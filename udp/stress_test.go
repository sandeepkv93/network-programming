package udp

import (
	"testing"
	"time"
)

func TestServerLoad(t *testing.T) {
	t.Skip("Skipping this test because it takes too long to run in github actions")
	serverAddr := "localhost:9002"
	server := NewUDPServer(serverAddr)

	// Start the server in a goroutine
	go server.Start()

	// Give some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clients := 10000 // Number of clients sending messages concurrently
	doneCh := make(chan bool, clients)

	for i := 0; i < clients; i++ {
		go func() {
			client := NewUDPClient(serverAddr)
			client.SendMessage("Stress Test Message")
			doneCh <- true
		}()
	}

	// Wait for all clients to finish
	for i := 0; i < clients; i++ {
		<-doneCh
	}
}
