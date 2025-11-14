package ftp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Server represents an FTP server
type Server struct {
	address string
	rootDir string
	listener net.Listener
	quit     chan bool
	wg       sync.WaitGroup
}

// NewServer creates a new FTP server
func NewServer(address string, rootDir string) *Server {
	return &Server{
		address: address,
		rootDir: rootDir,
		quit:    make(chan bool),
	}
}

// Start starts the FTP server
func (s *Server) Start() error {
	// Create root directory if it doesn't exist
	if err := os.MkdirAll(s.rootDir, 0755); err != nil {
		return fmt.Errorf("failed to create root directory: %v", err)
	}

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start FTP server: %v", err)
	}
	s.listener = listener
	log.Printf("FTP Server listening on %s\n", s.address)
	log.Printf("Root directory: %s\n", s.rootDir)

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

	// Send welcome message
	s.sendResponse(conn, 220, "FTP Server Ready")

	reader := bufio.NewReader(conn)
	currentDir := "/"

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
			s.sendResponse(conn, 230, "User logged in")
		case "PASS":
			s.sendResponse(conn, 230, "Password accepted")
		case "PWD":
			s.sendResponse(conn, 257, fmt.Sprintf("\"%s\" is current directory", currentDir))
		case "CWD":
			if arg != "" {
				currentDir = s.normalizePath(currentDir, arg)
				s.sendResponse(conn, 250, "Directory changed")
			} else {
				s.sendResponse(conn, 501, "No directory specified")
			}
		case "LIST":
			s.handleList(conn, currentDir)
		case "TYPE":
			s.sendResponse(conn, 200, "Type set")
		case "SYST":
			s.sendResponse(conn, 215, "UNIX Type: L8")
		case "QUIT":
			s.sendResponse(conn, 221, "Goodbye")
			return
		default:
			s.sendResponse(conn, 502, "Command not implemented")
		}
	}
}

// handleList handles the LIST command
func (s *Server) handleList(conn net.Conn, dir string) {
	fullPath := filepath.Join(s.rootDir, dir)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		s.sendResponse(conn, 550, "Failed to list directory")
		return
	}

	s.sendResponse(conn, 150, "Opening data connection")

	var listing strings.Builder
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		var perms string
		if entry.IsDir() {
			perms = "drwxr-xr-x"
		} else {
			perms = "-rw-r--r--"
		}

		listing.WriteString(fmt.Sprintf("%s 1 owner group %10d %s %s\r\n",
			perms,
			info.Size(),
			info.ModTime().Format("Jan 02 15:04"),
			entry.Name()))
	}

	// For simplicity, send listing on control connection
	conn.Write([]byte(listing.String()))
	s.sendResponse(conn, 226, "Transfer complete")
}

// sendResponse sends an FTP response
func (s *Server) sendResponse(conn net.Conn, code int, message string) {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	conn.Write([]byte(response))
}

// normalizePath normalizes a path
func (s *Server) normalizePath(current, target string) string {
	if strings.HasPrefix(target, "/") {
		return filepath.Clean(target)
	}
	return filepath.Clean(filepath.Join(current, target))
}

// Stop stops the FTP server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Println("FTP Server stopped")
}
