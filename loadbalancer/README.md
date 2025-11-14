## Load Balancer

This package implements an HTTP load balancer that distributes requests across multiple backend servers.

### Features

- **Load Distribution**: Distributes traffic across multiple backends
- **Multiple Strategies**:
  - Round-robin: Cycles through backends equally
  - Least-connections: Routes to backend with fewest active connections
- **Health Checks**: Automatically detects and removes unhealthy backends
- **Connection Tracking**: Monitors active connections per backend
- **Reverse Proxy**: Transparent request forwarding

### Usage

#### Basic Load Balancer

```go
package main

import (
	"log"
	"time"
	"network-programming/loadbalancer"
)

func main() {
	// Create load balancer with round-robin strategy
	lb := loadbalancer.NewLoadBalancer(":8080", "round-robin")

	// Add backend servers
	lb.AddBackend("http://localhost:8001")
	lb.AddBackend("http://localhost:8002")
	lb.AddBackend("http://localhost:8003")

	// Start load balancer
	if err := lb.Start(); err != nil {
		log.Fatal(err)
	}

	// Run for some time
	time.Sleep(1 * time.Minute)

	// Get backend statistics
	stats := lb.GetBackendStats()
	for _, stat := range stats {
		log.Printf("Backend: %v\n", stat)
	}

	lb.Stop()
}
```

#### Least-Connections Strategy

```go
package main

import (
	"log"
	"network-programming/loadbalancer"
)

func main() {
	// Create load balancer with least-connections strategy
	lb := loadbalancer.NewLoadBalancer(":8080", "least-connections")

	lb.AddBackend("http://localhost:8001")
	lb.AddBackend("http://localhost:8002")
	lb.AddBackend("http://localhost:8003")

	if err := lb.Start(); err != nil {
		log.Fatal(err)
	}

	select {} // Run forever
}
```

### Load Balancing Strategies

#### Round-Robin
- Distributes requests equally across all backends
- Each backend receives requests in turn
- Simple and fair distribution
- Doesn't account for backend load or capacity

#### Least-Connections
- Routes requests to backend with fewest active connections
- Better for backends with varying capacities
- Automatically adapts to backend performance
- Tracks connections in real-time

### Health Checking

The load balancer automatically:
- Checks each backend every 10 seconds
- Makes HTTP GET request to backend
- Marks backend as DOWN if:
  - Request fails
  - Response status code >= 500
- Removes unhealthy backends from rotation
- Re-adds backends when they become healthy again

### How it Works

1. **Client** sends request to load balancer
2. **Load Balancer** selects a backend using chosen strategy
3. Load balancer forwards request to selected backend
4. **Backend** processes request and sends response
5. Load balancer forwards response back to client
6. Connection count is updated

### Architecture

```
Client Requests
      ↓
Load Balancer (:8080)
      ↓
   Strategy
      ↓
Backend Selection
      ↓
┌─────┬─────┬─────┐
│  B1 │  B2 │  B3 │
│8001 │8002 │8003 │
└─────┴─────┴─────┘
```

### Backend Statistics

Each backend tracks:
- **URL**: Backend server address
- **Alive**: Health status (true/false)
- **Connections**: Current active connections

### Use Cases

- **High Availability**: Distribute traffic across multiple servers
- **Scalability**: Add more backends to handle increased load
- **Fault Tolerance**: Automatically route around failed backends
- **Performance**: Distribute load to prevent server overload
- **Rolling Deployments**: Gradually shift traffic to new versions

### Example: Setting up Backends

Start multiple backend servers on different ports:

```bash
# Terminal 1
go run backend1.go :8001

# Terminal 2
go run backend2.go :8002

# Terminal 3
go run backend3.go :8003

# Terminal 4 - Start load balancer
go run loadbalancer.go
```

Then access the load balancer:
```bash
curl http://localhost:8080/
```

### Note

This is an educational load balancer. Production load balancers typically include:
- Session persistence (sticky sessions)
- SSL/TLS termination
- Request rate limiting
- Geographic routing
- Weighted backends
- Advanced health checks (custom endpoints, timeouts)
- Metrics and monitoring
- Dynamic backend registration/removal
- Circuit breakers
- Request queuing
- Connection pooling
- HTTP/2 and WebSocket support
