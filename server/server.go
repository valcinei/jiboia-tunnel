package server

import (
	"fmt"
	"net/http"
)

// Server represents a simple HTTP server for testing purposes.
type Server struct {
	addr string
}

// NewServer creates a new Server instance.
func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

// Start begins listening on the specified address.
func (s *Server) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Jiboia Server!")
	})
	return http.ListenAndServe(s.addr, nil)
}
