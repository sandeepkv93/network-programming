## FTP Client & Server

This package implements a simplified FTP (File Transfer Protocol) client and server.

### Features

- **FTP Server**: Handles basic FTP commands
- **FTP Client**: Connects and executes FTP commands
- Supported commands: USER, PASS, PWD, CWD, LIST, TYPE, SYST, QUIT
- Directory listing
- Authentication (simplified)

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"network-programming/ftp"
)

func main() {
	server := ftp.NewServer(":21", "/tmp/ftp")

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
	"network-programming/ftp"
)

func main() {
	client := ftp.NewClient("localhost:21")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	// Login
	if err := client.Login("user", "pass"); err != nil {
		log.Fatal(err)
	}

	// Get current directory
	dir, err := client.Pwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current directory:", dir)

	// List files
	listing, err := client.List()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Files:\n", listing)
}
```

### How it Works

1. **Server** listens on port 21 (or custom port)
2. Client connects and receives a welcome message
3. Client authenticates with USER and PASS commands
4. Client can navigate directories (PWD, CWD) and list files (LIST)
5. Commands are sent as text over the control connection
6. Server responds with numeric status codes (220, 230, 250, etc.)

### FTP Protocol

- FTP uses two connections:
  - **Control connection** (port 21): Commands and responses
  - **Data connection** (dynamic port): File transfers and listings
- This simplified implementation uses only the control connection
- Standard FTP status codes:
  - 220: Service ready
  - 230: User logged in
  - 250: Requested action completed
  - 550: Requested action not taken

### Note

This is a simplified educational FTP implementation. Production FTP servers support:
- Separate data connections (active and passive modes)
- File upload (STOR) and download (RETR)
- Binary and ASCII transfer modes
- Resume support (REST)
- Secure FTP (FTPS with TLS/SSL)
- Many more commands (DELE, MKD, RMD, RNFR, RNTO, etc.)
