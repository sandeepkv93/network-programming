package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Backend represents a backend server
type Backend struct {
	URL          *url.URL
	Alive        bool
	mu           sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
	Connections  int64
}

// SetAlive sets the alive status of the backend
func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.Alive = alive
	b.mu.Unlock()
}

// IsAlive returns true if backend is alive
func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	alive := b.Alive
	b.mu.RUnlock()
	return alive
}

// IncrementConnections increments the connection count
func (b *Backend) IncrementConnections() {
	atomic.AddInt64(&b.Connections, 1)
}

// DecrementConnections decrements the connection count
func (b *Backend) DecrementConnections() {
	atomic.AddInt64(&b.Connections, -1)
}

// GetConnections returns the current connection count
func (b *Backend) GetConnections() int64 {
	return atomic.LoadInt64(&b.Connections)
}

// LoadBalancer represents a load balancer
type LoadBalancer struct {
	address  string
	backends []*Backend
	current  uint64
	strategy string
	server   *http.Server
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(address string, strategy string) *LoadBalancer {
	return &LoadBalancer{
		address:  address,
		backends: []*Backend{},
		strategy: strategy,
	}
}

// AddBackend adds a backend server
func (lb *LoadBalancer) AddBackend(urlStr string) error {
	backendURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid backend URL: %v", err)
	}

	backend := &Backend{
		URL:          backendURL,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(backendURL),
	}

	lb.backends = append(lb.backends, backend)
	log.Printf("Added backend: %s\n", urlStr)
	return nil
}

// Start starts the load balancer
func (lb *LoadBalancer) Start() error {
	if len(lb.backends) == 0 {
		return fmt.Errorf("no backends configured")
	}

	// Start health checks
	go lb.healthCheck()

	handler := http.HandlerFunc(lb.handleRequest)

	lb.server = &http.Server{
		Addr:         lb.address,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Load Balancer listening on %s\n", lb.address)
	log.Printf("Strategy: %s\n", lb.strategy)

	go func() {
		if err := lb.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Load Balancer error: %v\n", err)
		}
	}()

	return nil
}

// handleRequest handles incoming requests
func (lb *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()
	if backend == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		log.Println("No available backends")
		return
	}

	log.Printf("Forwarding request to %s\n", backend.URL.String())

	backend.IncrementConnections()
	defer backend.DecrementConnections()

	backend.ReverseProxy.ServeHTTP(w, r)
}

// getNextBackend returns the next available backend
func (lb *LoadBalancer) getNextBackend() *Backend {
	switch lb.strategy {
	case "round-robin":
		return lb.roundRobin()
	case "least-connections":
		return lb.leastConnections()
	default:
		return lb.roundRobin()
	}
}

// roundRobin implements round-robin selection
func (lb *LoadBalancer) roundRobin() *Backend {
	n := len(lb.backends)
	if n == 0 {
		return nil
	}

	// Try all backends starting from current position
	for i := 0; i < n; i++ {
		idx := (atomic.AddUint64(&lb.current, 1) - 1) % uint64(n)
		backend := lb.backends[idx]
		if backend.IsAlive() {
			return backend
		}
	}

	return nil
}

// leastConnections implements least-connections selection
func (lb *LoadBalancer) leastConnections() *Backend {
	var selected *Backend
	minConnections := int64(-1)

	for _, backend := range lb.backends {
		if !backend.IsAlive() {
			continue
		}

		connections := backend.GetConnections()
		if minConnections == -1 || connections < minConnections {
			minConnections = connections
			selected = backend
		}
	}

	return selected
}

// healthCheck periodically checks backend health
func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, backend := range lb.backends {
			alive := lb.isBackendAlive(backend)
			backend.SetAlive(alive)

			status := "UP"
			if !alive {
				status = "DOWN"
			}
			log.Printf("Backend %s is %s (connections: %d)\n",
				backend.URL.String(), status, backend.GetConnections())
		}
	}
}

// isBackendAlive checks if a backend is alive
func (lb *LoadBalancer) isBackendAlive(backend *Backend) bool {
	timeout := 2 * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(backend.URL.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500
}

// GetBackendStats returns statistics for all backends
func (lb *LoadBalancer) GetBackendStats() []map[string]interface{} {
	stats := []map[string]interface{}{}

	for _, backend := range lb.backends {
		stats = append(stats, map[string]interface{}{
			"url":         backend.URL.String(),
			"alive":       backend.IsAlive(),
			"connections": backend.GetConnections(),
		})
	}

	return stats
}

// Stop stops the load balancer
func (lb *LoadBalancer) Stop() error {
	log.Println("Stopping Load Balancer...")
	if lb.server != nil {
		return lb.server.Close()
	}
	return nil
}
