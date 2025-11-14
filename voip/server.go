package voip

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
)

// VoIPServer represents a Voice over IP server
type VoIPServer struct {
	Address string
	conn    *net.UDPConn
	clients map[string]*net.UDPAddr
	mu      sync.Mutex
}

// AudioPacket represents a voice packet
type AudioPacket struct {
	Timestamp uint32
	Sequence  uint16
	Data      []byte
}

// NewServer creates a new VoIP server
func NewServer(address string) *VoIPServer {
	return &VoIPServer{
		Address: address,
		clients: make(map[string]*net.UDPAddr),
	}
}

// Start starts the VoIP server
func (s *VoIPServer) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.Address)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.conn = conn
	log.Printf("VoIP Server started on %s\n", s.Address)
	log.Println("Waiting for voice packets...")

	buffer := make([]byte, 1500) // MTU size

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading: %v\n", err)
			continue
		}

		// Register client
		s.mu.Lock()
		clientKey := clientAddr.String()
		if _, exists := s.clients[clientKey]; !exists {
			s.clients[clientKey] = clientAddr
			log.Printf("New client connected: %s\n", clientKey)
		}
		s.mu.Unlock()

		// Parse audio packet
		packet := s.parsePacket(buffer[:n])

		// Broadcast to all other clients
		s.broadcast(packet, clientAddr)
	}
}

// Stop stops the VoIP server
func (s *VoIPServer) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *VoIPServer) parsePacket(data []byte) *AudioPacket {
	if len(data) < 8 {
		return nil
	}

	return &AudioPacket{
		Timestamp: binary.BigEndian.Uint32(data[0:4]),
		Sequence:  binary.BigEndian.Uint16(data[4:6]),
		Data:      data[8:],
	}
}

func (s *VoIPServer) broadcast(packet *AudioPacket, sender *net.UDPAddr) {
	if packet == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Create packet data
	data := make([]byte, 8+len(packet.Data))
	binary.BigEndian.PutUint32(data[0:4], packet.Timestamp)
	binary.BigEndian.PutUint16(data[4:6], packet.Sequence)
	copy(data[8:], packet.Data)

	// Send to all clients except sender
	for clientKey, clientAddr := range s.clients {
		if clientAddr.String() != sender.String() {
			_, err := s.conn.WriteToUDP(data, clientAddr)
			if err != nil {
				log.Printf("Failed to send to %s: %v\n", clientKey, err)
			}
		}
	}
}

// GetConnectedClients returns the number of connected clients
func (s *VoIPServer) GetConnectedClients() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
}
