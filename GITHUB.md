# GitHub Workflow Guide

This guide explains how to use Git and GitHub with the Passkey Demo project.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Automated Push Script](#automated-push-script)
3. [Manual Git Commands](#manual-git-commands)
4. [.gitignore Rules](#gitignore-rules)
5. [Best Practices](#best-practices)
6. [Troubleshooting](#troubleshooting)

---

## Quick Start

### First Time Setup

1. **Create a GitHub repository**
   - Go to https://github.com/new
   - Create a new repository (e.g., `passkey-demo`)
   - Don't initialize with README (we already have one)

2. **Use the automated script**
   ```bash
   cd /Users/alainairom/Dev/passkey-demo
   ./push-to-github.sh
   ```

3. **Enter your repository URL when prompted**
   ```
   https://github.com/yourusername/passkey-demo.git
   ```

That's it! Your code is now on GitHub.

---

## Automated Push Script

The [`push-to-github.sh`](push-to-github.sh:1) script automates the entire Git workflow.

### Basic Usage

```bash
# Interactive mode (will prompt for commit message)
./push-to-github.sh

# With commit message
./push-to-github.sh "Add authentication feature"

# Show git status
./push-to-github.sh --status

# Show help
./push-to-github.sh --help
```

### What the Script Does

1. ✅ Checks if Git is installed
2. ✅ Initializes Git repository (if needed)
3. ✅ Checks/adds remote repository
4. ✅ Shows current changes
5. ✅ Stages all files (respecting .gitignore)
6. ✅ Commits with your message
7. ✅ Pushes to GitHub
8. ✅ Shows summary

### Script Features

- **Colored output** for better readability
- **Interactive prompts** for safety
- **Automatic remote setup** on first use
- **Branch detection** and creation
- **Error handling** with helpful messages
- **Confirmation prompts** before pushing

### Example Session

```bash
$ ./push-to-github.sh "Update documentation"

[INFO] === GitHub Push Automation ===
[INFO] Project: /Users/alainairom/Dev/passkey-demo

[INFO] Remote 'origin': https://github.com/username/passkey-demo.git
[INFO] Current branch: main

[INFO] Git status:
 M README.md
 M SETUP.md
?? GITHUB.md

[INFO] Changes detected
[INFO] Commit message: Update documentation

Proceed with commit and push? (Y/n): y

[INFO] Staging files...
[INFO] Files to be committed:
M       README.md
M       SETUP.md
A       GITHUB.md

[INFO] Committing changes...
[SUCCESS] Changes committed successfully

[INFO] Pushing to GitHub...
[SUCCESS] Successfully pushed to origin/main

[SUCCESS] === Push Summary ===
[INFO] Repository: https://github.com/username/passkey-demo.git
[INFO] Branch: main
[INFO] Commit: Update documentation
[INFO] Time: 2026-02-02 10:30:45

[SUCCESS] All done! 🚀
```

---

## Manual Git Commands

If you prefer manual control:

### Initial Setup

```bash
cd /Users/alainairom/Dev/passkey-demo

# Initialize Git (if not already done)
git init

# Add remote repository
git remote add origin https://github.com/username/passkey-demo.git

# Check remote
git remote -v
```

### Daily Workflow

```bash
# Check status
git status

# Stage all changes
git add .

# Or stage specific files
git add README.md SETUP.md

# Commit changes
git commit -m "Your commit message"

# Push to GitHub
git push origin main

# Or set upstream and push
git push -u origin main
```

### Viewing Changes

```bash
# See what changed
git diff

# See staged changes
git diff --cached

# View commit history
git log --oneline

# View specific file history
git log --follow README.md
```

### Branching

```bash
# Create new branch
git checkout -b feature-name

# Switch branches
git checkout main

# List branches
git branch -a

# Delete branch
git branch -d feature-name
```

---

## .gitignore Rules

The [`.gitignore`](.gitignore:1) file prevents certain files from being tracked.

### Current Rules

```gitignore
# Build artifacts
*.exe, *.dll, *.so, *.dylib
*.test, *.out

# IDE files
.vscode/, .idea/
*.swp, *.swo, *~

# OS files
.DS_Store, Thumbs.db

# Build output
server/passkey-demo
server/server

# Log files
*.log

# Folders starting with underscore
_*/
**/_*/
```

### Underscore Folder Rule

**Any folder starting with underscore will NOT be pushed to GitHub.**

Examples of ignored folders:
- `_temp/`
- `_backup/`
- `_private/`
- `_drafts/`
- `_local/`
- `server/_cache/`
- `client/_assets/`

This is useful for:
- Local development files
- Temporary backups
- Private notes
- Work-in-progress features
- Local configuration

### Testing .gitignore

```bash
# Check if a file would be ignored
git check-ignore -v _temp/file.txt

# List all ignored files
git status --ignored

# Force add an ignored file (not recommended)
git add -f _temp/file.txt
```

---

## Best Practices

### Commit Messages

**Good commit messages:**
```
✅ Add user authentication feature
✅ Fix login button styling
✅ Update API documentation
✅ Refactor storage layer
```

**Bad commit messages:**
```
❌ Update
❌ Fix stuff
❌ Changes
❌ WIP
```

### Commit Frequency

- **Commit often**: Small, focused commits
- **One feature per commit**: Easy to review and revert
- **Test before committing**: Ensure code works
- **Write clear messages**: Explain what and why

### Branch Strategy

```bash
# Main branch: stable, production-ready
main

# Feature branches: new features
feature/user-auth
feature/api-v2

# Fix branches: bug fixes
fix/login-error
fix/memory-leak

# Docs branches: documentation
docs/api-guide
docs/setup-instructions
```

### Before Pushing

1. **Review changes**: `git diff`
2. **Test locally**: Run the application
3. **Check status**: `git status`
4. **Read commit message**: Make sure it's clear
5. **Push**: Use script or manual command

---

## Troubleshooting

### "Permission denied (publickey)"

**Solution**: Set up SSH keys or use HTTPS with token

```bash
# Use HTTPS instead
git remote set-url origin https://github.com/username/passkey-demo.git

# Or set up SSH keys
ssh-keygen -t ed25519 -C "your_email@example.com"
# Add key to GitHub: Settings > SSH and GPG keys
```

### "Failed to push"

**Solution**: Pull first, then push

```bash
git pull origin main --rebase
git push origin main
```

### "Merge conflict"

**Solution**: Resolve conflicts manually

```bash
# Pull changes
git pull origin main

# Edit conflicted files (look for <<<<<<< markers)
# After resolving:
git add .
git commit -m "Resolve merge conflicts"
git push origin main
```

### "Detached HEAD state"

**Solution**: Create a branch or checkout existing one

```bash
# Create new branch from current state
git checkout -b recovery-branch

# Or go back to main
git checkout main
```

### "Large files"

**Solution**: Use Git LFS or remove from history

```bash
# Install Git LFS
brew install git-lfs  # macOS
git lfs install

# Track large files
git lfs track "*.psd"
git lfs track "*.mp4"

# Or remove from history
git filter-branch --tree-filter 'rm -f large-file.zip' HEAD
```

### "Accidentally committed sensitive data"

**Solution**: Remove from history immediately

```bash
# Remove file from all commits
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch path/to/sensitive-file" \
  --prune-empty --tag-name-filter cat -- --all

# Force push (WARNING: rewrites history)
git push origin --force --all

# Then change any exposed credentials!
```

---

## GitHub Features

### Pull Requests

1. Create feature branch
2. Push to GitHub
3. Open Pull Request on GitHub
4. Review and merge

### Issues

Track bugs and features:
- Go to repository > Issues
- Create new issue
- Assign, label, and track

### Actions (CI/CD)

Automate testing and deployment:
- Create `.github/workflows/` directory
- Add workflow YAML files
- Automatic builds on push

### Releases

Create versioned releases:
- Go to repository > Releases
- Create new release
- Tag version (e.g., v1.0.0)
- Add release notes

---

## Additional Resources

- [Git Documentation](https://git-scm.com/doc)
- [GitHub Guides](https://guides.github.com/)
- [Git Cheat Sheet](https://education.github.com/git-cheat-sheet-education.pdf)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

## Quick Reference

```bash
# Status and info
git status                    # Show working tree status
git log --oneline            # Show commit history
git diff                     # Show changes

# Staging and committing
git add .                    # Stage all changes
git commit -m "message"      # Commit with message
git commit --amend           # Modify last commit

# Pushing and pulling
git push origin main         # Push to remote
git pull origin main         # Pull from remote
git fetch origin             # Fetch without merging

# Branching
git branch                   # List branches
git checkout -b name         # Create and switch branch
git merge branch-name        # Merge branch

# Undoing changes
git reset HEAD file          # Unstage file
git checkout -- file         # Discard changes
git revert commit-hash       # Revert commit

# Remote management
git remote -v                # List remotes
git remote add origin url    # Add remote
git remote set-url origin url # Change remote URL
```

---

**Pro Tip**: Use the automated script [`push-to-github.sh`](push-to-github.sh:1) for hassle-free pushes! 🚀