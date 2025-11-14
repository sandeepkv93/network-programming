package echo

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

// Server represents an Echo server
type Server struct {
	address  string
	listener net.Listener
	quit     chan bool
	wg       sync.WaitGroup
}

// NewServer creates a new Echo server
func NewServer(address string) *Server {
	return &Server{
		address: address,
		quit:    make(chan bool),
	}
}

// Start starts the Echo server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.listener = listener
	log.Printf("Echo Server listening on %s\n", s.address)

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

// acceptConnections accepts incoming connections
func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					log.Printf("Error accepting connection: %v\n", err)
					continue
				}
			}
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

// handleConnection handles a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s\n", clientAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// Read message from client
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected\n", clientAddr)
			return
		}

		log.Printf("Received from %s: %s", clientAddr, message)

		// Echo back the exact message
		_, err = writer.WriteString(message)
		if err != nil {
			log.Printf("Error writing to client %s: %v\n", clientAddr, err)
			return
		}
		writer.Flush()
	}
}

// Stop stops the Echo server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Println("Echo Server stopped")
}
