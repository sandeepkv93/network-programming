## OAuth 2.0 Client & Server

OAuth 2.0 is an authorization framework that enables applications to obtain limited access to user accounts on an HTTP service.

## What is OAuth 2.0?

OAuth 2.0 is the industry-standard protocol for authorization, allowing third-party applications to access user data without exposing passwords.

**Grant Types**:
- **Authorization Code**: Most secure, used by server-side apps
- **Implicit**: For browser-based apps (deprecated)
- **Client Credentials**: Machine-to-machine
- **Resource Owner Password**: Direct login (not recommended)

**Flow** (Authorization Code):
1. Client redirects user to authorization server
2. User authenticates and grants permission
3. Authorization server redirects back with code
4. Client exchanges code for access token
5. Client uses token to access protected resources

**Tokens**:
- **Access Token**: Short-lived, used to access resources
- **Refresh Token**: Long-lived, used to get new access tokens

## Further Reading

- [OAuth 2.0 - RFC 6749](https://tools.ietf.org/html/rfc6749)
- [OAuth 2.0 Simplified](https://aaronparecki.com/oauth-2-simplified/)
