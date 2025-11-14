package portscanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ScanResult represents the result of a port scan
type ScanResult struct {
	Port   int
	Open   bool
	Banner string
}

// Scanner represents a port scanner
type Scanner struct {
	host    string
	timeout time.Duration
	threads int
}

// NewScanner creates a new port scanner
func NewScanner(host string) *Scanner {
	return &Scanner{
		host:    host,
		timeout: 1 * time.Second,
		threads: 100,
	}
}

// SetTimeout sets the connection timeout
func (s *Scanner) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// SetThreads sets the number of concurrent scanning threads
func (s *Scanner) SetThreads(threads int) {
	s.threads = threads
}

// ScanPort scans a single port
func (s *Scanner) ScanPort(port int) ScanResult {
	address := fmt.Sprintf("%s:%d", s.host, port)
	conn, err := net.DialTimeout("tcp", address, s.timeout)

	result := ScanResult{
		Port: port,
		Open: false,
	}

	if err != nil {
		return result
	}

	result.Open = true

	// Try to read banner
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err == nil && n > 0 {
		result.Banner = string(buffer[:n])
	}

	conn.Close()
	return result
}

// ScanRange scans a range of ports
func (s *Scanner) ScanRange(startPort, endPort int) []ScanResult {
	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create a channel to limit concurrent scans
	semaphore := make(chan struct{}, s.threads)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			result := s.ScanPort(p)
			if result.Open {
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return results
}

// ScanCommonPorts scans commonly used ports
func (s *Scanner) ScanCommonPorts() []ScanResult {
	commonPorts := []int{
		20, 21,   // FTP
		22,       // SSH
		23,       // Telnet
		25,       // SMTP
		53,       // DNS
		80, 443,  // HTTP/HTTPS
		110,      // POP3
		143,      // IMAP
		3306,     // MySQL
		5432,     // PostgreSQL
		6379,     // Redis
		8080,     // HTTP Alt
		8443,     // HTTPS Alt
		27017,    // MongoDB
	}

	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, s.threads)

	for _, port := range commonPorts {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			result := s.ScanPort(p)
			if result.Open {
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return results
}

// GetServiceName returns the common service name for a port
func GetServiceName(port int) string {
	services := map[int]string{
		20:    "FTP-DATA",
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		143:   "IMAP",
		443:   "HTTPS",
		3306:  "MySQL",
		5432:  "PostgreSQL",
		6379:  "Redis",
		8080:  "HTTP-Alt",
		8443:  "HTTPS-Alt",
		27017: "MongoDB",
	}

	if service, ok := services[port]; ok {
		return service
	}
	return "Unknown"
}
