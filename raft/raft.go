package raft

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// NodeState represents the state of a Raft node
type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

// RaftNode represents a node in a Raft cluster
type RaftNode struct {
	id              string
	state           NodeState
	currentTerm     int
	votedFor        string
	log             []LogEntry
	commitIndex     int
	lastApplied     int
	peers           []string
	votesReceived   int
	electionTimeout time.Duration
	mu              sync.RWMutex
	stopChan        chan bool
}

// LogEntry represents a log entry
type LogEntry struct {
	Term    int
	Command interface{}
}

// NewRaftNode creates a new Raft node
func NewRaftNode(id string) *RaftNode {
	return &RaftNode{
		id:              id,
		state:           Follower,
		currentTerm:     0,
		votedFor:        "",
		log:             make([]LogEntry, 0),
		commitIndex:     -1,
		lastApplied:     -1,
		peers:           make([]string, 0),
		electionTimeout: time.Duration(150+rand.Intn(150)) * time.Millisecond,
		stopChan:        make(chan bool),
	}
}

// AddPeer adds a peer to the cluster
func (n *RaftNode) AddPeer(peerID string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.peers = append(n.peers, peerID)
}

// Start starts the Raft node
func (n *RaftNode) Start() {
	log.Printf("[%s] Starting Raft node\n", n.id)
	go n.run()
}

// Stop stops the Raft node
func (n *RaftNode) Stop() {
	close(n.stopChan)
}

func (n *RaftNode) run() {
	for {
		select {
		case <-n.stopChan:
			return
		default:
			n.mu.RLock()
			state := n.state
			n.mu.RUnlock()

			switch state {
			case Follower:
				n.runFollower()
			case Candidate:
				n.runCandidate()
			case Leader:
				n.runLeader()
			}
		}
	}
}

func (n *RaftNode) runFollower() {
	timeout := time.After(n.electionTimeout)

	select {
	case <-timeout:
		// Election timeout - become candidate
		n.becomeCandidate()
	case <-n.stopChan:
		return
	}
}

func (n *RaftNode) runCandidate() {
	n.mu.Lock()
	n.currentTerm++
	n.votedFor = n.id
	n.votesReceived = 1 // Vote for self
	term := n.currentTerm
	n.mu.Unlock()

	log.Printf("[%s] Became candidate for term %d\n", n.id, term)

	// Request votes from peers
	votes := n.requestVotes(term)

	n.mu.Lock()
	n.votesReceived += votes
	totalVotes := n.votesReceived
	quorum := (len(n.peers) + 1) / 2 + 1
	n.mu.Unlock()

	if totalVotes >= quorum {
		n.becomeLeader()
	} else {
		// Election failed - go back to follower
		time.Sleep(n.electionTimeout)
		n.becomeFollower()
	}
}

func (n *RaftNode) runLeader() {
	log.Printf("[%s] Became leader for term %d\n", n.id, n.currentTerm)

	// Send heartbeats
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.sendHeartbeats()
		case <-n.stopChan:
			return
		}

		// Check if still leader (simplified)
		n.mu.RLock()
		if n.state != Leader {
			n.mu.RUnlock()
			return
		}
		n.mu.RUnlock()
	}
}

func (n *RaftNode) becomeFollower() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = Follower
	n.votedFor = ""
	log.Printf("[%s] Became follower\n", n.id)
}

func (n *RaftNode) becomeCandidate() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = Candidate
	log.Printf("[%s] Became candidate\n", n.id)
}

func (n *RaftNode) becomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = Leader
	log.Printf("[%s] Became leader\n", n.id)
}

func (n *RaftNode) requestVotes(term int) int {
	// In a real implementation, send RequestVote RPCs to all peers
	// Simplified: simulate some peers voting for us

	n.mu.RLock()
	peerCount := len(n.peers)
	n.mu.RUnlock()

	// Simulate majority voting for us
	return peerCount / 2
}

func (n *RaftNode) sendHeartbeats() {
	// In a real implementation, send AppendEntries RPCs to all peers
	// This prevents election timeouts

	n.mu.RLock()
	term := n.currentTerm
	n.mu.RUnlock()

	log.Printf("[%s] Sending heartbeats for term %d\n", n.id, term)
}

// AppendEntry appends a new entry to the log (only leader can do this)
func (n *RaftNode) AppendEntry(command interface{}) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.state != Leader {
		return fmt.Errorf("not the leader")
	}

	entry := LogEntry{
		Term:    n.currentTerm,
		Command: command,
	}

	n.log = append(n.log, entry)
	log.Printf("[%s] Appended entry: %v\n", n.id, command)

	return nil
}

// GetState returns the current state
func (n *RaftNode) GetState() (NodeState, int) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.state, n.currentTerm
}

// IsLeader returns whether this node is the leader
func (n *RaftNode) IsLeader() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.state == Leader
}

// String representation of NodeState
func (s NodeState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}
