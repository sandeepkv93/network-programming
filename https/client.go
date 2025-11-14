package https

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client represents an HTTPS client
type Client struct {
	HTTPClient       *http.Client
	SkipVerify       bool
	ClientCertFile   string
	ClientKeyFile    string
	CACertFile       string
}

// NewClient creates a new HTTPS client
func NewClient(skipVerify bool) *Client {
	return &Client{
		SkipVerify: skipVerify,
	}
}

// Initialize initializes the HTTPS client with TLS configuration
func (c *Client) Initialize() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.SkipVerify,
		MinVersion:         tls.VersionTLS12,
	}

	// Load CA certificate if provided
	if c.CACertFile != "" {
		caCert, err := os.ReadFile(c.CACertFile)
		if err != nil {
			return fmt.Errorf("failed to read CA certificate: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return fmt.Errorf("failed to parse CA certificate")
		}

		tlsConfig.RootCAs = caCertPool
	}

	// Load client certificate if provided
	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.ClientCertFile, c.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load client certificate: %v", err)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	c.HTTPClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:     tlsConfig,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return nil
}

// Get performs an HTTPS GET request
func (c *Client) Get(url string) (string, error) {
	if c.HTTPClient == nil {
		if err := c.Initialize(); err != nil {
			return "", err
		}
	}

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return string(body), nil
}

// Post performs an HTTPS POST request
func (c *Client) Post(url, contentType, body string) (string, error) {
	if c.HTTPClient == nil {
		if err := c.Initialize(); err != nil {
			return "", err
		}
	}

	resp, err := c.HTTPClient.Post(url, contentType, io.NopCloser(bytes.NewBufferString(body)))
	if err != nil {
		return "", fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return string(respBody), nil
}

// GetTLSInfo retrieves TLS connection information
func (c *Client) GetTLSInfo(url string) (map[string]string, error) {
	if c.HTTPClient == nil {
		if err := c.Initialize(); err != nil {
			return nil, err
		}
	}

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.TLS == nil {
		return nil, fmt.Errorf("no TLS connection info available")
	}

	info := map[string]string{
		"version":             getTLSVersion(resp.TLS.Version),
		"cipher_suite":        tls.CipherSuiteName(resp.TLS.CipherSuite),
		"server_name":         resp.TLS.ServerName,
		"negotiated_protocol": resp.TLS.NegotiatedProtocol,
	}

	if len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		info["subject"] = cert.Subject.String()
		info["issuer"] = cert.Issuer.String()
		info["not_before"] = cert.NotBefore.Format(time.RFC3339)
		info["not_after"] = cert.NotAfter.Format(time.RFC3339)
	}

	return info, nil
}
