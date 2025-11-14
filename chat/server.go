package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// Server represents a chat server
type Server struct {
	address string
	listener net.Listener
	clients map[net.Conn]*Client
	mutex   sync.RWMutex
	quit    chan bool
	wg      sync.WaitGroup
}

// Client represents a connected chat client
type Client struct {
	conn     net.Conn
	name     string
	outgoing chan string
}

// NewServer creates a new chat server
func NewServer(address string) *Server {
	return &Server{
		address: address,
		clients: make(map[net.Conn]*Client),
		quit:    make(chan bool),
	}
}

// Start starts the chat server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.listener = listener
	log.Printf("Chat Server listening on %s\n", s.address)

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

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Get client name
	writer.WriteString("Enter your name: ")
	writer.Flush()

	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)

	client := &Client{
		conn:     conn,
		name:     name,
		outgoing: make(chan string, 10),
	}

	// Add client to server
	s.mutex.Lock()
	s.clients[conn] = client
	s.mutex.Unlock()

	// Announce new client
	s.broadcast(fmt.Sprintf("%s joined the chat\n", name), conn)
	log.Printf("%s joined the chat\n", name)

	// Start message sender for this client
	s.wg.Add(1)
	go s.sendMessages(client)

	// Read messages from client
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Broadcast message to all clients
		fullMessage := fmt.Sprintf("%s: %s\n", name, message)
		s.broadcast(fullMessage, conn)
	}

	// Client disconnected
	s.mutex.Lock()
	delete(s.clients, conn)
	close(client.outgoing)
	s.mutex.Unlock()

	s.broadcast(fmt.Sprintf("%s left the chat\n", name), nil)
	log.Printf("%s left the chat\n", name)
}

// sendMessages sends messages to a client
func (s *Server) sendMessages(client *Client) {
	defer s.wg.Done()

	writer := bufio.NewWriter(client.conn)
	for message := range client.outgoing {
		_, err := writer.WriteString(message)
		if err != nil {
			return
		}
		writer.Flush()
	}
}

// broadcast sends a message to all clients except the sender
func (s *Server) broadcast(message string, sender net.Conn) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for conn, client := range s.clients {
		if conn != sender {
			select {
			case client.outgoing <- message:
			default:
				// Channel is full, skip
			}
		}
	}
}

// Stop stops the chat server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}

	// Close all client connections
	s.mutex.Lock()
	for conn, client := range s.clients {
		conn.Close()
		close(client.outgoing)
	}
	s.mutex.Unlock()

	s.wg.Wait()
	log.Println("Chat Server stopped")
}
