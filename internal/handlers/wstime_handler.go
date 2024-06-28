package handlers

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
)

type WSTimeHandler struct {
	conn map[*websocket.Conn]bool
	mu   sync.Mutex
}

func NewWSTimeHandler() *WSTimeHandler {
	return &WSTimeHandler{
		conn: make(map[*websocket.Conn]bool),
	}
}

func (w *WSTimeHandler) HandleWSTime(c *websocket.Conn) {
	log.Info("User connected Time")
	w.mu.Lock()
	w.conn[c] = true
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.conn, c)
		w.mu.Unlock()
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Info("Error reading message:", err)
			break
		}
	}
}

func (w *WSTimeHandler) StartBroadcast() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			w.mu.Lock()
			for conn := range w.conn {
				err := conn.WriteJSON(map[string]string{"time": t.Format(time.RFC3339)})
				if err != nil {
					log.Info("Error writing message:", err)
					conn.Close()
					delete(w.conn, conn)
				}
			}
			w.mu.Unlock()
		}
	}
}
