package dhcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Client represents a DHCP client
type Client struct {
	interfaceName string
	mac           net.HardwareAddr
	conn          *net.UDPConn
	assignedIP    net.IP
	serverIP      net.IP
	leaseTime     time.Duration
}

// NewClient creates a new DHCP client
func NewClient(interfaceName string) (*Client, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get interface: %v", err)
	}

	return &Client{
		interfaceName: interfaceName,
		mac:           iface.HardwareAddr,
	}, nil
}

// Discover sends a DHCP DISCOVER message
func (c *Client) Discover() error {
	// Create UDP connection
	addr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 68,
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %v", err)
	}
	c.conn = conn

	// Build DHCP DISCOVER packet
	packet := c.buildDiscoverPacket()

	// Send to broadcast address
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 67,
	}

	_, err = c.conn.WriteToUDP(packet, serverAddr)
	if err != nil {
		return fmt.Errorf("failed to send DISCOVER: %v", err)
	}

	return nil
}

// WaitForOffer waits for a DHCP OFFER message
func (c *Client) WaitForOffer(timeout time.Duration) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	c.conn.SetReadDeadline(time.Now().Add(timeout))

	buffer := make([]byte, 1500)
	n, _, err := c.conn.ReadFromUDP(buffer)
	if err != nil {
		return fmt.Errorf("failed to receive OFFER: %v", err)
	}

	if n < 240 {
		return fmt.Errorf("invalid DHCP packet")
	}

	// Parse OFFER (simplified)
	// Your IP address (yiaddr) at offset 16
	c.assignedIP = net.IP(buffer[16:20])

	// Server IP address (siaddr) at offset 20
	c.serverIP = net.IP(buffer[20:24])

	return nil
}

// Request sends a DHCP REQUEST message
func (c *Client) Request() error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	// Build DHCP REQUEST packet
	packet := c.buildRequestPacket()

	// Send to broadcast address
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 67,
	}

	_, err := c.conn.WriteToUDP(packet, serverAddr)
	if err != nil {
		return fmt.Errorf("failed to send REQUEST: %v", err)
	}

	return nil
}

// GetIP performs the full DHCP process and returns assigned IP
func (c *Client) GetIP(timeout time.Duration) (net.IP, error) {
	// Send DISCOVER
	if err := c.Discover(); err != nil {
		return nil, err
	}

	// Wait for OFFER
	if err := c.WaitForOffer(timeout); err != nil {
		return nil, err
	}

	// Send REQUEST
	if err := c.Request(); err != nil {
		return nil, err
	}

	return c.assignedIP, nil
}

// buildDiscoverPacket builds a DHCP DISCOVER packet
func (c *Client) buildDiscoverPacket() []byte {
	packet := make([]byte, 240)

	// Boot request
	packet[0] = 1

	// Hardware type (Ethernet)
	packet[1] = 1

	// Hardware address length
	packet[2] = 6

	// Transaction ID (random)
	binary.BigEndian.PutUint32(packet[4:8], uint32(time.Now().Unix()))

	// Copy MAC address (chaddr)
	copy(packet[28:34], c.mac)

	// Magic cookie would go here in full implementation
	// Options would follow in full implementation

	return packet
}

// buildRequestPacket builds a DHCP REQUEST packet
func (c *Client) buildRequestPacket() []byte {
	packet := make([]byte, 240)

	// Boot request
	packet[0] = 1

	// Hardware type (Ethernet)
	packet[1] = 1

	// Hardware address length
	packet[2] = 6

	// Transaction ID
	binary.BigEndian.PutUint32(packet[4:8], uint32(time.Now().Unix()))

	// Copy requested IP
	if c.assignedIP != nil {
		copy(packet[12:16], c.assignedIP.To4())
	}

	// Copy MAC address
	copy(packet[28:34], c.mac)

	return packet
}

// Release releases the DHCP lease
func (c *Client) Release() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.assignedIP = nil
	return nil
}

// GetAssignedIP returns the assigned IP address
func (c *Client) GetAssignedIP() net.IP {
	return c.assignedIP
}

// GetServerIP returns the DHCP server IP address
func (c *Client) GetServerIP() net.IP {
	return c.serverIP
}
