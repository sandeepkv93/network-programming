## TCP Client & Server

This package implements a basic TCP (Transmission Control Protocol) client and server.

### Features

- **TCP Server**: Accepts multiple client connections and echoes back received messages
- **TCP Client**: Connects to the server and sends messages
- Connection management with graceful shutdown
- Concurrent handling of multiple clients

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"network-programming/tcp"
)

func main() {
	server := tcp.NewServer(":8080")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Keep server running
	time.Sleep(1 * time.Minute)

	server.Stop()
}
```

#### Client

```go
package main

import (
	"log"
	"network-programming/tcp"
)

func main() {
	client := tcp.NewClient("localhost:8080")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	response, err := client.Send("Hello, Server!")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Response:", response)
}
```

### How it Works

1. **Server** listens on a specified port and accepts incoming TCP connections
2. Each client connection is handled in a separate goroutine
3. **Client** connects to the server and sends messages
4. Server receives messages and sends back responses
5. Communication uses newline-delimited messages

### TCP vs UDP

- **TCP** is connection-oriented and reliable (guaranteed delivery)
- **UDP** is connectionless and unreliable (no delivery guarantee)
- TCP provides ordering, error checking, and flow control
