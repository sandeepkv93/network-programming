package heartbeat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

// Client represents a heartbeat client
type Client struct {
	serverAddr string
	clientID   string
	interval   time.Duration
	conn       net.Conn
	running    bool
	stopChan   chan bool
}

// NewClient creates a new heartbeat client
func NewClient(serverAddr string, interval time.Duration) *Client {
	if interval == 0 {
		interval = 10 * time.Second
	}

	return &Client{
		serverAddr: serverAddr,
		clientID:   uuid.New().String(),
		interval:   interval,
		stopChan:   make(chan bool),
	}
}

// NewClientWithID creates a client with a specific ID
func NewClientWithID(serverAddr, clientID string, interval time.Duration) *Client {
	if interval == 0 {
		interval = 10 * time.Second
	}

	return &Client{
		serverAddr: serverAddr,
		clientID:   clientID,
		interval:   interval,
		stopChan:   make(chan bool),
	}
}

// Start starts sending heartbeats
func (c *Client) Start() error {
	var err error
	c.conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to heartbeat server: %v", err)
	}

	log.Printf("Connected to heartbeat server at %s\n", c.serverAddr)
	log.Printf("Client ID: %s\n", c.clientID)
	log.Printf("Heartbeat interval: %v\n", c.interval)

	c.running = true

	go c.heartbeatLoop()

	return nil
}

// heartbeatLoop sends periodic heartbeats
func (c *Client) heartbeatLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.sendHeartbeat(); err != nil {
				log.Printf("Failed to send heartbeat: %v\n", err)
				c.reconnect()
			}
		case <-c.stopChan:
			return
		}
	}
}

// sendHeartbeat sends a heartbeat message
func (c *Client) sendHeartbeat() error {
	msg := HeartbeatMessage{
		ClientID:  c.clientID,
		Timestamp: time.Now(),
		Status:    "alive",
	}

	encoder := json.NewEncoder(c.conn)
	if err := encoder.Encode(msg); err != nil {
		return err
	}

	// Read acknowledgment
	var ack map[string]interface{}
	decoder := json.NewDecoder(c.conn)
	if err := decoder.Decode(&ack); err != nil {
		return err
	}

	log.Printf("Heartbeat sent and acknowledged\n")
	return nil
}

// reconnect attempts to reconnect to the server
func (c *Client) reconnect() {
	log.Println("Attempting to reconnect...")

	if c.conn != nil {
		c.conn.Close()
	}

	for {
		conn, err := net.Dial("tcp", c.serverAddr)
		if err != nil {
			log.Printf("Reconnection failed: %v, retrying in 5s...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		c.conn = conn
		log.Println("Reconnected successfully")
		return
	}
}

// Stop stops sending heartbeats
func (c *Client) Stop() error {
	c.running = false
	close(c.stopChan)

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// GetClientID returns the client ID
func (c *Client) GetClientID() string {
	return c.clientID
}

// IsRunning returns whether the client is running
func (c *Client) IsRunning() bool {
	return c.running
}
