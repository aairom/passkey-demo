# API Documentation

This document describes the REST API endpoints for the Passkey Demo application.

## Base URL

```
http://localhost:8080/api
```

## Endpoints

### 1. Begin Registration

Initiates the passkey registration process for a new or existing user.

**Endpoint:** `POST /api/register/begin`

**Request Body:**
```json
{
  "username": "user@example.com"
}
```

**Response:** WebAuthn Credential Creation Options
```json
{
  "publicKey": {
    "challenge": "base64url-encoded-challenge",
    "rp": {
      "name": "Passkey Demo",
      "id": "localhost"
    },
    "user": {
      "id": "base64url-encoded-user-id",
      "name": "user@example.com",
      "displayName": "user@example.com"
    },
    "pubKeyCredParams": [
      {
        "type": "public-key",
        "alg": -7
      },
      {
        "type": "public-key",
        "alg": -257
      }
    ],
    "timeout": 60000,
    "attestation": "none",
    "authenticatorSelection": {
      "userVerification": "preferred"
    }
  }
}
```

**Status Codes:**
- `200 OK` - Registration options generated successfully
- `400 Bad Request` - Invalid request (missing username)
- `500 Internal Server Error` - Server error

---

### 2. Finish Registration

Completes the passkey registration process by verifying the credential.

**Endpoint:** `POST /api/register/finish?username={username}`

**Query Parameters:**
- `username` (required) - The username being registered

**Request Body:** WebAuthn Credential Creation Response
```json
{
  "id": "credential-id",
  "rawId": "base64url-encoded-raw-id",
  "type": "public-key",
  "response": {
    "attestationObject": "base64url-encoded-attestation",
    "clientDataJSON": "base64url-encoded-client-data"
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful"
}
```

**Status Codes:**
- `200 OK` - Registration completed successfully
- `400 Bad Request` - Invalid credential or session not found
- `404 Not Found` - User not found
- `500 Internal Server Error` - Server error

---

### 3. Begin Login

Initiates the passkey authentication process.

**Endpoint:** `POST /api/login/begin`

**Request Body:**
```json
{
  "username": "user@example.com"
}
```

**Response:** WebAuthn Credential Request Options
```json
{
  "publicKey": {
    "challenge": "base64url-encoded-challenge",
    "timeout": 60000,
    "rpId": "localhost",
    "allowCredentials": [
      {
        "type": "public-key",
        "id": "base64url-encoded-credential-id"
      }
    ],
    "userVerification": "preferred"
  }
}
```

**Status Codes:**
- `200 OK` - Login options generated successfully
- `400 Bad Request` - Invalid request or no credentials registered
- `404 Not Found` - User not found
- `500 Internal Server Error` - Server error

---

### 4. Finish Login

Completes the passkey authentication process by verifying the assertion.

**Endpoint:** `POST /api/login/finish?username={username}`

**Query Parameters:**
- `username` (required) - The username being authenticated

**Request Body:** WebAuthn Assertion Response
```json
{
  "id": "credential-id",
  "rawId": "base64url-encoded-raw-id",
  "type": "public-key",
  "response": {
    "authenticatorData": "base64url-encoded-authenticator-data",
    "clientDataJSON": "base64url-encoded-client-data",
    "signature": "base64url-encoded-signature",
    "userHandle": "base64url-encoded-user-handle"
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "username": "user@example.com",
  "userId": [1, 2, 3, 4, 5]
}
```

**Status Codes:**
- `200 OK` - Authentication successful
- `400 Bad Request` - Invalid assertion or session not found
- `404 Not Found` - User not found
- `500 Internal Server Error` - Server error

---

### 5. Health Check

Simple health check endpoint to verify the server is running.

**Endpoint:** `GET /health`

**Response:**
```
OK
```

**Status Codes:**
- `200 OK` - Server is healthy

---

## Authentication Flow

### Registration Flow

```
Client                          Server
  |                               |
  |---(1) POST /register/begin--->|
  |    {username}                 |
  |                               |
  |<--(2) Credential Options------|
  |                               |
  |---(3) Create Credential------>|
  |    (via WebAuthn API)         |
  |                               |
  |---(4) POST /register/finish-->|
  |    {credential}               |
  |                               |
  |<--(5) Success Response--------|
```

### Login Flow

```
Client                          Server
  |                               |
  |---(1) POST /login/begin------>|
  |    {username}                 |
  |                               |
  |<--(2) Assertion Options-------|
  |                               |
  |---(3) Get Credential--------->|
  |    (via WebAuthn API)         |
  |                               |
  |---(4) POST /login/finish----->|
  |    {assertion}                |
  |                               |
  |<--(5) Success + User Data-----|
```

## Error Handling

All endpoints return appropriate HTTP status codes and error messages:

**Error Response Format:**
```
Plain text error message
```

**Common Error Codes:**
- `400 Bad Request` - Invalid input or request format
- `404 Not Found` - Resource not found (user, session, etc.)
- `500 Internal Server Error` - Server-side error

## CORS

The server is configured to allow CORS from all origins for development purposes:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

**Note:** In production, restrict CORS to specific origins.

## Data Encoding

All binary data (challenges, credential IDs, etc.) is encoded using **base64url** encoding:
- No padding (`=`)
- URL-safe characters (`-` instead of `+`, `_` instead of `/`)

## Session Management

The server maintains temporary sessions during registration and login:
- Sessions are stored in-memory
- Sessions are deleted after successful completion
- Sessions timeout after 60 seconds (configurable)

## Security Considerations

### Development vs Production

This API is designed for **demonstration purposes**. For production:

1. **Use HTTPS** - Required for WebAuthn
2. **Implement Rate Limiting** - Prevent brute force attacks
3. **Add CSRF Protection** - Protect against cross-site attacks
4. **Validate Origins** - Restrict CORS to known origins
5. **Use Database** - Replace in-memory storage
6. **Session Management** - Implement proper session handling
7. **Logging** - Add comprehensive logging
8. **Error Messages** - Don't expose internal details

### WebAuthn Security

The implementation follows WebAuthn best practices:
- Challenge-response authentication
- Origin validation
- Attestation verification
- Signature verification
- Counter validation (prevents replay attacks)

## Testing with cURL

### Register a User

```bash
# Step 1: Begin registration
curl -X POST http://localhost:8080/api/register/begin \
  -H "Content-Type: application/json" \
  -d '{"username":"test@example.com"}'

# Step 2: Complete registration (requires WebAuthn credential)
# This step must be done through a browser with WebAuthn support
```

### Login

```bash
# Step 1: Begin login
curl -X POST http://localhost:8080/api/login/begin \
  -H "Content-Type: application/json" \
  -d '{"username":"test@example.com"}'

# Step 2: Complete login (requires WebAuthn credential)
# This step must be done through a browser with WebAuthn support
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Client Implementation

See the HTML/JavaScript files in the `client/` directory for complete client-side implementation examples:

- `register.html` - Registration flow
- `login.html` - Authentication flow
- Helper functions for base64url encoding/decoding
- WebAuthn API usage examples

## Resources

- [WebAuthn Specification](https://www.w3.org/TR/webauthn/)
- [go-webauthn Library](https://github.com/go-webauthn/webauthn)
- [WebAuthn Guide](https://webauthn.guide/)