package ipscanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ScanResult represents the result of scanning an IP address
type ScanResult struct {
	IP        string
	IsAlive   bool
	Latency   time.Duration
	Hostname  string
	Timestamp time.Time
}

// Scanner represents an IP scanner
type Scanner struct {
	timeout     time.Duration
	concurrency int
}

// NewScanner creates a new IP scanner
func NewScanner(timeout time.Duration, concurrency int) *Scanner {
	if timeout == 0 {
		timeout = 1 * time.Second
	}
	if concurrency == 0 {
		concurrency = 100
	}
	return &Scanner{
		timeout:     timeout,
		concurrency: concurrency,
	}
}

// ScanIP scans a single IP address to check if it's alive
func (s *Scanner) ScanIP(ip string) ScanResult {
	result := ScanResult{
		IP:        ip,
		IsAlive:   false,
		Timestamp: time.Now(),
	}

	start := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", ip, s.timeout)
	if err != nil {
		// Fallback to TCP ping on port 80
		conn, err = net.DialTimeout("tcp", net.JoinHostPort(ip, "80"), s.timeout)
		if err != nil {
			return result
		}
	}
	defer conn.Close()

	result.IsAlive = true
	result.Latency = time.Since(start)

	// Try to resolve hostname
	if names, err := net.LookupAddr(ip); err == nil && len(names) > 0 {
		result.Hostname = names[0]
	}

	return result
}

// ScanRange scans a range of IP addresses
func (s *Scanner) ScanRange(startIP, endIP string) []ScanResult {
	start := parseIP(startIP)
	end := parseIP(endIP)

	if start == 0 || end == 0 || start > end {
		return nil
	}

	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create a semaphore channel to limit concurrency
	semaphore := make(chan struct{}, s.concurrency)

	for ip := start; ip <= end; ip++ {
		wg.Add(1)
		go func(ipNum uint32) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			ipStr := formatIP(ipNum)
			result := s.ScanIP(ipStr)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(ip)
	}

	wg.Wait()
	return results
}

// ScanSubnet scans all hosts in a subnet (CIDR notation)
func (s *Scanner) ScanSubnet(cidr string) ([]ScanResult, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %v", err)
	}

	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, s.concurrency)

	// Iterate through all IPs in the subnet
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		ipStr := ip.String()
		wg.Add(1)
		go func(ipAddr string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := s.ScanIP(ipAddr)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(ipStr)
	}

	wg.Wait()
	return results, nil
}

// parseIP converts an IP string to a uint32
func parseIP(ipStr string) uint32 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// formatIP converts a uint32 to an IP string
func formatIP(ipNum uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ipNum>>24)&0xFF,
		(ipNum>>16)&0xFF,
		(ipNum>>8)&0xFF,
		ipNum&0xFF)
}

// incIP increments an IP address
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// GetAliveHosts filters and returns only alive hosts from scan results
func GetAliveHosts(results []ScanResult) []ScanResult {
	var alive []ScanResult
	for _, result := range results {
		if result.IsAlive {
			alive = append(alive, result)
		}
	}
	return alive
}

// PrintResults prints scan results in a formatted way
func PrintResults(results []ScanResult) {
	fmt.Println("\n=== IP Scanner Results ===")
	fmt.Printf("Total scanned: %d\n", len(results))

	aliveCount := 0
	for _, result := range results {
		if result.IsAlive {
			aliveCount++
			hostname := result.Hostname
			if hostname == "" {
				hostname = "N/A"
			}
			fmt.Printf("[ALIVE] %s - Latency: %v - Hostname: %s\n",
				result.IP, result.Latency, hostname)
		}
	}

	fmt.Printf("\nAlive hosts: %d/%d\n", aliveCount, len(results))
}
