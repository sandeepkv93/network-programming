## SFTP Client & Server

SFTP (SSH File Transfer Protocol) is a network protocol that provides file access, file transfer, and file management over a reliable data stream using SSH.

## Table of Contents

1. [What is SFTP?](#what-is-sftp)
2. [How Does SFTP Work?](#how-does-sftp-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is SFTP?

SFTP is a secure file transfer protocol that runs over SSH (Secure Shell). Unlike FTPS, SFTP is not FTP over SSH but a completely separate protocol designed from the ground up.

**Key Features**:
- Single secure connection (port 22 by default)
- Strong encryption via SSH
- Authentication via SSH (passwords, keys)
- File operations (upload, download, delete, rename, etc.)

**SFTP vs FTP vs FTPS**:
- **FTP**: Plain text, insecure
- **FTPS**: FTP with TLS/SSL added
- **SFTP**: Completely different protocol using SSH

### How Does SFTP Work?

1. **SSH Connection**: Establish SSH connection
2. **Authentication**: Authenticate via password or key
3. **SFTP Subsystem**: Start SFTP subsystem over SSH
4. **File Operations**: Perform secure file operations
5. **Close**: Terminate SSH connection

**Advantages**:
- Only one port needed (SSH port)
- Firewall-friendly
- Platform-independent
- Comprehensive file management

### Understanding the Code

This is a simplified SFTP implementation. Production systems should use `github.com/pkg/sftp` package.

**Server**: SSH server with SFTP subsystem
**Client**: SSH client with SFTP support

**Authentication**: Password-based (key-based auth also possible)

### Further Reading

- [SFTP - Wikipedia](https://en.wikipedia.org/wiki/SSH_File_Transfer_Protocol)
- [SSH File Transfer Protocol](https://datatracker.ietf.org/doc/html/draft-ietf-secsh-filexfer)
- [pkg/sftp - Go SFTP implementation](https://github.com/pkg/sftp)
