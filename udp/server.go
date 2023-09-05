package udp

import (
	"fmt"
	"net"
	"sync"
)

// UDPServer represents a UDP server that listens for incoming client requests
type UDPServer struct {
	address string        // The address to listen on (e.g. "localhost:8080")
	mutex   sync.Mutex   // A mutex to synchronize access to shared resources
	bufPool *sync.Pool   // A pool of reusable byte buffers
}

// NewUDPServer creates a new instance of UDPServer with the given address
func NewUDPServer(address string) *UDPServer {
	// Initialize a buffer pool with 1KB buffers
	pool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}

	return &UDPServer{
		address: address,
		bufPool: pool,
	}
}

// handleClientRequest handles an incoming client request by sending a response back to the client
func (s *UDPServer) handleClientRequest(clientAddr *net.UDPAddr, message []byte, conn *net.UDPConn) {
	defer s.bufPool.Put(message) // Return the buffer to the pool once done

	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Printf("Received %s from %s\n", string(message), clientAddr)

	// Send response to client
	_, err := conn.WriteToUDP([]byte("Hello, client!"), clientAddr)
	if err != nil {
		fmt.Println("Error sending response:", err)
	}
}

// Start starts the UDP server and listens for incoming client requests
func (s *UDPServer) Start() {
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	for {
		buf := s.bufPool.Get().([]byte)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error:", err)
			s.bufPool.Put(buf) // Don't forget to return the buffer to the pool if an error occurs
			continue
		}

		go s.handleClientRequest(clientAddr, buf[:n], conn)
	}
}
