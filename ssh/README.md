## SSH Client & Server

This package implements a basic SSH (Secure Shell) client and server for secure remote access.

### Features

- **SSH Server**: Accepts secure client connections with password authentication
- **SSH Client**: Connects to SSH servers and executes commands
- Encrypted communication using SSH protocol
- Password-based authentication
- Command execution support
- Interactive shell support

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"networkprogramming/ssh"
)

func main() {
	server, err := ssh.NewServer(":2222")
	if err != nil {
		log.Fatal(err)
	}

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
	"networkprogramming/ssh"
)

func main() {
	client := ssh.NewClient("localhost:2222", "user", "password")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Execute a command
	output, err := client.ExecuteCommand("ls -la")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Command output:", output)
}
```

### How it Works

1. **Server** generates an RSA key pair for encryption
2. Client connects and authenticates using username/password
3. All communication is encrypted using SSH protocol
4. Commands can be executed remotely
5. Interactive shell sessions are supported

### Security Notes

- This is a simplified implementation for educational purposes
- Uses password authentication (demo password: "password")
- Generates a new server key on each start
- Production servers should:
  - Use key-based authentication
  - Persist server keys
  - Implement proper user management
  - Use strong password policies

### Dependencies

Requires `golang.org/x/crypto/ssh` package:
```
go get golang.org/x/crypto/ssh
```
