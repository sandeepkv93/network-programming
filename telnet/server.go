package telnet

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// Server represents a telnet server
type Server struct {
	address  string
	listener net.Listener
	quit     chan bool
	wg       sync.WaitGroup
}

// NewServer creates a new telnet server
func NewServer(address string) *Server {
	return &Server{
		address: address,
		quit:    make(chan bool),
	}
}

// Start starts the telnet server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.listener = listener
	log.Printf("Telnet Server listening on %s\n", s.address)

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

// handleConnection handles a single telnet connection
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s\n", clientAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Send welcome message
	welcome := "Welcome to Telnet Server\r\n"
	welcome += "Type 'help' for available commands\r\n"
	welcome += "> "
	writer.WriteString(welcome)
	writer.Flush()

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected\n", clientAddr)
			return
		}

		command = strings.TrimSpace(command)
		response := s.processCommand(command)

		writer.WriteString(response + "\r\n> ")
		writer.Flush()
	}
}

// processCommand processes telnet commands
func (s *Server) processCommand(command string) string {
	switch strings.ToLower(command) {
	case "help":
		return "Available commands:\r\n" +
			"  help  - Show this help message\r\n" +
			"  time  - Show current server time\r\n" +
			"  echo <message> - Echo back the message\r\n" +
			"  quit  - Disconnect from server"
	case "time":
		return fmt.Sprintf("Server time: %s", "2025-11-14 04:48:00")
	case "quit":
		return "Goodbye!"
	default:
		if strings.HasPrefix(command, "echo ") {
			return strings.TrimPrefix(command, "echo ")
		}
		return fmt.Sprintf("Unknown command: %s (type 'help' for available commands)", command)
	}
}

// Stop stops the telnet server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Println("Telnet Server stopped")
}
