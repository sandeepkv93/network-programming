## Tunneling

TCP tunneling implementation for forwarding traffic between networks, useful for bypassing firewalls and accessing services through intermediary servers.

## Features

- TCP port forwarding
- Bidirectional data transfer
- Connection tracking
- Multiple concurrent tunnels

## Usage

### Forward Tunnel (Local to Remote)
```go
// Forward local port 8080 to remote service
server := tunneling.NewServer(":8080", "remote.example.com:80")
server.Start()
```

### Reverse Tunnel
```go
// Expose local service through remote server
client := tunneling.NewClient("tunnel.example.com:9000", "localhost:3000")
client.Connect()
```

## Use Cases

- Access services behind NAT/firewall
- Secure service exposure
- Port forwarding
- Network bridging
