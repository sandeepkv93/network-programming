package vpn

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// Client represents a VPN client
type Client struct {
	serverAddr string
	conn       net.Conn
	aead       cipher.AEAD
	assignedIP net.IP
	netmask    net.IPMask
	connected  bool
}

// NewClient creates a new VPN client
func NewClient(serverAddr, key string) (*Client, error) {
	// Create AES-GCM cipher for encryption
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	return &Client{
		serverAddr: serverAddr,
		aead:       aead,
	}, nil
}

// Connect connects to the VPN server
func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to VPN server: %v", err)
	}

	log.Printf("Connected to VPN server at %s\n", c.serverAddr)

	// Receive VPN configuration
	if err := c.receiveConfig(); err != nil {
		c.conn.Close()
		return fmt.Errorf("failed to receive VPN config: %v", err)
	}

	c.connected = true
	log.Printf("VPN tunnel established. Assigned IP: %s/%d\n",
		c.assignedIP.String(), c.maskSize())

	return nil
}

// receiveConfig receives VPN configuration from server
func (c *Client) receiveConfig() error {
	buffer := make([]byte, 1024)
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := c.conn.Read(buffer)
	c.conn.SetReadDeadline(time.Time{})

	if err != nil {
		return err
	}

	// Decrypt configuration
	config, err := c.decrypt(buffer[:n])
	if err != nil {
		return err
	}

	if len(config) < 8 {
		return fmt.Errorf("invalid config size")
	}

	c.assignedIP = net.IP(config[0:4])
	c.netmask = net.IPMask(config[4:8])

	return nil
}

// SendPacket sends an IP packet through the VPN tunnel
func (c *Client) SendPacket(packet []byte) error {
	if !c.connected {
		return fmt.Errorf("not connected to VPN server")
	}

	// Encrypt packet
	encrypted, err := c.encrypt(packet)
	if err != nil {
		return fmt.Errorf("failed to encrypt packet: %v", err)
	}

	_, err = c.conn.Write(encrypted)
	if err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
	}

	return nil
}

// ReceivePacket receives an IP packet from the VPN tunnel
func (c *Client) ReceivePacket() ([]byte, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to VPN server")
	}

	buffer := make([]byte, 2048)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Decrypt packet
	packet, err := c.decrypt(buffer[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt packet: %v", err)
	}

	return packet, nil
}

// StartTunnel starts the VPN tunnel and handles packet forwarding
func (c *Client) StartTunnel() error {
	if !c.connected {
		return fmt.Errorf("not connected to VPN server")
	}

	log.Println("VPN tunnel active - receiving packets...")

	for {
		packet, err := c.ReceivePacket()
		if err != nil {
			if err == io.EOF {
				log.Println("VPN server closed connection")
				break
			}
			log.Printf("Error receiving packet: %v\n", err)
			continue
		}

		// In a real implementation, this would write to a TUN device
		// For demo purposes, just log the packet
		if len(packet) >= 20 {
			srcIP := net.IP(packet[12:16])
			dstIP := net.IP(packet[16:20])
			log.Printf("Received packet: %s -> %s (%d bytes)\n",
				srcIP.String(), dstIP.String(), len(packet))
		}
	}

	return nil
}

// encrypt encrypts data using AES-GCM
func (c *Client) encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := c.aead.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func (c *Client) decrypt(data []byte) ([]byte, error) {
	if len(data) < c.aead.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:c.aead.NonceSize()], data[c.aead.NonceSize():]
	return c.aead.Open(nil, nonce, ciphertext, nil)
}

// Disconnect disconnects from the VPN server
func (c *Client) Disconnect() error {
	c.connected = false
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetAssignedIP returns the assigned VPN IP address
func (c *Client) GetAssignedIP() string {
	if c.assignedIP == nil {
		return ""
	}
	return c.assignedIP.String()
}

// GetNetmask returns the VPN netmask
func (c *Client) GetNetmask() string {
	if c.netmask == nil {
		return ""
	}
	return net.IP(c.netmask).String()
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	return c.connected
}

// maskSize returns the number of bits in the netmask
func (c *Client) maskSize() int {
	if c.netmask == nil {
		return 0
	}
	ones, _ := c.netmask.Size()
	return ones
}

// Ping sends a ping through the VPN tunnel to test connectivity
func (c *Client) Ping(destIP string) error {
	if !c.connected {
		return fmt.Errorf("not connected to VPN server")
	}

	// Create a simple ICMP echo request packet
	// This is a simplified version - real implementation would use proper ICMP
	dip := net.ParseIP(destIP).To4()
	if dip == nil {
		return fmt.Errorf("invalid destination IP")
	}

	// Simplified IP packet (normally would include proper ICMP payload)
	packet := make([]byte, 20)
	packet[0] = 0x45 // Version 4, header length 5
	copy(packet[12:16], c.assignedIP.To4())
	copy(packet[16:20], dip)

	return c.SendPacket(packet)
}
