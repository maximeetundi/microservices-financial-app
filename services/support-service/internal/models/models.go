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
	MessageTypeUser  MessageType = "user"
	MessageTypeAgent MessageType = "agent"
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
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email,omitempty" db:"email"`
	Type        AgentType `json:"type" db:"type"`
	Avatar      string    `json:"avatar,omitempty" db:"avatar"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	MaxChats    int       `json:"max_chats" db:"max_chats"`
	ActiveChats int       `json:"active_chats" db:"active_chats"`
	Skills      []string  `json:"skills,omitempty"`
	Rating      float64   `json:"rating" db:"rating"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Conversation represents a support conversation
type Conversation struct {
	ID             string             `json:"id" db:"id"`
	UserID         string             `json:"user_id" db:"user_id"`
	UserName       string             `json:"user_name" db:"user_name"`
	UserEmail      string             `json:"user_email,omitempty" db:"user_email"`
	AgentID        *string            `json:"agent_id,omitempty" db:"agent_id"`
	AgentType      AgentType          `json:"agent_type" db:"agent_type"`
	Subject        string             `json:"subject" db:"subject"`
	Category       string             `json:"category" db:"category"`
	Status         ConversationStatus `json:"status" db:"status"`
	Priority       Priority           `json:"priority" db:"priority"`
	LastMessage    *string            `json:"last_message,omitempty" db:"last_message"`
	LastMessageAt  *time.Time         `json:"last_message_at,omitempty" db:"last_message_at"`
	UnreadCount    int                `json:"unread_count" db:"unread_count"`
	MessageCount   int                `json:"message_count" db:"message_count"`
	Rating         *int               `json:"rating,omitempty" db:"rating"`
	Feedback       *string            `json:"feedback,omitempty" db:"feedback"`
	ResolvedAt     *time.Time         `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt      time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" db:"updated_at"`
	
	// Relations (not stored in DB)
	Agent    *Agent    `json:"agent,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

// Message represents a single message in a conversation
type Message struct {
	ID             string      `json:"id" db:"id"`
	ConversationID string      `json:"conversation_id" db:"conversation_id"`
	SenderID       string      `json:"sender_id" db:"sender_id"`
	SenderName     string      `json:"sender_name" db:"sender_name"`
	SenderType     MessageType `json:"sender_type" db:"sender_type"`
	Content        string      `json:"content" db:"content"`
	ContentType    string      `json:"content_type" db:"content_type"` // text, image, file
	Attachments    []string    `json:"attachments,omitempty" db:"attachments"`
	IsRead         bool        `json:"is_read" db:"is_read"`
	ReadAt         *time.Time  `json:"read_at,omitempty" db:"read_at"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
}

// QuickReply represents a pre-defined quick reply option
type QuickReply struct {
	ID       string `json:"id" db:"id"`
	Label    string `json:"label" db:"label"`
	Response string `json:"response" db:"response"`
	Category string `json:"category" db:"category"`
	IsActive bool   `json:"is_active" db:"is_active"`
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
	Content     string   `json:"content" binding:"required,min=1,max=2000"`
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
