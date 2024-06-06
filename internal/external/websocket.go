package external

import (
	"encoding/json"
	"errors"
	"log"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type WebSocketHandler struct {
	clients   map[*websocket.Conn]*Client
	broadcast chan map[string]map[string]float32
	cfg       *config.Config
	mu        sync.Mutex
	userRepo  *repositories.UserRepository
}

type Client struct {
	conn   *websocket.Conn
	userID string
}

func NewWebSocketHandler(userRepo *repositories.UserRepository, cfg *config.Config) *WebSocketHandler {
	return &WebSocketHandler{
		clients:   make(map[*websocket.Conn]*Client),
		broadcast: make(chan map[string]map[string]float32),
		cfg:       cfg,
		userRepo:  userRepo,
	}
}

func (w *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	token := c.Query("token")
	if token == "" {
		log.Println("Token not provided")
		c.Close()
		return
	}

	claims, err := w.ValidateToken(token)
	if err != nil {
		log.Println("Invalid token:", err)
		c.Close()
		return
	}

	client := &Client{
		conn:   c,
		userID: claims.UserID,
	}

	w.mu.Lock()
	w.clients[c] = client
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.clients, c)
		w.mu.Unlock()
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
	}
}

func (w *WebSocketHandler) Start() {
	for {
		msg := <-w.broadcast
		w.mu.Lock()
		for _, client := range w.clients {
			data := w.filterDataForDevices(msg, client.userID)
			if data == nil {
				continue
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Println("Error marshaling data:", err)
				continue
			}

			err = client.conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Error writing message:", err)
				client.conn.Close()
				delete(w.clients, client.conn)
			}

		}
		w.mu.Unlock()
	}
}

func (w *WebSocketHandler) Broadcast(data map[string]map[string]float32) {
	w.broadcast <- data
}

func (w *WebSocketHandler) filterDataForDevices(data map[string]map[string]float32, user_id string) map[string]map[string]float32 {
	userId, err := uuid.Parse(user_id)
	if err != nil {
		log.Println("Error parsing user_id:", err)
		return nil
	}

	userDevice, err := w.userRepo.FindUserDeviceById(userId)
	if err != nil {
		log.Println("Error finding user device:", err)
		return nil
	}

	devices := make([]string, 0)
	for _, device := range userDevice {
		devices = append(devices, device.DeviceID)
	}

	filteredData := make(map[string]map[string]float32)
	for _, device := range devices {
		if d, ok := data[device]; ok {
			filteredData[device] = d
		}
	}
	return filteredData
}

func (w *WebSocketHandler) ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(w.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Exp.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
