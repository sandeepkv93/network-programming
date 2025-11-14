package mail

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Message represents an email message
type Message struct {
	ID        string
	From      string
	To        string
	Subject   string
	Body      string
	Timestamp time.Time
}

// Mailbox represents a user's mailbox
type Mailbox struct {
	Messages []Message
	mu       sync.RWMutex
}

// Server represents a mail server
type Server struct {
	address  string
	mailboxes map[string]*Mailbox
	listener net.Listener
	quit     chan bool
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

// NewServer creates a new mail server
func NewServer(address string) *Server {
	return &Server{
		address:   address,
		mailboxes: make(map[string]*Mailbox),
		quit:      make(chan bool),
	}
}

// Start starts the mail server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start mail server: %v", err)
	}
	s.listener = listener
	log.Printf("Mail Server listening on %s\n", s.address)

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

// acceptConnections accepts incoming connections
func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					log.Printf("Error accepting connection: %v\n", err)
					continue
				}
			}
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

// handleConnection handles a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s\n", clientAddr)

	s.sendResponse(conn, "+OK Mail Server Ready")

	reader := bufio.NewReader(conn)
	var currentUser string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected\n", clientAddr)
			return
		}

		command := strings.TrimSpace(line)
		if command == "" {
			continue
		}

		log.Printf("Command from %s: %s\n", clientAddr, command)

		parts := strings.SplitN(command, " ", 2)
		cmd := strings.ToUpper(parts[0])
		var arg string
		if len(parts) > 1 {
			arg = parts[1]
		}

		switch cmd {
		case "USER":
			if arg != "" {
				currentUser = arg
				s.ensureMailbox(currentUser)
				s.sendResponse(conn, "+OK User accepted")
			} else {
				s.sendResponse(conn, "-ERR No username provided")
			}
		case "STAT":
			if currentUser == "" {
				s.sendResponse(conn, "-ERR Not logged in")
				continue
			}
			count := s.getMessageCount(currentUser)
			s.sendResponse(conn, fmt.Sprintf("+OK %d messages", count))
		case "LIST":
			if currentUser == "" {
				s.sendResponse(conn, "-ERR Not logged in")
				continue
			}
			s.listMessages(conn, currentUser)
		case "RETR":
			if currentUser == "" {
				s.sendResponse(conn, "-ERR Not logged in")
				continue
			}
			s.retrieveMessage(conn, currentUser, arg)
		case "SEND":
			if arg != "" {
				s.handleSend(conn, reader, currentUser, arg)
			} else {
				s.sendResponse(conn, "-ERR Invalid SEND command")
			}
		case "QUIT":
			s.sendResponse(conn, "+OK Goodbye")
			return
		default:
			s.sendResponse(conn, "-ERR Unknown command")
		}
	}
}

// handleSend handles sending a message
func (s *Server) handleSend(conn net.Conn, reader *bufio.Reader, from string, to string) {
	s.sendResponse(conn, "+OK Send message (end with single '.' on a line)")

	var subject, body strings.Builder
	lineCount := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "." {
			break
		}

		if lineCount == 0 {
			subject.WriteString(trimmed)
		} else {
			body.WriteString(line)
		}
		lineCount++
	}

	msg := Message{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		From:      from,
		To:        to,
		Subject:   subject.String(),
		Body:      body.String(),
		Timestamp: time.Now(),
	}

	s.storeMessage(to, msg)
	s.sendResponse(conn, "+OK Message sent")
}

// ensureMailbox ensures a mailbox exists for a user
func (s *Server) ensureMailbox(user string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.mailboxes[user]; !exists {
		s.mailboxes[user] = &Mailbox{
			Messages: []Message{},
		}
	}
}

// storeMessage stores a message in a user's mailbox
func (s *Server) storeMessage(user string, msg Message) {
	s.ensureMailbox(user)

	s.mu.RLock()
	mailbox := s.mailboxes[user]
	s.mu.RUnlock()

	mailbox.mu.Lock()
	mailbox.Messages = append(mailbox.Messages, msg)
	mailbox.mu.Unlock()

	log.Printf("Message stored for %s from %s\n", user, msg.From)
}

// getMessageCount returns the number of messages for a user
func (s *Server) getMessageCount(user string) int {
	s.mu.RLock()
	mailbox, exists := s.mailboxes[user]
	s.mu.RUnlock()

	if !exists {
		return 0
	}

	mailbox.mu.RLock()
	defer mailbox.mu.RUnlock()
	return len(mailbox.Messages)
}

// listMessages lists all messages for a user
func (s *Server) listMessages(conn net.Conn, user string) {
	s.mu.RLock()
	mailbox, exists := s.mailboxes[user]
	s.mu.RUnlock()

	if !exists {
		s.sendResponse(conn, "+OK 0 messages")
		return
	}

	mailbox.mu.RLock()
	defer mailbox.mu.RUnlock()

	s.sendResponse(conn, fmt.Sprintf("+OK %d messages", len(mailbox.Messages)))
	for i, msg := range mailbox.Messages {
		info := fmt.Sprintf("%d: From=%s Subject=%s", i+1, msg.From, msg.Subject)
		s.sendResponse(conn, info)
	}
	s.sendResponse(conn, ".")
}

// retrieveMessage retrieves a specific message
func (s *Server) retrieveMessage(conn net.Conn, user string, indexStr string) {
	s.mu.RLock()
	mailbox, exists := s.mailboxes[user]
	s.mu.RUnlock()

	if !exists {
		s.sendResponse(conn, "-ERR No mailbox")
		return
	}

	var index int
	fmt.Sscanf(indexStr, "%d", &index)
	index-- // Convert to 0-based

	mailbox.mu.RLock()
	defer mailbox.mu.RUnlock()

	if index < 0 || index >= len(mailbox.Messages) {
		s.sendResponse(conn, "-ERR Invalid message number")
		return
	}

	msg := mailbox.Messages[index]
	s.sendResponse(conn, "+OK Message follows")
	s.sendResponse(conn, fmt.Sprintf("From: %s", msg.From))
	s.sendResponse(conn, fmt.Sprintf("To: %s", msg.To))
	s.sendResponse(conn, fmt.Sprintf("Subject: %s", msg.Subject))
	s.sendResponse(conn, fmt.Sprintf("Date: %s", msg.Timestamp.Format(time.RFC822)))
	s.sendResponse(conn, "")
	s.sendResponse(conn, msg.Body)
	s.sendResponse(conn, ".")
}

// sendResponse sends a response to the client
func (s *Server) sendResponse(conn net.Conn, message string) {
	conn.Write([]byte(message + "\r\n"))
}

// Stop stops the mail server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Println("Mail Server stopped")
}
