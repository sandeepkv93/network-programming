## Traceroute Implementation

This package implements the traceroute network diagnostic tool that tracks the path packets take to a destination.

### Features

- ICMP-based traceroute implementation
- Track all hops to destination
- Measure round-trip time (RTT) for each hop
- Handle timeouts gracefully
- Support for retry attempts
- IPv4 support

### Usage

#### Basic Traceroute

```go
package main

import (
	"fmt"
	"log"
	"time"
	"networkprogramming/traceroute"
)

func main() {
	result, err := traceroute.Traceroute("8.8.8.8", 30, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(traceroute.FormatResult(result))
}
```

#### Traceroute with Retries

```go
package main

import (
	"fmt"
	"log"
	"time"
	"networkprogramming/traceroute"
)

func main() {
	// 30 max hops, 3 second timeout, 3 retries per hop
	result, err := traceroute.TraceWithRetries("google.com", 30, 3*time.Second, 3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(traceroute.FormatResult(result))
}
```

### Output Format

```
traceroute to 8.8.8.8, 10 hops max
 1  192.168.1.1  2.345ms
 2  10.0.0.1  5.123ms
 3  172.16.0.1  8.456ms
 4  * * * Request timeout
 5  8.8.8.8  15.789ms
```

### How it Works

1. Sends ICMP Echo Request packets with incrementing TTL (Time To Live)
2. Each router along the path decrements TTL
3. When TTL reaches 0, router sends back ICMP Time Exceeded message
4. This reveals the router's IP address
5. Process continues until destination is reached or max hops exceeded
6. Measures RTT for each hop

### Parameters

- **destination**: Target hostname or IP address
- **maxHops**: Maximum number of hops to try (default: 30)
- **timeout**: How long to wait for each response
- **retries**: Number of attempts per hop (for TraceWithRetries)

### Requirements

- Requires raw socket access (may need root/admin privileges)
- Uses `golang.org/x/net/icmp` and `golang.org/x/net/ipv4` packages

```
go get golang.org/x/net/icmp
go get golang.org/x/net/ipv4
```

### Use Cases

- Network path diagnostics
- Identifying routing issues
- Measuring network latency at each hop
- Debugging connectivity problems
- Network topology discovery

### Limitations

- Requires elevated privileges to create raw ICMP sockets
- Some routers may not respond to ICMP (firewalls)
- IPv4 only in this implementation
- May be blocked by security policies
