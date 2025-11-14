## Two-Factor Authentication (2FA)

Two-Factor Authentication adds an extra layer of security by requiring a second form of verification beyond just a password.

## What is 2FA?

2FA requires two different authentication factors:
1. **Something you know**: Password
2. **Something you have**: Phone, hardware token, authenticator app

**Common Methods**:
- **TOTP** (Time-based One-Time Password): Google Authenticator, Authy
- **SMS**: Text message codes
- **Hardware tokens**: YubiKey, security keys
- **Biometric**: Fingerprint, face recognition

## TOTP (RFC 6238)

Time-based One-Time Password algorithm:
1. Shared secret stored on server and client device
2. Current time divided into 30-second intervals
3. HMAC-SHA1 of secret and time counter
4. Generate 6-digit code

**Features**:
- Works offline (no network required)
- Synchronized by time
- New code every 30 seconds
- Backup codes for recovery

## Implementation

This package provides:
- Secret generation
- QR code URL generation
- TOTP generation and validation
- Backup code generation
- Time window validation

## Further Reading

- [RFC 6238 - TOTP](https://tools.ietf.org/html/rfc6238)
- [Two-Factor Authentication](https://en.wikipedia.org/wiki/Multi-factor_authentication)
