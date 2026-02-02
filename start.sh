#!/bin/bash

# Passkey Demo - Start Script
# This script starts the application in detached mode

set -e

# Configuration
APP_NAME="passkey-demo"
PID_FILE="/tmp/${APP_NAME}.pid"
LOG_FILE="/tmp/${APP_NAME}.log"
SERVER_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/server" && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        print_error "Application is already running with PID $PID"
        print_info "Use './stop.sh' to stop it first"
        exit 1
    else
        print_warning "Stale PID file found. Removing..."
        rm -f "$PID_FILE"
    fi
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"
if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    print_error "Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or higher."
    exit 1
fi

print_info "Starting $APP_NAME..."

# Navigate to server directory
cd "$SERVER_DIR"

# Download dependencies if needed
if [ ! -f "go.sum" ]; then
    print_info "Downloading dependencies..."
    go mod download
fi

# Build the application
print_info "Building application..."
go build -o "${APP_NAME}" .

if [ $? -ne 0 ]; then
    print_error "Build failed"
    exit 1
fi

# Start the application in background
print_info "Starting server in detached mode..."
nohup ./"${APP_NAME}" > "$LOG_FILE" 2>&1 &
APP_PID=$!

# Save PID
echo $APP_PID > "$PID_FILE"

# Wait a moment and check if it's still running
sleep 2
if ps -p $APP_PID > /dev/null 2>&1; then
    print_info "✓ Application started successfully!"
    print_info "  PID: $APP_PID"
    print_info "  Log file: $LOG_FILE"
    print_info "  URL: http://localhost:8080"
    print_info ""
    print_info "To view logs: tail -f $LOG_FILE"
    print_info "To stop: ./stop.sh"
else
    print_error "Application failed to start. Check logs at $LOG_FILE"
    rm -f "$PID_FILE"
    exit 1
fi

# Made with Bob
