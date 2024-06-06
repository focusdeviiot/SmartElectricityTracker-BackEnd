package external

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

type WebSocketHandler struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (w *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	w.clients[c] = true
	defer func() {
		delete(w.clients, c)
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		w.broadcast <- msg
	}
}

func (w *WebSocketHandler) Start() {
	for {
		msg := <-w.broadcast
		for client := range w.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Error writing message:", err)
				client.Close()
				delete(w.clients, client)
			}
		}
	}
}

func (w *WebSocketHandler) Broadcast(data []byte) {
	w.broadcast <- data
}
