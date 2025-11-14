package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Server represents an SSH server
type Server struct {
	address  string
	config   *ssh.ServerConfig
	listener net.Listener
	quit     chan bool
	wg       sync.WaitGroup
}

// NewServer creates a new SSH server
func NewServer(address string) (*Server, error) {
	// Generate server key
	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// Simple authentication - accept any username/password for demo
			if string(password) == "password" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", conn.User())
		},
	}

	config.AddHostKey(privateKey)

	return &Server{
		address: address,
		config:  config,
		quit:    make(chan bool),
	}, nil
}

// Start starts the SSH server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.listener = listener
	log.Printf("SSH Server listening on %s\n", s.address)

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

// handleConnection handles a single SSH connection
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	// Upgrade to SSH connection
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, s.config)
	if err != nil {
		log.Printf("Failed to handshake: %v\n", err)
		return
	}
	defer sshConn.Close()

	log.Printf("New SSH connection from %s (%s)\n", sshConn.RemoteAddr(), sshConn.User())

	// Discard all global requests
	go ssh.DiscardRequests(reqs)

	// Handle channels
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v\n", err)
			continue
		}

		s.wg.Add(1)
		go s.handleChannel(channel, requests)
	}
}

// handleChannel handles an SSH channel
func (s *Server) handleChannel(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer s.wg.Done()
	defer channel.Close()

	for req := range requests {
		switch req.Type {
		case "shell":
			if req.WantReply {
				req.Reply(true, nil)
			}
			// Send welcome message
			io.WriteString(channel, "Welcome to SSH Server\r\n")
			io.WriteString(channel, "Type 'exit' to disconnect\r\n")
			io.WriteString(channel, "$ ")

			// Echo back everything
			buf := make([]byte, 1024)
			for {
				n, err := channel.Read(buf)
				if err != nil {
					return
				}
				if n > 0 {
					data := buf[:n]
					// Echo back
					channel.Write(data)
					if string(data) == "\r" || string(data) == "\n" {
						channel.Write([]byte("$ "))
					}
				}
			}
		case "exec":
			if req.WantReply {
				req.Reply(true, nil)
			}
			// Execute command (simplified)
			command := string(req.Payload[4:])
			response := fmt.Sprintf("Executed: %s\r\n", command)
			io.WriteString(channel, response)
			channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
			return
		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}
}

// Stop stops the SSH server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Println("SSH Server stopped")
}

// generatePrivateKey generates an RSA private key for the server
func generatePrivateKey() (ssh.Signer, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)
	return ssh.ParsePrivateKey(privateKeyBytes)
}
