## WebRTC Client & Server

The WebRTC package provides peer-to-peer real-time communication capabilities, including video, audio, and data channels, with a signaling server to facilitate connection establishment.

## Table of Contents

1. [What is WebRTC?](#what-is-webrtc)
2. [How Does It Work?](#how-does-it-work)
3. [Understanding the Code](#understanding-the-code)
4. [Usage Examples](#usage-examples)
5. [Further Reading](#further-reading)

### What is WebRTC?

WebRTC (Web Real-Time Communication) is an open-source project that enables real-time communication of audio, video, and data in web browsers and mobile applications via simple APIs. It allows peer-to-peer connections with minimal latency.

**Key Features**:
- Peer-to-peer communication (reduces server load)
- Low latency real-time data transfer
- Built-in encryption (DTLS and SRTP)
- NAT traversal using ICE, STUN, and TURN
- Support for audio, video, and arbitrary data

**Common Use Cases**:
- Video conferencing
- Voice calling (VoIP)
- Screen sharing
- File transfer
- Gaming
- Live broadcasting
- IoT device communication

### How Does It Work?

WebRTC connection establishment follows the JSEP (JavaScript Session Establishment Protocol):

1. **Signaling**: Peers exchange connection metadata (not part of WebRTC spec)
2. **Offer/Answer**: SDP (Session Description Protocol) exchange
3. **ICE Candidates**: Network path discovery using STUN/TURN servers
4. **Connection**: Direct peer-to-peer connection established
5. **Data/Media Transfer**: Real-time communication begins

**Connection Flow**:
```
Peer A                  Signaling Server              Peer B
  |                            |                         |
  |--- Create Offer ---------->|                         |
  |                            |--- Relay Offer -------->|
  |                            |<-- Create Answer -------|
  |<-- Relay Answer -----------|                         |
  |                            |                         |
  |<------- ICE Candidates Exchange -------------------->|
  |                            |                         |
  |<=============== Direct P2P Connection ==============>|
```

### Understanding the Code

#### Signaling Server Components:

- `SignalingServer`: WebSocket-based server to relay signaling messages
- `SignalMessage`: Structure for offer, answer, and ICE candidate messages
- `handleWebSocket`: Manages peer connections and message relay
- `relayMessage`: Forwards signaling messages between peers

**The signaling server does NOT handle media/data - it only facilitates connection setup.**

#### Peer Components:

- `Peer`: Manages WebRTC peer connection and signaling
- `CreatePeerConnection`: Initializes RTCPeerConnection with ICE servers
- `CreateOffer`: Initiates connection by creating SDP offer
- `HandleSignaling`: Processes incoming signaling messages
- `SendMessage`: Sends data through established data channel

**Key WebRTC Elements**:
- **ICE (Interactive Connectivity Establishment)**: NAT traversal
- **SDP (Session Description Protocol)**: Media and connection metadata
- **Data Channels**: Bidirectional data transfer
- **STUN Server**: Public IP discovery
- **TURN Server**: Relay when direct connection fails (not implemented in basic example)

### Usage Examples

#### Starting the Signaling Server:
```go
server := webrtc.NewSignalingServer(":8080")
if err := server.Start(); err != nil {
    log.Fatal(err)
}
```

#### Creating a Peer (Caller):
```go
peer := webrtc.NewPeer("ws://localhost:8080/ws")

// Connect to signaling server
if err := peer.Connect(); err != nil {
    log.Fatal(err)
}

// Create peer connection
if err := peer.CreatePeerConnection(); err != nil {
    log.Fatal(err)
}

// Create and send offer
if err := peer.CreateOffer(); err != nil {
    log.Fatal(err)
}

// Handle incoming signaling messages
go peer.HandleSignaling()

// Send a message
time.Sleep(5 * time.Second) // Wait for connection
peer.SendMessage("Hello, WebRTC!")
```

#### Creating a Peer (Receiver):
```go
peer := webrtc.NewPeer("ws://localhost:8080/ws")

// Connect to signaling server
peer.Connect()

// Create peer connection
peer.CreatePeerConnection()

// Set message handler
peer.SetOnMessage(func(msg string) {
    log.Printf("Received: %s\n", msg)
})

// Handle signaling (will automatically respond to offers)
peer.HandleSignaling()
```

#### Checking Connection State:
```go
state := peer.GetConnectionState()
fmt.Printf("Connection state: %s\n", state)

// Get detailed statistics
stats, err := peer.GetStats()
if err == nil {
    fmt.Println(stats)
}
```

### Architecture

**Signaling Server**:
- Does not see or handle actual media/data
- Only relays connection setup messages
- Can be implemented with any protocol (WebSocket, HTTP, etc.)

**Peer Connection**:
- Direct connection between peers
- All data/media flows peer-to-peer
- Server only used for initial handshake

**NAT Traversal**:
- STUN: Discovers public IP and port
- TURN: Relays traffic when direct connection impossible
- ICE: Coordinates the process

### Security

WebRTC includes built-in security:
- **DTLS**: Encryption for data channels
- **SRTP**: Encryption for media streams
- **Secure origins**: HTTPS required in browsers
- **Permission prompts**: User consent for camera/microphone

### Further Reading

- [WebRTC Specification](https://www.w3.org/TR/webrtc/)
- [Pion WebRTC Documentation](https://github.com/pion/webrtc)
- [WebRTC for the Curious](https://webrtcforthecurious.com/)
- [ICE RFC 8445](https://datatracker.ietf.org/doc/html/rfc8445)
- [SDP RFC 4566](https://datatracker.ietf.org/doc/html/rfc4566)
