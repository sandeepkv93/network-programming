package twofa

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// TwoFactorAuth represents a 2FA authentication system
type TwoFactorAuth struct {
	Issuer string
}

// NewTwoFactorAuth creates a new 2FA instance
func NewTwoFactorAuth(issuer string) *TwoFactorAuth {
	return &TwoFactorAuth{
		Issuer: issuer,
	}
}

// GenerateSecret generates a new random secret key
func (t *TwoFactorAuth) GenerateSecret() (string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateQRCodeURL generates a URL for QR code
func (t *TwoFactorAuth) GenerateQRCodeURL(accountName, secret string) string {
	// Format: otpauth://totp/Issuer:AccountName?secret=SECRET&issuer=Issuer
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		t.Issuer, accountName, secret, t.Issuer)
}

// GenerateTOTP generates a Time-based One-Time Password
func (t *TwoFactorAuth) GenerateTOTP(secret string) (string, error) {
	return t.GenerateTOTPAtTime(secret, time.Now())
}

// GenerateTOTPAtTime generates TOTP for a specific time
func (t *TwoFactorAuth) GenerateTOTPAtTime(secret string, timestamp time.Time) (string, error) {
	// Decode secret
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", fmt.Errorf("invalid secret: %v", err)
	}

	// Calculate time counter (30 second intervals)
	counter := uint64(timestamp.Unix()) / 30

	// Convert counter to bytes
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	// Calculate HMAC-SHA1
	h := hmac.New(sha1.New, key)
	h.Write(buf)
	hash := h.Sum(nil)

	// Dynamic truncation
	offset := hash[len(hash)-1] & 0x0F
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

	// Generate 6-digit code
	code := truncated % 1000000

	return fmt.Sprintf("%06d", code), nil
}

// ValidateTOTP validates a TOTP code
func (t *TwoFactorAuth) ValidateTOTP(secret, code string) (bool, error) {
	return t.ValidateTOTPWithWindow(secret, code, 1)
}

// ValidateTOTPWithWindow validates TOTP with time window
func (t *TwoFactorAuth) ValidateTOTPWithWindow(secret, code string, window int) (bool, error) {
	now := time.Now()

	// Check current time and windows before/after
	for i := -window; i <= window; i++ {
		testTime := now.Add(time.Duration(i*30) * time.Second)
		expectedCode, err := t.GenerateTOTPAtTime(secret, testTime)
		if err != nil {
			return false, err
		}

		if code == expectedCode {
			return true, nil
		}
	}

	return false, nil
}

// GenerateBackupCodes generates backup recovery codes
func (t *TwoFactorAuth) GenerateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)

	for i := 0; i < count; i++ {
		code := make([]byte, 8)
		_, err := rand.Read(code)
		if err != nil {
			return nil, err
		}

		// Format as XXXX-XXXX
		codes[i] = fmt.Sprintf("%04X-%04X",
			binary.BigEndian.Uint16(code[0:2]),
			binary.BigEndian.Uint16(code[2:4]))
	}

	return codes, nil
}

// Example user registration flow

// UserSetup represents user 2FA setup
type UserSetup struct {
	AccountName string
	Secret      string
	QRCodeURL   string
	BackupCodes []string
}

// SetupUser sets up 2FA for a new user
func (t *TwoFactorAuth) SetupUser(accountName string) (*UserSetup, error) {
	// Generate secret
	secret, err := t.GenerateSecret()
	if err != nil {
		return nil, err
	}

	// Generate QR code URL
	qrURL := t.GenerateQRCodeURL(accountName, secret)

	// Generate backup codes
	backupCodes, err := t.GenerateBackupCodes(10)
	if err != nil {
		return nil, err
	}

	return &UserSetup{
		AccountName: accountName,
		Secret:      secret,
		QRCodeURL:   qrURL,
		BackupCodes: backupCodes,
	}, nil
}
