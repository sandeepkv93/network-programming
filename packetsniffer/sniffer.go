package packetsniffer

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

// PacketInfo represents information about a captured packet
type PacketInfo struct {
	Timestamp   time.Time
	Protocol    string
	SrcIP       string
	DstIP       string
	SrcPort     uint16
	DstPort     uint16
	Length      int
	PayloadSize int
}

// Sniffer represents a packet sniffer
type Sniffer struct {
	Interface   string
	Filter      string
	PacketCount int
	OnPacket    func(info PacketInfo)
	packets     []PacketInfo
	mu          sync.Mutex
	running     bool
}

// NewSniffer creates a new packet sniffer
func NewSniffer(iface string) *Sniffer {
	return &Sniffer{
		Interface: iface,
		packets:   make([]PacketInfo, 0),
	}
}

// Start starts the packet sniffer
func (s *Sniffer) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("sniffer is already running")
	}
	s.running = true
	s.mu.Unlock()

	fmt.Printf("Starting packet sniffer on interface: %s\n", s.Interface)
	if s.Interface == "" {
		s.Interface = "any"
	}

	// Create a raw socket to capture packets
	// Note: This requires elevated privileges (root/admin)
	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		return fmt.Errorf("failed to create raw socket: %v (requires elevated privileges)", err)
	}
	defer conn.Close()

	fmt.Println("Packet sniffer started. Capturing packets...")

	buffer := make([]byte, 65535)
	count := 0

	for s.running {
		// Set read deadline to check if we should stop
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			fmt.Printf("Error reading packet: %v\n", err)
			continue
		}

		// Parse the packet
		info := s.parsePacket(buffer[:n], addr)
		if info != nil {
			count++
			fmt.Printf("[%d] %s\n", count, s.formatPacketInfo(*info))

			s.mu.Lock()
			s.packets = append(s.packets, *info)
			s.mu.Unlock()

			if s.OnPacket != nil {
				s.OnPacket(*info)
			}

			// Check if we've reached the packet count limit
			if s.PacketCount > 0 && count >= s.PacketCount {
				break
			}
		}
	}

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	fmt.Printf("\nPacket capture stopped. Total packets captured: %d\n", count)
	return nil
}

// Stop stops the packet sniffer
func (s *Sniffer) Stop() {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
}

func (s *Sniffer) parsePacket(data []byte, addr net.Addr) *PacketInfo {
	if len(data) < 20 {
		return nil
	}

	info := &PacketInfo{
		Timestamp: time.Now(),
		Length:    len(data),
	}

	// Parse IP header (simplified)
	version := data[0] >> 4
	if version != 4 {
		return nil // Only IPv4 for now
	}

	ihl := (data[0] & 0x0F) * 4
	if len(data) < int(ihl) {
		return nil
	}

	protocol := data[9]
	switch protocol {
	case 1:
		info.Protocol = "ICMP"
	case 6:
		info.Protocol = "TCP"
	case 17:
		info.Protocol = "UDP"
	default:
		info.Protocol = fmt.Sprintf("Unknown(%d)", protocol)
	}

	// Source and destination IP
	info.SrcIP = fmt.Sprintf("%d.%d.%d.%d", data[12], data[13], data[14], data[15])
	info.DstIP = fmt.Sprintf("%d.%d.%d.%d", data[16], data[17], data[18], data[19])

	// Parse TCP/UDP ports if applicable
	if (protocol == 6 || protocol == 17) && len(data) >= int(ihl)+4 {
		info.SrcPort = binary.BigEndian.Uint16(data[ihl : ihl+2])
		info.DstPort = binary.BigEndian.Uint16(data[ihl+2 : ihl+4])
	}

	info.PayloadSize = len(data) - int(ihl)

	return info
}

func (s *Sniffer) formatPacketInfo(info PacketInfo) string {
	timestamp := info.Timestamp.Format("15:04:05.000")
	if info.SrcPort > 0 && info.DstPort > 0 {
		return fmt.Sprintf("%s | %-8s | %15s:%-5d -> %15s:%-5d | %d bytes",
			timestamp, info.Protocol, info.SrcIP, info.SrcPort, info.DstIP, info.DstPort, info.Length)
	}
	return fmt.Sprintf("%s | %-8s | %15s -> %15s | %d bytes",
		timestamp, info.Protocol, info.SrcIP, info.DstIP, info.Length)
}

// GetPackets returns all captured packets
func (s *Sniffer) GetPackets() []PacketInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	packets := make([]PacketInfo, len(s.packets))
	copy(packets, s.packets)
	return packets
}

// GetStatistics returns statistics about captured packets
func (s *Sniffer) GetStatistics() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := make(map[string]int)
	for _, packet := range s.packets {
		stats[packet.Protocol]++
	}
	return stats
}
