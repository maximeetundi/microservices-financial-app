package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (configure properly in production)
	},
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 4096
)

// IncomingMessage represents a message from the client
type IncomingMessage struct {
	Type           string `json:"type"` // typing, read, presence
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id,omitempty"`
	IsTyping       bool   `json:"is_typing,omitempty"`
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from query or auth
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		// Create client
		client := &Client{
			Hub:    hub,
			Conn:   conn,
			UserID: userID,
			Send:   make(chan []byte, 256),
		}

		// Register client
		hub.Register(client)

		// Broadcast online presence
		hub.BroadcastPresence(userID, "online")

		// Start goroutines for reading and writing
		go client.writePump()
		go client.readPump()
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		// Broadcast offline presence
		c.Hub.BroadcastPresence(c.UserID, "offline")
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse and handle incoming messages
		c.handleMessage(message)
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(data []byte) {
	var msg IncomingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Failed to parse message from %s: %v", c.UserID, err)
		return
	}

	switch msg.Type {
	case "typing":
		c.handleTyping(msg)
	case "read":
		c.handleRead(msg)
	case "presence":
		// Handled automatically on connect/disconnect
	default:
		log.Printf("Unknown message type from %s: %s", c.UserID, msg.Type)
	}
}

// handleTyping broadcasts typing indicator to conversation participants
func (c *Client) handleTyping(msg IncomingMessage) {
	if msg.ConversationID == "" {
		return
	}

	// Get participants for this conversation
	participants := c.Hub.GetConversationParticipants(msg.ConversationID)

	// Broadcast to all other participants
	for _, pid := range participants {
		if pid != c.UserID {
			c.Hub.broadcast <- &Message{
				Type:           "typing",
				ConversationID: msg.ConversationID,
				SenderID:       c.UserID,
				RecipientID:    pid,
				Content: map[string]interface{}{
					"user_id":   c.UserID,
					"is_typing": msg.IsTyping,
				},
			}
		}
	}
}

// handleRead broadcasts read receipt to conversation participants
func (c *Client) handleRead(msg IncomingMessage) {
	if msg.MessageID == "" || msg.ConversationID == "" {
		return
	}

	// Get participants for this conversation
	participants := c.Hub.GetConversationParticipants(msg.ConversationID)

	// Broadcast to all other participants
	for _, pid := range participants {
		if pid != c.UserID {
			c.Hub.broadcast <- &Message{
				Type:           "read",
				ConversationID: msg.ConversationID,
				SenderID:       c.UserID,
				RecipientID:    pid,
				Content: map[string]interface{}{
					"message_id": msg.MessageID,
					"read_by":    c.UserID,
					"read_at":    time.Now().Format(time.RFC3339),
				},
			}
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
