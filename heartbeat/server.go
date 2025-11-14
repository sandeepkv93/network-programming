package heartbeat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Server represents a heartbeat monitoring server
type Server struct {
	address     string
	listener    net.Listener
	clients     map[string]*ClientStatus
	mu          sync.RWMutex
	timeout     time.Duration
	onDead      func(string)
	onAlive     func(string)
}

// ClientStatus represents the status of a monitored client
type ClientStatus struct {
	ID           string
	Addr         string
	LastHeartbeat time.Time
	Alive        bool
	HeartbeatCount int
	RegisteredAt time.Time
}

// HeartbeatMessage represents a heartbeat message
type HeartbeatMessage struct {
	ClientID  string    `json:"client_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// NewServer creates a new heartbeat server
func NewServer(address string, timeout time.Duration) *Server {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Server{
		address: address,
		clients: make(map[string]*ClientStatus),
		timeout: timeout,
	}
}

// SetCallbacks sets callbacks for client status changes
func (s *Server) SetCallbacks(onAlive, onDead func(string)) {
	s.onAlive = onAlive
	s.onDead = onDead
}

// Start starts the heartbeat server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start heartbeat server: %v", err)
	}

	log.Printf("Heartbeat Server listening on %s\n", s.address)
	log.Printf("Client timeout: %v\n", s.timeout)

	// Start monitoring goroutine
	go s.monitorClients()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleClient(conn)
	}
}

// handleClient handles a client connection
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var msg HeartbeatMessage
		if err := decoder.Decode(&msg); err != nil {
			break
		}

		s.processHeartbeat(&msg, conn.RemoteAddr().String())

		// Send acknowledgment
		ack := map[string]interface{}{
			"status": "ok",
			"time":   time.Now(),
		}
		json.NewEncoder(conn).Encode(ack)
	}
}

// processHeartbeat processes a heartbeat message
func (s *Server) processHeartbeat(msg *HeartbeatMessage, addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.clients[msg.ClientID]
	if !exists {
		client = &ClientStatus{
			ID:           msg.ClientID,
			Addr:         addr,
			RegisteredAt: time.Now(),
			Alive:        true,
		}
		s.clients[msg.ClientID] = client
		log.Printf("New client registered: %s (%s)\n", msg.ClientID, addr)

		if s.onAlive != nil {
			go s.onAlive(msg.ClientID)
		}
	}

	wasAlive := client.Alive
	client.LastHeartbeat = time.Now()
	client.HeartbeatCount++
	client.Alive = true

	if !wasAlive && s.onAlive != nil {
		log.Printf("Client %s is back online\n", msg.ClientID)
		go s.onAlive(msg.ClientID)
	}
}

// monitorClients monitors client heartbeats for timeouts
func (s *Server) monitorClients() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()

		for id, client := range s.clients {
			if client.Alive && now.Sub(client.LastHeartbeat) > s.timeout {
				client.Alive = false
				log.Printf("Client %s timed out (last heartbeat: %v ago)\n",
					id, now.Sub(client.LastHeartbeat))

				if s.onDead != nil {
					go s.onDead(id)
				}
			}
		}
		s.mu.Unlock()
	}
}

// GetClientStatus returns the status of a specific client
func (s *Server) GetClientStatus(clientID string) (*ClientStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	client, exists := s.clients[clientID]
	return client, exists
}

// GetAllClients returns all client statuses
func (s *Server) GetAllClients() map[string]*ClientStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*ClientStatus)
	for id, client := range s.clients {
		result[id] = client
	}
	return result
}

// GetAliveClients returns a list of alive clients
func (s *Server) GetAliveClients() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var alive []string
	for id, client := range s.clients {
		if client.Alive {
			alive = append(alive, id)
		}
	}
	return alive
}

// GetDeadClients returns a list of dead clients
func (s *Server) GetDeadClients() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var dead []string
	for id, client := range s.clients {
		if !client.Alive {
			dead = append(dead, id)
		}
	}
	return dead
}

// RemoveClient removes a client from tracking
func (s *Server) RemoveClient(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, clientID)
	log.Printf("Client %s removed from tracking\n", clientID)
}

// Stop stops the heartbeat server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
