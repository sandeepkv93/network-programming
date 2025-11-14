package portforwarding

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Forwarder represents a port forwarder
type Forwarder struct {
	LocalAddr   string
	RemoteAddr  string
	listener    net.Listener
	connections map[string]net.Conn
	mu          sync.Mutex
	running     bool
}

// NewForwarder creates a new port forwarder
func NewForwarder(localAddr, remoteAddr string) *Forwarder {
	return &Forwarder{
		LocalAddr:   localAddr,
		RemoteAddr:  remoteAddr,
		connections: make(map[string]net.Conn),
	}
}

// Start starts the port forwarder
func (f *Forwarder) Start() error {
	f.mu.Lock()
	if f.running {
		f.mu.Unlock()
		return fmt.Errorf("forwarder is already running")
	}
	f.running = true
	f.mu.Unlock()

	listener, err := net.Listen("tcp", f.LocalAddr)
	if err != nil {
		f.mu.Lock()
		f.running = false
		f.mu.Unlock()
		return fmt.Errorf("failed to listen on %s: %v", f.LocalAddr, err)
	}

	f.listener = listener
	log.Printf("Port forwarder started: %s -> %s\n", f.LocalAddr, f.RemoteAddr)

	go f.acceptConnections()

	return nil
}

// Stop stops the port forwarder
func (f *Forwarder) Stop() error {
	f.mu.Lock()
	if !f.running {
		f.mu.Unlock()
		return fmt.Errorf("forwarder is not running")
	}
	f.running = false
	f.mu.Unlock()

	// Close listener
	if f.listener != nil {
		f.listener.Close()
	}

	// Close all active connections
	f.mu.Lock()
	for _, conn := range f.connections {
		conn.Close()
	}
	f.connections = make(map[string]net.Conn)
	f.mu.Unlock()

	log.Println("Port forwarder stopped")
	return nil
}

func (f *Forwarder) acceptConnections() {
	for {
		conn, err := f.listener.Accept()
		if err != nil {
			f.mu.Lock()
			running := f.running
			f.mu.Unlock()

			if !running {
				return
			}

			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go f.handleConnection(conn)
	}
}

func (f *Forwarder) handleConnection(clientConn net.Conn) {
	clientAddr := clientConn.RemoteAddr().String()
	log.Printf("New connection from %s\n", clientAddr)

	// Store the connection
	f.mu.Lock()
	f.connections[clientAddr] = clientConn
	f.mu.Unlock()

	// Clean up on exit
	defer func() {
		clientConn.Close()
		f.mu.Lock()
		delete(f.connections, clientAddr)
		f.mu.Unlock()
		log.Printf("Connection closed: %s\n", clientAddr)
	}()

	// Connect to remote server
	remoteConn, err := net.DialTimeout("tcp", f.RemoteAddr, 10*time.Second)
	if err != nil {
		log.Printf("Failed to connect to remote %s: %v\n", f.RemoteAddr, err)
		return
	}
	defer remoteConn.Close()

	log.Printf("Forwarding %s -> %s\n", clientAddr, f.RemoteAddr)

	// Create channels for errors from each goroutine
	errChan := make(chan error, 2)

	// Forward data from client to remote
	go func() {
		_, err := io.Copy(remoteConn, clientConn)
		errChan <- err
	}()

	// Forward data from remote to client
	go func() {
		_, err := io.Copy(clientConn, remoteConn)
		errChan <- err
	}()

	// Wait for either direction to finish
	err = <-errChan

	if err != nil && err != io.EOF {
		log.Printf("Error forwarding data for %s: %v\n", clientAddr, err)
	}
}

// GetActiveConnections returns the number of active connections
func (f *Forwarder) GetActiveConnections() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.connections)
}

// IsRunning returns whether the forwarder is running
func (f *Forwarder) IsRunning() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.running
}
