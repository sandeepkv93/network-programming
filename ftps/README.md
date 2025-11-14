## FTPS Client & Server

FTPS (FTP Secure) is an extension to the File Transfer Protocol that adds support for TLS (Transport Layer Security) and SSL (Secure Sockets Layer) cryptographic protocols.

## Table of Contents

1. [What is FTPS?](#what-is-ftps)
2. [How Does FTPS Work?](#how-does-ftps-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is FTPS?

FTPS is FTP with added security through TLS/SSL encryption. It provides secure file transfer capabilities while maintaining compatibility with the FTP protocol.

**FTPS vs SFTP**:
- **FTPS**: FTP over TLS/SSL (extends FTP)
- **SFTP**: SSH File Transfer Protocol (completely different protocol)

**Security Features**:
- Encrypted control and data connections
- Server and client authentication
- Data integrity verification

### How Does FTPS Work?

FTPS can operate in two modes:
1. **Explicit FTPS** (FTPES): Starts as plain FTP, upgrades to TLS
2. **Implicit FTPS**: TLS from the start (this implementation)

**Connection Process**:
1. TLS handshake
2. FTP authentication
3. Encrypted command/response exchange
4. Secure data transfer

### Understanding the Code

**Server**: Accepts TLS-encrypted FTP connections
**Client**: Connects to FTPS servers with TLS

**Commands Supported**:
- USER, PASS: Authentication
- PWD: Print working directory
- CWD: Change directory
- LIST: List files
- RETR: Retrieve file
- STOR: Store file
- QUIT: Disconnect

### Further Reading

- [FTPS - Wikipedia](https://en.wikipedia.org/wiki/FTPS)
- [RFC 4217 - FTP over TLS](https://tools.ietf.org/html/rfc4217)
