package udp

import (
	"fmt"
	"net"
	"sync"
)

type UDPServer struct {
	address string
	mutex   sync.Mutex
}

func NewUDPServer(address string) *UDPServer {
	return &UDPServer{
		address: address,
	}
}

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

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Protect shared resources using Mutex
		s.mutex.Lock()
		fmt.Printf("Received %s from %s\n", string(buf[:n]), clientAddr)

		// Send response to client
		_, err = conn.WriteToUDP([]byte("Hello, client!"), clientAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
		s.mutex.Unlock()
	}
}
