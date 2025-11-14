package webrtc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

// SignalingServer handles WebRTC signaling
type SignalingServer struct {
	address  string
	upgrader websocket.Upgrader
	peers    map[string]*websocket.Conn
	mu       sync.Mutex
}

// SignalMessage represents a signaling message
type SignalMessage struct {
	Type      string                    `json:"type"`
	From      string                    `json:"from"`
	To        string                    `json:"to"`
	SDP       *webrtc.SessionDescription `json:"sdp,omitempty"`
	Candidate *webrtc.ICECandidateInit   `json:"candidate,omitempty"`
}

// NewSignalingServer creates a new signaling server
func NewSignalingServer(address string) *SignalingServer {
	return &SignalingServer{
		address: address,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		peers: make(map[string]*websocket.Conn),
	}
}

// Start starts the signaling server
func (s *SignalingServer) Start() error {
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/", s.handleHome)

	log.Printf("WebRTC Signaling Server listening on %s\n", s.address)
	log.Printf("WebSocket endpoint: ws://%s/ws\n", s.address)
	return http.ListenAndServe(s.address, nil)
}

// handleHome serves a simple HTML page for testing
func (s *SignalingServer) handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>WebRTC Test</title>
</head>
<body>
    <h1>WebRTC Peer Connection Test</h1>
    <div>
        <button onclick="createOffer()">Create Offer</button>
        <button onclick="createAnswer()">Create Answer</button>
    </div>
    <div>
        <textarea id="localSDP" rows="10" cols="50" placeholder="Local SDP"></textarea>
        <textarea id="remoteSDP" rows="10" cols="50" placeholder="Remote SDP"></textarea>
    </div>
    <div id="status"></div>
    <script>
        let pc = null;
        const ws = new WebSocket('ws://' + window.location.host + '/ws');

        ws.onmessage = function(event) {
            const msg = JSON.parse(event.data);
            handleSignalMessage(msg);
        };

        async function createOffer() {
            pc = new RTCPeerConnection({
                iceServers: [{urls: 'stun:stun.l.google.com:19302'}]
            });

            const offer = await pc.createOffer();
            await pc.setLocalDescription(offer);

            document.getElementById('localSDP').value = JSON.stringify(offer);
            ws.send(JSON.stringify({type: 'offer', sdp: offer}));
        }

        async function createAnswer() {
            pc = new RTCPeerConnection({
                iceServers: [{urls: 'stun:stun.l.google.com:19302'}]
            });

            const remoteSDP = JSON.parse(document.getElementById('remoteSDP').value);
            await pc.setRemoteDescription(new RTCSessionDescription(remoteSDP));

            const answer = await pc.createAnswer();
            await pc.setLocalDescription(answer);

            document.getElementById('localSDP').value = JSON.stringify(answer);
            ws.send(JSON.stringify({type: 'answer', sdp: answer}));
        }

        async function handleSignalMessage(msg) {
            if (msg.type === 'offer') {
                document.getElementById('remoteSDP').value = JSON.stringify(msg.sdp);
            } else if (msg.type === 'answer') {
                await pc.setRemoteDescription(new RTCSessionDescription(msg.sdp));
            }
        }
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// handleWebSocket handles WebSocket connections for signaling
func (s *SignalingServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	peerID := fmt.Sprintf("peer_%p", conn)
	s.mu.Lock()
	s.peers[peerID] = conn
	s.mu.Unlock()

	log.Printf("Peer connected: %s (Total peers: %d)\n", peerID, len(s.peers))

	defer func() {
		s.mu.Lock()
		delete(s.peers, peerID)
		s.mu.Unlock()
		conn.Close()
		log.Printf("Peer disconnected: %s (Total peers: %d)\n", peerID, len(s.peers))
	}()

	for {
		var msg SignalMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v\n", err)
			}
			break
		}

		msg.From = peerID
		log.Printf("Received %s message from %s\n", msg.Type, peerID)

		// Relay the message to other peers
		s.relayMessage(msg, peerID)
	}
}

// relayMessage relays a signaling message to other peers
func (s *SignalingServer) relayMessage(msg SignalMessage, senderID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, conn := range s.peers {
		if id == senderID {
			continue
		}

		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error relaying message to %s: %v\n", id, err)
		}
	}
}

// GetPeerCount returns the number of connected peers
func (s *SignalingServer) GetPeerCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.peers)
}
