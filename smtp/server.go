package smtp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Email represents an email message
type Email struct {
	From      string
	To        []string
	Data      string
	Timestamp time.Time
}

// Server represents an SMTP server
type Server struct {
	address  string
	hostname string
	emails   []Email
	listener net.Listener
	quit     chan bool
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

// NewServer creates a new SMTP server
func NewServer(address, hostname string) *Server {
	return &Server{
		address:  address,
		hostname: hostname,
		emails:   []Email{},
		quit:     make(chan bool),
	}
}

// Start starts the SMTP server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start SMTP server: %v", err)
	}
	s.listener = listener
	log.Printf("SMTP Server listening on %s\n", s.address)

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

// handleConnection handles a single SMTP session
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("SMTP client connected: %s\n", clientAddr)

	// Send greeting
	s.sendResponse(conn, 220, fmt.Sprintf("%s SMTP Service Ready", s.hostname))

	reader := bufio.NewReader(conn)
	var email Email
	var inData bool

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected\n", clientAddr)
			return
		}

		line = strings.TrimSpace(line)

		// Handle DATA mode
		if inData {
			if line == "." {
				s.storeEmail(email)
				s.sendResponse(conn, 250, "OK: Message accepted")
				inData = false
				email = Email{} // Reset
				continue
			}
			email.Data += line + "\n"
			continue
		}

		log.Printf("Command from %s: %s\n", clientAddr, line)

		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(parts[0])
		var arg string
		if len(parts) > 1 {
			arg = parts[1]
		}

		switch cmd {
		case "HELO", "EHLO":
			s.sendResponse(conn, 250, fmt.Sprintf("%s Hello", s.hostname))

		case "MAIL":
			if strings.HasPrefix(strings.ToUpper(arg), "FROM:") {
				from := strings.TrimPrefix(strings.ToUpper(arg), "FROM:")
				from = strings.TrimSpace(from)
				from = strings.Trim(from, "<>")
				email.From = from
				s.sendResponse(conn, 250, "OK")
			} else {
				s.sendResponse(conn, 501, "Syntax error in parameters")
			}

		case "RCPT":
			if strings.HasPrefix(strings.ToUpper(arg), "TO:") {
				to := strings.TrimPrefix(strings.ToUpper(arg), "TO:")
				to = strings.TrimSpace(to)
				to = strings.Trim(to, "<>")
				email.To = append(email.To, to)
				s.sendResponse(conn, 250, "OK")
			} else {
				s.sendResponse(conn, 501, "Syntax error in parameters")
			}

		case "DATA":
			if email.From == "" || len(email.To) == 0 {
				s.sendResponse(conn, 503, "Bad sequence of commands")
			} else {
				s.sendResponse(conn, 354, "Start mail input; end with <CRLF>.<CRLF>")
				inData = true
			}

		case "RSET":
			email = Email{}
			s.sendResponse(conn, 250, "OK")

		case "NOOP":
			s.sendResponse(conn, 250, "OK")

		case "QUIT":
			s.sendResponse(conn, 221, fmt.Sprintf("%s closing connection", s.hostname))
			return

		default:
			s.sendResponse(conn, 502, "Command not implemented")
		}
	}
}

// storeEmail stores an email
func (s *Server) storeEmail(email Email) {
	email.Timestamp = time.Now()

	s.mu.Lock()
	s.emails = append(s.emails, email)
	s.mu.Unlock()

	log.Printf("Email stored: From=%s To=%v\n", email.From, email.To)
}

// GetEmails returns all stored emails
func (s *Server) GetEmails() []Email {
	s.mu.RLock()
	defer s.mu.RUnlock()

	emails := make([]Email, len(s.emails))
	copy(emails, s.emails)
	return emails
}

// sendResponse sends an SMTP response
func (s *Server) sendResponse(conn net.Conn, code int, message string) {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	conn.Write([]byte(response))
}

// Stop stops the SMTP server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Println("SMTP Server stopped")
}
