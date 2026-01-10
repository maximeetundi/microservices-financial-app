package models

import "time"

// AgentType represents the type of support agent
type AgentType string

const (
	AgentTypeAI    AgentType = "ai"
	AgentTypeHuman AgentType = "human"
)

// ConversationStatus represents the current status of a conversation
type ConversationStatus string

const (
	ConversationStatusOpen      ConversationStatus = "open"
	ConversationStatusPending   ConversationStatus = "pending"    // Waiting for agent
	ConversationStatusActive    ConversationStatus = "active"     // Agent is responding
	ConversationStatusResolved  ConversationStatus = "resolved"
	ConversationStatusClosed    ConversationStatus = "closed"
	ConversationStatusEscalated ConversationStatus = "escalated"  // Escalated to human
)

// MessageType represents who sent the message
type MessageType string

const (
	MessageTypeUser   MessageType = "user"
	MessageTypeAgent  MessageType = "agent"
	MessageTypeSystem MessageType = "system"
)

// Priority represents ticket/conversation priority
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Agent represents a support agent (AI or human)
type Agent struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email,omitempty" bson:"email"`
	Type        AgentType `json:"type" bson:"type"`
	Avatar      string    `json:"avatar,omitempty" bson:"avatar"`
	IsAvailable bool      `json:"is_available" bson:"is_available"`
	MaxChats    int       `json:"max_chats" bson:"max_chats"`
	ActiveChats int       `json:"active_chats" bson:"active_chats"`
	Skills      []string  `json:"skills,omitempty" bson:"skills"`
	Rating      float64   `json:"rating" bson:"rating"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Conversation represents a support conversation
type Conversation struct {
	ID             string             `json:"id" bson:"_id,omitempty"`
	UserID         string             `json:"user_id" bson:"user_id"`
	UserName       string             `json:"user_name" bson:"user_name"`
	UserEmail      string             `json:"user_email,omitempty" bson:"user_email"`
	AgentID        *string            `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentType      AgentType          `json:"agent_type" bson:"agent_type"`
	Subject        string             `json:"subject" bson:"subject"`
	Category       string             `json:"category" bson:"category"`
	Status         ConversationStatus `json:"status" bson:"status"`
	Priority       Priority           `json:"priority" bson:"priority"`
	LastMessage    *string            `json:"last_message,omitempty" bson:"last_message,omitempty"`
	LastMessageAt  *time.Time         `json:"last_message_at,omitempty" bson:"last_message_at,omitempty"`
	UnreadCount    int                `json:"unread_count" bson:"unread_count"`
	MessageCount   int                `json:"message_count" bson:"message_count"`
	Rating         *int               `json:"rating,omitempty" bson:"rating,omitempty"`
	Feedback       *string            `json:"feedback,omitempty" bson:"feedback,omitempty"`
	ResolvedAt     *time.Time         `json:"resolved_at,omitempty" bson:"resolved_at,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`

	// Relations (not stored in DB, populated if needed)
	Agent    *Agent    `json:"agent,omitempty" bson:"-"`
	Messages []Message `json:"messages,omitempty" bson:"-"` // Or maybe store truncated list?
}

// Message represents a single message in a conversation
type Message struct {
	ID             string      `json:"id" bson:"_id,omitempty"`
	ConversationID string      `json:"conversation_id" bson:"conversation_id"`
	SenderID       string      `json:"sender_id" bson:"sender_id"`
	SenderName     string      `json:"sender_name" bson:"sender_name"`
	SenderType     MessageType `json:"sender_type" bson:"sender_type"`
	Content        string      `json:"content" bson:"content"`
	ContentType    string      `json:"content_type" bson:"content_type"` // text, image, file
	Attachments    []string    `json:"attachments" bson:"attachments"`
	IsRead         bool        `json:"is_read" bson:"is_read"`
	ReadAt         *time.Time  `json:"read_at,omitempty" bson:"read_at,omitempty"`
	CreatedAt      time.Time   `json:"created_at" bson:"created_at"`
}

// QuickReply represents a pre-defined quick reply option
type QuickReply struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Label    string `json:"label" bson:"label"`
	Response string `json:"response" bson:"response"`
	Category string `json:"category" bson:"category"`
	IsActive bool   `json:"is_active" bson:"is_active"`
}

// SupportStats represents support statistics
type SupportStats struct {
	TotalConversations    int     `json:"total_conversations"`
	OpenConversations     int     `json:"open_conversations"`
	ResolvedToday         int     `json:"resolved_today"`
	AvgResponseTime       float64 `json:"avg_response_time_minutes"`
	AvgResolutionTime     float64 `json:"avg_resolution_time_hours"`
	CustomerSatisfaction  float64 `json:"customer_satisfaction"`
	ActiveAgents          int     `json:"active_agents"`
	PendingConversations  int     `json:"pending_conversations"`
}

// --- Request/Response DTOs ---

type CreateConversationRequest struct {
	AgentType AgentType `json:"agent_type" binding:"required,oneof=ai human"`
	Subject   string    `json:"subject" binding:"required,min=5,max=200"`
	Category  string    `json:"category" binding:"required"`
	Message   string    `json:"message" binding:"required,min=1,max=2000"`
	Priority  Priority  `json:"priority,omitempty"`
}

type SendMessageRequest struct {
	Content     string   `json:"content" binding:"max=2000"`
	ContentType string   `json:"content_type,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type AssignAgentRequest struct {
	AgentID string `json:"agent_id" binding:"required"`
}

type CloseConversationRequest struct {
	Rating   int    `json:"rating,omitempty"`
	Feedback string `json:"feedback,omitempty"`
}

type EscalateRequest struct {
	Reason   string   `json:"reason" binding:"required"`
	Priority Priority `json:"priority,omitempty"`
}
