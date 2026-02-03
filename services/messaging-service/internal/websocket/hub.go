package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID string
	Send   chan []byte
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

	// Conversation participants cache (conversation_id -> []user_ids)
	conversations map[string][]string

	mu sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type           string      `json:"type"` // new_message, typing, read, presence
	ConversationID string      `json:"conversation_id,omitempty"`
	SenderID       string      `json:"sender_id,omitempty"`
	RecipientID    string      `json:"recipient_id,omitempty"`
	Content        interface{} `json:"content,omitempty"`
	UserID         string      `json:"user_id,omitempty"`
	IsTyping       bool        `json:"is_typing,omitempty"`
	Status         string      `json:"status,omitempty"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:       make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan *Message, 256),
		conversations: make(map[string][]string),
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
			log.Printf("WebSocket: Client connected: %s (total: %d)", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket: Client disconnected: %s (total: %d)", client.UserID, len(h.clients))

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
			log.Printf("WebSocket: Sent %s to %s", msg.Type, msg.RecipientID)
		default:
			// Client buffer full, disconnect
			h.mu.Lock()
			delete(h.clients, msg.RecipientID)
			close(client.Send)
			h.mu.Unlock()
			log.Printf("WebSocket: Client buffer full, disconnected: %s", msg.RecipientID)
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

// SendNewMessage notifies a user of a new message
func (h *Hub) SendNewMessage(recipientID, senderID, conversationID string, messageContent interface{}) {
	h.broadcast <- &Message{
		Type:           "new_message",
		ConversationID: conversationID,
		SenderID:       senderID,
		RecipientID:    recipientID,
		Content:        messageContent,
	}
}

// BroadcastPresence broadcasts presence status to all connected users who share conversations
func (h *Hub) BroadcastPresence(userID, status string) {
	h.mu.RLock()
	// For simplicity, broadcast to all connected users
	// In production, you'd want to filter to only users who share conversations
	for id := range h.clients {
		if id != userID {
			h.broadcast <- &Message{
				Type:        "presence",
				RecipientID: id,
				UserID:      userID,
				Status:      status,
			}
		}
	}
	h.mu.RUnlock()
}

// BroadcastTyping sends typing indicator to conversation participants
func (h *Hub) BroadcastTyping(conversationID, userID string, isTyping bool, participantIDs []string) {
	for _, pid := range participantIDs {
		if pid != userID {
			h.broadcast <- &Message{
				Type:           "typing",
				ConversationID: conversationID,
				RecipientID:    pid,
				UserID:         userID,
				IsTyping:       isTyping,
			}
		}
	}
}

// SetConversationParticipants caches conversation participants
func (h *Hub) SetConversationParticipants(conversationID string, participants []string) {
	h.mu.Lock()
	h.conversations[conversationID] = participants
	h.mu.Unlock()
}

// GetConversationParticipants returns cached conversation participants
func (h *Hub) GetConversationParticipants(conversationID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.conversations[conversationID]
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

// GetOnlineStatus returns online status for multiple users
func (h *Hub) GetOnlineStatus(userIDs []string) map[string]string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]string)
	for _, id := range userIDs {
		if _, ok := h.clients[id]; ok {
			result[id] = "online"
		} else {
			result[id] = "offline"
		}
	}
	return result
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
