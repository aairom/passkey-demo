# Setup Guide

This guide will help you set up and run the Passkey Demo application in various environments.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Running Options](#running-options)
4. [Docker Deployment](#docker-deployment)
5. [Kubernetes Deployment](#kubernetes-deployment)
6. [Configuration](#configuration)
7. [Testing](#testing)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required

1. **Go 1.21 or higher** (for direct execution)
   - Check your version: `go version`
   - Download from: https://golang.org/dl/

2. **Modern Web Browser**
   - Chrome 67+
   - Firefox 60+
   - Safari 13+
   - Edge 18+

3. **Device with Passkey Support**
   - Biometric sensor (fingerprint, Face ID, etc.)
   - Security key (YubiKey, etc.)
   - Device PIN/password

### Optional (for containerization)

4. **Docker** (for container deployment)
   - Check: `docker --version`
   - Download from: https://www.docker.com/

5. **Kubernetes** (for K8s deployment)
   - kubectl: `kubectl version`
   - Cluster access configured

---

## Installation

### 1. Navigate to Project Directory

```bash
cd /Users/alainairom/Dev/passkey-demo
```

### 2. Install Dependencies

```bash
cd server
go mod download
go mod verify
```

---

## Running Options

### Option 1: Using Start/Stop Scripts (Recommended)

The easiest way to run the application in detached mode:

```bash
# Start server in background
./start.sh

# View logs
tail -f /tmp/passkey-demo.log

# Stop server
./stop.sh
```

**Features:**
- Runs in background
- Automatic PID management
- Log file at `/tmp/passkey-demo.log`
- Graceful shutdown

### Option 2: Direct Execution

For development and debugging:

```bash
cd server
go run .
```

Or build first:

```bash
cd server
go build -o passkey-demo
./passkey-demo
```

**Features:**
- Immediate feedback
- Easy debugging
- Foreground process

### Option 3: Docker

See [Docker Deployment](#docker-deployment) section below.

### Option 4: Kubernetes

See [Kubernetes Deployment](#kubernetes-deployment) section below.

---

## Docker Deployment

### Using Docker Compose (Recommended)

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Using Docker Directly

```bash
# Build image
docker build -t passkey-demo:latest .

# Run container
docker run -d \
  --name passkey-demo \
  -p 8080:8080 \
  -e RP_ID=localhost \
  -e RP_ORIGIN=http://localhost:8080 \
  passkey-demo:latest

# View logs
docker logs -f passkey-demo

# Stop and remove
docker stop passkey-demo
docker rm passkey-demo
```

### Docker Configuration

The [`Dockerfile`](Dockerfile:1) uses multi-stage build:
- **Builder stage**: Compiles Go application
- **Final stage**: Minimal Alpine image (~20MB)
- **Security**: Non-root user, dropped capabilities
- **Health check**: Built-in health monitoring

---

## Kubernetes Deployment

### Quick Deploy

```bash
# Deploy all resources
kubectl apply -f k8s/

# Check status
kubectl get all -n passkey-demo

# View logs
kubectl logs -n passkey-demo -l app=passkey-demo -f
```

### Step-by-Step Deploy

```bash
# 1. Create namespace
kubectl apply -f k8s/namespace.yaml

# 2. Create configuration
kubectl apply -f k8s/configmap.yaml

# 3. Deploy application
kubectl apply -f k8s/deployment.yaml

# 4. Create service
kubectl apply -f k8s/service.yaml

# 5. Setup ingress
kubectl apply -f k8s/ingress.yaml

# 6. Enable autoscaling
kubectl apply -f k8s/hpa.yaml
```

### Kubernetes Resources

- **Namespace**: `passkey-demo`
- **Deployment**: 3 replicas with health checks
- **Service**: ClusterIP on port 80
- **Ingress**: HTTPS with TLS
- **HPA**: Auto-scale 3-10 pods based on CPU/memory

See [k8s/README.md](k8s/README.md) for detailed Kubernetes documentation.

---

You should see output like:

```
Starting server on port 8080
RP ID: localhost
RP Origin: http://localhost:8080
API endpoints:
  POST /api/register/begin
  POST /api/register/finish
  POST /api/login/begin
  POST /api/login/finish
  GET  /health
Serving static files from ../client
```

### Access the Application

Open your browser and navigate to:

```
http://localhost:8080
```

Or directly open the HTML files:

```
passkey-demo/client/index.html
```

## Configuration

You can configure the server using environment variables:

```bash
# Set the Relying Party ID (default: localhost)
export RP_ID=localhost

# Set the Relying Party Origin (default: http://localhost:8080)
export RP_ORIGIN=http://localhost:8080

# Set the server port (default: 8080)
export PORT=8080

# Run the server
go run .
```

## Testing the Application

### 1. Register a New User

1. Click "Register New Account"
2. Enter a username (e.g., `user@example.com`)
3. Click "Register with Passkey"
4. Follow your device prompts to create a passkey
5. You should see "Registration successful!"

### 2. Login with Passkey

1. Click "Login with Passkey"
2. Enter the same username
3. Click "Login with Passkey"
4. Follow your device prompts to authenticate
5. You should see "Login Successful!"

## Troubleshooting

### "WebAuthn not supported"

**Solution**: Ensure you're using a modern browser and accessing via `localhost` or `https://`

### "Failed to begin registration"

**Possible causes**:
- Server not running
- CORS issues
- Check browser console for errors

**Solution**: 
- Verify server is running on port 8080
- Check server logs for errors
- Open browser developer tools (F12) and check console

### "User not found" during login

**Solution**: Make sure you registered the user first

### "No credentials registered"

**Solution**: Complete the registration process before attempting to login

### Port already in use

**Solution**: Change the port:

```bash
PORT=3000 go run .
```

Then update the API_BASE in the HTML files to match.

## Development Tips

### Enable Verbose Logging

The server logs all requests. Check the terminal where the server is running for detailed logs.

### Testing with Different Users

Simply use different usernames to create multiple accounts.

### Resetting Data

Since the application uses in-memory storage, simply restart the server to clear all data.

## Security Notes for Production

This is a **demonstration application**. For production use:

1. **Use HTTPS**: Required for WebAuthn in production
2. **Implement Database**: Replace in-memory storage with PostgreSQL, MySQL, etc.
3. **Add Session Management**: Implement proper session handling with secure cookies
4. **Rate Limiting**: Prevent brute force attacks
5. **CSRF Protection**: Add CSRF tokens
6. **Input Validation**: Validate and sanitize all inputs
7. **Error Handling**: Don't expose internal errors to clients
8. **Logging**: Implement proper logging and monitoring
9. **Backup**: Regular backups of credential data

## Next Steps

- Explore the code in `server/` directory
- Customize the UI in `client/` directory
- Add database integration
- Implement session management
- Deploy to production with HTTPS

## Resources

- [WebAuthn Guide](https://webauthn.guide/)
- [go-webauthn Documentation](https://github.com/go-webauthn/webauthn)
- [FIDO Alliance](https://fidoalliance.org/)
- [W3C WebAuthn Spec](https://www.w3.org/TR/webauthn/)