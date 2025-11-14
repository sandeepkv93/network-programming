package vpn

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// Server represents a VPN server
type Server struct {
	address    string
	listener   net.Listener
	clients    map[string]*VPNClient
	mu         sync.RWMutex
	aead       cipher.AEAD
	subnet     *net.IPNet
	nextIP     uint32
	routeTable map[string]string // IP -> client ID mapping
}

// VPNClient represents a connected VPN client
type VPNClient struct {
	conn       net.Conn
	id         string
	assignedIP net.IP
	mu         sync.Mutex
}

// NewServer creates a new VPN server
func NewServer(address, subnet, key string) (*Server, error) {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet: %v", err)
	}

	// Create AES-GCM cipher for encryption
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Calculate starting IP (skip network address)
	startIP := binary.BigEndian.Uint32(ipnet.IP) + 1

	return &Server{
		address:    address,
		clients:    make(map[string]*VPNClient),
		aead:       aead,
		subnet:     ipnet,
		nextIP:     startIP,
		routeTable: make(map[string]string),
	}, nil
}

// Start starts the VPN server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start VPN server: %v", err)
	}

	log.Printf("VPN Server listening on %s\n", s.address)
	log.Printf("VPN Subnet: %s\n", s.subnet.String())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleClient(conn)
	}
}

// handleClient handles a VPN client connection
func (s *Server) handleClient(conn net.Conn) {
	clientID := fmt.Sprintf("client_%s", conn.RemoteAddr().String())

	// Assign IP address to client
	assignedIP := s.assignIP()
	if assignedIP == nil {
		log.Printf("No available IP addresses\n")
		conn.Close()
		return
	}

	client := &VPNClient{
		conn:       conn,
		id:         clientID,
		assignedIP: assignedIP,
	}

	s.mu.Lock()
	s.clients[clientID] = client
	s.routeTable[assignedIP.String()] = clientID
	s.mu.Unlock()

	log.Printf("Client %s connected, assigned IP: %s\n", clientID, assignedIP.String())

	// Send assigned IP to client
	if err := s.sendConfig(conn, assignedIP); err != nil {
		log.Printf("Failed to send config: %v\n", err)
		s.removeClient(clientID)
		return
	}

	defer s.removeClient(clientID)

	// Handle client packets
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from client %s: %v\n", clientID, err)
			}
			break
		}

		// Decrypt packet
		packet, err := s.decrypt(buffer[:n])
		if err != nil {
			log.Printf("Failed to decrypt packet: %v\n", err)
			continue
		}

		// Route packet to destination
		s.routePacket(packet, clientID)
	}
}

// sendConfig sends VPN configuration to client
func (s *Server) sendConfig(conn net.Conn, ip net.IP) error {
	config := make([]byte, 8)
	copy(config[0:4], ip.To4())
	copy(config[4:8], s.subnet.Mask)

	encrypted, err := s.encrypt(config)
	if err != nil {
		return err
	}

	_, err = conn.Write(encrypted)
	return err
}

// assignIP assigns a new IP address to a client
func (s *Server) assignIP() net.IP {
	s.mu.Lock()
	defer s.mu.Unlock()

	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, s.nextIP)

	if !s.subnet.Contains(ip) {
		return nil
	}

	s.nextIP++
	return ip
}

// routePacket routes a packet to the appropriate client
func (s *Server) routePacket(packet []byte, senderID string) {
	if len(packet) < 20 {
		return // Invalid IP packet
	}

	// Extract destination IP (bytes 16-19 in IP header)
	destIP := net.IP(packet[16:20]).String()

	s.mu.RLock()
	clientID, exists := s.routeTable[destIP]
	s.mu.RUnlock()

	if !exists {
		log.Printf("No route to destination: %s\n", destIP)
		return
	}

	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()

	if !exists {
		return
	}

	// Encrypt and send packet
	encrypted, err := s.encrypt(packet)
	if err != nil {
		log.Printf("Failed to encrypt packet: %v\n", err)
		return
	}

	client.mu.Lock()
	_, err = client.conn.Write(encrypted)
	client.mu.Unlock()

	if err != nil {
		log.Printf("Failed to send packet to %s: %v\n", clientID, err)
	}
}

// removeClient removes a client from the server
func (s *Server) removeClient(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, exists := s.clients[clientID]; exists {
		client.conn.Close()
		delete(s.routeTable, client.assignedIP.String())
		delete(s.clients, clientID)
		log.Printf("Client %s disconnected\n", clientID)
	}
}

// encrypt encrypts data using AES-GCM
func (s *Server) encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := s.aead.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func (s *Server) decrypt(data []byte) ([]byte, error) {
	if len(data) < s.aead.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:s.aead.NonceSize()], data[s.aead.NonceSize():]
	return s.aead.Open(nil, nonce, ciphertext, nil)
}

// Stop stops the VPN server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// GetConnectedClients returns the number of connected clients
func (s *Server) GetConnectedClients() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

// GetClientList returns a list of connected clients and their IPs
func (s *Server) GetClientList() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make(map[string]string)
	for id, client := range s.clients {
		list[id] = client.assignedIP.String()
	}
	return list
}
