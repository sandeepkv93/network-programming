package gossip

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

// Node represents a gossip protocol node
type Node struct {
	ID      string
	Address string
	peers   map[string]string // ID -> Address
	data    map[string]interface{}
	mu      sync.RWMutex
	conn    *net.UDPConn
}

// Message represents a gossip message
type Message struct {
	Type      string                 `json:"type"`
	SenderID  string                 `json:"sender_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewNode creates a new gossip node
func NewNode(id, address string) *Node {
	return &Node{
		ID:      id,
		Address: address,
		peers:   make(map[string]string),
		data:    make(map[string]interface{}),
	}
}

// Start starts the gossip node
func (n *Node) Start() error {
	addr, err := net.ResolveUDPAddr("udp", n.Address)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	n.conn = conn
	log.Printf("Gossip node %s started on %s\n", n.ID, n.Address)

	// Start listening for messages
	go n.listen()

	// Start periodic gossip
	go n.periodicGossip()

	return nil
}

// AddPeer adds a peer to the node
func (n *Node) AddPeer(id, address string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.peers[id] = address
	log.Printf("Node %s added peer %s at %s\n", n.ID, id, address)
}

// SetData sets a key-value pair in the node's data
func (n *Node) SetData(key string, value interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.data[key] = value
	log.Printf("Node %s set %s = %v\n", n.ID, key, value)
}

// GetData gets a value by key
func (n *Node) GetData(key string) (interface{}, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	val, exists := n.data[key]
	return val, exists
}

func (n *Node) listen() {
	buffer := make([]byte, 65535)

	for {
		bytesRead, _, err := n.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP: %v\n", err)
			continue
		}

		var msg Message
		if err := json.Unmarshal(buffer[:bytesRead], &msg); err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		go n.handleMessage(msg)
	}
}

func (n *Node) handleMessage(msg Message) {
	if msg.SenderID == n.ID {
		return // Ignore messages from self
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	log.Printf("Node %s received gossip from %s\n", n.ID, msg.SenderID)

	// Merge received data
	for key, value := range msg.Data {
		// Simple last-write-wins strategy
		if existing, exists := n.data[key]; !exists {
			n.data[key] = value
			log.Printf("Node %s learned %s = %v from %s\n", n.ID, key, value, msg.SenderID)
		} else {
			// Could implement vector clocks or other conflict resolution here
			_ = existing
		}
	}

	// Add sender as peer if not already known
	if _, exists := n.peers[msg.SenderID]; !exists {
		// We would need sender's address in message for this
	}
}

func (n *Node) periodicGossip() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		n.gossip()
	}
}

func (n *Node) gossip() {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if len(n.peers) == 0 {
		return
	}

	// Select random peer(s) to gossip to
	fanout := 3 // Number of peers to gossip to
	if fanout > len(n.peers) {
		fanout = len(n.peers)
	}

	// Get random peers
	peerList := make([]string, 0, len(n.peers))
	for _, addr := range n.peers {
		peerList = append(peerList, addr)
	}

	// Shuffle and select
	rand.Shuffle(len(peerList), func(i, j int) {
		peerList[i], peerList[j] = peerList[j], peerList[i]
	})

	// Create gossip message
	msg := Message{
		Type:      "gossip",
		SenderID:  n.ID,
		Data:      n.data,
		Timestamp: time.Now(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v\n", err)
		return
	}

	// Send to selected peers
	for i := 0; i < fanout; i++ {
		peerAddr, err := net.ResolveUDPAddr("udp", peerList[i])
		if err != nil {
			log.Printf("Error resolving peer address: %v\n", err)
			continue
		}

		_, err = n.conn.WriteToUDP(msgBytes, peerAddr)
		if err != nil {
			log.Printf("Error sending to peer: %v\n", err)
		}
	}

	log.Printf("Node %s gossiped to %d peers\n", n.ID, fanout)
}

// Stop stops the gossip node
func (n *Node) Stop() error {
	if n.conn != nil {
		return n.conn.Close()
	}
	return nil
}

// GetAllData returns all data stored in the node
func (n *Node) GetAllData() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()

	data := make(map[string]interface{})
	for k, v := range n.data {
		data[k] = v
	}
	return data
}
