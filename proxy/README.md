## Proxy Server

This package implements an HTTP/HTTPS proxy server that forwards requests to target servers.

### Features

- **HTTP Proxy**: Forwards HTTP requests to target servers
- **HTTPS Support**: Handles CONNECT method for HTTPS tunneling
- **Request/Response forwarding**: Preserves headers and body
- **Statistics tracking**: Tracks requests and data transfer
- **Connection hijacking**: Supports HTTPS tunneling

### Usage

#### Server

```go
package main

import (
	"fmt"
	"log"
	"time"
	"network-programming/proxy"
)

func main() {
	server := proxy.NewServer(":8888")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Run for some time
	time.Sleep(1 * time.Minute)

	// Get statistics
	totalReqs, bytesIn, bytesOut := server.GetStats()
	fmt.Printf("Total requests: %d\n", totalReqs)
	fmt.Printf("Bytes in: %d, Bytes out: %d\n", bytesIn, bytesOut)

	server.Stop()
}
```

#### Using the Proxy

Configure your HTTP client to use the proxy:

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	// Configure proxy
	proxyURL, _ := url.Parse("http://localhost:8888")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// Make request through proxy
	resp, err := client.Get("http://example.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

Or use environment variables:
```bash
export HTTP_PROXY=http://localhost:8888
export HTTPS_PROXY=http://localhost:8888
curl http://example.com
```

### How it Works

#### HTTP Requests
1. Client sends HTTP request to the proxy
2. Proxy creates a new request to the target server
3. Proxy forwards headers and body to the target
4. Target server responds to the proxy
5. Proxy forwards the response back to the client

#### HTTPS Requests (CONNECT Tunneling)
1. Client sends CONNECT request to the proxy
2. Proxy establishes a TCP connection to the target server
3. Proxy hijacks the client connection
4. Proxy responds with "200 Connection Established"
5. Proxy sets up bidirectional forwarding between client and target
6. All data is forwarded transparently (encrypted end-to-end)

### Features

- **Transparent Proxying**: Forwards requests without modification
- **Header Preservation**: Maintains all HTTP headers
- **HTTPS Tunneling**: Supports encrypted HTTPS connections via CONNECT
- **Statistics**: Tracks number of requests and data transfer
- **Timeout Handling**: Configurable timeouts for connections

### Use Cases

- Web traffic monitoring and logging
- Content filtering and blocking
- Caching (not implemented in this basic version)
- Load distribution
- Security scanning
- Development and debugging

### Note

This is a basic educational proxy server. Production proxy servers typically include:
- Authentication and authorization
- Content caching
- Access control lists (ACLs)
- SSL/TLS interception (for HTTPS inspection)
- Bandwidth throttling
- Request/response modification
- Logging and analytics
