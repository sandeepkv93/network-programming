## File Transfer over Network

This package implements a simple file transfer protocol for sending files over TCP.

### Features

- **File Transfer Server**: Receives files from clients and saves them to disk
- **File Transfer Client**: Sends files to the server
- Binary protocol with file metadata (name, size)
- Acknowledgment mechanism
- Support for large files
- Progress tracking capability

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"networkprogramming/filetransfer"
)

func main() {
	// Create server that saves files to ./uploads directory
	server := filetransfer.NewServer(":9999", "./uploads")

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
	"log"
	"networkprogramming/filetransfer"
)

func main() {
	client := filetransfer.NewClient("localhost:9999")

	// Send a file
	if err := client.SendFile("./myfile.txt"); err != nil {
		log.Fatal(err)
	}

	log.Println("File sent successfully!")
}
```

### Protocol Format

The file transfer protocol uses the following format:

1. **Filename Length** (4 bytes, big-endian uint32)
2. **Filename** (variable length string)
3. **File Size** (8 bytes, big-endian int64)
4. **File Data** (variable length, based on file size)
5. **Acknowledgment** (2 bytes, "OK" from server)

### How it Works

1. Client opens a file and gets its size
2. Client connects to server
3. Client sends filename length, filename, and file size
4. Client streams file data
5. Server receives and saves the file
6. Server sends acknowledgment ("OK")
7. Connection closes

### Use Cases

- File backup systems
- Content distribution
- Log file collection
- Data synchronization
- Remote file uploads

### Notes

- Files are saved with their base filename (path is stripped)
- Server creates the upload directory if it doesn't exist
- Binary protocol for efficient transfer
- Suitable for transferring any file type (text, images, videos, etc.)
