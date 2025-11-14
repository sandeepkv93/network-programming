package chat

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// ChatClient represents a chat client
type ChatClient struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	wg      sync.WaitGroup
	quit    chan bool
}

// NewClient creates a new chat client
func NewClient(address string) *ChatClient {
	return &ChatClient{
		address: address,
		quit:    make(chan bool),
	}
}

// Connect connects to the chat server
func (c *ChatClient) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)

	return nil
}

// SetName sets the client's name
func (c *ChatClient) SetName(name string) error {
	// Read the prompt
	_, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Send name
	_, err = c.writer.WriteString(name + "\n")
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

// SendMessage sends a message to the chat
func (c *ChatClient) SendMessage(message string) error {
	_, err := c.writer.WriteString(message + "\n")
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

// ReceiveMessages receives messages from the server
func (c *ChatClient) ReceiveMessages(callback func(string)) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-c.quit:
				return
			default:
				message, err := c.reader.ReadString('\n')
				if err != nil {
					return
				}
				callback(message)
			}
		}
	}()
}

// Close closes the connection
func (c *ChatClient) Close() {
	close(c.quit)
	if c.conn != nil {
		c.conn.Close()
	}
	c.wg.Wait()
}
