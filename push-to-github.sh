#!/bin/bash

# GitHub Push Automation Script
# This script automates the process of committing and pushing changes to GitHub

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEFAULT_BRANCH="main"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to check if git is installed
check_git() {
    if ! command -v git &> /dev/null; then
        print_error "Git is not installed. Please install Git first."
        exit 1
    fi
}

# Function to check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_warning "Not a git repository. Initializing..."
        git init
        print_success "Git repository initialized"
    fi
}

# Function to check if remote exists
check_remote() {
    if ! git remote get-url origin > /dev/null 2>&1; then
        print_warning "No remote 'origin' found."
        read -p "Enter your GitHub repository URL (e.g., https://github.com/username/repo.git): " REPO_URL
        
        if [ -z "$REPO_URL" ]; then
            print_error "Repository URL cannot be empty"
            exit 1
        fi
        
        git remote add origin "$REPO_URL"
        print_success "Remote 'origin' added: $REPO_URL"
    else
        REMOTE_URL=$(git remote get-url origin)
        print_info "Remote 'origin': $REMOTE_URL"
    fi
}

# Function to get current branch
get_current_branch() {
    CURRENT_BRANCH=$(git branch --show-current)
    if [ -z "$CURRENT_BRANCH" ]; then
        CURRENT_BRANCH="$DEFAULT_BRANCH"
        print_info "No branch found, will use: $CURRENT_BRANCH"
    else
        print_info "Current branch: $CURRENT_BRANCH"
    fi
}

# Function to check for changes
check_changes() {
    if git diff-index --quiet HEAD -- 2>/dev/null; then
        print_warning "No changes to commit"
        
        read -p "Do you want to push anyway? (y/N): " PUSH_ANYWAY
        if [[ ! "$PUSH_ANYWAY" =~ ^[Yy]$ ]]; then
            print_info "Aborted by user"
            exit 0
        fi
    else
        print_info "Changes detected"
    fi
}

# Function to show status
show_status() {
    print_info "Git status:"
    git status --short
    echo ""
}

# Function to get commit message
get_commit_message() {
    if [ -n "$1" ]; then
        COMMIT_MSG="$1"
    else
        read -p "Enter commit message (or press Enter for default): " COMMIT_MSG
        
        if [ -z "$COMMIT_MSG" ]; then
            COMMIT_MSG="Update: $(date '+%Y-%m-%d %H:%M:%S')"
        fi
    fi
    
    print_info "Commit message: $COMMIT_MSG"
}

# Function to stage files
stage_files() {
    print_info "Staging files..."
    
    # Add all files except those in .gitignore
    git add -A
    
    # Show what will be committed
    print_info "Files to be committed:"
    git --no-pager diff --cached --name-status
    echo ""
}

# Function to commit changes
commit_changes() {
    print_info "Committing changes..."
    
    if git commit -m "$COMMIT_MSG"; then
        print_success "Changes committed successfully"
    else
        print_error "Commit failed"
        exit 1
    fi
}

# Function to push to GitHub
push_to_github() {
    print_info "Pushing to GitHub..."
    
    # Check if branch exists on remote
    if git ls-remote --heads origin "$CURRENT_BRANCH" | grep -q "$CURRENT_BRANCH"; then
        # Branch exists, normal push
        if git push origin "$CURRENT_BRANCH"; then
            print_success "Successfully pushed to origin/$CURRENT_BRANCH"
        else
            print_error "Push failed"
            print_info "You may need to pull first: git pull origin $CURRENT_BRANCH"
            exit 1
        fi
    else
        # Branch doesn't exist, push with --set-upstream
        print_info "Branch doesn't exist on remote, creating..."
        if git push --set-upstream origin "$CURRENT_BRANCH"; then
            print_success "Successfully pushed and set upstream to origin/$CURRENT_BRANCH"
        else
            print_error "Push failed"
            exit 1
        fi
    fi
}

# Function to show summary
show_summary() {
    echo ""
    print_success "=== Push Summary ==="
    print_info "Repository: $(git remote get-url origin)"
    print_info "Branch: $CURRENT_BRANCH"
    print_info "Commit: $COMMIT_MSG"
    print_info "Time: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
}

# Main execution
main() {
    cd "$PROJECT_DIR"
    
    print_info "=== GitHub Push Automation ==="
    print_info "Project: $PROJECT_DIR"
    echo ""
    
    # Run checks
    check_git
    check_git_repo
    check_remote
    get_current_branch
    
    # Show current status
    show_status
    
    # Check for changes
    check_changes
    
    # Get commit message from argument or prompt
    get_commit_message "$1"
    
    # Confirm before proceeding
    echo ""
    read -p "Proceed with commit and push? (Y/n): " CONFIRM
    if [[ "$CONFIRM" =~ ^[Nn]$ ]]; then
        print_info "Aborted by user"
        exit 0
    fi
    
    # Execute git operations
    stage_files
    commit_changes
    push_to_github
    
    # Show summary
    show_summary
    
    print_success "All done! 🚀"
}

# Handle script arguments
case "${1:-}" in
    -h|--help)
        echo "Usage: $0 [commit-message]"
        echo ""
        echo "Options:"
        echo "  -h, --help     Show this help message"
        echo "  -s, --status   Show git status and exit"
        echo ""
        echo "Examples:"
        echo "  $0                           # Interactive mode"
        echo "  $0 \"Fix bug in auth\"        # With commit message"
        echo "  $0 --status                  # Show status only"
        exit 0
        ;;
    -s|--status)
        cd "$PROJECT_DIR"
        git status
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac

# Made with Bob
