package ftps

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Server represents an FTPS server (FTP over TLS)
type Server struct {
	Address  string
	CertFile string
	KeyFile  string
	RootDir  string
	listener net.Listener
}

// NewServer creates a new FTPS server
func NewServer(address, certFile, keyFile, rootDir string) *Server {
	return &Server{
		Address:  address,
		CertFile: certFile,
		KeyFile:  keyFile,
		RootDir:  rootDir,
	}
}

// Start starts the FTPS server
func (s *Server) Start() error {
	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to load certificate: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Create TLS listener
	listener, err := tls.Listen("tcp", s.Address, config)
	if err != nil {
		return fmt.Errorf("failed to start FTPS server: %v", err)
	}

	s.listener = listener
	log.Printf("FTPS Server started on %s\n", s.Address)
	log.Printf("Root directory: %s\n", s.RootDir)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// Stop stops the FTPS server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	log.Printf("New FTPS connection from %s\n", clientAddr)

	// Send welcome message
	conn.Write([]byte("220 Welcome to FTPS Server\r\n"))

	buffer := make([]byte, 4096)
	currentDir := "/"

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from client: %v\n", err)
			}
			break
		}

		command := strings.TrimSpace(string(buffer[:n]))
		log.Printf("[%s] Command: %s\n", clientAddr, command)

		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}

		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "USER":
			conn.Write([]byte("331 Username OK, need password\r\n"))
		case "PASS":
			conn.Write([]byte("230 User logged in\r\n"))
		case "PWD":
			response := fmt.Sprintf("257 \"%s\" is current directory\r\n", currentDir)
			conn.Write([]byte(response))
		case "CWD":
			if len(parts) > 1 {
				currentDir = parts[1]
				conn.Write([]byte("250 Directory changed\r\n"))
			} else {
				conn.Write([]byte("501 Syntax error\r\n"))
			}
		case "LIST":
			s.handleList(conn, currentDir)
		case "RETR":
			if len(parts) > 1 {
				s.handleRetrieve(conn, filepath.Join(currentDir, parts[1]))
			} else {
				conn.Write([]byte("501 Syntax error\r\n"))
			}
		case "STOR":
			if len(parts) > 1 {
				s.handleStore(conn, filepath.Join(currentDir, parts[1]))
			} else {
				conn.Write([]byte("501 Syntax error\r\n"))
			}
		case "QUIT":
			conn.Write([]byte("221 Goodbye\r\n"))
			return
		default:
			conn.Write([]byte("502 Command not implemented\r\n"))
		}
	}
}

func (s *Server) handleList(conn net.Conn, dir string) {
	fullPath := filepath.Join(s.RootDir, dir)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		conn.Write([]byte("550 Failed to list directory\r\n"))
		return
	}

	conn.Write([]byte("150 Opening data connection\r\n"))

	for _, entry := range entries {
		info, _ := entry.Info()
		line := fmt.Sprintf("%s %10d %s\r\n",
			entry.Name(),
			info.Size(),
			info.ModTime().Format("Jan 02 15:04"))
		conn.Write([]byte(line))
	}

	conn.Write([]byte("226 Transfer complete\r\n"))
}

func (s *Server) handleRetrieve(conn net.Conn, filename string) {
	fullPath := filepath.Join(s.RootDir, filename)
	file, err := os.Open(fullPath)
	if err != nil {
		conn.Write([]byte("550 File not found\r\n"))
		return
	}
	defer file.Close()

	conn.Write([]byte("150 Opening data connection\r\n"))
	io.Copy(conn, file)
	conn.Write([]byte("226 Transfer complete\r\n"))
}

func (s *Server) handleStore(conn net.Conn, filename string) {
	// Simplified - would need data connection in full implementation
	conn.Write([]byte("150 Ready to receive file\r\n"))
	conn.Write([]byte("226 Transfer complete\r\n"))
}
