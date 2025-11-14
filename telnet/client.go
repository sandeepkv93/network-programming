package telnet

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Client represents a telnet client
type Client struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	wg      sync.WaitGroup
	quit    chan bool
}

// NewClient creates a new telnet client
func NewClient(address string) *Client {
	return &Client{
		address: address,
		quit:    make(chan bool),
	}
}

// Connect connects to the telnet server
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)

	return nil
}

// ReadWelcome reads the welcome message from the server
func (c *Client) ReadWelcome() (string, error) {
	var response strings.Builder

	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		response.WriteString(line)

		if strings.HasSuffix(line, "> ") {
			break
		}
	}

	return response.String(), nil
}

// SendCommand sends a command to the server
func (c *Client) SendCommand(command string) error {
	_, err := c.writer.WriteString(command + "\n")
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

// ReadResponse reads the response from the server
func (c *Client) ReadResponse() (string, error) {
	var response strings.Builder

	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// Check if this is the prompt
		if strings.HasSuffix(line, "> ") {
			// Remove the prompt from response
			line = strings.TrimSuffix(line, "> ")
			if line != "" {
				response.WriteString(line)
			}
			break
		}

		response.WriteString(line)
	}

	return response.String(), nil
}

// ReceiveData receives data from the server continuously
func (c *Client) ReceiveData(callback func(string)) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-c.quit:
				return
			default:
				line, err := c.reader.ReadString('\n')
				if err != nil {
					return
				}
				callback(line)
			}
		}
	}()
}

// Close closes the connection
func (c *Client) Close() {
	close(c.quit)
	if c.conn != nil {
		c.conn.Close()
	}
	c.wg.Wait()
}
