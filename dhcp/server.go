package dhcp

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// MessageType represents DHCP message types
type MessageType byte

const (
	DHCPDiscover MessageType = 1
	DHCPOffer    MessageType = 2
	DHCPRequest  MessageType = 3
	DHCPAck      MessageType = 4
	DHCPNak      MessageType = 5
	DHCPRelease  MessageType = 6
)

// Server represents a DHCP server
type Server struct {
	address     string
	ipPool      []net.IP
	leases      map[string]Lease
	mutex       sync.RWMutex
	subnetMask  net.IP
	router      net.IP
	dnsServer   net.IP
	leaseTime   uint32
	conn        *net.UDPConn
	quit        chan bool
	wg          sync.WaitGroup
}

// Lease represents a DHCP lease
type Lease struct {
	IP         net.IP
	MAC        net.HardwareAddr
	Expiry     time.Time
}

// NewServer creates a new DHCP server
func NewServer(address string, ipPoolStart, ipPoolEnd string, subnetMask, router, dnsServer string) (*Server, error) {
	start := net.ParseIP(ipPoolStart)
	end := net.ParseIP(ipPoolEnd)

	if start == nil || end == nil {
		return nil, fmt.Errorf("invalid IP pool range")
	}

	// Generate IP pool
	ipPool := generateIPPool(start, end)

	return &Server{
		address:    address,
		ipPool:     ipPool,
		leases:     make(map[string]Lease),
		subnetMask: net.ParseIP(subnetMask),
		router:     net.ParseIP(router),
		dnsServer:  net.ParseIP(dnsServer),
		leaseTime:  3600, // 1 hour
		quit:       make(chan bool),
	}, nil
}

// Start starts the DHCP server
func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp4", s.address)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.conn = conn
	log.Printf("DHCP Server listening on %s\n", s.address)

	s.wg.Add(1)
	go s.handleRequests()

	return nil
}

// handleRequests handles incoming DHCP requests
func (s *Server) handleRequests() {
	defer s.wg.Done()

	buffer := make([]byte, 1500)

	for {
		select {
		case <-s.quit:
			return
		default:
			n, clientAddr, err := s.conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}

			// Parse DHCP packet (simplified)
			if n < 240 {
				continue
			}

			// Extract MAC address (chaddr field at offset 28)
			mac := net.HardwareAddr(buffer[28:34])

			// Determine message type (simplified - would need to parse options)
			// For this simplified implementation, we'll treat all requests as DISCOVER
			s.handleDiscover(mac, clientAddr)
		}
	}
}

// handleDiscover handles DHCP DISCOVER messages
func (s *Server) handleDiscover(mac net.HardwareAddr, clientAddr *net.UDPAddr) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	macStr := mac.String()

	// Check if client already has a lease
	if lease, exists := s.leases[macStr]; exists && time.Now().Before(lease.Expiry) {
		s.sendOffer(mac, lease.IP, clientAddr)
		return
	}

	// Allocate new IP from pool
	if len(s.ipPool) == 0 {
		log.Println("No available IPs in pool")
		return
	}

	ip := s.ipPool[0]
	s.ipPool = s.ipPool[1:]

	// Create lease
	lease := Lease{
		IP:     ip,
		MAC:    mac,
		Expiry: time.Now().Add(time.Duration(s.leaseTime) * time.Second),
	}

	s.leases[macStr] = lease
	log.Printf("Allocated IP %s to MAC %s\n", ip, mac)

	s.sendOffer(mac, ip, clientAddr)
}

// sendOffer sends a DHCP OFFER message
func (s *Server) sendOffer(mac net.HardwareAddr, ip net.IP, clientAddr *net.UDPAddr) {
	// Build simplified DHCP OFFER packet
	packet := make([]byte, 240)

	// Boot reply
	packet[0] = 2

	// Hardware type (Ethernet)
	packet[1] = 1

	// Hardware address length
	packet[2] = 6

	// Copy MAC address
	copy(packet[28:34], mac)

	// Your IP address (yiaddr)
	copy(packet[16:20], ip.To4())

	// Server IP (siaddr)
	serverIP := net.ParseIP("192.168.1.1").To4()
	copy(packet[20:24], serverIP)

	// Send packet
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 68,
	}

	s.conn.WriteToUDP(packet, broadcastAddr)
	log.Printf("Sent OFFER for IP %s to MAC %s\n", ip, mac)
}

// GetLeases returns all active leases
func (s *Server) GetLeases() []Lease {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var leases []Lease
	for _, lease := range s.leases {
		if time.Now().Before(lease.Expiry) {
			leases = append(leases, lease)
		}
	}

	return leases
}

// Stop stops the DHCP server
func (s *Server) Stop() {
	close(s.quit)
	if s.conn != nil {
		s.conn.Close()
	}
	s.wg.Wait()
	log.Println("DHCP Server stopped")
}

// generateIPPool generates a pool of IP addresses
func generateIPPool(start, end net.IP) []net.IP {
	var pool []net.IP

	startInt := binary.BigEndian.Uint32(start.To4())
	endInt := binary.BigEndian.Uint32(end.To4())

	for i := startInt; i <= endInt; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		pool = append(pool, ip)
	}

	return pool
}
