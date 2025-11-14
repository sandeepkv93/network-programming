package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Server represents a WebSocket server
type Server struct {
	address  string
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mu       sync.Mutex
	broadcast chan []byte
	quit     chan bool
}

// NewServer creates a new WebSocket server
func NewServer(address string) *Server {
	return &Server{
		address: address,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo purposes
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
		quit:      make(chan bool),
	}
}

// Start starts the WebSocket server
func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/", s.handleHome)

	// Start the broadcast handler
	go s.handleBroadcast()

	log.Printf("WebSocket Server listening on %s\n", s.address)
	log.Printf("WebSocket endpoint: ws://%s/ws\n", s.address)
	log.Printf("Web interface: http://%s/\n", s.address)

	return http.ListenAndServe(s.address, nil)
}

// handleHome serves a simple HTML page for testing
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Test</title>
</head>
<body>
    <h1>WebSocket Test Client</h1>
    <div>
        <input type="text" id="message" placeholder="Enter message">
        <button onclick="sendMessage()">Send</button>
    </div>
    <div id="messages"></div>
    <script>
        const ws = new WebSocket('ws://' + window.location.host + '/ws');

        ws.onopen = function() {
            addMessage('Connected to WebSocket server');
        };

        ws.onmessage = function(event) {
            addMessage('Received: ' + event.data);
        };

        ws.onclose = function() {
            addMessage('Disconnected from WebSocket server');
        };

        function sendMessage() {
            const input = document.getElementById('message');
            ws.send(input.value);
            addMessage('Sent: ' + input.value);
            input.value = '';
        }

        function addMessage(msg) {
            const div = document.createElement('div');
            div.textContent = msg;
            document.getElementById('messages').appendChild(div);
        }
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	log.Printf("Client connected from %s (Total clients: %d)\n", conn.RemoteAddr(), len(s.clients))

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		log.Printf("Client disconnected (Total clients: %d)\n", len(s.clients))
	}()

	// Read messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v\n", err)
			}
			break
		}

		log.Printf("Received message from %s: %s\n", conn.RemoteAddr(), message)

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Error writing message: %v\n", err)
			break
		}

		// Broadcast to all other clients
		s.broadcast <- message
	}
}

// handleBroadcast broadcasts messages to all connected clients
func (s *Server) handleBroadcast() {
	for {
		select {
		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Broadcast: %s", message)))
				if err != nil {
					log.Printf("Error broadcasting to client: %v\n", err)
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		case <-s.quit:
			return
		}
	}
}

// Stop stops the WebSocket server
func (s *Server) Stop() {
	close(s.quit)
	s.mu.Lock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[*websocket.Conn]bool)
	s.mu.Unlock()
}

// GetClientCount returns the number of connected clients
func (s *Server) GetClientCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
}
