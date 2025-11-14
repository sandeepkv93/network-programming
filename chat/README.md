## Simple Chat App

This package implements a simple multi-client chat application using TCP.

### Features

- **Chat Server**: Accepts multiple client connections and broadcasts messages
- **Chat Client**: Connects to the server and sends/receives messages
- Real-time message broadcasting to all connected clients
- Client name registration
- Join/leave notifications
- Concurrent handling of multiple clients

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"networkprogramming/chat"
)

func main() {
	server := chat.NewServer(":9000")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Keep server running
	time.Sleep(10 * time.Minute)

	server.Stop()
}
```

#### Client

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/chat"
)

func main() {
	client := chat.NewClient("localhost:9000")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Set client name
	if err := client.SetName("Alice"); err != nil {
		log.Fatal(err)
	}

	// Start receiving messages
	client.ReceiveMessages(func(message string) {
		fmt.Print(message)
	})

	// Send a message
	if err := client.SendMessage("Hello, everyone!"); err != nil {
		log.Fatal(err)
	}
}
```

### How it Works

1. **Server** listens on a specified port and accepts incoming TCP connections
2. Each client provides a name upon connection
3. Messages from any client are broadcast to all other connected clients
4. Server notifies all clients when someone joins or leaves
5. Each client has a dedicated goroutine for sending messages

### Features

- Multi-client support with concurrent connections
- Real-time message broadcasting
- Client identification by name
- Join/leave notifications
