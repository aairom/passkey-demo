# Passkey Authentication Demo

A comprehensive demonstration of passkey (WebAuthn) implementation with a Go server and HTML/JavaScript client.

## Features

- **Passwordless Authentication**: Uses WebAuthn/FIDO2 passkeys
- **Go Backend**: RESTful API server with WebAuthn support
- **Simple Client**: HTML/JavaScript frontend
- **In-Memory Storage**: User and credential storage (can be replaced with a database)
- **CORS Support**: Configured for local development

## Architecture

```
passkey-demo/
├── server/              # Go backend server
│   ├── main.go          # Server entry point
│   ├── handlers.go      # HTTP handlers
│   ├── storage.go       # In-memory storage
│   ├── go.mod           # Go dependencies
│   └── go.sum           # Dependency checksums
├── client/              # HTML/JavaScript client
│   ├── index.html       # Main page
│   ├── register.html    # Registration page
│   └── login.html       # Login page
├── k8s/                 # Kubernetes manifests
│   ├── namespace.yaml   # Namespace definition
│   ├── configmap.yaml   # Configuration
│   ├── deployment.yaml  # Deployment spec
│   ├── service.yaml     # Service definition
│   ├── ingress.yaml     # Ingress rules
│   ├── hpa.yaml         # Horizontal Pod Autoscaler
│   └── README.md        # K8s deployment guide
├── Dockerfile           # Container image definition
├── .dockerignore        # Docker ignore rules
├── docker-compose.yml   # Docker Compose configuration
├── start.sh             # Start script (detached mode)
├── stop.sh              # Stop script
├── ARCHITECTURE.md      # Architecture diagrams (Mermaid)
├── QUICKSTART.md        # Quick start guide
├── SETUP.md             # Detailed setup instructions
├── API.md               # API documentation
└── PROJECT_OVERVIEW.md  # Comprehensive overview
```

For detailed architecture diagrams, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Prerequisites

- Go 1.21 or higher
   go run .
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

### Option 2: Using Start/Stop Scripts (Detached Mode)

```bash
# Start the server in background
./start.sh

# View logs
tail -f /tmp/passkey-demo.log

# Stop the server
./stop.sh
```

### Option 3: Using Docker

```bash
# Build and run with Docker
docker build -t passkey-demo .
docker run -p 8080:8080 passkey-demo

# Or use Docker Compose
docker-compose up -d
docker-compose logs -f
docker-compose down
```

### Option 4: Deploy to Kubernetes

```bash
# Build and push image
docker build -t your-registry/passkey-demo:latest .
docker push your-registry/passkey-demo:latest

# Deploy to Kubernetes
kubectl apply -f k8s/

# Check status
kubectl get all -n passkey-demo
```

See [k8s/README.md](k8s/README.md) for detailed Kubernetes deployment instructions.

## Quick Links

- 📖 [Quick Start Guide](QUICKSTART.md) - Get started in 5 minutes
- 🏗️ [Architecture Diagrams](ARCHITECTURE.md) - Mermaid flowcharts and diagrams
- 🔧 [Setup Guide](SETUP.md) - Detailed installation and configuration
- 📡 [API Documentation](API.md) - Complete API reference
- 📋 [Project Overview](PROJECT_OVERVIEW.md) - Comprehensive technical overview
- ☸️ [Kubernetes Guide](k8s/README.md) - K8s deployment instructions
- 🔀 [GitHub Workflow](GITHUB.md) - Git and GitHub usage guide

## Deployment Options

| Method | Use Case | Command |
|--------|----------|---------|
| **Direct** | Development | `cd server && go run .` |
| **Scripts** | Background process | `./start.sh` |
| **Docker** | Containerized | `docker-compose up` |
| **Kubernetes** | Production | `kubectl apply -f k8s/` |

## Project Location

This project is located at: `/Users/alainairom/Dev/passkey-demo`

## Usage

### Registration Flow

1. Go to the registration page
2. Enter a username
3. Click "Register with Passkey"
4. Follow your browser/device prompts to create a passkey
5. Your passkey is now registered!

### Authentication Flow

1. Go to the login page
2. Enter your username
3. Click "Login with Passkey"
4. Follow your browser/device prompts to authenticate
5. You're logged in!

## API Endpoints

### POST `/api/register/begin`
Initiates passkey registration
- Request: `{"username": "user@example.com"}`
- Response: WebAuthn credential creation options

### POST `/api/register/finish`
Completes passkey registration
- Request: WebAuthn credential + username
- Response: Success confirmation

### POST `/api/login/begin`
Initiates passkey authentication
- Request: `{"username": "user@example.com"}`
- Response: WebAuthn credential request options

### POST `/api/login/finish`
Completes passkey authentication
- Request: WebAuthn assertion
- Response: Authentication token/session

See [API.md](API.md) for complete API documentation.

## Development

### Building

```bash
cd server
go build -o passkey-demo
./passkey-demo
```

### Testing

```bash
# Run tests

## Version Control & GitHub

### Push to GitHub

Use the automated push script:

```bash
# Interactive mode (prompts for commit message)
./push-to-github.sh

# With commit message
./push-to-github.sh "Add new feature"

# Show status only
./push-to-github.sh --status

# Show help
./push-to-github.sh --help
```

The script will:
1. Check if git is initialized
2. Check/add remote repository
3. Show current changes
4. Stage all files (respecting .gitignore)
5. Commit with your message
6. Push to GitHub

### First Time Setup

If this is your first push:

```bash
# The script will prompt for your GitHub repository URL
./push-to-github.sh

# Or set it up manually first
git remote add origin https://github.com/username/passkey-demo.git
```

### .gitignore Rules

The project ignores:
- Build artifacts (*.exe, *.out, binaries)
- IDE files (.vscode/, .idea/)
- OS files (.DS_Store, Thumbs.db)
- Log files (*.log)
- **Folders starting with underscore** (`_*/`, `**/_*/`)

This means any folder like `_temp/`, `_backup/`, `_private/` will not be pushed to GitHub.

cd server
go test ./...

# With coverage
go test -cover ./...
```

### Docker Build

```bash
# Build image
docker build -t passkey-demo:latest .

# Run container
docker run -p 8080:8080 \
  -e RP_ID=localhost \
  -e RP_ORIGIN=http://localhost:8080 \
  passkey-demo:latest
```

## Configuration

Environment variables:

- `RP_ID` - Relying Party ID (default: `localhost`)
- `RP_ORIGIN` - Relying Party Origin (default: `http://localhost:8080`)
- `PORT` - Server port (default: `8080`)

Example:

```bash
export RP_ID=myapp.com
export RP_ORIGIN=https://myapp.com
export PORT=3000
go run .
- Modern web browser with WebAuthn support (Chrome, Firefox, Safari, Edge)
- HTTPS or localhost (required for WebAuthn)

## Installation

1. Clone or download this project
2. Navigate to the server directory:
   ```bash
   cd passkey-demo/server
   ```
3. Install Go dependencies:
   ```bash
   go mod download
   ```

## Running the Application

### Option 1: Direct Execution

1. Start the Go server:
   ```bash
   cd server
   go run .
   ```
   The server will start on `http://localhost:8080`

2. Open the client in your browser:
   - Open `client/index.html` in your browser, or
   - Navigate to `http://localhost:8080` (the server serves static files)

## Usage

### Registration Flow

1. Go to the registration page
2. Enter a username
3. Click "Register with Passkey"
4. Follow your browser/device prompts to create a passkey
5. Your passkey is now registered!

### Authentication Flow

1. Go to the login page
2. Enter your username
3. Click "Login with Passkey"
4. Follow your browser/device prompts to authenticate
5. You're logged in!

## API Endpoints

### POST `/api/register/begin`
Initiates passkey registration
- Request: `{"username": "user@example.com"}`
- Response: WebAuthn credential creation options

### POST `/api/register/finish`
Completes passkey registration
- Request: WebAuthn credential + username
- Response: Success confirmation

### POST `/api/login/begin`
Initiates passkey authentication
- Request: `{"username": "user@example.com"}`
- Response: WebAuthn credential request options

### POST `/api/login/finish`
Completes passkey authentication
- Request: WebAuthn assertion
- Response: Authentication token/session

## Security Notes

- This is a **demonstration** application
- Uses in-memory storage (data lost on restart)
- For production:
  - Use a proper database
  - Implement session management
  - Add rate limiting
  - Use HTTPS
  - Add CSRF protection
  - Implement proper error handling

## Technology Stack

- **Backend**: Go with `go-webauthn/webauthn` library
- **Frontend**: Vanilla JavaScript with WebAuthn API
- **Protocol**: WebAuthn/FIDO2

## Browser Support

Passkeys are supported in:
- Chrome 67+
- Firefox 60+
- Safari 13+
- Edge 18+

## Troubleshooting

**"WebAuthn not supported"**: Ensure you're using HTTPS or localhost

**Registration fails**: Check browser console for errors, ensure your device supports passkeys

**Authentication fails**: Verify the username matches a registered user

## License

MIT License - Feel free to use for learning and development