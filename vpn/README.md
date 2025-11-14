## VPN Client & Server

The VPN package provides a Virtual Private Network implementation with encrypted tunneling, allowing secure communication between clients over untrusted networks.

## Table of Contents

1. [What is a VPN?](#what-is-a-vpn)
2. [How Does It Work?](#how-does-it-work)
3. [Understanding the Code](#understanding-the-code)
4. [Usage Examples](#usage-examples)
5. [Further Reading](#further-reading)

### What is a VPN?

A VPN (Virtual Private Network) extends a private network across a public network, enabling users to send and receive data as if their devices were directly connected to the private network. VPNs encrypt traffic to ensure privacy and security.

**Key Features**:
- Encrypted communication (AES-GCM)
- IP address assignment and management
- Packet routing between clients
- Network isolation
- Secure tunneling

**Common Use Cases**:
- Remote access to corporate networks
- Secure browsing on public Wi-Fi
- Bypassing geographic restrictions
- Privacy protection
- Site-to-site connectivity

### How Does It Work?

The VPN implementation uses the following architecture:

1. **Connection Establishment**: Client connects to VPN server via TCP
2. **Configuration**: Server assigns IP address from VPN subnet
3. **Encryption**: All packets encrypted using AES-GCM
4. **Tunneling**: IP packets encapsulated and routed through VPN server
5. **Routing**: Server maintains routing table for packet forwarding

**VPN Connection Flow**:
```
Client A (10.0.0.2)          VPN Server           Client B (10.0.0.3)
       |                          |                         |
       |--- TCP Connect --------->|                         |
       |<-- Assign IP 10.0.0.2 ---|                         |
       |                          |<-- TCP Connect ---------|
       |                          |--- Assign IP 10.0.0.3 ->|
       |                          |                         |
       |-- Encrypted Packet ----->|                         |
       |    (Dest: 10.0.0.3)      |--- Encrypted Packet --->|
       |                          |    (Decrypted & Routed) |
```

### Understanding the Code

#### Server Components:

- `Server`: Manages VPN server, client connections, and routing
- `VPNClient`: Represents a connected client with assigned IP
- `assignIP()`: Allocates IP addresses from subnet pool
- `routePacket()`: Routes packets to appropriate destination client
- `encrypt/decrypt()`: AES-GCM encryption for all traffic

**Key Features**:
- Automatic IP address management
- Thread-safe client handling
- Routing table for packet forwarding
- AES-GCM encryption (authenticated encryption)
- Graceful client disconnection

#### Client Components:

- `Client`: Manages connection to VPN server
- `Connect()`: Establishes VPN connection and receives config
- `SendPacket()`: Sends encrypted IP packets through tunnel
- `ReceivePacket()`: Receives and decrypts packets from tunnel
- `StartTunnel()`: Continuously processes incoming packets

### Usage Examples

#### Starting a VPN Server:
```go
// Key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256
key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256

server, err := vpn.NewServer(":1194", "10.0.0.0/24", key)
if err != nil {
    log.Fatal(err)
}

if err := server.Start(); err != nil {
    log.Fatal(err)
}
```

#### Connecting a VPN Client:
```go
key := "0123456789abcdef0123456789abcdef" // Must match server key

client, err := vpn.NewClient("server.example.com:1194", key)
if err != nil {
    log.Fatal(err)
}

// Connect to VPN server
if err := client.Connect(); err != nil {
    log.Fatal(err)
}
defer client.Disconnect()

fmt.Printf("Assigned IP: %s\n", client.GetAssignedIP())
fmt.Printf("Netmask: %s\n", client.GetNetmask())

// Start receiving packets
go client.StartTunnel()
```

#### Sending Packets Through VPN:
```go
// Create an IP packet (simplified example)
packet := make([]byte, 64)
// ... fill packet with IP header and data ...

if err := client.SendPacket(packet); err != nil {
    log.Printf("Failed to send packet: %v\n", err)
}
```

#### Checking Server Status:
```go
// Get number of connected clients
count := server.GetConnectedClients()
fmt.Printf("Connected clients: %d\n", count)

// Get client list
clients := server.GetClientList()
for id, ip := range clients {
    fmt.Printf("%s: %s\n", id, ip)
}
```

### Security Features

**Encryption**:
- **AES-GCM**: Authenticated encryption with associated data
- **256-bit keys**: Strong encryption (configurable to 128 or 192)
- **Unique nonces**: Every packet encrypted with unique nonce
- **Authentication**: GCM mode provides integrity verification

**Additional Security Considerations** (not implemented in basic version):
- Perfect Forward Secrecy with key rotation
- Certificate-based authentication
- Multi-factor authentication
- IP filtering and access control lists
- DDoS protection

### Architecture Notes

**This is a simplified VPN implementation. Production VPNs require**:
- TUN/TAP device integration for actual packet injection
- Proper IP routing and forwarding
- DNS configuration
- MTU handling and packet fragmentation
- Reconnection logic
- Bandwidth management
- Logging and monitoring

**Protocol Alternatives**:
- **OpenVPN**: SSL/TLS-based VPN
- **WireGuard**: Modern, fast VPN protocol
- **IPSec**: Industry-standard VPN protocol
- **L2TP**: Layer 2 tunneling

### Further Reading

- [Virtual Private Networks - RFC 4026](https://datatracker.ietf.org/doc/html/rfc4026)
- [AES-GCM - NIST SP 800-38D](https://csrc.nist.gov/publications/detail/sp/800-38d/final)
- [WireGuard Protocol](https://www.wireguard.com/protocol/)
- [OpenVPN Documentation](https://openvpn.net/community-resources/)
- [TUN/TAP Interfaces](https://www.kernel.org/doc/Documentation/networking/tuntap.txt)
