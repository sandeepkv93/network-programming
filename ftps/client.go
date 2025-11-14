package ftps

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
)

// Client represents an FTPS client
type Client struct {
	Address    string
	SkipVerify bool
	conn       net.Conn
	reader     *bufio.Reader
}

// NewClient creates a new FTPS client
func NewClient(address string, skipVerify bool) *Client {
	return &Client{
		Address:    address,
		SkipVerify: skipVerify,
	}
}

// Connect connects to the FTPS server
func (c *Client) Connect() error {
	config := &tls.Config{
		InsecureSkipVerify: c.SkipVerify,
		MinVersion:         tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", c.Address, config)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)

	// Read welcome message
	response, err := c.readResponse()
	if err != nil {
		return err
	}

	fmt.Printf("Server: %s\n", response)
	return nil
}

// Login authenticates with the server
func (c *Client) Login(username, password string) error {
	// Send USER command
	if err := c.sendCommand(fmt.Sprintf("USER %s", username)); err != nil {
		return err
	}

	response, err := c.readResponse()
	if err != nil {
		return err
	}
	fmt.Printf("Server: %s\n", response)

	// Send PASS command
	if err := c.sendCommand(fmt.Sprintf("PASS %s", password)); err != nil {
		return err
	}

	response, err = c.readResponse()
	if err != nil {
		return err
	}
	fmt.Printf("Server: %s\n", response)

	return nil
}

// List lists files in the current directory
func (c *Client) List() (string, error) {
	if err := c.sendCommand("LIST"); err != nil {
		return "", err
	}

	response, err := c.readResponse()
	if err != nil {
		return "", err
	}

	return response, nil
}

// Close closes the connection
func (c *Client) Close() error {
	if c.conn != nil {
		c.sendCommand("QUIT")
		return c.conn.Close()
	}
	return nil
}

func (c *Client) sendCommand(command string) error {
	_, err := fmt.Fprintf(c.conn, "%s\r\n", command)
	return err
}

func (c *Client) readResponse() (string, error) {
	var response strings.Builder
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		response.WriteString(line)

		// Simple response parsing - check if line starts with digit
		if len(line) >= 3 && line[3] == ' ' {
			break
		}
	}
	return response.String(), nil
}
