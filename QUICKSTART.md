# Quick Start Guide

Get up and running with the Passkey Demo in 5 minutes!

## Prerequisites

- Go 1.21+ installed
- Modern web browser (Chrome, Firefox, Safari, or Edge)

## Choose Your Method

### Method 1: Quick Start (Recommended)

```bash
cd /Users/alainairom/Dev/passkey-demo
./start.sh
```

That's it! The server is now running in the background.

### Method 2: Direct Execution

```bash
cd /Users/alainairom/Dev/passkey-demo/server
go mod download
go run .
```

### Method 3: Docker

```bash
cd /Users/alainairom/Dev/passkey-demo
docker-compose up -d
```

### Method 4: Kubernetes

```bash
cd /Users/alainairom/Dev/passkey-demo
kubectl apply -f k8s/
```

## Step 2: Verify Server is Running

Check the server status:

```bash
# If using start.sh
tail -f /tmp/passkey-demo.log

# If using Docker
docker-compose logs -f

# If using Kubernetes
kubectl logs -n passkey-demo -l app=passkey-demo -f
```

## Step 3: Access the Application

Open your browser and navigate to:
```
http://localhost:8080
```

Or directly open:
```
passkey-demo/client/index.html
```

## Step 4: Register a User

1. Click **"Register New Account"**
2. Enter a username (e.g., `alice@example.com`)
3. Click **"Register with Passkey"**
4. Follow your device prompts:
   - **Touch ID** on Mac
   - **Windows Hello** on Windows
   - **Fingerprint** on mobile
   - **Security Key** if you have one

✅ Registration complete!

## Step 5: Login

1. Click **"Login with Passkey"**
2. Enter the same username
3. Click **"Login with Passkey"**
4. Authenticate with your device

✅ You're logged in!

## What Just Happened?

You just experienced **passwordless authentication**:

1. **No password to remember** - Your device stores the passkey
2. **Phishing-resistant** - Passkeys are bound to the website
3. **Biometric security** - Uses your device's built-in security
4. **Fast & convenient** - One touch to authenticate

## Troubleshooting

### "WebAuthn not supported"
- Use a modern browser
- Access via `localhost` or `https://`

### "User not found"
- Register first before logging in

### Port 8080 in use
```bash
PORT=3000 go run .
```
Then update `API_BASE` in HTML files to `http://localhost:3000/api`

## Next Steps

- Read [SETUP.md](SETUP.md) for detailed setup
- Check [API.md](API.md) for API documentation
- Explore the code in `server/` and `client/`
- Try registering multiple users
- Test on different devices

## Architecture Overview

```
┌─────────────┐         ┌─────────────┐
│   Browser   │         │  Go Server  │
│             │         │             │
│  WebAuthn   │◄───────►│  WebAuthn   │
│    API      │  HTTPS  │   Library   │
│             │         │             │

## Stopping the Application

### If using start.sh
```bash
./stop.sh
```

### If using Docker
```bash
docker-compose down
```

### If using Kubernetes
```bash
kubectl delete -f k8s/
```

### If using direct execution
Press `Ctrl+C` in the terminal

## Deployment Comparison

| Method | Startup Time | Best For | Persistence |
|--------|--------------|----------|-------------|
| **start.sh** | Fast | Development | Background process |
| **Direct** | Fast | Development | Foreground only |
| **Docker** | Medium | Testing/Staging | Container lifecycle |
| **Kubernetes** | Slow | Production | Cluster managed |

## Additional Resources

- 📖 [Full Documentation](README.md)
- 🏗️ [Architecture Diagrams](ARCHITECTURE.md)
- 🔧 [Detailed Setup](SETUP.md)
- 📡 [API Reference](API.md)
- ☸️ [Kubernetes Guide](k8s/README.md)

## Project Location

```
/Users/alainairom/Dev/passkey-demo
```
│  Passkey    │         │  Storage    │
│  Storage    │         │  (Memory)   │
└─────────────┘         └─────────────┘
```

## Key Features

✨ **Passwordless** - No passwords to manage
🔒 **Secure** - FIDO2/WebAuthn standard
🚀 **Fast** - One-touch authentication
📱 **Cross-platform** - Works on all devices
🎯 **Simple** - Easy to implement

## Learn More

- [WebAuthn Guide](https://webauthn.guide/)
- [FIDO Alliance](https://fidoalliance.org/)
- [Passkeys.dev](https://passkeys.dev/)

---

**Ready to build?** Check out the code and customize it for your needs!