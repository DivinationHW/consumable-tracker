package handlers

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WSHub struct {
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

func NewWSHub() *WSHub {
	return &WSHub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *WSHub) HandleWS(c *websocket.Conn) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.clients, c)
		h.mu.Unlock()
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *WSHub) Broadcast(msg interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("[WS] send error: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}

func WSUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return c.Status(400).JSON(fiber.Map{"error": "需要WebSocket连接"})
}

func BroadcastFunc(hub *WSHub) func(msg interface{}) {
	return func(msg interface{}) {
		hub.Broadcast(msg)
	}
}
