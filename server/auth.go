package server

import (
	"encoding/json"
	"net/http"
)

// AuthHandler handles authentication requests.
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Dummy authentication logic
	if credentials.Username == "admin" && credentials.Password == "password" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authentication successful"))
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// SetupAuthRoutes sets up the authentication routes.
func SetupAuthRoutes() {
	http.HandleFunc("/auth", AuthHandler)
}
