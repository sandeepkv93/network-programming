## Telnet Client & Server

This package implements a basic Telnet protocol client and server for remote terminal access.

### Features

- **Telnet Server**: Accepts client connections and provides a command-line interface
- **Telnet Client**: Connects to the server and sends commands
- Interactive command processing
- Built-in commands (help, time, echo, quit)
- Support for custom command handlers

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"networkprogramming/telnet"
)

func main() {
	server := telnet.NewServer(":23")

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
	"networkprogramming/telnet"
)

func main() {
	client := telnet.NewClient("localhost:23")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Read welcome message
	welcome, err := client.ReadWelcome()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(welcome)

	// Send command
	if err := client.SendCommand("help"); err != nil {
		log.Fatal(err)
	}

	// Read response
	response, err := client.ReadResponse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(response)
}
```

### Available Commands

- `help` - Show available commands
- `time` - Show current server time
- `echo <message>` - Echo back the message
- `quit` - Disconnect from server

### How it Works

1. **Server** listens on port 23 (traditional Telnet port)
2. Client connects and receives a welcome message
3. Server provides an interactive prompt (>)
4. Client sends commands, server processes and responds
5. Communication uses CRLF line endings (\r\n)

### Note

This is a simplified Telnet implementation for educational purposes. Production Telnet servers should implement the full Telnet protocol with option negotiation (RFC 854).
