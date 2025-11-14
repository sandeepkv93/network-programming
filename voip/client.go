package voip

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// VoIPClient represents a Voice over IP client
type VoIPClient struct {
	ServerAddr string
	conn       *net.UDPConn
	sequence   uint16
}

// NewClient creates a new VoIP client
func NewClient(serverAddr string) *VoIPClient {
	return &VoIPClient{
		ServerAddr: serverAddr,
		sequence:   0,
	}
}

// Connect connects to the VoIP server
func (c *VoIPClient) Connect() error {
	addr, err := net.ResolveUDPAddr("udp", c.ServerAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	fmt.Printf("Connected to VoIP server at %s\n", c.ServerAddr)
	return nil
}

// SendAudio sends audio data to the server
func (c *VoIPClient) SendAudio(audioData []byte) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	// Create packet
	packet := make([]byte, 8+len(audioData))
	binary.BigEndian.PutUint32(packet[0:4], uint32(time.Now().Unix()))
	binary.BigEndian.PutUint16(packet[4:6], c.sequence)
	copy(packet[8:], audioData)

	c.sequence++

	// Send packet
	_, err := c.conn.Write(packet)
	return err
}

// ReceiveAudio receives audio data from the server
func (c *VoIPClient) ReceiveAudio() ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	buffer := make([]byte, 1500)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	if n < 8 {
		return nil, fmt.Errorf("invalid packet")
	}

	// Extract audio data (skip header)
	return buffer[8:n], nil
}

// Close closes the connection
func (c *VoIPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
