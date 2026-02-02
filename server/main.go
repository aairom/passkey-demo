package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/mux"
)

func main() {
	// Get configuration from environment or use defaults
	rpID := getEnv("RP_ID", "localhost")
	rpOrigin := getEnv("RP_ORIGIN", "http://localhost:8080")
	port := getEnv("PORT", "8080")

	// Initialize WebAuthn
	wconfig := &webauthn.Config{
		RPDisplayName: "Passkey Demo",
		RPID:          rpID,
		RPOrigins:     []string{rpOrigin},
		// Timeout for registration/authentication (in milliseconds)
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    60000,
				TimeoutUVD: 60000,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    60000,
				TimeoutUVD: 60000,
			},
		},
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		log.Fatalf("Failed to create WebAuthn instance: %v", err)
	}

	// Initialize storage
	storage := NewStorage()

	// Initialize handlers
	handlers := NewHandlers(webAuthn, storage)

	// Setup router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register/begin", handlers.RegisterBegin).Methods("POST", "OPTIONS")
	api.HandleFunc("/register/finish", handlers.RegisterFinish).Methods("POST", "OPTIONS")
	api.HandleFunc("/login/begin", handlers.LoginBegin).Methods("POST", "OPTIONS")
	api.HandleFunc("/login/finish", handlers.LoginFinish).Methods("POST", "OPTIONS")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Serve static files from client directory
	clientDir := "../client"
	if _, err := os.Stat(clientDir); err == nil {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(clientDir)))
		log.Printf("Serving static files from %s", clientDir)
	} else {
		log.Printf("Client directory not found at %s, skipping static file serving", clientDir)
	}

	// Apply middleware
	handler := LoggingMiddleware(EnableCORS(r))

	// Start server
	log.Printf("Starting server on port %s", port)
	log.Printf("RP ID: %s", rpID)
	log.Printf("RP Origin: %s", rpOrigin)
	log.Printf("API endpoints:")
	log.Printf("  POST /api/register/begin")
	log.Printf("  POST /api/register/finish")
	log.Printf("  POST /api/login/begin")
	log.Printf("  POST /api/login/finish")
	log.Printf("  GET  /health")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Made with Bob
