## DHCP Implementation

This package implements a simplified DHCP (Dynamic Host Configuration Protocol) server and client for automatic IP address assignment.

### Features

- **DHCP Server**: Assigns IP addresses from a configured pool
- **DHCP Client**: Requests IP address from DHCP server
- IP address pool management
- Lease management with expiration
- Support for DISCOVER, OFFER, REQUEST, ACK messages
- Subnet mask, router, and DNS server configuration

### Usage

#### DHCP Server

```go
package main

import (
	"log"
	"time"
	"networkprogramming/dhcp"
)

func main() {
	// Create server with IP pool 192.168.1.100-200
	server, err := dhcp.NewServer(
		":67",                  // Address
		"192.168.1.100",       // Pool start
		"192.168.1.200",       // Pool end
		"255.255.255.0",       // Subnet mask
		"192.168.1.1",         // Router
		"8.8.8.8",             // DNS server
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Keep server running
	time.Sleep(10 * time.Minute)

	server.Stop()
}
```

#### DHCP Client

```go
package main

import (
	"fmt"
	"log"
	"time"
	"networkprogramming/dhcp"
)

func main() {
	client, err := dhcp.NewClient("eth0")
	if err != nil {
		log.Fatal(err)
	}

	// Get IP address
	ip, err := client.GetIP(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Assigned IP: %s\n", ip)
	fmt.Printf("DHCP Server: %s\n", client.GetServerIP())

	// Release when done
	defer client.Release()
}
```

#### Get Active Leases

```go
package main

import (
	"fmt"
	"networkprogramming/dhcp"
)

func main() {
	// Assuming server is already running
	server, _ := dhcp.NewServer(":67", "192.168.1.100", "192.168.1.200",
		"255.255.255.0", "192.168.1.1", "8.8.8.8")

	leases := server.GetLeases()

	fmt.Println("Active leases:")
	for _, lease := range leases {
		fmt.Printf("IP: %s, MAC: %s, Expires: %s\n",
			lease.IP, lease.MAC, lease.Expiry)
	}
}
```

### DHCP Process

1. **DISCOVER**: Client broadcasts request for IP address
2. **OFFER**: Server offers an available IP address
3. **REQUEST**: Client requests the offered IP address
4. **ACK**: Server acknowledges and assigns the IP address

```
Client                          Server
  |                               |
  |-------- DISCOVER ------------>|
  |                               |
  |<-------- OFFER --------------|
  |                               |
  |-------- REQUEST ------------>|
  |                               |
  |<-------- ACK ----------------|
  |                               |
```

### Server Configuration

The DHCP server manages:
- **IP Pool**: Range of available IP addresses
- **Leases**: Active IP assignments with expiration
- **Network Config**: Subnet mask, router, DNS servers
- **Lease Time**: Duration of IP assignment (default: 1 hour)

### Client Operations

The DHCP client can:
- Send DISCOVER messages
- Receive and parse OFFER messages
- Send REQUEST messages
- Retrieve assigned IP and server information
- Release leases

### Lease Management

- Leases are stored with MAC address as key
- Automatic lease expiration after configured time
- IP addresses returned to pool when leases expire
- Existing clients can renew leases

### Message Types

- `DHCPDiscover (1)`: Client searches for servers
- `DHCPOffer (2)`: Server offers configuration
- `DHCPRequest (3)`: Client requests configuration
- `DHCPAck (4)`: Server confirms configuration
- `DHCPNak (5)`: Server denies configuration
- `DHCPRelease (6)`: Client releases IP address

### Notes

- This is a simplified educational implementation
- Production DHCP servers should implement full RFC 2131
- Requires elevated privileges to bind to port 67
- Uses UDP broadcast for discovery
- Does not implement all DHCP options
- No persistent lease storage

### Requirements

- Port 67 (server) and 68 (client) access
- May require root/administrator privileges
- Network interface configuration

### Limitations

- Simplified packet format (educational purposes)
- Limited DHCP options support
- No DHCP relay agent support
- No conflict detection
- Basic lease management
- IPv4 only

### Security Considerations

- DHCP has no built-in authentication
- Vulnerable to rogue DHCP servers
- Clients trust any DHCP response
- Production use requires additional security measures:
  - DHCP snooping on switches
  - Port security
  - Network segmentation
  - Monitoring for rogue servers
