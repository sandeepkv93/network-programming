package ipspoofing

import (
	"encoding/binary"
	"fmt"
	"net"
)

// SECURITY WARNING: IP Spoofing is illegal in many jurisdictions when used maliciously.
// This implementation is for educational purposes only.
// Only use in authorized testing environments (pentesting, CTF, security research).

// IPSpoofing represents an IP spoofing utility
type IPSpoofing struct {
	SourceIP string
	DestIP   string
}

// NewIPSpoofing creates a new IP spoofing utility
func NewIPSpoofing(sourceIP, destIP string) *IPSpoofing {
	return &IPSpoofing{
		SourceIP: sourceIP,
		DestIP:   destIP,
	}
}

// CreateSpoofedPacket creates a packet with spoofed source IP
// Note: This requires raw socket access and elevated privileges
func (s *IPSpoofing) CreateSpoofedPacket(payload []byte) ([]byte, error) {
	// Parse IP addresses
	srcIP := net.ParseIP(s.SourceIP)
	dstIP := net.ParseIP(s.DestIP)

	if srcIP == nil || dstIP == nil {
		return nil, fmt.Errorf("invalid IP address")
	}

	// Convert to 4-byte representation
	srcIP = srcIP.To4()
	dstIP = dstIP.To4()

	if srcIP == nil || dstIP == nil {
		return nil, fmt.Errorf("only IPv4 is supported")
	}

	// Create IP header (20 bytes minimum)
	ipHeader := make([]byte, 20)

	// Version (4) and IHL (5) = 0x45
	ipHeader[0] = 0x45

	// Type of Service
	ipHeader[1] = 0

	// Total Length (header + payload)
	totalLen := 20 + len(payload)
	binary.BigEndian.PutUint16(ipHeader[2:4], uint16(totalLen))

	// Identification
	binary.BigEndian.PutUint16(ipHeader[4:6], 54321)

	// Flags and Fragment Offset
	binary.BigEndian.PutUint16(ipHeader[6:8], 0)

	// TTL
	ipHeader[8] = 64

	// Protocol (6 = TCP, 17 = UDP, 1 = ICMP)
	ipHeader[9] = 17 // UDP

	// Header Checksum (will be calculated)
	ipHeader[10] = 0
	ipHeader[11] = 0

	// Source IP
	copy(ipHeader[12:16], srcIP)

	// Destination IP
	copy(ipHeader[16:20], dstIP)

	// Calculate checksum
	checksum := calculateChecksum(ipHeader)
	binary.BigEndian.PutUint16(ipHeader[10:12], checksum)

	// Combine header and payload
	packet := append(ipHeader, payload...)

	return packet, nil
}

// SendSpoofedPacket sends a packet with spoofed source IP
// WARNING: Requires raw socket privileges (root/admin)
func (s *IPSpoofing) SendSpoofedPacket(payload []byte) error {
	packet, err := s.CreateSpoofedPacket(payload)
	if err != nil {
		return err
	}

	// Create raw socket
	// Note: This will fail without proper privileges
	conn, err := net.Dial("ip4:udp", s.DestIP)
	if err != nil {
		return fmt.Errorf("failed to create socket (requires elevated privileges): %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
	}

	return nil
}

func calculateChecksum(data []byte) uint16 {
	sum := uint32(0)

	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}

	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}

	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}

	return ^uint16(sum)
}

// DetectSpoofing provides basic IP spoofing detection
func DetectSpoofing(srcIP string, expectedNetwork string) bool {
	// Simple check: is the source IP from an unexpected network?
	ip := net.ParseIP(srcIP)
	if ip == nil {
		return false
	}

	_, network, err := net.ParseCIDR(expectedNetwork)
	if err != nil {
		return false
	}

	// If IP is not in expected network, might be spoofed
	return !network.Contains(ip)
}
