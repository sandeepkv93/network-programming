package dht

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"sort"
	"sync"
)

// DHT represents a Distributed Hash Table
type DHT struct {
	nodeID   *big.Int
	address  string
	fingers  []*Node // Finger table for Chord-like routing
	data     map[string]string
	peers    map[string]*Node
	mu       sync.RWMutex
	keyBits  int
}

// Node represents a DHT node
type Node struct {
	ID      *big.Int
	Address string
}

// NewDHT creates a new DHT node
func NewDHT(address string, keyBits int) *DHT {
	nodeID := hashAddress(address)

	return &DHT{
		nodeID:  nodeID,
		address: address,
		fingers: make([]*Node, keyBits),
		data:    make(map[string]string),
		peers:   make(map[string]*Node),
		keyBits: keyBits,
	}
}

// hashAddress generates a consistent hash for an address
func hashAddress(address string) *big.Int {
	h := sha1.New()
	h.Write([]byte(address))
	hashBytes := h.Sum(nil)

	// Convert to big.Int
	id := new(big.Int)
	id.SetBytes(hashBytes)

	return id
}

// hashKey generates a hash for a key
func hashKey(key string) *big.Int {
	return hashAddress(key)
}

// Join joins the DHT network through a known node
func (d *DHT) Join(knownAddress string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// In a real implementation, this would contact the known node
	// and initialize finger table
	knownNode := &Node{
		ID:      hashAddress(knownAddress),
		Address: knownAddress,
	}

	d.peers[knownAddress] = knownNode
	log.Printf("Node %s joined DHT via %s\n", d.address, knownAddress)

	return nil
}

// Put stores a key-value pair in the DHT
func (d *DHT) Put(key, value string) error {
	keyHash := hashKey(key)

	// Find responsible node
	responsible := d.findSuccessor(keyHash)

	if d.isResponsible(keyHash) {
		// We are responsible
		d.mu.Lock()
		d.data[key] = value
		d.mu.Unlock()
		log.Printf("Node %s stored key %s locally\n", d.address, key)
	} else {
		// Forward to responsible node
		log.Printf("Node %s forwarding key %s to %s\n", d.address, key, responsible.Address)
		// In a real implementation, make RPC to responsible node
	}

	return nil
}

// Get retrieves a value by key from the DHT
func (d *DHT) Get(key string) (string, error) {
	keyHash := hashKey(key)

	if d.isResponsible(keyHash) {
		// We are responsible
		d.mu.RLock()
		value, exists := d.data[key]
		d.mu.RUnlock()

		if !exists {
			return "", fmt.Errorf("key not found")
		}

		return value, nil
	}

	// Forward to responsible node
	responsible := d.findSuccessor(keyHash)
	log.Printf("Node %s forwarding lookup of %s to %s\n", d.address, key, responsible.Address)

	// In a real implementation, make RPC to responsible node
	return "", fmt.Errorf("key lookup not fully implemented in this example")
}

// isResponsible checks if this node is responsible for a key
func (d *DHT) isResponsible(keyHash *big.Int) bool {
	// Simplified: if no peers, we're responsible
	if len(d.peers) == 0 {
		return true
	}

	// In Chord, a node is responsible for keys between predecessor and itself
	// This is a simplified check
	return true
}

// findSuccessor finds the node responsible for a key (Chord algorithm)
func (d *DHT) findSuccessor(keyHash *big.Int) *Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Simplified successor lookup
	// In real Chord, this would use finger table for O(log N) lookups

	if len(d.peers) == 0 {
		return &Node{ID: d.nodeID, Address: d.address}
	}

	// Find closest node
	var closest *Node
	minDist := new(big.Int)

	for _, peer := range d.peers {
		dist := new(big.Int).Sub(keyHash, peer.ID)
		dist.Abs(dist)

		if closest == nil || dist.Cmp(minDist) < 0 {
			closest = peer
			minDist = dist
		}
	}

	return closest
}

// AddPeer adds a peer to the DHT
func (d *DHT) AddPeer(address string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	node := &Node{
		ID:      hashAddress(address),
		Address: address,
	}

	d.peers[address] = node
	log.Printf("Node %s added peer %s\n", d.address, address)

	// Update finger table (simplified)
	d.updateFingerTable()
}

// updateFingerTable updates the Chord finger table
func (d *DHT) updateFingerTable() {
	// Simplified finger table update
	// In real Chord, finger[i] points to successor of (n + 2^i) mod 2^m

	peerList := make([]*Node, 0, len(d.peers))
	for _, node := range d.peers {
		peerList = append(peerList, node)
	}

	// Sort peers by ID
	sort.Slice(peerList, func(i, j int) bool {
		return peerList[i].ID.Cmp(peerList[j].ID) < 0
	})

	// Fill finger table
	for i := 0; i < d.keyBits && i < len(peerList); i++ {
		d.fingers[i] = peerList[i]
	}
}

// GetLocalData returns all data stored locally
func (d *DHT) GetLocalData() map[string]string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	data := make(map[string]string)
	for k, v := range d.data {
		data[k] = v
	}
	return data
}

// GetNodeInfo returns information about this node
func (d *DHT) GetNodeInfo() string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return fmt.Sprintf("Node ID: %s, Address: %s, Peers: %d, Keys: %d",
		d.nodeID.String()[:16]+"...", d.address, len(d.peers), len(d.data))
}

// Helper to convert hash to uint64 for display
func hashToUint64(hash *big.Int) uint64 {
	bytes := hash.Bytes()
	if len(bytes) < 8 {
		padded := make([]byte, 8)
		copy(padded[8-len(bytes):], bytes)
		bytes = padded
	}
	return binary.BigEndian.Uint64(bytes[:8])
}
