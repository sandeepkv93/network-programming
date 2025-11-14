package cdn

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// OriginServer represents a CDN origin server
type OriginServer struct {
	Address string
	server  *http.Server
}

// NewOriginServer creates a new origin server
func NewOriginServer(address string) *OriginServer {
	return &OriginServer{
		Address: address,
	}
}

// Start starts the origin server
func (o *OriginServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", o.handleRequest)

	o.server = &http.Server{
		Addr:    o.Address,
		Handler: mux,
	}

	log.Printf("CDN Origin Server starting on %s\n", o.Address)
	return o.server.ListenAndServe()
}

// Stop stops the origin server
func (o *OriginServer) Stop() error {
	if o.server != nil {
		return o.server.Close()
	}
	return nil
}

func (o *OriginServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Origin] Request: %s from %s\n", r.URL.Path, r.RemoteAddr)

	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	// Serve content based on path
	path := r.URL.Path

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Generate simple content
	content := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>CDN Test Page</title>
</head>
<body>
    <h1>CDN Origin Server</h1>
    <p>Path: %s</p>
    <p>Time: %s</p>
    <p>This content is served from the origin server.</p>
</body>
</html>`, path, time.Now().Format(time.RFC3339))

	fmt.Fprint(w, content)
}
