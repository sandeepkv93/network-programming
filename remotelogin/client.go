package remotelogin

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// Client represents a remote login client
type Client struct {
	serverAddr string
	conn       net.Conn
	connected  bool
}

// NewClient creates a new remote login client
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
	}
}

// Connect connects to the remote login server
func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.connected = true
	log.Printf("Connected to %s\n", c.serverAddr)
	return nil
}

// Login performs interactive login
func (c *Client) Login(username, password string) error {
	if !c.connected {
		return fmt.Errorf("not connected to server")
	}

	reader := bufio.NewReader(c.conn)
	scanner := bufio.NewScanner(os.Stdin)

	// Read welcome banner
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print(line)

		if strings.Contains(line, "Username:") {
			break
		}
	}

	// Send username
	if username == "" {
		fmt.Print("") // Prompt already shown by server
		if scanner.Scan() {
			username = scanner.Text()
		}
		c.conn.Write([]byte(username + "\n"))
	} else {
		c.conn.Write([]byte(username + "\n"))
	}

	// Wait for password prompt
	line, _ := reader.ReadString('\n')
	fmt.Print(line)

	// Send password
	if password == "" {
		if scanner.Scan() {
			password = scanner.Text()
		}
		c.conn.Write([]byte(password + "\n"))
	} else {
		c.conn.Write([]byte(password + "\n"))
	}

	// Read authentication result
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Print(line)

	if strings.Contains(line, "failed") {
		return fmt.Errorf("authentication failed")
	}

	return nil
}

// StartInteractive starts an interactive session
func (c *Client) StartInteractive() error {
	if !c.connected {
		return fmt.Errorf("not connected to server")
	}

	// Start goroutine to read from server
	go func() {
		reader := bufio.NewReader(c.conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Connection error: %v\n", err)
				}
				fmt.Println("\nConnection closed by server")
				os.Exit(0)
			}
			fmt.Print(line)
		}
	}()

	// Read from stdin and send to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		_, err := c.conn.Write([]byte(line + "\n"))
		if err != nil {
			return fmt.Errorf("failed to send command: %v", err)
		}

		if line == "exit" || line == "logout" {
			break
		}
	}

	return nil
}

// Disconnect disconnects from the server
func (c *Client) Disconnect() error {
	c.connected = false
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected returns whether connected to server
func (c *Client) IsConnected() bool {
	return c.connected
}
