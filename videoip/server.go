package videoip

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
)

// VideoServer represents a Video over IP server
type VideoServer struct {
	Address string
	conn    *net.UDPConn
	clients map[string]*net.UDPAddr
	mu      sync.Mutex
}

// VideoPacket represents a video frame packet
type VideoPacket struct {
	Timestamp  uint32
	Sequence   uint16
	FrameType  uint8 // I-frame, P-frame, B-frame
	FragmentID uint8
	Data       []byte
}

// NewServer creates a new Video over IP server
func NewServer(address string) *VideoServer {
	return &VideoServer{
		Address: address,
		clients: make(map[string]*net.UDPAddr),
	}
}

// Start starts the video server
func (s *VideoServer) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.Address)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.conn = conn
	log.Printf("Video over IP Server started on %s\n", s.Address)
	log.Println("Waiting for video packets...")

	buffer := make([]byte, 65535) // Max UDP packet size

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

		// Parse video packet
		packet := s.parsePacket(buffer[:n])
		if packet != nil {
			// Broadcast to all other clients
			s.broadcast(buffer[:n], clientAddr)
		}
	}
}

// Stop stops the video server
func (s *VideoServer) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *VideoServer) parsePacket(data []byte) *VideoPacket {
	if len(data) < 10 {
		return nil
	}

	return &VideoPacket{
		Timestamp:  binary.BigEndian.Uint32(data[0:4]),
		Sequence:   binary.BigEndian.Uint16(data[4:6]),
		FrameType:  data[6],
		FragmentID: data[7],
		Data:       data[10:],
	}
}

func (s *VideoServer) broadcast(data []byte, sender *net.UDPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
func (s *VideoServer) GetConnectedClients() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
}
