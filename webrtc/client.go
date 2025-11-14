package webrtc

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

// Peer represents a WebRTC peer
type Peer struct {
	signalingURL string
	pc           *webrtc.PeerConnection
	ws           *websocket.Conn
	dataChannel  *webrtc.DataChannel
	mu           sync.Mutex
	onMessage    func(string)
}

// NewPeer creates a new WebRTC peer
func NewPeer(signalingURL string) *Peer {
	return &Peer{
		signalingURL: signalingURL,
	}
}

// Connect connects to the signaling server
func (p *Peer) Connect() error {
	var err error
	p.ws, _, err = websocket.DefaultDialer.Dial(p.signalingURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to signaling server: %v", err)
	}

	log.Println("Connected to signaling server")
	return nil
}

// CreatePeerConnection creates a new WebRTC peer connection
func (p *Peer) CreatePeerConnection() error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	var err error
	p.pc, err = webrtc.NewPeerConnection(config)
	if err != nil {
		return fmt.Errorf("failed to create peer connection: %v", err)
	}

	// Set up ICE candidate handler
	p.pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		candidateInit := candidate.ToJSON()
		msg := SignalMessage{
			Type:      "candidate",
			Candidate: &candidateInit,
		}

		p.mu.Lock()
		defer p.mu.Unlock()
		if err := p.ws.WriteJSON(msg); err != nil {
			log.Printf("Failed to send ICE candidate: %v\n", err)
		}
	})

	// Set up connection state handler
	p.pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("Peer connection state changed: %s\n", state.String())
	})

	log.Println("Peer connection created")
	return nil
}

// CreateOffer creates and sends an offer to establish a connection
func (p *Peer) CreateOffer() error {
	// Create a data channel
	var err error
	p.dataChannel, err = p.pc.CreateDataChannel("data", nil)
	if err != nil {
		return fmt.Errorf("failed to create data channel: %v", err)
	}

	p.setupDataChannel(p.dataChannel)

	// Create offer
	offer, err := p.pc.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("failed to create offer: %v", err)
	}

	// Set local description
	if err := p.pc.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("failed to set local description: %v", err)
	}

	// Send offer through signaling server
	msg := SignalMessage{
		Type: "offer",
		SDP:  &offer,
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ws.WriteJSON(msg); err != nil {
		return fmt.Errorf("failed to send offer: %v", err)
	}

	log.Println("Offer created and sent")
	return nil
}

// HandleSignaling handles incoming signaling messages
func (p *Peer) HandleSignaling() error {
	for {
		var msg SignalMessage
		if err := p.ws.ReadJSON(&msg); err != nil {
			return fmt.Errorf("failed to read signaling message: %v", err)
		}

		log.Printf("Received %s message\n", msg.Type)

		switch msg.Type {
		case "offer":
			if err := p.handleOffer(msg.SDP); err != nil {
				log.Printf("Failed to handle offer: %v\n", err)
			}
		case "answer":
			if err := p.handleAnswer(msg.SDP); err != nil {
				log.Printf("Failed to handle answer: %v\n", err)
			}
		case "candidate":
			if err := p.handleCandidate(msg.Candidate); err != nil {
				log.Printf("Failed to handle ICE candidate: %v\n", err)
			}
		}
	}
}

// handleOffer handles incoming offer messages
func (p *Peer) handleOffer(sdp *webrtc.SessionDescription) error {
	if sdp == nil {
		return fmt.Errorf("received nil SDP in offer")
	}

	// Set remote description
	if err := p.pc.SetRemoteDescription(*sdp); err != nil {
		return fmt.Errorf("failed to set remote description: %v", err)
	}

	// Set up data channel handler for the answering peer
	p.pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.Printf("Data channel '%s' opened\n", dc.Label())
		p.dataChannel = dc
		p.setupDataChannel(dc)
	})

	// Create answer
	answer, err := p.pc.CreateAnswer(nil)
	if err != nil {
		return fmt.Errorf("failed to create answer: %v", err)
	}

	// Set local description
	if err := p.pc.SetLocalDescription(answer); err != nil {
		return fmt.Errorf("failed to set local description: %v", err)
	}

	// Send answer
	msg := SignalMessage{
		Type: "answer",
		SDP:  &answer,
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ws.WriteJSON(msg); err != nil {
		return fmt.Errorf("failed to send answer: %v", err)
	}

	log.Println("Answer created and sent")
	return nil
}

// handleAnswer handles incoming answer messages
func (p *Peer) handleAnswer(sdp *webrtc.SessionDescription) error {
	if sdp == nil {
		return fmt.Errorf("received nil SDP in answer")
	}

	if err := p.pc.SetRemoteDescription(*sdp); err != nil {
		return fmt.Errorf("failed to set remote description: %v", err)
	}

	log.Println("Remote description set from answer")
	return nil
}

// handleCandidate handles incoming ICE candidate messages
func (p *Peer) handleCandidate(candidate *webrtc.ICECandidateInit) error {
	if candidate == nil {
		return nil
	}

	if err := p.pc.AddICECandidate(*candidate); err != nil {
		return fmt.Errorf("failed to add ICE candidate: %v", err)
	}

	return nil
}

// setupDataChannel sets up event handlers for a data channel
func (p *Peer) setupDataChannel(dc *webrtc.DataChannel) {
	dc.OnOpen(func() {
		log.Printf("Data channel '%s' opened\n", dc.Label())
	})

	dc.OnClose(func() {
		log.Printf("Data channel '%s' closed\n", dc.Label())
	})

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Printf("Received message: %s\n", string(msg.Data))
		if p.onMessage != nil {
			p.onMessage(string(msg.Data))
		}
	})
}

// SendMessage sends a message through the data channel
func (p *Peer) SendMessage(message string) error {
	if p.dataChannel == nil {
		return fmt.Errorf("data channel not available")
	}

	if p.dataChannel.ReadyState() != webrtc.DataChannelStateOpen {
		return fmt.Errorf("data channel not open")
	}

	return p.dataChannel.SendText(message)
}

// SetOnMessage sets the callback for received messages
func (p *Peer) SetOnMessage(handler func(string)) {
	p.onMessage = handler
}

// Close closes the peer connection and signaling connection
func (p *Peer) Close() error {
	if p.dataChannel != nil {
		p.dataChannel.Close()
	}

	if p.pc != nil {
		if err := p.pc.Close(); err != nil {
			return err
		}
	}

	if p.ws != nil {
		return p.ws.Close()
	}

	return nil
}

// GetConnectionState returns the current connection state
func (p *Peer) GetConnectionState() webrtc.PeerConnectionState {
	if p.pc == nil {
		return webrtc.PeerConnectionStateNew
	}
	return p.pc.ConnectionState()
}

// GetStats returns statistics about the peer connection
func (p *Peer) GetStats() (string, error) {
	if p.pc == nil {
		return "", fmt.Errorf("peer connection not initialized")
	}

	stats := p.pc.GetStats()
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
