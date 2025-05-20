package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valcinei/jiboia-tunnel/shared"
)

type Client struct {
	localURL   string
	relayURL   string
	name       string
	baseDomain string
}

func NewClient(localURL, relayURL, name, baseDomain string) *Client {
	return &Client{
		localURL:   localURL,
		relayURL:   relayURL,
		name:       name,
		baseDomain: baseDomain,
	}
}

func (c *Client) Start() error {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s?id=%s", c.relayURL, c.name), nil)
	if err != nil {
		return fmt.Errorf("falha ao conectar ao relay: %w", err)
	}
	defer conn.Close()

	fmt.Printf("Túnel disponível em: http://%s.%s\n", c.name, c.baseDomain)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error to read message from relay:", err)
			break
		}

		var tm shared.TunnelMessage
		if err := json.Unmarshal(msg, &tm); err != nil {
			log.Println("Error to decode message:", err)
			continue
		}

		req, err := http.NewRequest(tm.Method, c.localURL+tm.Path, bytes.NewReader(tm.Body))
		if err != nil {
			log.Println("Error to create local request:", err)
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		var response shared.TunnelResponse

		if err != nil {
			log.Println("Error sending to local service:", err)
			response = shared.TunnelResponse{
				StatusCode: 502,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       []byte("Error connecting to local service"),
			}
		} else {
			defer resp.Body.Close()
			data, _ := io.ReadAll(resp.Body)

			headers := map[string]string{}
			for k, v := range resp.Header {
				kl := strings.ToLower(k)
				if kl == "transfer-encoding" || kl == "connection" || kl == "keep-alive" || kl == "content-length" {
					continue
				}
				headers[k] = v[0]
			}

			response = shared.TunnelResponse{
				StatusCode: resp.StatusCode,
				Headers:    headers,
				Body:       data,
			}
		}

		jsonResp, err := json.Marshal(response)
		if err != nil {
			log.Println("Error serializing response:", err)
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, jsonResp); err != nil {
			log.Println("Error sending response to relay:", err)
		} else {
			log.Printf("Responsed: %s %s → %d", tm.Method, tm.Path, response.StatusCode)
		}
	}

	return nil
}

func GenerateRandomName() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
