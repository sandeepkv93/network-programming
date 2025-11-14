## IP Spoofing

IP Spoofing is a technique where an attacker sends packets with a forged source IP address to hide their identity or impersonate another system.

**⚠️ SECURITY WARNING**: IP Spoofing is illegal when used maliciously. This implementation is for:
- Authorized penetration testing
- CTF competitions
- Security research
- Educational purposes only

## What is IP Spoofing?

IP Spoofing involves creating IP packets with a false source IP address. This can be used for both legitimate testing and malicious attacks.

**Use Cases**:
- **Security Testing**: Testing firewall and IDS effectiveness
- **DDoS Attacks**: Amplification and reflection attacks (illegal)
- **Bypassing Filters**: Evading IP-based access controls (illegal)
- **Research**: Understanding network security

**Detection**:
- Ingress filtering (RFC 2827)
- Egress filtering
- Deep packet inspection
- Anomaly detection

## How It Works

1. Create IP packet with custom header
2. Set source IP to spoofed address
3. Calculate checksums
4. Send via raw socket (requires privileges)
5. Responses go to spoofed IP (not attacker)

**Limitations**:
- Cannot receive responses (goes to spoofed IP)
- ISPs often filter spoofed packets
- Requires raw socket access
- Traceable through network forensics

## Further Reading

- [IP Spoofing - Wikipedia](https://en.wikipedia.org/wiki/IP_address_spoofing)
- [RFC 2827 - Network Ingress Filtering](https://tools.ietf.org/html/rfc2827)
