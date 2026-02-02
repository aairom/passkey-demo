#!/bin/bash

# Passkey Demo - Stop Script
# This script stops the application running in detached mode

set -e

# Configuration
APP_NAME="passkey-demo"
PID_FILE="/tmp/${APP_NAME}.pid"
LOG_FILE="/tmp/${APP_NAME}.log"

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

# Check if PID file exists
if [ ! -f "$PID_FILE" ]; then
    # Try to find the process by name and port
    print_warning "PID file not found. Checking for running process..."
    
    # Find process listening on port 8080
    PID=$(lsof -ti :8080 2>/dev/null | head -n 1)
    
    if [ -z "$PID" ]; then
        print_error "Application is not running"
        exit 1
    fi
    
    print_info "Found process with PID $PID listening on port 8080"
else
    # Read PID from file
    PID=$(cat "$PID_FILE")
    
    # Check if process is running
    if ! ps -p "$PID" > /dev/null 2>&1; then
        print_warning "Process with PID $PID is not running"
        print_info "Cleaning up PID file..."
        rm -f "$PID_FILE"
        
        # Check if another process is using the port
        PORT_PID=$(lsof -ti :8080 2>/dev/null | head -n 1)
        if [ -n "$PORT_PID" ]; then
            print_warning "Found different process (PID $PORT_PID) using port 8080"
            PID=$PORT_PID
        else
            exit 0
        fi
    fi
fi

print_info "Stopping $APP_NAME (PID: $PID)..."

# Try graceful shutdown first (SIGTERM)
kill -TERM "$PID" 2>/dev/null

# Wait for process to stop (max 10 seconds)
TIMEOUT=10
COUNTER=0
while ps -p "$PID" > /dev/null 2>&1; do
    if [ $COUNTER -ge $TIMEOUT ]; then
        print_warning "Graceful shutdown timed out. Forcing shutdown..."
        kill -KILL "$PID" 2>/dev/null
        sleep 1
        break
    fi
    sleep 1
    COUNTER=$((COUNTER + 1))
done

# Verify process is stopped
if ps -p "$PID" > /dev/null 2>&1; then
    print_error "Failed to stop application"
    exit 1
else
    print_info "✓ Application stopped successfully"
    rm -f "$PID_FILE"
    
    # Show log file location
    if [ -f "$LOG_FILE" ]; then
        print_info "Log file available at: $LOG_FILE"
        print_info "Last 10 lines of log:"
        echo "----------------------------------------"
        tail -n 10 "$LOG_FILE"
        echo "----------------------------------------"
    fi
fi

# Made with Bob
