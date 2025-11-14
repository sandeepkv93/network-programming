package filetransfer

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
)

// Server represents a file transfer server
type Server struct {
	address   string
	directory string
	listener  net.Listener
	quit      chan bool
	wg        sync.WaitGroup
}

// NewServer creates a new file transfer server
func NewServer(address, directory string) *Server {
	return &Server{
		address:   address,
		directory: directory,
		quit:      make(chan bool),
	}
}

// Start starts the file transfer server
func (s *Server) Start() error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(s.directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.listener = listener
	log.Printf("File Transfer Server listening on %s\n", s.address)
	log.Printf("Storing files in: %s\n", s.directory)

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

// handleConnection handles a single file transfer connection
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s\n", clientAddr)

	// Read filename length
	var filenameLen uint32
	if err := binary.Read(conn, binary.BigEndian, &filenameLen); err != nil {
		log.Printf("Error reading filename length: %v\n", err)
		return
	}

	// Read filename
	filenameBytes := make([]byte, filenameLen)
	if _, err := io.ReadFull(conn, filenameBytes); err != nil {
		log.Printf("Error reading filename: %v\n", err)
		return
	}
	filename := string(filenameBytes)

	// Read file size
	var fileSize int64
	if err := binary.Read(conn, binary.BigEndian, &fileSize); err != nil {
		log.Printf("Error reading file size: %v\n", err)
		return
	}

	log.Printf("Receiving file: %s (%d bytes)\n", filename, fileSize)

	// Create file
	filepath := filepath.Join(s.directory, filepath.Base(filename))
	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Receive file data
	written, err := io.CopyN(file, conn, fileSize)
	if err != nil {
		log.Printf("Error receiving file: %v\n", err)
		return
	}

	log.Printf("File received successfully: %s (%d bytes)\n", filename, written)

	// Send acknowledgment
	conn.Write([]byte("OK"))
}

// Stop stops the file transfer server
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	log.Println("File Transfer Server stopped")
}
