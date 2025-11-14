package sftp

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// Server represents an SFTP server (SSH File Transfer Protocol)
type Server struct {
	Address    string
	PrivateKey string
	RootDir    string
	listener   net.Listener
}

// NewServer creates a new SFTP server
func NewServer(address, privateKey, rootDir string) *Server {
	return &Server{
		Address:    address,
		PrivateKey: privateKey,
		RootDir:    rootDir,
	}
}

// Start starts the SFTP server
func (s *Server) Start() error {
	// Configure SSH server
	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// Simple authentication - in production, verify against a user database
			if string(password) == "password" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", conn.User())
		},
	}

	// Load private key
	privateBytes, err := os.ReadFile(s.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	config.AddHostKey(private)

	// Start listening
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.Address, err)
	}

	s.listener = listener
	log.Printf("SFTP Server started on %s\n", s.Address)
	log.Printf("Root directory: %s\n", s.RootDir)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn, config)
	}
}

// Stop stops the SFTP server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handleConnection(netConn net.Conn, config *ssh.ServerConfig) {
	defer netConn.Close()

	// Perform SSH handshake
	sshConn, chans, reqs, err := ssh.NewServerConn(netConn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v\n", err)
		return
	}
	defer sshConn.Close()

	log.Printf("New SSH connection from %s (%s)\n", sshConn.RemoteAddr(), sshConn.ClientVersion())

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

		go s.handleChannel(channel, requests)
	}
}

func (s *Server) handleChannel(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	for req := range requests {
		switch req.Type {
		case "subsystem":
			if len(req.Payload) > 4 && string(req.Payload[4:]) == "sftp" {
				req.Reply(true, nil)
				s.handleSFTP(channel)
				return
			}
			req.Reply(false, nil)
		default:
			req.Reply(false, nil)
		}
	}
}

func (s *Server) handleSFTP(channel ssh.Channel) {
	log.Println("Starting SFTP subsystem")

	// Simplified SFTP implementation
	// In a real implementation, you would use github.com/pkg/sftp
	channel.Write([]byte("SFTP subsystem started\n"))

	// Echo back any data (placeholder for real SFTP protocol)
	io.Copy(channel, channel)
}
