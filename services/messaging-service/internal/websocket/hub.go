package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	UserID   string
	Send     chan []byte
}

// Hub maintains active client connections and broadcasts messages
type Hub struct {
	// Registered clients by user ID
	clients map[string]*Client
	
	// Register requests from clients
	register chan *Client
	
	// Unregister requests from clients
	unregister chan *Client
	
	// Broadcast message to specific user
	broadcast chan *Message
	
	mu sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type           string      `json:"type"` // new_message, typing, read, presence
	ConversationID string      `json:"conversation_id,omitempty"`
	SenderID       string      `json:"sender_id,omitempty"`
	RecipientID    string      `json:"recipient_id,omitempty"`
	Content        interface{} `json:"content,omitempty"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("WebSocket: Client connected: %s", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket: Client disconnected: %s", client.UserID)

		case message := <-h.broadcast:
			h.sendToUser(message)
		}
	}
}

// sendToUser sends a message to a specific user
func (h *Hub) sendToUser(msg *Message) {
	h.mu.RLock()
	client, ok := h.clients[msg.RecipientID]
	h.mu.RUnlock()

	if ok {
		data, _ := json.Marshal(msg)
		select {
		case client.Send <- data:
		default:
			// Client buffer full, disconnect
			h.mu.Lock()
			delete(h.clients, msg.RecipientID)
			close(client.Send)
			h.mu.Unlock()
		}
	}
}

// SendMessage broadcasts a message to a user
func (h *Hub) SendMessage(recipientID string, msgType string, content interface{}) {
	h.broadcast <- &Message{
		Type:        msgType,
		RecipientID: recipientID,
		Content:     content,
	}
}

// SendToConversation sends a message to all participants in a conversation
func (h *Hub) SendToConversation(participantIDs []string, senderID string, msg *Message) {
	for _, pid := range participantIDs {
		if pid != senderID { // Don't send to sender
			msgCopy := *msg
			msgCopy.RecipientID = pid
			h.broadcast <- &msgCopy
		}
	}
}

// IsOnline checks if a user is connected
func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// GetOnlineUsers returns list of online user IDs
func (h *Hub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	users := make([]string, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
