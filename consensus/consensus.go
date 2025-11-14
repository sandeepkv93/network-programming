package consensus

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ConsensusProtocol represents a generic consensus protocol
type ConsensusProtocol interface {
	Propose(value interface{}) error
	GetValue() (interface{}, bool)
	GetState() string
}

// SimpleConsensus implements a simple majority-based consensus
type SimpleConsensus struct {
	nodeID    string
	peers     []string
	proposals map[string]int // value -> count
	decided   bool
	value     interface{}
	mu        sync.RWMutex
}

// NewSimpleConsensus creates a new simple consensus instance
func NewSimpleConsensus(nodeID string) *SimpleConsensus {
	return &SimpleConsensus{
		nodeID:    nodeID,
		peers:     make([]string, 0),
		proposals: make(map[string]int),
		decided:   false,
	}
}

// AddPeer adds a peer to the consensus group
func (c *SimpleConsensus) AddPeer(peerID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.peers = append(c.peers, peerID)
}

// Propose proposes a value for consensus
func (c *SimpleConsensus) Propose(value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.decided {
		return fmt.Errorf("consensus already reached")
	}

	valueStr := fmt.Sprintf("%v", value)
	c.proposals[valueStr]++

	log.Printf("[%s] Proposed value: %v (count: %d)\n", c.nodeID, value, c.proposals[valueStr])

	// Check if majority reached
	totalNodes := len(c.peers) + 1
	majority := totalNodes/2 + 1

	if c.proposals[valueStr] >= majority {
		c.decided = true
		c.value = value
		log.Printf("[%s] ✓ Consensus reached on value: %v\n", c.nodeID, value)
	}

	return nil
}

// GetValue returns the consensus value if decided
func (c *SimpleConsensus) GetValue() (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value, c.decided
}

// GetState returns the current consensus state
func (c *SimpleConsensus) GetState() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.decided {
		return fmt.Sprintf("Decided: %v", c.value)
	}
	return "Undecided"
}

// QuorumConsensus implements quorum-based consensus
type QuorumConsensus struct {
	nodeID      string
	peers       []string
	votes       map[string]map[string]bool // value -> (nodeID -> voted)
	quorumSize  int
	decided     bool
	value       interface{}
	mu          sync.RWMutex
}

// NewQuorumConsensus creates a new quorum-based consensus
func NewQuorumConsensus(nodeID string, quorumSize int) *QuorumConsensus {
	return &QuorumConsensus{
		nodeID:     nodeID,
		peers:      make([]string, 0),
		votes:      make(map[string]map[string]bool),
		quorumSize: quorumSize,
		decided:    false,
	}
}

// AddPeer adds a peer
func (q *QuorumConsensus) AddPeer(peerID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.peers = append(q.peers, peerID)
}

// Vote votes for a value
func (q *QuorumConsensus) Vote(value interface{}, voterID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.decided {
		return fmt.Errorf("consensus already reached")
	}

	valueStr := fmt.Sprintf("%v", value)

	if _, exists := q.votes[valueStr]; !exists {
		q.votes[valueStr] = make(map[string]bool)
	}

	q.votes[valueStr][voterID] = true

	voteCount := len(q.votes[valueStr])
	log.Printf("[%s] Vote for %v from %s (count: %d/%d)\n",
		q.nodeID, value, voterID, voteCount, q.quorumSize)

	// Check if quorum reached
	if voteCount >= q.quorumSize {
		q.decided = true
		q.value = value
		log.Printf("[%s] ✓ Quorum reached on value: %v\n", q.nodeID, value)
	}

	return nil
}

// GetValue returns the consensus value if decided
func (q *QuorumConsensus) GetValue() (interface{}, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.value, q.decided
}

// GetState returns the current state
func (q *QuorumConsensus) GetState() string {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.decided {
		return fmt.Sprintf("Decided: %v", q.value)
	}
	return "Undecided"
}

// TimedConsensus implements time-based consensus with timeout
type TimedConsensus struct {
	*SimpleConsensus
	timeout  time.Duration
	deadline time.Time
}

// NewTimedConsensus creates a timed consensus instance
func NewTimedConsensus(nodeID string, timeout time.Duration) *TimedConsensus {
	return &TimedConsensus{
		SimpleConsensus: NewSimpleConsensus(nodeID),
		timeout:         timeout,
		deadline:        time.Now().Add(timeout),
	}
}

// Propose proposes with timeout check
func (t *TimedConsensus) Propose(value interface{}) error {
	if time.Now().After(t.deadline) {
		return fmt.Errorf("consensus timeout expired")
	}

	return t.SimpleConsensus.Propose(value)
}

// IsExpired checks if consensus has expired
func (t *TimedConsensus) IsExpired() bool {
	return time.Now().After(t.deadline)
}

// ConsensusBuilder helps build different consensus configurations
type ConsensusBuilder struct {
	nodeID string
	peers  []string
}

// NewConsensusBuilder creates a new builder
func NewConsensusBuilder(nodeID string) *ConsensusBuilder {
	return &ConsensusBuilder{
		nodeID: nodeID,
		peers:  make([]string, 0),
	}
}

// WithPeers sets the peers
func (b *ConsensusBuilder) WithPeers(peers []string) *ConsensusBuilder {
	b.peers = peers
	return b
}

// BuildSimple builds a simple consensus
func (b *ConsensusBuilder) BuildSimple() *SimpleConsensus {
	c := NewSimpleConsensus(b.nodeID)
	for _, peer := range b.peers {
		c.AddPeer(peer)
	}
	return c
}

// BuildQuorum builds a quorum consensus
func (b *ConsensusBuilder) BuildQuorum(quorumSize int) *QuorumConsensus {
	c := NewQuorumConsensus(b.nodeID, quorumSize)
	for _, peer := range b.peers {
		c.AddPeer(peer)
	}
	return c
}

// BuildTimed builds a timed consensus
func (b *ConsensusBuilder) BuildTimed(timeout time.Duration) *TimedConsensus {
	c := NewTimedConsensus(b.nodeID, timeout)
	for _, peer := range b.peers {
		c.AddPeer(peer)
	}
	return c
}
