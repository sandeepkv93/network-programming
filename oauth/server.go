package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// OAuthServer represents an OAuth 2.0 authorization server
type OAuthServer struct {
	Address      string
	clients      map[string]*Client
	tokens       map[string]*Token
	authCodes    map[string]*AuthCode
	mu           sync.RWMutex
	server       *http.Server
}

// Client represents an OAuth client
type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Name         string
}

// Token represents an access token
type Token struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int64
	RefreshToken string
	Scope        string
	CreatedAt    time.Time
}

// AuthCode represents an authorization code
type AuthCode struct {
	Code        string
	ClientID    string
	RedirectURI string
	Scope       string
	ExpiresAt   time.Time
}

// NewOAuthServer creates a new OAuth server
func NewOAuthServer(address string) *OAuthServer {
	return &OAuthServer{
		Address:   address,
		clients:   make(map[string]*Client),
		tokens:    make(map[string]*Token),
		authCodes: make(map[string]*AuthCode),
	}
}

// RegisterClient registers a new OAuth client
func (s *OAuthServer) RegisterClient(name, redirectURI string) (*Client, error) {
	clientID := generateRandomString(32)
	clientSecret := generateRandomString(64)

	client := &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Name:         name,
	}

	s.mu.Lock()
	s.clients[clientID] = client
	s.mu.Unlock()

	return client, nil
}

// Start starts the OAuth server
func (s *OAuthServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/authorize", s.handleAuthorize)
	mux.HandleFunc("/token", s.handleToken)
	mux.HandleFunc("/userinfo", s.handleUserInfo)

	s.server = &http.Server{
		Addr:    s.Address,
		Handler: mux,
	}

	log.Printf("OAuth 2.0 Server started on %s\n", s.Address)
	return s.server.ListenAndServe()
}

func (s *OAuthServer) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")
	responseType := r.URL.Query().Get("response_type")

	// Validate client
	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()

	if !exists || client.RedirectURI != redirectURI {
		http.Error(w, "Invalid client", http.StatusBadRequest)
		return
	}

	if responseType != "code" {
		http.Error(w, "Unsupported response type", http.StatusBadRequest)
		return
	}

	// Generate authorization code
	code := generateRandomString(32)
	authCode := &AuthCode{
		Code:        code,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	s.mu.Lock()
	s.authCodes[code] = authCode
	s.mu.Unlock()

	// Redirect back to client
	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *OAuthServer) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	grantType := r.FormValue("grant_type")
	code := r.FormValue("code")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	if grantType != "authorization_code" {
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
		return
	}

	// Validate client credentials
	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()

	if !exists || client.ClientSecret != clientSecret {
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	// Validate authorization code
	s.mu.Lock()
	authCode, exists := s.authCodes[code]
	if exists {
		delete(s.authCodes, code) // One-time use
	}
	s.mu.Unlock()

	if !exists || authCode.ExpiresAt.Before(time.Now()) || authCode.ClientID != clientID {
		http.Error(w, "Invalid authorization code", http.StatusBadRequest)
		return
	}

	// Generate access token
	token := &Token{
		AccessToken:  generateRandomString(64),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: generateRandomString(64),
		Scope:        authCode.Scope,
		CreatedAt:    time.Now(),
	}

	s.mu.Lock()
	s.tokens[token.AccessToken] = token
	s.mu.Unlock()

	// Return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  token.AccessToken,
		"token_type":    token.TokenType,
		"expires_in":    token.ExpiresIn,
		"refresh_token": token.RefreshToken,
		"scope":         token.Scope,
	})
}

func (s *OAuthServer) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}

	// Extract token
	var accessToken string
	fmt.Sscanf(authHeader, "Bearer %s", &accessToken)

	// Validate token
	s.mu.RLock()
	token, exists := s.tokens[accessToken]
	s.mu.RUnlock()

	if !exists || time.Since(token.CreatedAt) > time.Duration(token.ExpiresIn)*time.Second {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sub":   "user123",
		"name":  "Test User",
		"email": "user@example.com",
	})
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

// Stop stops the OAuth server
func (s *OAuthServer) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// ParseRedirectURL helper for client
func ParseRedirectURL(redirectURL string) (code string, state string, err error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return "", "", err
	}

	code = u.Query().Get("code")
	state = u.Query().Get("state")
	return code, state, nil
}
