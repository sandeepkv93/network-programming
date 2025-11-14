package websocket

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	url  string
	conn *websocket.Conn
	done chan bool
}

// NewClient creates a new WebSocket client
func NewClient(url string) *Client {
	return &Client{
		url:  url,
		done: make(chan bool),
	}
}

// Connect connects to the WebSocket server
func (c *Client) Connect() error {
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket server: %v", err)
	}

	log.Printf("Connected to WebSocket server at %s\n", c.url)
	return nil
}

// SendMessage sends a message to the WebSocket server
func (c *Client) SendMessage(message string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}

	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Sent: %s\n", message)
	return nil
}

// ReceiveMessages listens for messages from the server
func (c *Client) ReceiveMessages() {
	defer close(c.done)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v\n", err)
			}
			return
		}

		log.Printf("Received: %s\n", message)
	}
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Printf("Error sending close message: %v\n", err)
	}

	time.Sleep(100 * time.Millisecond)
	return c.conn.Close()
}

// RunInteractive runs an interactive WebSocket client session
func (c *Client) RunInteractive() error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	// Handle interrupt signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start receiving messages
	go c.ReceiveMessages()

	// Read user input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages to send (Ctrl+C to quit):")

	go func() {
		for scanner.Scan() {
			message := scanner.Text()
			if message == "" {
				continue
			}

			if err := c.SendMessage(message); err != nil {
				log.Printf("Error sending message: %v\n", err)
				return
			}
		}
	}()

	// Wait for interrupt or connection close
	select {
	case <-interrupt:
		log.Println("Interrupt received, closing connection...")
	case <-c.done:
		log.Println("Connection closed by server")
	}

	return nil
}

// SendPing sends a ping message to the server
func (c *Client) SendPing() error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}

	err := c.conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		return fmt.Errorf("failed to send ping: %v", err)
	}

	log.Println("Ping sent")
	return nil
}

// SetReadDeadline sets a deadline for read operations
func (c *Client) SetReadDeadline(deadline time.Time) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}
	return c.conn.SetReadDeadline(deadline)
}

// SetWriteDeadline sets a deadline for write operations
func (c *Client) SetWriteDeadline(deadline time.Time) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}
	return c.conn.SetWriteDeadline(deadline)
}
