package ratelimiter

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	requests map[string]*ClientLimitInfo
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
	cleanup  time.Duration // cleanup interval
}

// ClientLimitInfo tracks rate limit info for a client
type ClientLimitInfo struct {
	Count      int
	ResetTime  time.Time
	Blocked    bool
	TotalReqs  int
	BlockedReqs int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerWindow int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*ClientLimitInfo),
		rate:     requestsPerWindow,
		window:   window,
		cleanup:  window * 2,
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// Allow checks if a request from the client should be allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.requests[clientID]

	if !exists {
		rl.requests[clientID] = &ClientLimitInfo{
			Count:     1,
			ResetTime: now.Add(rl.window),
			TotalReqs: 1,
		}
		return true
	}

	info.TotalReqs++

	// Reset window if expired
	if now.After(info.ResetTime) {
		info.Count = 1
		info.ResetTime = now.Add(rl.window)
		info.Blocked = false
		return true
	}

	// Check rate limit
	if info.Count >= rl.rate {
		info.Blocked = true
		info.BlockedReqs++
		return false
	}

	info.Count++
	return true
}

// GetClientInfo returns rate limit info for a client
func (rl *RateLimiter) GetClientInfo(clientID string) *ClientLimitInfo {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.requests[clientID]
}

// GetAllClients returns all tracked clients
func (rl *RateLimiter) GetAllClients() map[string]*ClientLimitInfo {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	result := make(map[string]*ClientLimitInfo)
	for k, v := range rl.requests {
		result[k] = v
	}
	return result
}

// ResetClient resets rate limit for a specific client
func (rl *RateLimiter) ResetClient(clientID string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.requests, clientID)
}

// cleanupLoop periodically cleans up old entries
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup removes expired entries
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for clientID, info := range rl.requests {
		if now.After(info.ResetTime.Add(rl.window)) {
			delete(rl.requests, clientID)
		}
	}
}

// HTTPMiddleware returns an HTTP middleware for rate limiting
func (rl *RateLimiter) HTTPMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use IP address as client ID
		clientID := getClientIP(r)

		if !rl.Allow(clientID) {
			info := rl.GetClientInfo(clientID)
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.rate))
			w.Header().Set("X-RateLimit-Reset", info.ResetTime.Format(time.RFC3339))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded\n"))
			log.Printf("Rate limit exceeded for %s\n", clientID)
			return
		}

		next(w, r)
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Use remote address
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// TCPRateLimiter wraps a TCP listener with rate limiting
type TCPRateLimiter struct {
	listener net.Listener
	limiter  *RateLimiter
}

// NewTCPRateLimiter creates a rate-limited TCP listener
func NewTCPRateLimiter(listener net.Listener, requestsPerWindow int, window time.Duration) *TCPRateLimiter {
	return &TCPRateLimiter{
		listener: listener,
		limiter:  NewRateLimiter(requestsPerWindow, window),
	}
}

// Accept accepts a connection with rate limiting
func (trl *TCPRateLimiter) Accept() (net.Conn, error) {
	conn, err := trl.listener.Accept()
	if err != nil {
		return nil, err
	}

	// Extract IP from connection
	ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

	if !trl.limiter.Allow(ip) {
		log.Printf("Rate limit exceeded for TCP connection from %s\n", ip)
		conn.Close()
		return nil, fmt.Errorf("rate limit exceeded")
	}

	return conn, nil
}

// Close closes the listener
func (trl *TCPRateLimiter) Close() error {
	return trl.listener.Close()
}

// Addr returns the listener's network address
func (trl *TCPRateLimiter) Addr() net.Addr {
	return trl.listener.Addr()
}

// TokenBucket implements token bucket rate limiting algorithm
type TokenBucket struct {
	capacity  int
	tokens    int
	refillRate int // tokens per second
	mu        sync.Mutex
	lastRefill time.Time
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	tb := &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}

	// Start refill goroutine
	go tb.refillLoop()

	return tb
}

// Take attempts to take n tokens from the bucket
func (tb *TokenBucket) Take(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}

	return false
}

// refill refills the bucket based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
}

// refillLoop continuously refills the bucket
func (tb *TokenBucket) refillLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		tb.mu.Lock()
		tb.refill()
		tb.mu.Unlock()
	}
}

// GetTokens returns the current number of tokens
func (tb *TokenBucket) GetTokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}
