## Ifconfig Implementation

This package implements functionality similar to the `ifconfig` command for displaying network interface information.

### Features

- List all network interfaces
- Get specific interface details
- Display IP addresses and netmasks
- Show hardware (MAC) addresses
- Check interface flags (up, loopback, broadcast, etc.)
- Format output similar to traditional ifconfig

### Usage

#### Get All Interfaces

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/ifconfig"
)

func main() {
	interfaces, err := ifconfig.GetInterfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range interfaces {
		fmt.Print(ifconfig.FormatInterface(iface))
		fmt.Println()
	}
}
```

#### Get Specific Interface

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/ifconfig"
)

func main() {
	info, err := ifconfig.GetInterface("eth0")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(ifconfig.FormatInterface(*info))
}
```

#### Check Interface Status

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/ifconfig"
)

func main() {
	isUp, err := ifconfig.IsUp("eth0")
	if err != nil {
		log.Fatal(err)
	}

	if isUp {
		fmt.Println("Interface is up")
	} else {
		fmt.Println("Interface is down")
	}
}
```

### Output Format

```
eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>
    ether 02:42:ac:11:00:02
    inet 172.17.0.2 netmask 255.255.0.0

lo: flags=73<UP,LOOPBACK,RUNNING>
    inet 127.0.0.1 netmask 255.0.0.0
    inet ::1 netmask ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff
```

### Interface Information

Each interface includes:
- **Name**: Interface name (e.g., eth0, lo, wlan0)
- **Flags**: Status flags (UP, BROADCAST, RUNNING, MULTICAST, LOOPBACK)
- **Hardware Address**: MAC address (for physical interfaces)
- **IP Addresses**: IPv4 and IPv6 addresses
- **Netmask**: Network mask for each address

### Common Flags

- `UP`: Interface is active
- `BROADCAST`: Interface supports broadcasting
- `LOOPBACK`: Interface is a loopback interface
- `RUNNING`: Interface is operational
- `MULTICAST`: Interface supports multicast

### Use Cases

- Network diagnostics
- Interface monitoring
- Configuration verification
- Network troubleshooting
- System administration tools
