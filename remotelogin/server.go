package remotelogin

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Server represents a remote login server
type Server struct {
	address  string
	listener net.Listener
	users    map[string]string // username -> password hash
	sessions map[string]*Session
	mu       sync.RWMutex
}

// Session represents an active login session
type Session struct {
	ID        string
	Username  string
	RemoteAddr string
	LoginTime time.Time
	conn      net.Conn
}

// NewServer creates a new remote login server
func NewServer(address string) *Server {
	return &Server{
		address:  address,
		users:    make(map[string]string),
		sessions: make(map[string]*Session),
	}
}

// AddUser adds a user with hashed password
func (s *Server) AddUser(username, password string) {
	hash := sha256.Sum256([]byte(password))
	s.mu.Lock()
	s.users[username] = hex.EncodeToString(hash[:])
	s.mu.Unlock()
	log.Printf("User added: %s\n", username)
}

// Start starts the remote login server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	log.Printf("Remote Login Server listening on %s\n", s.address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	// Send welcome banner
	conn.Write([]byte("Remote Login Server v1.0\n\n"))

	// Authenticate user
	username, authenticated := s.authenticate(conn, reader)
	if !authenticated {
		conn.Write([]byte("Authentication failed. Goodbye.\n"))
		return
	}

	// Create session
	sessionID := s.generateSessionID()
	session := &Session{
		ID:         sessionID,
		Username:   username,
		RemoteAddr: conn.RemoteAddr().String(),
		LoginTime:  time.Now(),
		conn:       conn,
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	log.Printf("User %s logged in from %s (Session: %s)\n", username, session.RemoteAddr, sessionID)

	defer func() {
		s.mu.Lock()
		delete(s.sessions, sessionID)
		s.mu.Unlock()
		log.Printf("User %s logged out (Session: %s)\n", username, sessionID)
	}()

	// Send welcome message
	conn.Write([]byte(fmt.Sprintf("Welcome %s!\nSession ID: %s\n\n", username, sessionID)))
	conn.Write([]byte("Type 'help' for available commands, 'exit' to logout\n\n"))

	// Handle user commands
	s.handleSession(session, reader)
}

// authenticate authenticates a user
func (s *Server) authenticate(conn net.Conn, reader *bufio.Reader) (string, bool) {
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Get username
		conn.Write([]byte("Username: "))
		username, err := reader.ReadString('\n')
		if err != nil {
			return "", false
		}
		username = strings.TrimSpace(username)

		// Get password
		conn.Write([]byte("Password: "))
		password, err := reader.ReadString('\n')
		if err != nil {
			return "", false
		}
		password = strings.TrimSpace(password)

		// Verify credentials
		hash := sha256.Sum256([]byte(password))
		passwordHash := hex.EncodeToString(hash[:])

		s.mu.RLock()
		storedHash, exists := s.users[username]
		s.mu.RUnlock()

		if exists && storedHash == passwordHash {
			return username, true
		}

		if attempt < maxAttempts {
			conn.Write([]byte(fmt.Sprintf("Invalid credentials. Attempt %d/%d\n\n", attempt, maxAttempts)))
		}
	}

	return "", false
}

// handleSession handles an authenticated session
func (s *Server) handleSession(session *Session, reader *bufio.Reader) {
	for {
		session.conn.Write([]byte(fmt.Sprintf("%s@remote$ ", session.Username)))

		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading command: %v\n", err)
			}
			break
		}

		command := strings.TrimSpace(line)
		if command == "" {
			continue
		}

		if command == "exit" || command == "logout" {
			session.conn.Write([]byte("Goodbye!\n"))
			break
		}

		s.executeCommand(session, command)
	}
}

// executeCommand executes a command in the session
func (s *Server) executeCommand(session *Session, command string) {
	switch {
	case command == "help":
		help := `Available commands:
  help      - Show this help message
  whoami    - Show current username
  session   - Show session information
  users     - List active users
  uptime    - Show login time
  pwd       - Print working directory
  ls        - List directory contents
  date      - Show current date/time
  echo      - Echo arguments
  exit      - Logout
`
		session.conn.Write([]byte(help + "\n"))

	case command == "whoami":
		session.conn.Write([]byte(session.Username + "\n"))

	case command == "session":
		info := fmt.Sprintf("Session ID: %s\nUsername: %s\nRemote Address: %s\nLogin Time: %s\n",
			session.ID, session.Username, session.RemoteAddr, session.LoginTime.Format(time.RFC1123))
		session.conn.Write([]byte(info))

	case command == "users":
		s.mu.RLock()
		session.conn.Write([]byte(fmt.Sprintf("Active sessions: %d\n", len(s.sessions))))
		for _, sess := range s.sessions {
			session.conn.Write([]byte(fmt.Sprintf("  %s (%s) - %s\n",
				sess.Username, sess.RemoteAddr, sess.LoginTime.Format(time.RFC1123))))
		}
		s.mu.RUnlock()

	case command == "uptime":
		duration := time.Since(session.LoginTime)
		session.conn.Write([]byte(fmt.Sprintf("Logged in for: %s\n", duration.Round(time.Second))))

	case strings.HasPrefix(command, "echo "):
		message := strings.TrimPrefix(command, "echo ")
		session.conn.Write([]byte(message + "\n"))

	default:
		// Try to execute as system command
		parts := strings.Fields(command)
		if len(parts) == 0 {
			return
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			session.conn.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
		} else {
			session.conn.Write(output)
		}
	}
}

// generateSessionID generates a unique session ID
func (s *Server) generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Stop stops the server
func (s *Server) Stop() error {
	// Close all sessions
	s.mu.Lock()
	for _, session := range s.sessions {
		session.conn.Close()
	}
	s.sessions = make(map[string]*Session)
	s.mu.Unlock()

	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// GetActiveSessions returns the number of active sessions
func (s *Server) GetActiveSessions() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

// KickUser forcibly disconnects a user session
func (s *Server) KickUser(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.conn.Write([]byte("\nYou have been disconnected by administrator.\n"))
	session.conn.Close()
	delete(s.sessions, sessionID)

	return nil
}
