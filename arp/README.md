## ARP Implementation

This package implements functionality for working with the ARP (Address Resolution Protocol) table, which maps IP addresses to MAC addresses.

### Features

- Read system ARP table
- Get specific ARP entries by IP address
- Display ARP cache in formatted table
- Filter complete/incomplete entries
- Get entries by network device
- Resolve MAC addresses from IP addresses
- Get MAC address of network interfaces

### Usage

#### Get Full ARP Table

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/arp"
)

func main() {
	entries, err := arp.GetARPTable()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(arp.FormatTable(entries))
}
```

#### Get Specific ARP Entry

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/arp"
)

func main() {
	entry, err := arp.GetARPEntry("192.168.1.1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(arp.FormatEntry(*entry))
}
```

#### Resolve MAC Address

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/arp"
)

func main() {
	mac, err := arp.ResolveMAC("192.168.1.1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MAC Address: %s\n", mac)
}
```

#### Get Complete Entries Only

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/arp"
)

func main() {
	entries, err := arp.GetCompleteEntries()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(arp.FormatTable(entries))
}
```

#### Get MAC of Network Interface

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/arp"
)

func main() {
	mac, err := arp.GetMACByInterface("eth0")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("eth0 MAC: %s\n", mac)
}
```

### Output Format

```
IP Address      HW Address        HW Type  Flags  Mask     Device
--------------------------------------------------------------------------------
192.168.1.1     aa:bb:cc:dd:ee:ff 0x1      0x2    *        eth0
192.168.1.100   11:22:33:44:55:66 0x1      0x2    *        eth0
10.0.0.1        ff:ee:dd:cc:bb:aa 0x1      0x2    *        wlan0
```

### ARP Entry Fields

- **IP Address**: IPv4 address
- **HW Address**: MAC (hardware) address
- **HW Type**: Hardware type (0x1 = Ethernet)
- **Flags**: Entry flags (0x2 = complete, 0x0 = incomplete)
- **Mask**: Published proxy ARP entry mask
- **Device**: Network interface device

### ARP Flags

- `0x2` (C): Complete entry with valid MAC address
- `0x4` (M): Permanent entry
- `0x8` (P): Published entry
- `0x0`: Incomplete entry

### How it Works

1. Reads `/proc/net/arp` on Linux systems
2. Parses each line into structured ARP entries
3. Provides filtering and lookup capabilities
4. Maps IP addresses to MAC addresses
5. Shows which network interface each entry uses

### Platform Support

This implementation reads from `/proc/net/arp`, which is Linux-specific. For other platforms:
- **Windows**: Use `arp -a` command or Windows APIs
- **macOS**: Use `arp -a` command or BSD socket APIs

### Use Cases

- Network discovery
- MAC address lookup
- Network diagnostics
- Security monitoring
- Duplicate IP detection
- Network mapping
- Device identification

### Understanding ARP

ARP (Address Resolution Protocol) is used to:
1. Map IP addresses (Layer 3) to MAC addresses (Layer 2)
2. Enable communication within a local network
3. Cache mappings to reduce network traffic
4. Automatically update when devices communicate

When a device wants to communicate with another device on the same network, it uses ARP to find the target's MAC address.

### Limitations

- Read-only access to ARP table
- Cannot add/delete ARP entries (requires privileges)
- Platform-specific implementation
- Static ARP entries may require special handling
