package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// Server represents a proxy server
type Server struct {
	address string
	server  *http.Server
	stats   *Stats
}

// Stats tracks proxy statistics
type Stats struct {
	mu            sync.RWMutex
	TotalRequests int64
	BytesIn       int64
	BytesOut      int64
}

// NewServer creates a new proxy server
func NewServer(address string) *Server {
	return &Server{
		address: address,
		stats:   &Stats{},
	}
}

// Start starts the proxy server
func (s *Server) Start() error {
	handler := http.HandlerFunc(s.handleRequest)

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Proxy Server listening on %s\n", s.address)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Proxy Server error: %v\n", err)
		}
	}()

	return nil
}

// handleRequest handles HTTP requests
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.stats.mu.Lock()
	s.stats.TotalRequests++
	s.stats.mu.Unlock()

	log.Printf("Proxy request: %s %s from %s\n", r.Method, r.URL.String(), r.RemoteAddr)

	// Handle CONNECT method for HTTPS
	if r.Method == http.MethodConnect {
		s.handleConnect(w, r)
		return
	}

	// Create a new request to the target server
	targetURL := r.URL.String()
	if r.URL.Scheme == "" {
		targetURL = "http://" + r.Host + r.URL.Path
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		log.Printf("Error creating proxy request: %v\n", err)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Make the request
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to reach target server", http.StatusBadGateway)
		log.Printf("Error proxying request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	bytesWritten, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response: %v\n", err)
		return
	}

	s.stats.mu.Lock()
	s.stats.BytesOut += bytesWritten
	s.stats.mu.Unlock()
}

// handleConnect handles CONNECT requests for HTTPS tunneling
func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	// Connect to the target server
	targetConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, "Failed to connect to target", http.StatusBadGateway)
		log.Printf("Error connecting to target: %v\n", err)
		return
	}
	defer targetConn.Close()

	// Hijack the client connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Failed to hijack connection", http.StatusInternalServerError)
		log.Printf("Error hijacking connection: %v\n", err)
		return
	}
	defer clientConn.Close()

	// Send 200 Connection Established
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// Bidirectional copy
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(targetConn, clientConn)
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientConn, targetConn)
	}()

	wg.Wait()
}

// GetStats returns proxy statistics
func (s *Server) GetStats() (totalRequests int64, bytesIn int64, bytesOut int64) {
	s.stats.mu.RLock()
	defer s.stats.mu.RUnlock()
	return s.stats.TotalRequests, s.stats.BytesIn, s.stats.BytesOut
}

// Stop stops the proxy server
func (s *Server) Stop() error {
	log.Println("Stopping Proxy Server...")
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}
