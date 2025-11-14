package remoteexec

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Client represents a remote execution client
type Client struct {
	serverAddr string
	conn       net.Conn
	authToken  string
	reader     *bufio.Reader
	mu         sync.Mutex
	connected  bool
}

// NewClient creates a new remote execution client
func NewClient(serverAddr, authToken string) *Client {
	return &Client{
		serverAddr: serverAddr,
		authToken:  authToken,
	}
}

// Connect connects to the remote execution server
func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	c.reader = bufio.NewReader(c.conn)
	c.connected = true

	log.Printf("Connected to remote execution server at %s\n", c.serverAddr)
	return nil
}

// Execute executes a command on the remote server
func (c *Client) Execute(command string, args ...string) (*CommandResponse, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	// Create request
	req := CommandRequest{
		ID:      uuid.New().String(),
		Command: command,
		Args:    args,
		Token:   c.authToken,
	}

	// Send request
	c.mu.Lock()
	reqJSON, err := json.Marshal(req)
	if err != nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	_, err = c.conn.Write(append(reqJSON, '\n'))
	if err != nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	// Read response
	line, err := c.reader.ReadString('\n')
	c.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var response CommandResponse
	if err := json.Unmarshal([]byte(line), &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &response, nil
}

// ExecuteAsync executes a command asynchronously
func (c *Client) ExecuteAsync(command string, args ...string) (chan *CommandResponse, error) {
	responseChan := make(chan *CommandResponse, 1)

	go func() {
		response, err := c.Execute(command, args...)
		if err != nil {
			responseChan <- &CommandResponse{
				Success: false,
				Error:   err.Error(),
			}
		} else {
			responseChan <- response
		}
		close(responseChan)
	}()

	return responseChan, nil
}

// ExecuteWithTimeout executes a command with a timeout
func (c *Client) ExecuteWithTimeout(timeout time.Duration, command string, args ...string) (*CommandResponse, error) {
	responseChan, err := c.ExecuteAsync(command, args...)
	if err != nil {
		return nil, err
	}

	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("command execution timed out after %v", timeout)
	}
}

// ExecuteScript executes a shell script on the remote server
func (c *Client) ExecuteScript(script string) (*CommandResponse, error) {
	return c.Execute("sh", "-c", script)
}

// ExecuteBatch executes multiple commands in sequence
func (c *Client) ExecuteBatch(commands [][]string) ([]*CommandResponse, error) {
	var responses []*CommandResponse

	for i, cmd := range commands {
		if len(cmd) == 0 {
			continue
		}

		var response *CommandResponse
		var err error

		if len(cmd) == 1 {
			response, err = c.Execute(cmd[0])
		} else {
			response, err = c.Execute(cmd[0], cmd[1:]...)
		}

		if err != nil {
			return responses, fmt.Errorf("failed at command %d: %v", i, err)
		}

		responses = append(responses, response)

		// Stop on first failure
		if !response.Success {
			return responses, fmt.Errorf("command %d failed: %s", i, response.Error)
		}
	}

	return responses, nil
}

// RunInteractive runs an interactive command execution session
func (c *Client) RunInteractive() error {
	if !c.connected {
		return fmt.Errorf("not connected to server")
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Remote Execution Client - Interactive Mode")
	fmt.Println("Enter commands to execute remotely (or 'quit' to exit):")
	fmt.Println()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "quit" || input == "exit" {
			break
		}

		// Parse command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}

		// Execute command
		fmt.Printf("Executing: %s %v\n", command, args)
		response, err := c.Execute(command, args...)

		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}

		// Display results
		fmt.Printf("Status: %v\n", response.Success)
		fmt.Printf("Exit Code: %d\n", response.ExitCode)
		fmt.Printf("Duration: %s\n", response.Duration)

		if response.Output != "" {
			fmt.Printf("\nOutput:\n%s\n", response.Output)
		}

		if response.Error != "" {
			fmt.Printf("\nError:\n%s\n", response.Error)
		}

		fmt.Println()
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

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	return c.connected
}

// Ping tests the connection to the server
func (c *Client) Ping() error {
	response, err := c.ExecuteWithTimeout(5*time.Second, "echo", "ping")
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("ping failed: %s", response.Error)
	}

	return nil
}

// GetServerInfo retrieves information about the remote server
func (c *Client) GetServerInfo() (map[string]string, error) {
	info := make(map[string]string)

	// Get hostname
	response, err := c.Execute("hostname")
	if err == nil && response.Success {
		info["hostname"] = strings.TrimSpace(response.Output)
	}

	// Get OS info
	response, err = c.Execute("uname", "-a")
	if err == nil && response.Success {
		info["os"] = strings.TrimSpace(response.Output)
	}

	// Get uptime
	response, err = c.Execute("uptime")
	if err == nil && response.Success {
		info["uptime"] = strings.TrimSpace(response.Output)
	}

	// Get user
	response, err = c.Execute("whoami")
	if err == nil && response.Success {
		info["user"] = strings.TrimSpace(response.Output)
	}

	return info, nil
}
