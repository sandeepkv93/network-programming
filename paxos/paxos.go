package paxos

import (
	"fmt"
	"log"
	"sync"
)

// Paxos implements the Paxos consensus algorithm
type Paxos struct {
	nodeID           string
	proposalNumber   int
	promisedNumber   int
	acceptedNumber   int
	acceptedValue    interface{}
	peers            []string
	mu               sync.RWMutex
}

// Message represents a Paxos message
type Message struct {
	Type     string      // Prepare, Promise, Accept, Accepted
	From     string
	Number   int         // Proposal number
	Value    interface{} // Proposed value
}

// NewPaxos creates a new Paxos instance
func NewPaxos(nodeID string) *Paxos {
	return &Paxos{
		nodeID:           nodeID,
		proposalNumber:   0,
		promisedNumber:   -1,
		acceptedNumber:   -1,
		peers:            make([]string, 0),
	}
}

// AddPeer adds a peer to the Paxos group
func (p *Paxos) AddPeer(peerID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = append(p.peers, peerID)
}

// Propose proposes a value (Proposer role)
func (p *Paxos) Propose(value interface{}) (interface{}, error) {
	p.mu.Lock()
	p.proposalNumber++
	proposalNum := p.proposalNumber
	p.mu.Unlock()

	log.Printf("[%s] Proposing value %v with proposal number %d\n",
		p.nodeID, value, proposalNum)

	// Phase 1: Prepare
	promises := p.sendPrepare(proposalNum)

	if !p.hasQuorum(promises) {
		return nil, fmt.Errorf("failed to get quorum in prepare phase")
	}

	// Check if any acceptor has accepted a value
	highestNum := -1
	var highestValue interface{}

	for _, promise := range promises {
		if promise.Number > highestNum {
			highestNum = promise.Number
			highestValue = promise.Value
		}
	}

	// Use highest accepted value or our proposed value
	proposalValue := value
	if highestValue != nil {
		proposalValue = highestValue
		log.Printf("[%s] Using previously accepted value: %v\n", p.nodeID, proposalValue)
	}

	// Phase 2: Accept
	accepted := p.sendAccept(proposalNum, proposalValue)

	if !p.hasQuorum(accepted) {
		return nil, fmt.Errorf("failed to get quorum in accept phase")
	}

	log.Printf("[%s] Consensus reached on value: %v\n", p.nodeID, proposalValue)
	return proposalValue, nil
}

// HandlePrepare handles a prepare request (Acceptor role)
func (p *Paxos) HandlePrepare(proposalNum int, from string) *Message {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Promise not to accept proposals numbered less than proposalNum
	if proposalNum > p.promisedNumber {
		p.promisedNumber = proposalNum

		log.Printf("[%s] Promised to proposal %d from %s\n", p.nodeID, proposalNum, from)

		return &Message{
			Type:   "Promise",
			From:   p.nodeID,
			Number: p.acceptedNumber,
			Value:  p.acceptedValue,
		}
	}

	log.Printf("[%s] Rejected prepare %d (already promised %d)\n",
		p.nodeID, proposalNum, p.promisedNumber)

	return nil
}

// HandleAccept handles an accept request (Acceptor role)
func (p *Paxos) HandleAccept(proposalNum int, value interface{}, from string) *Message {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Accept if we haven't promised to ignore this proposal
	if proposalNum >= p.promisedNumber {
		p.promisedNumber = proposalNum
		p.acceptedNumber = proposalNum
		p.acceptedValue = value

		log.Printf("[%s] Accepted proposal %d with value %v from %s\n",
			p.nodeID, proposalNum, value, from)

		return &Message{
			Type:   "Accepted",
			From:   p.nodeID,
			Number: proposalNum,
			Value:  value,
		}
	}

	log.Printf("[%s] Rejected accept %d (promised %d)\n",
		p.nodeID, proposalNum, p.promisedNumber)

	return nil
}

// sendPrepare sends prepare requests to all acceptors
func (p *Paxos) sendPrepare(proposalNum int) []*Message {
	// In a real implementation, this would send network messages
	// For this simplified version, we simulate responses

	responses := make([]*Message, 0)

	// Simulate majority responding with promises
	quorumSize := (len(p.peers) + 1) / 2 + 1
	for i := 0; i < quorumSize; i++ {
		responses = append(responses, &Message{
			Type:   "Promise",
			From:   fmt.Sprintf("peer-%d", i),
			Number: -1,
			Value:  nil,
		})
	}

	return responses
}

// sendAccept sends accept requests to all acceptors
func (p *Paxos) sendAccept(proposalNum int, value interface{}) []*Message {
	// In a real implementation, this would send network messages
	// For this simplified version, we simulate responses

	responses := make([]*Message, 0)

	// Simulate majority accepting
	quorumSize := (len(p.peers) + 1) / 2 + 1
	for i := 0; i < quorumSize; i++ {
		responses = append(responses, &Message{
			Type:   "Accepted",
			From:   fmt.Sprintf("peer-%d", i),
			Number: proposalNum,
			Value:  value,
		})
	}

	return responses
}

// hasQuorum checks if we have a quorum of responses
func (p *Paxos) hasQuorum(responses []*Message) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	totalNodes := len(p.peers) + 1 // peers + self
	quorumSize := totalNodes/2 + 1

	return len(responses) >= quorumSize
}

// GetAcceptedValue returns the currently accepted value
func (p *Paxos) GetAcceptedValue() (interface{}, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.acceptedNumber >= 0 {
		return p.acceptedValue, true
	}
	return nil, false
}
