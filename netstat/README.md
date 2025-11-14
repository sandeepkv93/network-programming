## Netstat Implementation

This package implements functionality similar to the `netstat` command for displaying network connections and statistics.

### Features

- List all active network connections
- Display TCP connections
- Display UDP connections
- Show connection states (ESTABLISHED, LISTEN, etc.)
- Show local and remote addresses/ports
- Filter listening ports
- Parse /proc/net/tcp and /proc/net/udp (Linux)

### Usage

#### Get All Connections

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/netstat"
)

func main() {
	connections, err := netstat.GetConnections()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(netstat.FormatConnections(connections))
}
```

#### Get TCP Connections Only

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/netstat"
)

func main() {
	connections, err := netstat.GetTCPConnections()
	if err != nil {
		log.Fatal(err)
	}

	for _, conn := range connections {
		fmt.Println(netstat.FormatConnection(conn))
	}
}
```

#### Get Listening Ports

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/netstat"
)

func main() {
	listening, err := netstat.GetListeningPorts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening ports:")
	for _, conn := range listening {
		fmt.Println(netstat.FormatConnection(conn))
	}
}
```

### Output Format

```
Proto  Local Address           Foreign Address         State
tcp    0.0.0.0:22              0.0.0.0:0               LISTEN
tcp    127.0.0.1:3306          0.0.0.0:0               LISTEN
tcp    192.168.1.10:52341      93.184.216.34:443       ESTABLISHED
udp    0.0.0.0:68              0.0.0.0:0               LISTEN
```

### Connection States

TCP connections can be in various states:

- **ESTABLISHED**: Connection is established
- **LISTEN**: Server is waiting for connections
- **SYN_SENT**: Attempting to establish connection
- **SYN_RECV**: Connection request received
- **FIN_WAIT1**: Connection is closing
- **FIN_WAIT2**: Connection is closed, waiting for remote shutdown
- **TIME_WAIT**: Waiting to ensure remote received shutdown
- **CLOSE**: Connection is closed
- **CLOSE_WAIT**: Remote has closed connection
- **LAST_ACK**: Waiting for connection termination acknowledgment
- **CLOSING**: Both sides closing simultaneously

### Connection Information

Each connection includes:
- **Protocol**: tcp or udp
- **Local Address**: Local IP address
- **Local Port**: Local port number
- **Remote Address**: Remote IP address
- **Remote Port**: Remote port number
- **State**: Connection state (TCP only)
- **UID**: User ID of the process
- **Inode**: Socket inode number

### Platform Support

This implementation reads from `/proc/net/tcp` and `/proc/net/udp`, which is Linux-specific. For cross-platform support, platform-specific implementations would be needed for:
- Windows: Use netstat command or Windows APIs
- macOS: Use netstat command or BSD socket APIs

### Use Cases

- Monitor active connections
- Identify listening services
- Debug network applications
- Security auditing
- Network troubleshooting
- Port usage verification
