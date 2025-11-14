package videoconference

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ConferenceServer represents a video conferencing server
type ConferenceServer struct {
	Address string
	rooms   map[string]*Room
	mu      sync.RWMutex
	upgrader websocket.Upgrader
	server  *http.Server
}

// Room represents a conference room
type Room struct {
	ID      string
	clients map[*websocket.Conn]string
	mu      sync.Mutex
}

// Message represents a signaling message
type Message struct {
	Type     string          `json:"type"`
	From     string          `json:"from"`
	To       string          `json:"to"`
	RoomID   string          `json:"roomId"`
	Payload  json.RawMessage `json:"payload"`
}

// NewServer creates a new video conferencing server
func NewServer(address string) *ConferenceServer {
	return &ConferenceServer{
		Address: address,
		rooms:   make(map[string]*Room),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// Start starts the conference server
func (s *ConferenceServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/", s.handleIndex)

	s.server = &http.Server{
		Addr:    s.Address,
		Handler: mux,
	}

	log.Printf("Video Conference Server started on %s\n", s.Address)
	return s.server.ListenAndServe()
}

// Stop stops the conference server
func (s *ConferenceServer) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func (s *ConferenceServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head><title>Video Conference</title></head>
<body>
    <h1>Video Conference Server</h1>
    <p>WebSocket endpoint: ws://` + s.Address + `/ws</p>
    <p>Connect your WebRTC client to join a conference room.</p>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func (s *ConferenceServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v\n", err)
		return
	}
	defer conn.Close()

	log.Printf("New WebSocket connection from %s\n", conn.RemoteAddr())

	var currentRoom *Room
	var clientID string

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			break
		}

		switch msg.Type {
		case "join":
			clientID = msg.From
			currentRoom = s.joinRoom(msg.RoomID, conn, clientID)
			log.Printf("Client %s joined room %s\n", clientID, msg.RoomID)

		case "offer", "answer", "ice-candidate":
			if currentRoom != nil {
				currentRoom.broadcast(msg, conn)
			}
		}
	}

	// Cleanup
	if currentRoom != nil {
		currentRoom.removeClient(conn)
		log.Printf("Client %s left room\n", clientID)
	}
}

func (s *ConferenceServer) joinRoom(roomID string, conn *websocket.Conn, clientID string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[roomID]
	if !exists {
		room = &Room{
			ID:      roomID,
			clients: make(map[*websocket.Conn]string),
		}
		s.rooms[roomID] = room
	}

	room.addClient(conn, clientID)
	return room
}

func (r *Room) addClient(conn *websocket.Conn, clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[conn] = clientID
}

func (r *Room) removeClient(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, conn)
}

func (r *Room) broadcast(msg Message, sender *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for conn, _ := range r.clients {
		if conn != sender {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error broadcasting: %v\n", err)
			}
		}
	}
}
