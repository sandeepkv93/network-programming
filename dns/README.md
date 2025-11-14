## DNS Server

This package implements a simplified DNS (Domain Name System) server and client that handles A record queries.

### Features

- **DNS Server**: Responds to DNS queries for configured domains
- **DNS Client**: Performs DNS lookups
- A record (IPv4 address) support
- UDP-based communication (standard DNS protocol)
- Configurable domain-to-IP mappings

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"network-programming/dns"
)

func main() {
	server := dns.NewServer(":53")

	// Add DNS records
	server.AddRecord("example.com", "93.184.216.34")
	server.AddRecord("localhost", "127.0.0.1")
	server.AddRecord("myapp.local", "192.168.1.100")

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
	"network-programming/dns"
)

func main() {
	client := dns.NewClient("localhost:53")

	ip, err := client.Query("example.com")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("IP address: %s\n", ip)
}
```

### How it Works

1. **Server** listens on UDP port 53 (or custom port)
2. Client sends a DNS query packet with a domain name
3. Server parses the query and looks up the domain in its records
4. Server builds a DNS response with the IP address (or error)
5. Client receives and parses the response to extract the IP

### DNS Protocol

- DNS uses UDP on port 53 for queries (TCP for zone transfers)
- Query/Response format includes headers and resource records
- This implementation supports:
  - A records (IPv4 addresses)
  - Standard query/response format
  - Basic error handling
- Production DNS servers support many more record types (AAAA, MX, CNAME, etc.)

### Note

This is a simplified educational DNS server. Production DNS servers are much more complex and handle:
- Multiple record types (AAAA, MX, CNAME, NS, SOA, etc.)
- Recursive queries
- Caching
- Zone transfers
- DNSSEC
- Load distribution
