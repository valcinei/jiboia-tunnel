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
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("erro upgrade:", err)
		return
	}
	log.Printf("Novo cliente: %s", id)
	s.clients.Store(id, &ClientConn{Conn: conn})
	defer func() {
		s.clients.Delete(id)
		conn.Close()
		log.Printf("Cliente %s desconectado", id)
	}()
	select {}
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requisição recebida do navegador: Host=%s, URL=%s", r.Host, r.URL.String())
	id := s.extractSubdomain(r.Host)
	clientRaw, ok := s.clients.Load(id)
	if !ok {
		http.Error(w, "cliente offline", http.StatusBadGateway)
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
		http.Error(w, "erro do cliente", http.StatusInternalServerError)
		return
	}
	var tr shared.TunnelResponse
	if err := json.Unmarshal(rawResp, &tr); err != nil {
		http.Error(w, "resposta inválida", http.StatusInternalServerError)
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
	log.Printf("Relay ouvindo em http://%s", addr)
	return http.ListenAndServe(addr, nil)
} 