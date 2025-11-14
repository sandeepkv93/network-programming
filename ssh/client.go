package ssh

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

// Client represents an SSH client
type Client struct {
	address string
	config  *ssh.ClientConfig
	client  *ssh.Client
}

// NewClient creates a new SSH client
func NewClient(address, username, password string) *Client {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For demo purposes only
	}

	return &Client{
		address: address,
		config:  config,
	}
}

// Connect connects to the SSH server
func (c *Client) Connect() error {
	client, err := ssh.Dial("tcp", c.address, c.config)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.client = client
	return nil
}

// ExecuteCommand executes a command on the remote server
func (c *Client) ExecuteCommand(command string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// StartShell starts an interactive shell session
func (c *Client) StartShell() error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		return err
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		return err
	}

	// Wait for session to finish
	return session.Wait()
}

// RunInteractive runs an interactive session with custom I/O
func (c *Client) RunInteractive(stdin io.Reader, stdout, stderr io.Writer) error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdin = stdin
	session.Stdout = stdout
	session.Stderr = stderr

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		return err
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		return err
	}

	// Wait for session to finish
	return session.Wait()
}

// Close closes the SSH connection
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
