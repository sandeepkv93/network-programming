## Heartbeat Server

A heartbeat monitoring system for tracking the health and availability of distributed clients/services.

## Features

- Periodic heartbeat monitoring
- Client timeout detection
- Status callbacks (onAlive, onDead)
- Client registration and tracking
- Automatic reconnection (client)
- Configurable intervals and timeouts

## Usage

### Server
```go
server := heartbeat.NewServer(":5000", 30*time.Second)

// Set callbacks for status changes
server.SetCallbacks(
    func(clientID string) {
        log.Printf("Client %s is alive\n", clientID)
    },
    func(clientID string) {
        log.Printf("Client %s is dead\n", clientID)
    },
)

server.Start()
```

### Client
```go
client := heartbeat.NewClient("localhost:5000", 10*time.Second)
client.Start()

// Client will send heartbeats every 10 seconds
// Server will mark it dead if no heartbeat for 30 seconds
```

## Use Cases

- Service health monitoring
- Distributed system coordination
- Failure detection
- Load balancer health checks
