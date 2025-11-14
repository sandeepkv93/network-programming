package sftp

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// Client represents an SFTP client
type Client struct {
	Address  string
	Username string
	Password string
	sshClient *ssh.Client
}

// NewClient creates a new SFTP client
func NewClient(address, username, password string) *Client {
	return &Client{
		Address:  address,
		Username: username,
		Password: password,
	}
}

// Connect connects to the SFTP server
func (c *Client) Connect() error {
	config := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", c.Address, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}

	c.sshClient = client
	fmt.Printf("Connected to SFTP server at %s\n", c.Address)
	return nil
}

// Upload uploads a file to the server
func (c *Client) Upload(localPath, remotePath string) error {
	if c.sshClient == nil {
		return fmt.Errorf("not connected")
	}

	// Open SFTP session
	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Read local file
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// In a real implementation, use SFTP protocol to upload
	fmt.Printf("Would upload %d bytes from %s to %s\n", len(data), localPath, remotePath)
	return nil
}

// Download downloads a file from the server
func (c *Client) Download(remotePath, localPath string) error {
	if c.sshClient == nil {
		return fmt.Errorf("not connected")
	}

	// In a real implementation, use SFTP protocol to download
	fmt.Printf("Would download from %s to %s\n", remotePath, localPath)
	return nil
}

// Close closes the connection
func (c *Client) Close() error {
	if c.sshClient != nil {
		return c.sshClient.Close()
	}
	return nil
}

// ListenAndDial is a helper for SFTP operations
func listenAndDial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}
