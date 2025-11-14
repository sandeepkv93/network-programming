package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server represents an HTTP server
type Server struct {
	address string
	server  *http.Server
	mux     *http.ServeMux
}

// Response represents a JSON response
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
}

// NewServer creates a new HTTP server
func NewServer(address string) *Server {
	mux := http.NewServeMux()
	return &Server{
		address: address,
		mux:     mux,
		server: &http.Server{
			Addr:         address,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupRoutes()
	log.Printf("HTTP Server listening on %s\n", s.address)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP Server error: %v\n", err)
		}
	}()

	return nil
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", s.handleRoot)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/echo", s.handleEcho)
	s.mux.HandleFunc("/time", s.handleTime)
}

// handleRoot handles requests to the root path
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	response := Response{
		Message:   "Welcome to the HTTP Server!",
		Timestamp: time.Now(),
		Path:      r.URL.Path,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Server is healthy",
		Timestamp: time.Now(),
		Path:      r.URL.Path,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleEcho echoes back the request information
func (s *Server) handleEcho(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message == "" {
		message = "No message provided"
	}

	response := Response{
		Message:   fmt.Sprintf("Echo: %s", message),
		Timestamp: time.Now(),
		Path:      r.URL.Path,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleTime returns the current server time
func (s *Server) handleTime(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   time.Now().Format(time.RFC3339),
		Timestamp: time.Now(),
		Path:      r.URL.Path,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// writeJSON writes a JSON response
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	log.Println("Stopping HTTP Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	log.Println("HTTP Server stopped")
	return nil
}

// RequestCounter middleware for counting requests
var (
	requestCount int
	countMutex   sync.Mutex
)

// GetRequestCount returns the current request count
func GetRequestCount() int {
	countMutex.Lock()
	defer countMutex.Unlock()
	return requestCount
}

// IncrementRequestCount increments the request counter
func IncrementRequestCount() {
	countMutex.Lock()
	defer countMutex.Unlock()
	requestCount++
}
