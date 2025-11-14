package videoip

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// VideoClient represents a Video over IP client
type VideoClient struct {
	ServerAddr string
	conn       *net.UDPConn
	sequence   uint16
}

// NewClient creates a new video client
func NewClient(serverAddr string) *VideoClient {
	return &VideoClient{
		ServerAddr: serverAddr,
		sequence:   0,
	}
}

// Connect connects to the video server
func (c *VideoClient) Connect() error {
	addr, err := net.ResolveUDPAddr("udp", c.ServerAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	fmt.Printf("Connected to Video server at %s\n", c.ServerAddr)
	return nil
}

// SendVideoFrame sends a video frame to the server
func (c *VideoClient) SendVideoFrame(frameData []byte, frameType uint8) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	// Create packet
	packet := make([]byte, 10+len(frameData))
	binary.BigEndian.PutUint32(packet[0:4], uint32(time.Now().UnixMilli()))
	binary.BigEndian.PutUint16(packet[4:6], c.sequence)
	packet[6] = frameType
	packet[7] = 0 // Fragment ID
	copy(packet[10:], frameData)

	c.sequence++

	// Send packet
	_, err := c.conn.Write(packet)
	return err
}

// ReceiveVideoFrame receives a video frame from the server
func (c *VideoClient) ReceiveVideoFrame() (*VideoPacket, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	buffer := make([]byte, 65535)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	if n < 10 {
		return nil, fmt.Errorf("invalid packet")
	}

	return &VideoPacket{
		Timestamp:  binary.BigEndian.Uint32(buffer[0:4]),
		Sequence:   binary.BigEndian.Uint16(buffer[4:6]),
		FrameType:  buffer[6],
		FragmentID: buffer[7],
		Data:       buffer[10:n],
	}, nil
}

// Close closes the connection
func (c *VideoClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
