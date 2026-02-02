package main

import (
	"sync"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// User represents a user in the system
type User struct {
	ID          []byte
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	return u.ID
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.Name
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnCredentials returns the user's credentials
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

// WebAuthnIcon returns the user's icon URL (optional)
func (u User) WebAuthnIcon() string {
	return ""
}

// Storage provides in-memory storage for users and sessions
type Storage struct {
	users    map[string]*User // username -> User
	sessions map[string]*webauthn.SessionData
	mu       sync.RWMutex
}

// NewStorage creates a new storage instance
func NewStorage() *Storage {
	return &Storage{
		users:    make(map[string]*User),
		sessions: make(map[string]*webauthn.SessionData),
	}
}

// GetUser retrieves a user by username
func (s *Storage) GetUser(username string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[username]
	return user, exists
}

// CreateUser creates a new user
func (s *Storage) CreateUser(username string, userID []byte) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &User{
		ID:          userID,
		Name:        username,
		DisplayName: username,
		Credentials: []webauthn.Credential{},
	}
	s.users[username] = user
	return user
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(user *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.Name] = user
}

// SaveSession stores a WebAuthn session
func (s *Storage) SaveSession(username string, session *webauthn.SessionData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[username] = session
}

// GetSession retrieves a WebAuthn session
func (s *Storage) GetSession(username string) (*webauthn.SessionData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, exists := s.sessions[username]
	return session, exists
}

// DeleteSession removes a WebAuthn session
func (s *Storage) DeleteSession(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, username)
}

// AddCredential adds a credential to a user
func (s *Storage) AddCredential(username string, credential webauthn.Credential) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[username]
	if !exists {
		return protocol.ErrBadRequest.WithDetails("User not found")
	}

	user.Credentials = append(user.Credentials, credential)
	return nil
}

// UpdateCredential updates a user's credential (e.g., sign count)
func (s *Storage) UpdateCredential(username string, credential webauthn.Credential) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[username]
	if !exists {
		return protocol.ErrBadRequest.WithDetails("User not found")
	}

	for i, cred := range user.Credentials {
		if string(cred.ID) == string(credential.ID) {
			user.Credentials[i] = credential
			return nil
		}
	}

	return protocol.ErrBadRequest.WithDetails("Credential not found")
}

// Made with Bob
