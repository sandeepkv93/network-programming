## IP Scanner

The IP Scanner package provides functionality to scan IP addresses and subnets to discover active hosts on a network.

## Table of Contents

1. [What is IP Scanning?](#what-is-ip-scanning)
2. [How Does It Work?](#how-does-it-work)
3. [Understanding the Code](#understanding-the-code)
4. [Usage Examples](#usage-examples)
5. [Further Reading](#further-reading)

### What is IP Scanning?

IP scanning is a network reconnaissance technique used to identify active hosts on a network. It involves sending packets to IP addresses and analyzing responses to determine if a host is alive and reachable.

**Common Use Cases**:
- Network inventory and asset discovery
- Network troubleshooting
- Security auditing
- Monitoring network availability

### How Does It Work?

The IP scanner uses multiple techniques to detect live hosts:

1. **ICMP Ping**: Attempts to establish an ICMP connection to the target IP
2. **TCP Ping**: Falls back to TCP connection on port 80 if ICMP fails
3. **Hostname Resolution**: Attempts reverse DNS lookup for discovered hosts
4. **Concurrent Scanning**: Uses goroutines to scan multiple IPs in parallel

**Scanning Process**:
1. Parse the IP range or subnet to scan
2. Create a pool of workers limited by concurrency settings
3. For each IP, attempt connection with timeout
4. Record results including latency and hostname
5. Return summary of alive and dead hosts

### Understanding the Code

#### Data Structures:

- `Scanner`: Main scanner with configurable timeout and concurrency
- `ScanResult`: Contains scan results including IP, status, latency, and hostname

#### Key Functions:

- `NewScanner(timeout, concurrency)`: Creates a scanner with specified settings
- `ScanIP(ip)`: Scans a single IP address
- `ScanRange(startIP, endIP)`: Scans a range of IP addresses
- `ScanSubnet(cidr)`: Scans all hosts in a CIDR subnet
- `GetAliveHosts(results)`: Filters results to show only active hosts
- `PrintResults(results)`: Displays formatted scan results

### Usage Examples

#### Scanning a Single IP:
```go
scanner := ipscanner.NewScanner(1*time.Second, 100)
result := scanner.ScanIP("192.168.1.1")
if result.IsAlive {
    fmt.Printf("Host %s is alive\n", result.IP)
}
```

#### Scanning an IP Range:
```go
scanner := ipscanner.NewScanner(1*time.Second, 100)
results := scanner.ScanRange("192.168.1.1", "192.168.1.254")
ipscanner.PrintResults(results)
```

#### Scanning a Subnet:
```go
scanner := ipscanner.NewScanner(1*time.Second, 100)
results, err := scanner.ScanSubnet("192.168.1.0/24")
if err == nil {
    aliveHosts := ipscanner.GetAliveHosts(results)
    fmt.Printf("Found %d alive hosts\n", len(aliveHosts))
}
```

### Further Reading

- [ICMP Protocol - RFC 792](https://datatracker.ietf.org/doc/html/rfc792)
- [Network Scanning Techniques](https://nmap.org/book/man-host-discovery.html)
- [Go net package](https://pkg.go.dev/net)
