## Packet Sniffer

A packet sniffer (also known as a packet analyzer or network analyzer) is a tool that captures and analyzes network packets. It is commonly used for network troubleshooting, security analysis, and protocol development.

## Table of Contents

1. [What is a Packet Sniffer?](#what-is-a-packet-sniffer)
2. [How Does Packet Sniffing Work?](#how-does-packet-sniffing-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is a Packet Sniffer?

A packet sniffer is a program that intercepts and logs network traffic that passes over a digital network. It captures each packet and decodes the packet's raw data, showing the values of various fields in the packet, and analyzes its content according to the appropriate protocol specifications.

**Common Uses**:
- **Network Troubleshooting**: Diagnosing network problems
- **Security Analysis**: Detecting unauthorized network access or attacks
- **Protocol Analysis**: Understanding how network protocols work
- **Performance Monitoring**: Analyzing network bandwidth usage

**Note**: Packet sniffing requires administrative/root privileges to create raw sockets.

### How Does Packet Sniffing Work?

1. **Capture**: The sniffer places the network interface in promiscuous mode to capture all packets
2. **Filter**: Apply filters to capture only specific types of packets (optional)
3. **Decode**: Parse the packet headers (Ethernet, IP, TCP/UDP, etc.)
4. **Display**: Present the packet information in a readable format
5. **Store**: Save captured packets for later analysis

**Network Interface Modes**:
- **Normal Mode**: The NIC only accepts packets destined for its MAC address
- **Promiscuous Mode**: The NIC accepts all packets on the network segment

**Packet Structure**:
```
[Ethernet Header][IP Header][TCP/UDP Header][Payload]
```

**Security Considerations**:
- Packet sniffing can capture sensitive information (passwords, etc.)
- Should only be used on networks you own or have permission to monitor
- Many networks use encryption (HTTPS, VPN) to protect against sniffing

### Understanding the Code

#### Data Structures:

- `PacketInfo`: Information about a captured packet:
  - `Timestamp`: When the packet was captured
  - `Protocol`: Protocol type (TCP, UDP, ICMP, etc.)
  - `SrcIP/DstIP`: Source and destination IP addresses
  - `SrcPort/DstPort`: Source and destination ports
  - `Length`: Total packet length
  - `PayloadSize`: Size of the packet payload

- `Sniffer`: The packet sniffer structure:
  - `Interface`: Network interface to sniff on
  - `Filter`: Packet filter (not fully implemented in basic version)
  - `PacketCount`: Maximum number of packets to capture
  - `OnPacket`: Callback function for each captured packet

#### Functions:

- `NewSniffer(iface string) *Sniffer`: Creates a new packet sniffer
- `Start() error`: Starts capturing packets
- `Stop()`: Stops the packet capture
- `parsePacket(data []byte, addr net.Addr) *PacketInfo`: Parses packet headers
- `GetPackets() []PacketInfo`: Returns all captured packets
- `GetStatistics() map[string]int`: Returns packet statistics by protocol

#### Features:

- Captures IPv4 packets (TCP, UDP, ICMP)
- Parses IP headers and transport layer headers
- Real-time packet display
- Packet statistics by protocol
- Configurable packet count limit
- Callback support for custom packet processing

#### Limitations:

- Requires elevated privileges (root on Linux, Administrator on Windows)
- Basic IPv4 support only (no IPv6 in this simple implementation)
- Limited protocol parsing (basic TCP/UDP/ICMP)
- No BPF (Berkeley Packet Filter) support

### Further Reading

- [Packet Analyzer - Wikipedia](https://en.wikipedia.org/wiki/Packet_analyzer)
- [Wireshark - Popular Packet Analyzer](https://www.wireshark.org/)
- [tcpdump - Command-line Packet Analyzer](https://www.tcpdump.org/)
- [Promiscuous Mode](https://en.wikipedia.org/wiki/Promiscuous_mode)
- [Raw Sockets](https://en.wikipedia.org/wiki/Network_socket#Raw_socket)
