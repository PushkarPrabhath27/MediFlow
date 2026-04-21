package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict to frontend domain
	},
}

type BroadcastMessage struct {
	TenantID string      `json:"tenant_id"`
	Type     string      `json:"type"`
	Payload  interface{} `json:"payload"`
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	tenantID string
	userID   string
	send     chan []byte
}

type Hub struct {
	// Registered clients by tenantID
	clients   map[string]map[*Client]bool
	broadcast chan BroadcastMessage
	register  chan *Client
	unregister chan *Client
	mu        sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.tenantID] == nil {
				h.clients[client.tenantID] = make(map[*Client]bool)
			}
			h.clients[client.tenantID][client] = true
			h.mu.Unlock()
			log.Info().Str("tenant_id", client.tenantID).Str("user_id", client.userID).Msg("Client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if tenantClients, ok := h.clients[client.tenantID]; ok {
				if _, ok := tenantClients[client]; ok {
					delete(tenantClients, client)
					close(client.send)
					if len(tenantClients) == 0 {
						delete(h.clients, client.tenantID)
					}
				}
			}
			h.mu.Unlock()
			log.Info().Str("tenant_id", client.tenantID).Str("user_id", client.userID).Msg("Client unregistered")

		case message := <-h.broadcast:
			h.mu.RLock()
			tenantClients := h.clients[message.TenantID]
			if tenantClients != nil {
				payload, _ := json.Marshal(message)
				for client := range tenantClients {
					select {
					case client.send <- payload:
					default:
						// If send buffer is full, unregister client
						go func(c *Client) { h.unregister <- c }(client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { 
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil 
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("WebSocket read error")
			}
			break
		}
		// We don't expect messages from clients in v1.0 (read-only real-time board)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, tenantID, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("WebSocket upgrade failed")
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		tenantID: tenantID,
		userID:   userID,
		send:     make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
