package tunneling

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// Server represents a tunneling server
type Server struct {
	listenAddr string
	targetAddr string
	listener   net.Listener
	mu         sync.Mutex
	tunnels    int
}

// NewServer creates a new tunneling server
func NewServer(listenAddr, targetAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		targetAddr: targetAddr,
	}
}

// Start starts the tunneling server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start tunnel server: %v", err)
	}

	log.Printf("Tunnel Server listening on %s\n", s.listenAddr)
	log.Printf("Forwarding to %s\n", s.targetAddr)

	for {
		clientConn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleTunnel(clientConn)
	}
}

// handleTunnel handles a tunnel connection
func (s *Server) handleTunnel(clientConn net.Conn) {
	s.mu.Lock()
	s.tunnels++
	tunnelID := s.tunnels
	s.mu.Unlock()

	log.Printf("Tunnel #%d: Connection from %s\n", tunnelID, clientConn.RemoteAddr())

	// Connect to target
	targetConn, err := net.Dial("tcp", s.targetAddr)
	if err != nil {
		log.Printf("Tunnel #%d: Failed to connect to target: %v\n", tunnelID, err)
		clientConn.Close()
		return
	}

	log.Printf("Tunnel #%d: Connected to target %s\n", tunnelID, s.targetAddr)

	// Start bidirectional forwarding
	go s.forward(clientConn, targetConn, tunnelID, "client->target")
	go s.forward(targetConn, clientConn, tunnelID, "target->client")
}

// forward forwards data between connections
func (s *Server) forward(src, dst net.Conn, tunnelID int, direction string) {
	defer src.Close()
	defer dst.Close()

	bytes, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("Tunnel #%d [%s]: Error forwarding: %v\n", tunnelID, direction, err)
	} else {
		log.Printf("Tunnel #%d [%s]: Forwarded %d bytes\n", tunnelID, direction, bytes)
	}
}

// Stop stops the tunneling server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// GetActiveTunnels returns the number of tunnels created
func (s *Server) GetActiveTunnels() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tunnels
}
