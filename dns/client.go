package dns

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents a DNS client
type Client struct {
	serverAddr string
	timeout    time.Duration
}

// NewClient creates a new DNS client
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
		timeout:    5 * time.Second,
	}
}

// Query performs a DNS query for a domain
func (c *Client) Query(domain string) (net.IP, error) {
	// Resolve server address
	addr, err := net.ResolveUDPAddr("udp", c.serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address: %v", err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DNS server: %v", err)
	}
	defer conn.Close()

	// Set deadline
	conn.SetDeadline(time.Now().Add(c.timeout))

	// Build query
	query := c.buildQuery(domain)

	// Send query
	_, err = conn.Write(query)
	if err != nil {
		return nil, fmt.Errorf("failed to send query: %v", err)
	}

	// Receive response
	buffer := make([]byte, 512)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response: %v", err)
	}

	// Parse response
	ip, err := c.parseResponse(buffer[:n])
	if err != nil {
		return nil, err
	}

	return ip, nil
}

// buildQuery builds a DNS query packet
func (c *Client) buildQuery(domain string) []byte {
	query := make([]byte, 512)

	// Transaction ID (random)
	binary.BigEndian.PutUint16(query[0:2], 0x1234)

	// Flags: Standard query
	binary.BigEndian.PutUint16(query[2:4], 0x0100)

	// Questions: 1
	binary.BigEndian.PutUint16(query[4:6], 1)

	// Answers: 0
	binary.BigEndian.PutUint16(query[6:8], 0)

	// Authority RRs: 0
	binary.BigEndian.PutUint16(query[8:10], 0)

	// Additional RRs: 0
	binary.BigEndian.PutUint16(query[10:12], 0)

	// Question
	offset := 12
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		query[offset] = byte(len(part))
		offset++
		copy(query[offset:], []byte(part))
		offset += len(part)
	}
	query[offset] = 0 // Null terminator
	offset++

	// Type A (1)
	binary.BigEndian.PutUint16(query[offset:offset+2], 1)
	offset += 2

	// Class IN (1)
	binary.BigEndian.PutUint16(query[offset:offset+2], 1)
	offset += 2

	return query[:offset]
}

// parseResponse parses a DNS response packet
func (c *Client) parseResponse(response []byte) (net.IP, error) {
	if len(response) < 12 {
		return nil, fmt.Errorf("response too short")
	}

	// Check response code
	flags := binary.BigEndian.Uint16(response[2:4])
	rcode := flags & 0x000F
	if rcode != 0 {
		return nil, fmt.Errorf("DNS error: code %d", rcode)
	}

	// Get answer count
	answerCount := binary.BigEndian.Uint16(response[6:8])
	if answerCount == 0 {
		return nil, fmt.Errorf("no answers in response")
	}

	// Skip question section
	offset := 12
	for offset < len(response) && response[offset] != 0 {
		length := int(response[offset])
		offset += length + 1
	}
	offset += 5 // Skip null terminator, type, and class

	// Parse answer
	if offset+12 > len(response) {
		return nil, fmt.Errorf("invalid response format")
	}

	// Skip name (assuming pointer)
	offset += 2

	// Check type
	recordType := binary.BigEndian.Uint16(response[offset : offset+2])
	offset += 2

	if recordType != 1 {
		return nil, fmt.Errorf("unexpected record type: %d", recordType)
	}

	// Skip class and TTL
	offset += 6

	// Get data length
	dataLen := binary.BigEndian.Uint16(response[offset : offset+2])
	offset += 2

	if dataLen != 4 {
		return nil, fmt.Errorf("unexpected IP length: %d", dataLen)
	}

	// Get IP address
	if offset+4 > len(response) {
		return nil, fmt.Errorf("invalid IP address in response")
	}

	ip := net.IPv4(response[offset], response[offset+1], response[offset+2], response[offset+3])
	return ip, nil
}

// SetTimeout sets the query timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}
