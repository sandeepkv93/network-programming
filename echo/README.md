## Echo Client & Server

This package implements an Echo server and client. The server echoes back exactly what it receives from clients.

### Features

- **Echo Server**: Accepts connections and echoes back received messages
- **Echo Client**: Sends messages and receives echoes
- Multiple concurrent client support
- Simple protocol using newline-delimited messages

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"network-programming/echo"
)

func main() {
	server := echo.NewServer(":7777")

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
	"network-programming/echo"
)

func main() {
	client := echo.NewClient("localhost:7777")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	echo, err := client.Echo("Hello, Echo Server!")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Echo:", echo)
}
```

### How it Works

1. **Server** listens on a TCP port for incoming connections
2. For each connection, it reads messages line by line
3. Each received message is immediately sent back (echoed) to the client
4. **Client** connects, sends a message, and receives the exact echo

### Echo Protocol

The Echo protocol is defined in RFC 862:
- Listens on port 7 (by convention, though any port can be used)
- Sends back any data it receives
- Can be implemented over TCP or UDP
- Useful for testing network connectivity and latency
