## WebSocket Client & Server

The WebSocket package provides a full-duplex communication channel over a single TCP connection, enabling real-time bidirectional data exchange between clients and servers.

## Table of Contents

1. [What is WebSocket?](#what-is-websocket)
2. [How Does It Work?](#how-does-it-work)
3. [Understanding the Code](#understanding-the-code)
4. [Usage Examples](#usage-examples)
5. [Further Reading](#further-reading)

### What is WebSocket?

WebSocket is a protocol that provides full-duplex communication channels over a single TCP connection. Unlike HTTP, which follows a request-response pattern, WebSocket allows the server to push data to clients without being polled.

**Key Features**:
- Full-duplex bidirectional communication
- Low latency real-time data exchange
- Persistent connection (no overhead of repeated HTTP handshakes)
- Event-driven messaging
- Support for both text and binary data

**Common Use Cases**:
- Real-time chat applications
- Live notifications and updates
- Collaborative editing tools
- Gaming applications
- Live streaming and dashboards
- IoT device communication

### How Does It Work?

WebSocket communication follows this flow:

1. **Handshake**: Client sends an HTTP upgrade request to the server
2. **Upgrade**: Server accepts and upgrades connection to WebSocket protocol
3. **Data Exchange**: Both parties can send messages at any time
4. **Close**: Either party can close the connection gracefully

**Message Flow**:
```
Client                          Server
  |                               |
  |--- HTTP Upgrade Request ----->|
  |<--- HTTP 101 Switching -------|
  |                               |
  |<-------- Messages ----------->|
  |<-------- Messages ----------->|
  |                               |
  |------- Close Frame ---------->|
  |<------ Close Frame -----------|
```

### Understanding the Code

#### Server Components:

- `Server`: Manages WebSocket connections and message broadcasting
- `handleWebSocket`: Handles individual client connections
- `handleBroadcast`: Broadcasts messages to all connected clients
- `handleHome`: Serves a test HTML page with WebSocket client

**Key Features**:
- Concurrent client handling with goroutines
- Thread-safe client management with mutex
- Echo and broadcast functionality
- Automatic client cleanup on disconnect

#### Client Components:

- `Client`: Manages connection to WebSocket server
- `Connect()`: Establishes WebSocket connection
- `SendMessage()`: Sends text messages to server
- `ReceiveMessages()`: Listens for incoming messages
- `RunInteractive()`: Provides interactive CLI interface

### Usage Examples

#### Starting a WebSocket Server:
```go
server := websocket.NewServer(":8080")
if err := server.Start(); err != nil {
    log.Fatal(err)
}
```

#### Connecting a Client:
```go
client := websocket.NewClient("ws://localhost:8080/ws")
if err := client.Connect(); err != nil {
    log.Fatal(err)
}
defer client.Close()

// Send a message
client.SendMessage("Hello, WebSocket!")

// Start receiving messages
go client.ReceiveMessages()
```

#### Interactive Client Session:
```go
client := websocket.NewClient("ws://localhost:8080/ws")
if err := client.RunInteractive(); err != nil {
    log.Fatal(err)
}
```

#### Broadcasting to All Clients:
The server automatically broadcasts received messages to all connected clients, making it easy to implement chat rooms or notification systems.

### Protocol Details

**WebSocket Frame Format**:
- FIN bit: Indicates final fragment
- Opcode: Message type (text, binary, close, ping, pong)
- Mask: Client-to-server messages must be masked
- Payload length: Size of the message data
- Payload data: The actual message content

**Message Types**:
- Text messages: UTF-8 encoded text
- Binary messages: Raw binary data
- Control frames: Ping, pong, close

### Further Reading

- [RFC 6455 - WebSocket Protocol](https://datatracker.ietf.org/doc/html/rfc6455)
- [Gorilla WebSocket Documentation](https://pkg.go.dev/github.com/gorilla/websocket)
- [WebSocket API - MDN](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
