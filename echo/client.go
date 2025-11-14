package echo

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents an Echo client
type Client struct {
	address string
	conn    net.Conn
}

// NewClient creates a new Echo client
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Connect connects to the Echo server
func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	c.conn = conn
	return nil
}

// Echo sends a message to the server and receives the echo
func (c *Client) Echo(message string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected to server")
	}

	// Send message
	_, err := c.conn.Write([]byte(message + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	// Receive echo
	reader := bufio.NewReader(c.conn)
	echo, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to receive echo: %v", err)
	}

	return strings.TrimSpace(echo), nil
}

// Close closes the connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
