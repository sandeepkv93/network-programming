## Port Scanner

This package implements a concurrent TCP port scanner that can scan single ports, ranges, or common ports.

### Features

- **Single port scanning**: Check if a specific port is open
- **Range scanning**: Scan a range of ports (e.g., 1-1000)
- **Common ports**: Scan well-known service ports
- **Concurrent scanning**: Multiple ports scanned simultaneously
- **Banner grabbing**: Attempts to read service banners
- **Configurable timeout and threads**

### Usage

#### Scan a Single Port

```go
package main

import (
	"fmt"
	"network-programming/portscanner"
)

func main() {
	scanner := portscanner.NewScanner("localhost")

	result := scanner.ScanPort(80)
	if result.Open {
		fmt.Printf("Port %d is OPEN\n", result.Port)
		if result.Banner != "" {
			fmt.Printf("Banner: %s\n", result.Banner)
		}
	} else {
		fmt.Printf("Port %d is CLOSED\n", result.Port)
	}
}
```

#### Scan a Range of Ports

```go
package main

import (
	"fmt"
	"network-programming/portscanner"
)

func main() {
	scanner := portscanner.NewScanner("192.168.1.1")
	scanner.SetTimeout(2 * time.Second)
	scanner.SetThreads(50)

	results := scanner.ScanRange(1, 1000)

	fmt.Printf("Found %d open ports:\n", len(results))
	for _, result := range results {
		service := portscanner.GetServiceName(result.Port)
		fmt.Printf("Port %d (%s) is OPEN\n", result.Port, service)
	}
}
```

#### Scan Common Ports

```go
package main

import (
	"fmt"
	"network-programming/portscanner"
)

func main() {
	scanner := portscanner.NewScanner("example.com")

	results := scanner.ScanCommonPorts()

	for _, result := range results {
		service := portscanner.GetServiceName(result.Port)
		fmt.Printf("Port %d (%s) is OPEN\n", result.Port, service)
	}
}
```

### How it Works

1. Scanner attempts to establish a TCP connection to each port
2. If connection succeeds, the port is considered open
3. Scanner tries to read a banner (initial response from the service)
4. Multiple ports are scanned concurrently using goroutines
5. A semaphore limits the number of concurrent connections

### Configuration

- **Timeout**: Connection timeout (default: 1 second)
- **Threads**: Number of concurrent scan threads (default: 100)

### Common Ports Scanned

- 20, 21: FTP
- 22: SSH
- 23: Telnet
- 25: SMTP
- 53: DNS
- 80, 443: HTTP/HTTPS
- 110: POP3
- 143: IMAP
- 3306: MySQL
- 5432: PostgreSQL
- 6379: Redis
- 8080, 8443: HTTP/HTTPS Alternatives
- 27017: MongoDB

### Note

Port scanning may be illegal or against terms of service when performed on systems you don't own or have explicit permission to test. This tool is for educational purposes and authorized security testing only.
