## HTTP Client & Server

This package implements an HTTP/1.1 client and server using Go's standard `net/http` package.

### Features

- **HTTP Server**: Serves HTTP requests with multiple endpoints
- **HTTP Client**: Makes HTTP requests (GET, POST)
- JSON response handling
- Request routing
- Graceful shutdown
- Configurable timeouts

### Server Endpoints

- `GET /` - Welcome message
- `GET /health` - Health check endpoint
- `GET /echo?message=<msg>` - Echo back a message
- `GET /time` - Returns current server time

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"network-programming/http"
)

func main() {
	server := http.NewServer(":8080")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Keep server running
	time.Sleep(1 * time.Minute)

	server.Stop()
}
```

#### Client

```go
package main

import (
	"log"
	"network-programming/http"
)

func main() {
	client := http.NewClient("http://localhost:8080")

	// Make a GET request
	response, err := client.Get("/health")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %+v\n", response)
}
```

### How it Works

1. **Server** listens for HTTP requests on a specified port
2. Requests are routed to appropriate handlers based on the path
3. **Client** makes HTTP requests using the standard library
4. Responses are formatted as JSON with timestamps
5. Both client and server support configurable timeouts

### HTTP Protocol

- HTTP is an application-layer protocol built on top of TCP
- It uses a request-response model
- Common methods: GET, POST, PUT, DELETE, PATCH
- Status codes indicate the result (200 OK, 404 Not Found, etc.)
