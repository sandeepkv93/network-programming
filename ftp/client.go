package ftp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents an FTP client
type Client struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
}

// NewClient creates a new FTP client
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Connect connects to the FTP server
func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, 10*time.Second)
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

// Login logs in to the FTP server
func (c *Client) Login(username, password string) error {
	// Send USER command
	if err := c.sendCommand(fmt.Sprintf("USER %s", username)); err != nil {
		return err
	}
	if _, err := c.readResponse(); err != nil {
		return err
	}

	// Send PASS command
	if err := c.sendCommand(fmt.Sprintf("PASS %s", password)); err != nil {
		return err
	}
	if _, err := c.readResponse(); err != nil {
		return err
	}

	return nil
}

// Pwd returns the current working directory
func (c *Client) Pwd() (string, error) {
	if err := c.sendCommand("PWD"); err != nil {
		return "", err
	}

	response, err := c.readResponse()
	if err != nil {
		return "", err
	}

	return response, nil
}

// Cwd changes the working directory
func (c *Client) Cwd(dir string) error {
	if err := c.sendCommand(fmt.Sprintf("CWD %s", dir)); err != nil {
		return err
	}

	_, err := c.readResponse()
	return err
}

// List lists files in the current directory
func (c *Client) List() (string, error) {
	if err := c.sendCommand("LIST"); err != nil {
		return "", err
	}

	// Read response header
	_, err := c.readResponse()
	if err != nil {
		return "", err
	}

	// Read directory listing
	var listing strings.Builder
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			break
		}

		// Check if this is a response code (end of listing)
		if len(line) > 0 && line[0] >= '0' && line[0] <= '9' {
			break
		}

		listing.WriteString(line)
	}

	return listing.String(), nil
}

// Quit disconnects from the FTP server
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
