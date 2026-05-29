package handlers

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"consumable-tracker/middleware"
)

type WebSocketHub struct {
	mu      sync.RWMutex
	clients map[int]map[*websocket.Conn]bool // userID -> connections
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[int]map[*websocket.Conn]bool),
	}
}

func (h *WebSocketHub) Register(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*websocket.Conn]bool)
	}
	h.clients[userID][conn] = true
	log.Printf("WebSocket client registered: user %d (total clients: %d)", userID, h.countAll())
}

func (h *WebSocketHub) Unregister(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] != nil {
		delete(h.clients[userID], conn)
		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
		}
	}
}

func (h *WebSocketHub) Broadcast(role string, message string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conns := range h.clients {
		for conn := range conns {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Printf("WebSocket write error: %v", err)
			}
		}
	}
}

func (h *WebSocketHub) countAll() int {
	count := 0
	for _, conns := range h.clients {
		count += len(conns)
	}
	return count
}

func WebSocketUpgrade(hub *WebSocketHub) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func WebSocketHandler(hub *WebSocketHub) websocket.Handler {
	return func(conn *websocket.Conn) {
		userID, ok := conn.Locals("user_id").(int)
		if !ok {
			conn.Close()
			return
		}

		role, ok := conn.Locals("role").(string)
		if !ok || role != "admin" {
			conn.Close()
			return
		}

		hub.Register(userID, conn)
		defer hub.Unregister(userID, conn)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}

func SetupWebSocketRoute(app *fiber.App, hub *WebSocketHub) {
	app.Use("/ws", middleware.AuthRequired())
	app.Use("/ws", middleware.RequireAdmin())
	app.Get("/ws", WebSocketUpgrade(hub))
	app.Get("/ws", websocket.New(WebSocketHandler(hub)))
}
