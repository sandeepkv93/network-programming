## SMTP Client & Server

This package implements an SMTP (Simple Mail Transfer Protocol) client and server for sending emails.

### Features

- **SMTP Server**: Receives and stores email messages
- **SMTP Client**: Sends emails via SMTP protocol
- Standard SMTP commands: HELO, MAIL FROM, RCPT TO, DATA, QUIT
- Multiple recipients support
- RFC 5321 compliant (simplified)

### Usage

#### Server

```go
package main

import (
	"log"
	"time"
	"network-programming/smtp"
)

func main() {
	server := smtp.NewServer(":25", "mail.example.com")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Keep server running
	time.Sleep(1 * time.Minute)

	// Get received emails
	emails := server.GetEmails()
	for _, email := range emails {
		log.Printf("Email: From=%s To=%v\n", email.From, email.To)
		log.Printf("Data: %s\n", email.Data)
	}

	server.Stop()
}
```

#### Client - Simple Email

```go
package main

import (
	"log"
	"network-programming/smtp"
)

func main() {
	client := smtp.NewClient("localhost:25")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	err := client.SendMail(
		"alice@example.com",
		[]string{"bob@example.com"},
		"Hello",
		"This is a test email.",
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully")
}
```

#### Client - Manual Commands

```go
package main

import (
	"log"
	"network-programming/smtp"
)

func main() {
	client := smtp.NewClient("localhost:25")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	// Send HELO
	client.Hello("myclient.com")

	// Set sender
	client.MailFrom("sender@example.com")

	// Set recipients
	client.RcptTo("recipient1@example.com")
	client.RcptTo("recipient2@example.com")

	// Send message
	message := "Subject: Test\r\n\r\nThis is the email body."
	client.Data(message)

	log.Println("Email sent")
}
```

### SMTP Protocol Flow

1. **Connection**: Server sends 220 greeting
2. **HELO/EHLO**: Client identifies itself
3. **MAIL FROM**: Sender's email address
4. **RCPT TO**: Recipient's email address (can be multiple)
5. **DATA**: Email content (headers + body)
6. **QUIT**: Close connection

### SMTP Commands

- **HELO <hostname>**: Identify client to server
- **MAIL FROM:<address>**: Specify sender
- **RCPT TO:<address>**: Specify recipient
- **DATA**: Start message content
- **RSET**: Reset session
- **NOOP**: No operation (keep-alive)
- **QUIT**: Close connection

### SMTP Response Codes

- **220**: Service ready
- **250**: Requested action completed
- **354**: Start mail input
- **221**: Service closing
- **500**: Syntax error
- **501**: Syntax error in parameters
- **502**: Command not implemented
- **503**: Bad sequence of commands

### Email Message Format

```
From: sender@example.com
To: recipient@example.com
Subject: Email Subject

Email body goes here.
Multiple lines are supported.
```

### How it Works

1. **Client** connects to SMTP server on port 25 (or custom)
2. Server greets with 220 response
3. Client sends commands in order: HELO, MAIL FROM, RCPT TO, DATA
4. Server validates and responds to each command
5. Client sends email content, ending with "." on a line by itself
6. Server stores the email
7. Client sends QUIT to close connection

### Note

This is a simplified educational SMTP implementation. Production SMTP servers include:
- Authentication (SMTP AUTH)
- Encryption (STARTTLS, SMTPS)
- SPF, DKIM, DMARC validation
- Message queue management
- Relay controls and anti-spam
- Delivery status notifications (DSN)
- Extended SMTP (ESMTP) features
- Message size limits
- Rate limiting
- Graylisting
- Virus scanning
