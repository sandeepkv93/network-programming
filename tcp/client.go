package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents a TCP client
type Client struct {
	address string
	conn    net.Conn
}

// NewClient creates a new TCP client
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Connect connects to the TCP server
func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	c.conn = conn
	return nil
}

// Send sends a message to the server and receives the response
func (c *Client) Send(message string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected to server")
	}

	// Send message
	_, err := c.conn.Write([]byte(message + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	// Receive response
	reader := bufio.NewReader(c.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to receive response: %v", err)
	}

	return strings.TrimSpace(response), nil
}

// Close closes the connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
