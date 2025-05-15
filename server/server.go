package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/valcinei/jiboia-tunnel/shared"
)

// Server represents a REST server for managing tunnels.
type Server struct {
	addr  string
	store *shared.InMemoryStore
}

// NewServer creates a new Server instance.
func NewServer(addr string) *Server {
	return &Server{addr: addr, store: shared.NewInMemoryStore()}
}

// Start begins listening on the specified address.
func (s *Server) Start() error {
	http.HandleFunc("/tunnels", s.handleTunnels)
	http.HandleFunc("/tunnels/", s.handleTunnel)
	return http.ListenAndServe(s.addr, nil)
}

// handleTunnels handles requests to the /tunnels endpoint.
func (s *Server) handleTunnels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tunnels := s.store.ListTunnels()
		json.NewEncoder(w).Encode(tunnels)
	case http.MethodPost:
		var tunnel struct {
			Name     string `json:"name"`
			LocalURL string `json:"local_url"`
			RelayURL string `json:"relay_url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&tunnel); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		s.store.AddTunnel(tunnel.Name, tunnel.LocalURL, tunnel.RelayURL)
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTunnel handles requests to the /tunnels/{id} endpoint.
func (s *Server) handleTunnel(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tunnels/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tunnel ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		tunnel, exists := s.store.GetTunnel(id)
		if !exists {
			http.Error(w, "Tunnel not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(tunnel)
	case http.MethodDelete:
		if !s.store.DeleteTunnel(id) {
			http.Error(w, "Tunnel not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
