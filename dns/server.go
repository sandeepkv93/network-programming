package dns

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// Record represents a DNS record
type Record struct {
	Domain string
	IP     net.IP
}

// Server represents a DNS server
type Server struct {
	address string
	conn    *net.UDPConn
	records map[string]net.IP
	quit    chan bool
	wg      sync.WaitGroup
	mu      sync.RWMutex
}

// NewServer creates a new DNS server
func NewServer(address string) *Server {
	return &Server{
		address: address,
		records: make(map[string]net.IP),
		quit:    make(chan bool),
	}
}

// AddRecord adds a DNS record
func (s *Server) AddRecord(domain string, ip string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}

	s.records[strings.ToLower(domain)] = parsedIP.To4()
	log.Printf("Added DNS record: %s -> %s\n", domain, ip)
	return nil
}

// Start starts the DNS server
func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to start DNS server: %v", err)
	}
	s.conn = conn
	log.Printf("DNS Server listening on %s\n", s.address)

	s.wg.Add(1)
	go s.handleQueries()

	return nil
}

// handleQueries handles incoming DNS queries
func (s *Server) handleQueries() {
	defer s.wg.Done()

	buffer := make([]byte, 512) // Standard DNS message size

	for {
		select {
		case <-s.quit:
			return
		default:
			n, clientAddr, err := s.conn.ReadFromUDP(buffer)
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					log.Printf("Error reading UDP: %v\n", err)
					continue
				}
			}

			if n > 0 {
				s.wg.Add(1)
				go s.processQuery(buffer[:n], clientAddr)
			}
		}
	}
}

// processQuery processes a single DNS query
func (s *Server) processQuery(query []byte, clientAddr *net.UDPAddr) {
	defer s.wg.Done()

	if len(query) < 12 {
		log.Println("Query too short")
		return
	}

	// Parse query
	transactionID := binary.BigEndian.Uint16(query[0:2])
	domain := s.parseDomainName(query[12:])

	log.Printf("DNS query from %s for domain: %s\n", clientAddr, domain)

	// Look up IP address
	s.mu.RLock()
	ip, exists := s.records[strings.ToLower(domain)]
	s.mu.RUnlock()

	// Build response
	response := s.buildResponse(transactionID, domain, ip, exists)

	// Send response
	_, err := s.conn.WriteToUDP(response, clientAddr)
	if err != nil {
		log.Printf("Error sending response: %v\n", err)
	}
}

// parseDomainName parses a domain name from DNS query format
func (s *Server) parseDomainName(data []byte) string {
	var domain strings.Builder
	i := 0

	for i < len(data) && data[i] != 0 {
		length := int(data[i])
		if length == 0 {
			break
		}
		i++

		if domain.Len() > 0 {
			domain.WriteByte('.')
		}

		if i+length <= len(data) {
			domain.Write(data[i : i+length])
			i += length
		} else {
			break
		}
	}

	return domain.String()
}

// buildResponse builds a DNS response
func (s *Server) buildResponse(transactionID uint16, domain string, ip net.IP, exists bool) []byte {
	response := make([]byte, 512)

	// Transaction ID
	binary.BigEndian.PutUint16(response[0:2], transactionID)

	// Flags: Standard query response
	if exists {
		binary.BigEndian.PutUint16(response[2:4], 0x8180) // Response, no error
	} else {
		binary.BigEndian.PutUint16(response[2:4], 0x8183) // Response, name error
	}

	// Questions: 1
	binary.BigEndian.PutUint16(response[4:6], 1)

	// Answers
	if exists {
		binary.BigEndian.PutUint16(response[6:8], 1)
	} else {
		binary.BigEndian.PutUint16(response[6:8], 0)
	}

	// Authority RRs: 0
	binary.BigEndian.PutUint16(response[8:10], 0)

	// Additional RRs: 0
	binary.BigEndian.PutUint16(response[10:12], 0)

	// Question section
	offset := 12
	offset = s.encodeDomainName(response, offset, domain)
	binary.BigEndian.PutUint16(response[offset:offset+2], 1)   // Type A
	binary.BigEndian.PutUint16(response[offset+2:offset+4], 1) // Class IN
	offset += 4

	// Answer section (if exists)
	if exists && ip != nil {
		// Name pointer to question
		binary.BigEndian.PutUint16(response[offset:offset+2], 0xC00C)
		offset += 2

		// Type A
		binary.BigEndian.PutUint16(response[offset:offset+2], 1)
		offset += 2

		// Class IN
		binary.BigEndian.PutUint16(response[offset:offset+2], 1)
		offset += 2

		// TTL (300 seconds)
		binary.BigEndian.PutUint32(response[offset:offset+4], 300)
		offset += 4

		// Data length (4 bytes for IPv4)
		binary.BigEndian.PutUint16(response[offset:offset+2], 4)
		offset += 2

		// IP address
		copy(response[offset:offset+4], ip)
		offset += 4
	}

	return response[:offset]
}

// encodeDomainName encodes a domain name in DNS format
func (s *Server) encodeDomainName(buffer []byte, offset int, domain string) int {
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		buffer[offset] = byte(len(part))
		offset++
		copy(buffer[offset:], []byte(part))
		offset += len(part)
	}
	buffer[offset] = 0 // Null terminator
	offset++
	return offset
}

// Stop stops the DNS server
func (s *Server) Stop() {
	close(s.quit)
	if s.conn != nil {
		s.conn.Close()
	}
	s.wg.Wait()
	log.Println("DNS Server stopped")
}
