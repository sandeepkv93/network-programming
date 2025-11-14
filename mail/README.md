## Mail Client & Server

This package implements a simplified mail server and client for basic email operations.

### Features

- **Mail Server**: Stores and delivers messages between users
- **Mail Client**: Sends and retrieves messages
- User mailboxes with message storage
- Basic commands: USER, STAT, LIST, RETR, SEND, QUIT
- In-memory message storage

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"network-programming/mail"
)

func main() {
	server := mail.NewServer(":2525")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Keep server running
	time.Sleep(1 * time.Minute)

	server.Stop()
}
```

#### Client - Sending a Message

```go
package main

import (
	"log"
	"network-programming/mail"
)

func main() {
	client := mail.NewClient("localhost:2525")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	// Login
	if err := client.Login("alice"); err != nil {
		log.Fatal(err)
	}

	// Send a message
	err := client.SendMessage(
		"bob",
		"Hello",
		"This is a test message.\nHow are you?",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Message sent successfully")
}
```

#### Client - Receiving Messages

```go
package main

import (
	"log"
	"network-programming/mail"
)

func main() {
	client := mail.NewClient("localhost:2525")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	// Login
	if err := client.Login("bob"); err != nil {
		log.Fatal(err)
	}

	// Get message count
	count, err := client.GetMessageCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("You have %d messages\n", count)

	// List messages
	messages, err := client.ListMessages()
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range messages {
		log.Println(msg)
	}

	// Retrieve first message
	if count > 0 {
		content, err := client.RetrieveMessage(1)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Message:\n", content)
	}
}
```

### Protocol Commands

- **USER <username>**: Log in as a user
- **STAT**: Get message count
- **LIST**: List all messages with headers
- **RETR <index>**: Retrieve a specific message
- **SEND <recipient>**: Send a message (followed by subject and body)
- **QUIT**: Disconnect from server

### Message Format

When sending a message:
1. First line: Subject
2. Following lines: Message body
3. End with a single "." on a line

### How it Works

1. **Server** maintains mailboxes for each user
2. Messages are stored in memory
3. **Client** connects and logs in with a username
4. Client can send messages to other users
5. Client can check and retrieve their messages
6. Each message includes: From, To, Subject, Body, Timestamp

### Response Format

- **+OK**: Successful command
- **-ERR**: Error response

### Note

This is a simplified educational mail server. Production mail servers (like those using SMTP/POP3/IMAP) include:
- Authentication and encryption
- Persistent storage
- Spam filtering
- Attachment support
- Multiple mailbox folders
- Message threading
- Search capabilities
- Quota management
- Delivery receipts
- MIME encoding
