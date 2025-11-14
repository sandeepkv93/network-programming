package mail

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents a mail client
type Client struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
}

// NewClient creates a new mail client
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Connect connects to the mail server
func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	c.conn = conn
	c.reader = bufio.NewReader(conn)

	// Read welcome message
	_, err = c.readResponse()
	if err != nil {
		return fmt.Errorf("failed to read welcome: %v", err)
	}

	return nil
}

// Login logs in as a user
func (c *Client) Login(username string) error {
	if err := c.sendCommand(fmt.Sprintf("USER %s", username)); err != nil {
		return err
	}

	response, err := c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "+OK") {
		return fmt.Errorf("login failed: %s", response)
	}

	return nil
}

// GetMessageCount returns the number of messages
func (c *Client) GetMessageCount() (int, error) {
	if err := c.sendCommand("STAT"); err != nil {
		return 0, err
	}

	response, err := c.readResponse()
	if err != nil {
		return 0, err
	}

	var count int
	fmt.Sscanf(response, "+OK %d", &count)
	return count, nil
}

// ListMessages lists all messages
func (c *Client) ListMessages() ([]string, error) {
	if err := c.sendCommand("LIST"); err != nil {
		return nil, err
	}

	// Read header
	_, err := c.readResponse()
	if err != nil {
		return nil, err
	}

	// Read message list
	var messages []string
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			break
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "." {
			break
		}

		messages = append(messages, trimmed)
	}

	return messages, nil
}

// RetrieveMessage retrieves a specific message
func (c *Client) RetrieveMessage(index int) (string, error) {
	if err := c.sendCommand(fmt.Sprintf("RETR %d", index)); err != nil {
		return "", err
	}

	// Read header
	header, err := c.readResponse()
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(header, "+OK") {
		return "", fmt.Errorf("failed to retrieve message: %s", header)
	}

	// Read message content
	var message strings.Builder
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			break
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "." {
			break
		}

		message.WriteString(line)
	}

	return message.String(), nil
}

// SendMessage sends a message to another user
func (c *Client) SendMessage(to, subject, body string) error {
	if err := c.sendCommand(fmt.Sprintf("SEND %s", to)); err != nil {
		return err
	}

	response, err := c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "+OK") {
		return fmt.Errorf("send failed: %s", response)
	}

	// Send subject (first line)
	c.conn.Write([]byte(subject + "\r\n"))

	// Send body
	c.conn.Write([]byte(body))

	// Send terminator
	c.conn.Write([]byte("\r\n.\r\n"))

	// Read confirmation
	response, err = c.readResponse()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "+OK") {
		return fmt.Errorf("send failed: %s", response)
	}

	return nil
}

// Quit disconnects from the mail server
func (c *Client) Quit() error {
	if c.conn == nil {
		return nil
	}

	c.sendCommand("QUIT")
	c.readResponse()

	return c.conn.Close()
}

// sendCommand sends a command to the server
func (c *Client) sendCommand(command string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	_, err := c.conn.Write([]byte(command + "\r\n"))
	return err
}

// readResponse reads a response from the server
func (c *Client) readResponse() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
