## Port Forwarding

Port forwarding is a technique that allows external devices to access services on a private network by redirecting communication requests from one address and port number combination to another. It's commonly used for accessing services behind NAT or firewalls.

## Table of Contents

1. [What is Port Forwarding?](#what-is-port-forwarding)
2. [How Does Port Forwarding Work?](#how-does-port-forwarding-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is Port Forwarding?

Port forwarding is the act of redirecting network traffic from one IP address and port number combination to another. It creates a mapping between a port on the local machine and a port on a remote machine, allowing traffic to be forwarded between them.

**Common Use Cases**:
- **Remote Access**: Access services behind a firewall or NAT
- **SSH Tunneling**: Securely forward traffic through an encrypted SSH connection
- **Load Distribution**: Distribute traffic across multiple backend servers
- **Development**: Access services running in containers or virtual machines
- **Bypassing Restrictions**: Access blocked services through an intermediary

**Types of Port Forwarding**:
- **Local Port Forwarding**: Forward local port to remote destination
- **Remote Port Forwarding**: Forward remote port to local destination
- **Dynamic Port Forwarding**: Create a SOCKS proxy for dynamic forwarding

### How Does Port Forwarding Work?

1. **Listen**: The forwarder listens on a local address and port
2. **Accept**: When a client connects to the local port, accept the connection
3. **Connect**: Establish a connection to the remote destination
4. **Forward**: Bidirectionally copy data between client and remote connections
5. **Close**: When either side closes, terminate both connections

**Traffic Flow**:
```
Client -> Local Port (Forwarder) -> Remote Port -> Service
       <-                       <-             <-
```

**Example Scenarios**:

1. **SSH Tunnel**: `localhost:8080 -> remote-server:80`
   - Access a web server running on remote-server:80 via localhost:8080

2. **Database Access**: `localhost:5432 -> db-server:5432`
   - Access a database server through a secure tunnel

3. **Port Mapping**: `0.0.0.0:8000 -> internal-service:3000`
   - Make an internal service accessible on a different port

### Understanding the Code

#### Data Structures:

- `Forwarder`: The main port forwarder structure:
  - `LocalAddr`: Address and port to listen on (e.g., "localhost:8080")
  - `RemoteAddr`: Destination address and port (e.g., "example.com:80")
  - `listener`: TCP listener for accepting connections
  - `connections`: Map of active client connections
  - `running`: Whether the forwarder is currently running

#### Functions:

- `NewForwarder(localAddr, remoteAddr string) *Forwarder`: Creates a new forwarder
- `Start() error`: Starts listening and forwarding
- `Stop() error`: Stops the forwarder and closes all connections
- `acceptConnections()`: Accepts incoming client connections
- `handleConnection(clientConn net.Conn)`: Handles a single connection
  - Connects to remote destination
  - Bidirectionally copies data using goroutines
  - Handles cleanup when either side closes
- `GetActiveConnections() int`: Returns number of active connections
- `IsRunning() bool`: Returns whether the forwarder is running

#### Features:

- Bidirectional data forwarding using `io.Copy`
- Concurrent connection handling with goroutines
- Automatic cleanup of closed connections
- Connection tracking and statistics
- Graceful shutdown of all connections

#### Usage Example:

```go
// Forward local port 8080 to google.com:80
forwarder := NewForwarder("localhost:8080", "google.com:80")
err := forwarder.Start()
if err != nil {
    log.Fatal(err)
}

// Now you can access google.com by visiting http://localhost:8080

// Later, stop the forwarder
forwarder.Stop()
```

### Further Reading

- [Port Forwarding - Wikipedia](https://en.wikipedia.org/wiki/Port_forwarding)
- [SSH Tunneling](https://www.ssh.com/academy/ssh/tunneling)
- [NAT (Network Address Translation)](https://en.wikipedia.org/wiki/Network_address_translation)
- [SOCKS Proxy](https://en.wikipedia.org/wiki/SOCKS)
