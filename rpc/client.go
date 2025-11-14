package rpc

import (
	"fmt"
	"log"
	"net/rpc"
)

// Client represents an RPC client
type Client struct {
	serverAddr string
	client     *rpc.Client
	connected  bool
}

// NewClient creates a new RPC client
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
	}
}

// Connect connects to the RPC server
func (c *Client) Connect() error {
	var err error
	c.client, err = rpc.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC server: %v", err)
	}

	c.connected = true
	log.Printf("Connected to RPC server at %s\n", c.serverAddr)
	return nil
}

// Call makes an RPC call
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if !c.connected {
		return fmt.Errorf("not connected to server")
	}

	return c.client.Call(serviceMethod, args, reply)
}

// Add calls the Add method on the server
func (c *Client) Add(a, b int) (int, error) {
	args := &Args{A: a, B: b}
	var result Result

	err := c.Call("MathService.Add", args, &result)
	if err != nil {
		return 0, err
	}

	return result.Value, nil
}

// Subtract calls the Subtract method on the server
func (c *Client) Subtract(a, b int) (int, error) {
	args := &Args{A: a, B: b}
	var result Result

	err := c.Call("MathService.Subtract", args, &result)
	if err != nil {
		return 0, err
	}

	return result.Value, nil
}

// Multiply calls the Multiply method on the server
func (c *Client) Multiply(a, b int) (int, error) {
	args := &Args{A: a, B: b}
	var result Result

	err := c.Call("MathService.Multiply", args, &result)
	if err != nil {
		return 0, err
	}

	return result.Value, nil
}

// Divide calls the Divide method on the server
func (c *Client) Divide(a, b int) (int, error) {
	args := &Args{A: a, B: b}
	var result Result

	err := c.Call("MathService.Divide", args, &result)
	if err != nil {
		return 0, err
	}

	return result.Value, nil
}

// Close closes the connection to the server
func (c *Client) Close() error {
	c.connected = false
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	return c.connected
}
