package tunneling

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Client represents a tunneling client (reverse tunnel)
type Client struct {
	serverAddr string
	localAddr  string
	conn       net.Conn
}

// NewClient creates a new tunneling client
func NewClient(serverAddr, localAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
		localAddr:  localAddr,
	}
}

// Connect establishes a reverse tunnel
func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to tunnel server: %v", err)
	}

	log.Printf("Connected to tunnel server at %s\n", c.serverAddr)
	log.Printf("Tunnel established for local service at %s\n", c.localAddr)

	// Connect to local service
	localConn, err := net.Dial("tcp", c.localAddr)
	if err != nil {
		c.conn.Close()
		return fmt.Errorf("failed to connect to local service: %v", err)
	}

	// Start bidirectional forwarding
	go c.forward(c.conn, localConn, "tunnel->local")
	go c.forward(localConn, c.conn, "local->tunnel")

	return nil
}

// forward forwards data between connections
func (c *Client) forward(src, dst net.Conn, direction string) {
	defer src.Close()
	defer dst.Close()

	bytes, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("[%s] Error forwarding: %v\n", direction, err)
	} else {
		log.Printf("[%s] Forwarded %d bytes\n", direction, bytes)
	}
}

// Close closes the tunnel
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
