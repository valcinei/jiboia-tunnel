package main

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

var upgrader = websocket.Upgrader{}

type ClientConn struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

var clients sync.Map // id -> *ClientConn

func extractSubdomain(host string) string {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "default"
	}
	return parts[0]
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = "default"
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("erro upgrade:", err)
		return
	}
	log.Printf("Novo cliente: %s", id)
	clients.Store(id, &ClientConn{Conn: conn})

	defer func() {
		clients.Delete(id)
		conn.Close()
		log.Printf("Cliente %s desconectado", id)
	}()

	select {} // bloqueia sem interferir na leitura do proxy
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requisição recebida do navegador: Host=%s, URL=%s", r.Host, r.URL.String())

	id := extractSubdomain(r.Host)
	clientRaw, ok := clients.Load(id)
	if !ok {
		http.Error(w, "cliente offline", http.StatusBadGateway)
		return
	}
	client := clientRaw.(*ClientConn)
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	msg := shared.TunnelMessage{
		Method: r.Method,
		Path:   r.URL.Path,
		Body:   body,
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "erro ao codificar mensagem", http.StatusInternalServerError)
		return
	}

	if err := client.Conn.WriteMessage(websocket.BinaryMessage, jsonMsg); err != nil {
		http.Error(w, "falha ao enviar para o cliente", http.StatusInternalServerError)
		return
	}

	_, rawResp, err := client.Conn.ReadMessage()
	if err != nil {
		log.Printf("Erro do cliente WebSocket: %v", err)
		http.Error(w, "erro do cliente", http.StatusInternalServerError)
		return
	}

	var tr shared.TunnelResponse
	if err := json.Unmarshal(rawResp, &tr); err != nil {
		http.Error(w, "resposta inválida do cliente", http.StatusInternalServerError)
		return
	}

	for k, v := range tr.Headers {
		k = strings.ToLower(k)
		if k == "transfer-encoding" || k == "content-encoding" || k == "connection" || k == "keep-alive" || k == "content-length" {
			continue
		}
		w.Header().Set(k, v)
	}

	w.WriteHeader(tr.StatusCode)
	w.Write(tr.Body)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", handleProxy)
	log.Println("Relay em http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
