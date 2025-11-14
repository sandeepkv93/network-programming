## HTTPS Client & Server

HTTPS (Hypertext Transfer Protocol Secure) is an extension of HTTP that uses TLS/SSL to encrypt communication between clients and servers, providing secure data transfer over the internet.

## Table of Contents

1. [What is HTTPS?](#what-is-https)
2. [How Does HTTPS Work?](#how-does-https-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is HTTPS?

HTTPS is the secure version of HTTP, using TLS (Transport Layer Security) or its predecessor SSL (Secure Sockets Layer) to encrypt data transmitted between a client and server. This encryption protects against eavesdropping, tampering, and message forgery.

**Key Features**:
- **Encryption**: Data is encrypted during transmission
- **Authentication**: Verifies the identity of the server (and optionally the client)
- **Integrity**: Ensures data hasn't been tampered with during transit
- **Privacy**: Protects user privacy and sensitive information

**Why HTTPS Matters**:
- Protects user credentials, payment information, and personal data
- SEO benefits (search engines prefer HTTPS sites)
- Browser trust indicators (padlock icon)
- Required for modern web features (geolocation, service workers, etc.)
- Compliance with security standards (PCI DSS, GDPR, etc.)

### How Does HTTPS Work?

HTTPS uses TLS/SSL to create a secure connection through a process called the TLS handshake:

1. **Client Hello**: Client sends supported TLS versions and cipher suites
2. **Server Hello**: Server selects TLS version and cipher suite
3. **Certificate Exchange**: Server sends its SSL/TLS certificate
4. **Certificate Verification**: Client verifies the certificate
5. **Key Exchange**: Client and server establish session keys
6. **Secure Communication**: All data is encrypted with session keys

**TLS Handshake Process**:
```
Client                                Server
  |                                     |
  |------- Client Hello --------------->|
  |                                     |
  |<------ Server Hello -----------------|
  |<------ Certificate ------------------|
  |<------ Server Key Exchange ----------|
  |<------ Server Hello Done ------------|
  |                                     |
  |------- Client Key Exchange -------->|
  |------- Change Cipher Spec --------->|
  |------- Finished ------------------->|
  |                                     |
  |<------ Change Cipher Spec -----------|
  |<------ Finished ---------------------|
  |                                     |
  |<====== Encrypted Data =============>|
```

**Certificates**:
- **Self-Signed**: Generated locally, not trusted by browsers (for testing)
- **CA-Signed**: Issued by Certificate Authority, trusted by browsers
- **Contains**: Public key, domain name, expiration date, issuer info

**Cipher Suites**:
- Combination of algorithms for key exchange, authentication, encryption, and MAC
- Example: `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384`
  - Key Exchange: ECDHE (Elliptic Curve Diffie-Hellman Ephemeral)
  - Authentication: RSA
  - Encryption: AES-256-GCM
  - MAC: SHA384

### Understanding the Code

#### Data Structures:

**Server**:
- `Address`: Server address and port
- `CertFile`: Path to TLS certificate file
- `KeyFile`: Path to private key file
- `server`: HTTP server instance

**Client**:
- `HTTPClient`: Configured HTTP client with TLS
- `SkipVerify`: Whether to skip certificate verification (for testing)
- `ClientCertFile`: Client certificate for mutual TLS
- `ClientKeyFile`: Client private key
- `CACertFile`: Custom CA certificate

#### Functions:

**Server**:
- `NewServer(address, certFile, keyFile string) *Server`: Creates HTTPS server
- `Start() error`: Starts the HTTPS server with TLS configuration
- `Stop() error`: Stops the server
- `handleRoot(w, r)`: Serves main page with connection info
- `handleStatus(w, r)`: Returns server status as JSON
- `handleSecure(w, r)`: Returns TLS connection details
- `GenerateSelfSignedCert(certFile, keyFile string) error`: Generates self-signed certificate for testing

**Client**:
- `NewClient(skipVerify bool) *Client`: Creates HTTPS client
- `Initialize() error`: Initializes client with TLS configuration
- `Get(url string) (string, error)`: Performs HTTPS GET request
- `Post(url, contentType, body string) (string, error)`: Performs HTTPS POST request
- `GetTLSInfo(url string) (map[string]string, error)`: Retrieves TLS connection information

#### TLS Configuration:

**Server**:
- Minimum TLS version: 1.2
- Preferred cipher suites (strong encryption)
- HSTS (HTTP Strict Transport Security) headers
- Read/write/idle timeouts

**Client**:
- Certificate verification (can be disabled for testing)
- Support for custom CA certificates
- Support for client certificates (mutual TLS)
- Connection pooling and timeouts

#### Features:

- Self-signed certificate generation for development
- TLS 1.2+ with strong cipher suites
- Connection information endpoints
- HSTS support
- Mutual TLS support (client certificates)
- Custom CA certificate support
- Comprehensive error handling

#### Example Usage:

```go
// Generate self-signed certificate (for testing)
err := GenerateSelfSignedCert("server.crt", "server.key")
if err != nil {
    log.Fatal(err)
}

// Start HTTPS server
server := NewServer(":8443", "server.crt", "server.key")
go server.Start()

// Create HTTPS client (skip verification for self-signed cert)
client := NewClient(true)
client.Initialize()

// Make secure request
response, err := client.Get("https://localhost:8443/api/status")
if err != nil {
    log.Fatal(err)
}
fmt.Println(response)

// Get TLS connection info
tlsInfo, err := client.GetTLSInfo("https://localhost:8443")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("TLS Info: %+v\n", tlsInfo)
```

### Further Reading

- [HTTPS - Wikipedia](https://en.wikipedia.org/wiki/HTTPS)
- [TLS - Transport Layer Security](https://en.wikipedia.org/wiki/Transport_Layer_Security)
- [How HTTPS Works](https://howhttps.works/)
- [SSL/TLS Certificates](https://www.cloudflare.com/learning/ssl/what-is-an-ssl-certificate/)
- [Mozilla TLS Configuration](https://wiki.mozilla.org/Security/Server_Side_TLS)
- [Let's Encrypt - Free SSL Certificates](https://letsencrypt.org/)
