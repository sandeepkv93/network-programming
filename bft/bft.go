package bft

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
)

// BFTNode represents a node in a Byzantine Fault Tolerant system
// This is a simplified PBFT (Practical Byzantine Fault Tolerance) implementation
type BFTNode struct {
	id            string
	isPrimary     bool
	view          int
	sequenceNum   int
	peers         []string
	prepareCount  map[string]int
	commitCount   map[string]int
	committed     map[string]bool
	mu            sync.RWMutex
}

// Message represents a BFT protocol message
type Message struct {
	Type        string // PrePrepare, Prepare, Commit
	View        int
	SequenceNum int
	Digest      string
	NodeID      string
}

// Request represents a client request
type Request struct {
	Operation string
	Timestamp int64
}

// NewBFTNode creates a new BFT node
func NewBFTNode(id string, isPrimary bool) *BFTNode {
	return &BFTNode{
		id:           id,
		isPrimary:    isPrimary,
		view:         0,
		sequenceNum:  0,
		peers:        make([]string, 0),
		prepareCount: make(map[string]int),
		commitCount:  make(map[string]int),
		committed:    make(map[string]bool),
	}
}

// AddPeer adds a peer to the network
func (n *BFTNode) AddPeer(peerID string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.peers = append(n.peers, peerID)
}

// ProcessRequest processes a client request (Primary only)
func (n *BFTNode) ProcessRequest(req Request) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.isPrimary {
		return fmt.Errorf("not the primary node")
	}

	// Generate digest
	digest := n.computeDigest(req)

	n.sequenceNum++
	seqNum := n.sequenceNum

	log.Printf("[Primary %s] Processing request, seq=%d, digest=%s\n",
		n.id, seqNum, digest)

	// Send Pre-Prepare to all replicas
	n.sendPrePrepare(seqNum, digest)

	return nil
}

// HandlePrePrepare handles Pre-Prepare message (Replica)
func (n *BFTNode) HandlePrePrepare(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.isPrimary {
		return // Primary doesn't handle Pre-Prepare
	}

	log.Printf("[Replica %s] Received Pre-Prepare seq=%d, digest=%s\n",
		n.id, msg.SequenceNum, msg.Digest)

	// Accept Pre-Prepare and send Prepare
	n.sendPrepare(msg.SequenceNum, msg.Digest)
}

// HandlePrepare handles Prepare message
func (n *BFTNode) HandlePrepare(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	key := fmt.Sprintf("%d:%s", msg.SequenceNum, msg.Digest)
	n.prepareCount[key]++

	count := n.prepareCount[key]
	totalNodes := len(n.peers) + 1
	quorum := 2*totalNodes/3 + 1

	log.Printf("[%s] Prepare count for seq=%d: %d/%d\n",
		n.id, msg.SequenceNum, count, quorum)

	// If prepared (received 2f+1 Prepare messages), send Commit
	if count >= quorum && !n.committed[key] {
		log.Printf("[%s] Prepared seq=%d, sending Commit\n", n.id, msg.SequenceNum)
		n.sendCommit(msg.SequenceNum, msg.Digest)
	}
}

// HandleCommit handles Commit message
func (n *BFTNode) HandleCommit(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	key := fmt.Sprintf("%d:%s", msg.SequenceNum, msg.Digest)
	n.commitCount[key]++

	count := n.commitCount[key]
	totalNodes := len(n.peers) + 1
	quorum := 2*totalNodes/3 + 1

	log.Printf("[%s] Commit count for seq=%d: %d/%d\n",
		n.id, msg.SequenceNum, count, quorum)

	// If committed (received 2f+1 Commit messages), execute
	if count >= quorum && !n.committed[key] {
		n.committed[key] = true
		log.Printf("[%s] ✓ Committed and executed seq=%d\n", n.id, msg.SequenceNum)
	}
}

func (n *BFTNode) sendPrePrepare(seqNum int, digest string) {
	// In real implementation, send to all replicas
	log.Printf("[Primary %s] Sent Pre-Prepare: seq=%d, digest=%s\n", n.id, seqNum, digest)
}

func (n *BFTNode) sendPrepare(seqNum int, digest string) {
	// In real implementation, send to all nodes
	log.Printf("[%s] Sent Prepare: seq=%d\n", n.id, seqNum)

	// Simulate receiving own prepare
	key := fmt.Sprintf("%d:%s", seqNum, digest)
	n.prepareCount[key]++
}

func (n *BFTNode) sendCommit(seqNum int, digest string) {
	// In real implementation, send to all nodes
	log.Printf("[%s] Sent Commit: seq=%d\n", n.id, seqNum)

	// Simulate receiving own commit
	key := fmt.Sprintf("%d:%s", seqNum, digest)
	n.commitCount[key]++
}

func (n *BFTNode) computeDigest(req Request) string {
	data := fmt.Sprintf("%s:%d", req.Operation, req.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}

// IsPrimary returns whether this node is the primary
func (n *BFTNode) IsPrimary() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.isPrimary
}

// GetView returns the current view number
func (n *BFTNode) GetView() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.view
}

// SimulateByzantineNode simulates a Byzantine (faulty/malicious) node
func (n *BFTNode) SimulateByzantineNode() {
	n.mu.Lock()
	defer n.mu.Unlock()

	log.Printf("[%s] ⚠️ Node turned Byzantine (malicious)\n", n.id)
	// Could send conflicting messages, remain silent, etc.
}
