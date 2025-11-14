## Remote Login

A secure remote login server and client implementation providing authenticated shell access to remote systems, similar to telnet but with password authentication.

## Features

- Password-based authentication with SHA-256 hashing
- Session management
- Interactive shell interface
- Built-in commands (help, whoami, session, users, etc.)
- System command execution
- Multiple concurrent sessions

## Usage

### Server
```go
server := remotelogin.NewServer(":2222")
server.AddUser("admin", "password123")
server.AddUser("user", "userpass")
server.Start()
```

### Client
```go
client := remotelogin.NewClient("localhost:2222")
client.Connect()
client.Login("admin", "password123")
client.StartInteractive()
```

## Security Notes

⚠️ This is a basic implementation. For production use:
- Use SSH instead of this implementation
- Implement TLS encryption
- Use stronger authentication (keys, 2FA)
- Add rate limiting
- Implement comprehensive logging
