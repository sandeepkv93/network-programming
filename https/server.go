package https

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

// Server represents an HTTPS server
type Server struct {
	Address  string
	CertFile string
	KeyFile  string
	server   *http.Server
}

// Response represents a JSON response
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Secure    bool      `json:"secure"`
}

// NewServer creates a new HTTPS server
func NewServer(address, certFile, keyFile string) *Server {
	return &Server{
		Address:  address,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}

// Start starts the HTTPS server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/secure", s.handleSecure)

	// Configure TLS
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	s.server = &http.Server{
		Addr:         s.Address,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("HTTPS Server starting on %s\n", s.Address)
	log.Printf("Using certificate: %s\n", s.CertFile)
	log.Printf("Using key: %s\n", s.KeyFile)

	return s.server.ListenAndServeTLS(s.CertFile, s.KeyFile)
}

// Stop stops the HTTPS server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<!DOCTYPE html>
<html>
<head>
    <title>HTTPS Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .secure { color: green; font-weight: bold; }
        .info { background: #f0f0f0; padding: 15px; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>HTTPS Server</h1>
    <p class="secure">🔒 Secure Connection Established</p>
    <div class="info">
        <h2>Connection Information</h2>
        <p><strong>Protocol:</strong> HTTPS (TLS)</p>
        <p><strong>Path:</strong> %s</p>
        <p><strong>Time:</strong> %s</p>
        <p><strong>Remote Address:</strong> %s</p>
    </div>
    <h2>API Endpoints</h2>
    <ul>
        <li><a href="/api/status">/api/status</a> - Server status</li>
        <li><a href="/api/secure">/api/secure</a> - Secure data endpoint</li>
    </ul>
</body>
</html>`
	fmt.Fprintf(w, html, r.URL.Path, time.Now().Format(time.RFC3339), r.RemoteAddr)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	response := Response{
		Message:   "HTTPS server is running",
		Timestamp: time.Now(),
		Path:      r.URL.Path,
		Secure:    r.TLS != nil,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleSecure(w http.ResponseWriter, r *http.Request) {
	// Check if connection is secure
	if r.TLS == nil {
		http.Error(w, "Secure connection required", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	tlsInfo := map[string]interface{}{
		"version":            getTLSVersion(r.TLS.Version),
		"cipher_suite":       tls.CipherSuiteName(r.TLS.CipherSuite),
		"server_name":        r.TLS.ServerName,
		"negotiated_protocol": r.TLS.NegotiatedProtocol,
		"secure":             true,
		"timestamp":          time.Now(),
	}

	json.NewEncoder(w).Encode(tlsInfo)
}

func getTLSVersion(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

// GenerateSelfSignedCert generates a self-signed certificate for testing
func GenerateSelfSignedCert(certFile, keyFile string) error {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Network Programming"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	// Create self-signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	// Write certificate to file
	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("failed to create cert file: %v", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return fmt.Errorf("failed to write certificate: %v", err)
	}

	// Write private key to file
	keyOut, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("failed to create key file: %v", err)
	}
	defer keyOut.Close()

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	log.Printf("Generated self-signed certificate: %s\n", certFile)
	log.Printf("Generated private key: %s\n", keyFile)

	return nil
}
