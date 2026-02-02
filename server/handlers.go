package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

// Handlers manages HTTP request handlers
type Handlers struct {
	webAuthn *webauthn.WebAuthn
	storage  *Storage
}

// NewHandlers creates a new handlers instance
func NewHandlers(webAuthn *webauthn.WebAuthn, storage *Storage) *Handlers {
	return &Handlers{
		webAuthn: webAuthn,
		storage:  storage,
	}
}

// RegisterBegin starts the registration process
func (h *Handlers) RegisterBegin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	user, exists := h.storage.GetUser(req.Username)
	if !exists {
		// Create new user
		userID := []byte(req.Username) // In production, use a proper UUID
		user = h.storage.CreateUser(req.Username, userID)
	}

	// Generate registration options
	options, session, err := h.webAuthn.BeginRegistration(user)
	if err != nil {
		log.Printf("Error beginning registration: %v", err)
		http.Error(w, "Failed to begin registration", http.StatusInternalServerError)
		return
	}

	// Store session
	h.storage.SaveSession(req.Username, session)

	// Return options to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// RegisterFinish completes the registration process
func (h *Handlers) RegisterFinish(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}

	// Parse username from query or body
	username := r.URL.Query().Get("username")
	if username == "" {
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			username = req.Username
		}
	}

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get user
	user, exists := h.storage.GetUser(username)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get session
	session, exists := h.storage.GetSession(username)
	if !exists {
		http.Error(w, "Session not found", http.StatusBadRequest)
		return
	}

	// Parse credential creation response
	credential, err := h.webAuthn.FinishRegistration(user, *session, r)
	if err != nil {
		log.Printf("Error finishing registration: %v", err)
		http.Error(w, "Failed to finish registration: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Add credential to user
	if err := h.storage.AddCredential(username, *credential); err != nil {
		log.Printf("Error adding credential: %v", err)
		http.Error(w, "Failed to save credential", http.StatusInternalServerError)
		return
	}

	// Clean up session
	h.storage.DeleteSession(username)

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Registration successful",
	})
}

// LoginBegin starts the authentication process
func (h *Handlers) LoginBegin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get user
	user, exists := h.storage.GetUser(req.Username)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has credentials
	if len(user.Credentials) == 0 {
		http.Error(w, "No credentials registered for this user", http.StatusBadRequest)
		return
	}

	// Generate authentication options
	options, session, err := h.webAuthn.BeginLogin(user)
	if err != nil {
		log.Printf("Error beginning login: %v", err)
		http.Error(w, "Failed to begin login", http.StatusInternalServerError)
		return
	}

	// Store session
	h.storage.SaveSession(req.Username, session)

	// Return options to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// LoginFinish completes the authentication process
func (h *Handlers) LoginFinish(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}

	// Parse username from query or body
	username := r.URL.Query().Get("username")
	if username == "" {
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			username = req.Username
		}
	}

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get user
	user, exists := h.storage.GetUser(username)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get session
	session, exists := h.storage.GetSession(username)
	if !exists {
		http.Error(w, "Session not found", http.StatusBadRequest)
		return
	}

	// Parse credential assertion response
	credential, err := h.webAuthn.FinishLogin(user, *session, r)
	if err != nil {
		log.Printf("Error finishing login: %v", err)
		http.Error(w, "Failed to finish login: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update credential (sign count, etc.)
	if err := h.storage.UpdateCredential(username, *credential); err != nil {
		log.Printf("Error updating credential: %v", err)
		// Don't fail the login, just log the error
	}

	// Clean up session
	h.storage.DeleteSession(username)

	// Return success with user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "Login successful",
		"username": username,
		"userId":   user.ID,
	})
}

// EnableCORS middleware to handle CORS
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Made with Bob
