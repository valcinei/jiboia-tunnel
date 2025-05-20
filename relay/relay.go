package relay

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/valcinei/jiboia-tunnel/shared"
)

type Server struct {
	upgrader websocket.Upgrader
	clients  sync.Map
}

type ClientConn struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{},
	}
}

func (s *Server) extractSubdomain(host string) string {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "default"
	}
	return parts[0]
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = "default"
	}

	// Check if the subdomain is already in use
	if _, exists := s.clients.Load(id); exists {
		http.Error(w, "Subdomain is already in use", http.StatusConflict)
		log.Printf("Connection rejected: %s is already in use", id)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Store client
	client := &ClientConn{Conn: conn}
	s.clients.Store(id, client)
	log.Printf("New client connected: %s", id)

	// Clean up on disconnect
	defer func() {
		s.clients.Delete(id)
		conn.Close()
		log.Printf("Client disconnected: %s (subdomain released)", id)
	}()

	// Loop to detect disconnection
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received from browser: Host=%s, URL=%s", r.Host, r.URL.String())
	id := s.extractSubdomain(r.Host)
	clientRaw, ok := s.clients.Load(id)
	if !ok {
		http.Error(w, "client offline", http.StatusBadGateway)
		return
	}
	client := clientRaw.(*ClientConn)
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	msg := shared.TunnelMessage{Method: r.Method, Path: r.URL.Path, Body: body}
	jsonMsg, _ := json.Marshal(msg)
	client.Conn.WriteMessage(websocket.BinaryMessage, jsonMsg)
	_, rawResp, err := client.Conn.ReadMessage()
	if err != nil {
		http.Error(w, "client error", http.StatusInternalServerError)
		return
	}
	var tr shared.TunnelResponse
	if err := json.Unmarshal(rawResp, &tr); err != nil {
		http.Error(w, "invalid response", http.StatusInternalServerError)
		return
	}
	for k, v := range tr.Headers {
		if strings.ToLower(k) == "content-length" {
			continue
		}
		w.Header().Set(k, v)
	}
	w.WriteHeader(tr.StatusCode)
	w.Write(tr.Body)
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/", s.handleProxy)
	log.Printf("Relay listening on http://%s", addr)
	return http.ListenAndServe(addr, nil)
}
