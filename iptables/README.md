## IP Tables

IPTables is a user-space utility program for configuring Linux kernel firewall. This implementation mimics iptables functionality in Go.

## What is IPTables?

IPTables uses chains and rules to filter network traffic. It's the standard firewall solution on Linux systems.

**Tables**:
- **filter**: Packet filtering (default)
- **nat**: Network Address Translation
- **mangle**: Packet alteration

**Chains**:
- **INPUT**: Incoming packets
- **OUTPUT**: Outgoing packets
- **FORWARD**: Routed packets
- **PREROUTING**: Before routing
- **POSTROUTING**: After routing

**Targets**:
- **ACCEPT**: Allow packet
- **DROP**: Silently discard
- **REJECT**: Discard with error
- **MASQUERADE**: NAT

## Further Reading

- [iptables - Wikipedia](https://en.wikipedia.org/wiki/Iptables)
- [Netfilter](https://www.netfilter.org/)
