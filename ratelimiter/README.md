## Rate Limiter

A flexible rate limiting implementation supporting multiple algorithms and protocols (HTTP, TCP) to control request rates and prevent abuse.

## Features

- Fixed window rate limiting
- Token bucket algorithm
- HTTP middleware support
- TCP connection rate limiting
- Per-client tracking
- Automatic cleanup
- Configurable limits and windows

## Algorithms

### Fixed Window
Limits requests per time window (e.g., 100 requests per minute).

### Token Bucket
Allows burst traffic while maintaining average rate (e.g., capacity of 100, refill 10/sec).

## Usage

### Fixed Window Rate Limiter
```go
// Allow 100 requests per minute per client
limiter := ratelimiter.NewRateLimiter(100, 1*time.Minute)

if limiter.Allow("client-123") {
    // Process request
} else {
    // Reject request - rate limit exceeded
}
```

### HTTP Middleware
```go
limiter := ratelimiter.NewRateLimiter(100, 1*time.Minute)

http.HandleFunc("/api", limiter.HTTPMiddleware(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API response"))
}))

http.ListenAndServe(":8080", nil)
```

### TCP Rate Limiting
```go
listener, _ := net.Listen("tcp", ":9000")
rateLimitedListener := ratelimiter.NewTCPRateLimiter(listener, 10, 1*time.Second)

for {
    conn, err := rateLimitedListener.Accept()
    if err != nil {
        continue // Rate limited or error
    }
    go handleConnection(conn)
}
```

### Token Bucket
```go
bucket := ratelimiter.NewTokenBucket(100, 10) // capacity 100, refill 10/sec

if bucket.Take(1) {
    // Request allowed
} else {
    // Rate limited
}
```

## Use Cases

- API rate limiting
- DDoS protection
- Resource usage control
- Fair usage enforcement
- Traffic shaping
