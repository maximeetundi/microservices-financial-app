package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"messaging-service/internal/services"
)

type MessageHandler struct {
	db      *mongo.Database
	storage *services.StorageService
}

func NewMessageHandler(db *mongo.Database, storage *services.StorageService) *MessageHandler {
	return &MessageHandler{
		db:      db,
		storage: storage,
	}
}

// Message models
type Attachment struct {
	FileName     string  `json:"file_name" bson:"file_name"`
	FileURL      string  `json:"file_url" bson:"file_url"`
	FileSize     int64   `json:"file_size" bson:"file_size"`
	MimeType     string  `json:"mime_type" bson:"mime_type"`
	Duration     float64 `json:"duration,omitempty" bson:"duration,omitempty"`
	ThumbnailURL string  `json:"thumbnail_url,omitempty" bson:"thumbnail_url,omitempty"`
	Width        int     `json:"width,omitempty" bson:"width,omitempty"`
	Height       int     `json:"height,omitempty" bson:"height,omitempty"`
}

type Message struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ConversationID string             `json:"conversation_id" bson:"conversation_id"`
	SenderID       string             `json:"sender_id" bson:"sender_id"`
	SenderName     string             `json:"sender_name,omitempty" bson:"sender_name,omitempty"`
	Content        string             `json:"content" bson:"content"`
	MessageType    string             `json:"message_type" bson:"message_type"`
	Attachment     *Attachment        `json:"attachment,omitempty" bson:"attachment,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	Read           bool               `json:"read" bson:"read"`
}

type Participant struct {
	UserID string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Email  string `json:"email,omitempty" bson:"email,omitempty"`
	Phone  string `json:"phone,omitempty" bson:"phone,omitempty"`
}

type Conversation struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Participants []Participant      `json:"participants" bson:"participants"`
	LastMessage  *Message           `json:"last_message,omitempty" bson:"last_message,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	// Helper fields for display
	Name         string             `json:"name,omitempty" bson:"-"`
}

// Helper to get user ID from context (set by auth middleware)
func getUserID(c *gin.Context) (string, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	userID, ok := userIDInterface.(string)
	return userID, ok
}

// GetConversations gets all conversations for the current user
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ctx := context.Background()
	conversations := []Conversation{}

	cursor, err := h.db.Collection("conversations").Find(ctx, bson.M{
		"participants.user_id": userID,
	}, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}))
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &conversations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode conversations"})
		return
	}

	// Set display name (the other participant's name)
	for i := range conversations {
		for _, p := range conversations[i].Participants {
			if p.UserID != userID {
				conversations[i].Name = p.Name
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

// SendMessage sends a message in a conversation
func (h *MessageHandler) SendMessage(c *gin.Context) {
	conversationID := c.Param("id")
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	var req struct {
		Content     string      `json:"content"`
		MessageType string      `json:"message_type"`
		Attachment  *Attachment `json:"attachment"`
		SenderName  string      `json:"sender_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.MessageType == "" {
		req.MessageType = "text"
	}

	message := Message{
		ID:             primitive.NewObjectID(),
		ConversationID: conversationID,
		SenderID:       userID,
		SenderName:     req.SenderName,
		Content:        req.Content,
		MessageType:    req.MessageType,
		Attachment:     req.Attachment,
		CreatedAt:      time.Now(),
		Read:           false,
	}

	ctx := context.Background()
	
	_, err := h.db.Collection("messages").InsertOne(ctx, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Update conversation's last_message and updated_at
	objID, _ := primitive.ObjectIDFromHex(conversationID)
	h.db.Collection("conversations").UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"last_message": message,
				"updated_at":   time.Now(),
			},
		},
	)

	c.JSON(http.StatusCreated, message)
}

// GetMessages gets messages for a conversation
func (h *MessageHandler) GetMessages(c *gin.Context) {
	conversationID := c.Param("id")
	
	ctx := context.Background()
	messages := []Message{}

	cursor, err := h.db.Collection("messages").Find(ctx, bson.M{
		"conversation_id": conversationID,
	}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// CreateConversation creates a new conversation or returns existing one
func (h *MessageHandler) CreateConversation(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	var req struct {
		ParticipantID    string `json:"participant_id"`
		ParticipantName  string `json:"participant_name"`
		ParticipantEmail string `json:"participant_email"`
		ParticipantPhone string `json:"participant_phone"`
		MyName           string `json:"my_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()

	// Check if conversation already exists between these two users
	var existingConv Conversation
	err := h.db.Collection("conversations").FindOne(ctx, bson.M{
		"$and": []bson.M{
			{"participants.user_id": userID},
			{"participants.user_id": req.ParticipantID},
		},
	}).Decode(&existingConv)

	if err == nil {
		// Conversation already exists, return it
		existingConv.Name = req.ParticipantName
		c.JSON(http.StatusOK, existingConv)
		return
	}

	// Create new conversation with participant details
	conversation := Conversation{
		ID: primitive.NewObjectID(),
		Participants: []Participant{
			{
				UserID: userID,
				Name:   req.MyName,
			},
			{
				UserID: req.ParticipantID,
				Name:   req.ParticipantName,
				Email:  req.ParticipantEmail,
				Phone:  req.ParticipantPhone,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.db.Collection("conversations").InsertOne(ctx, conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	conversation.Name = req.ParticipantName
	c.JSON(http.StatusCreated, conversation)
}

// MarkAsRead marks a message as read
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	messageID := c.Param("id")
	
	ctx := context.Background()
	objID, _ := primitive.ObjectIDFromHex(messageID)
	
	_, err := h.db.Collection("messages").UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"read": true}},
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteConversation deletes a conversation
func (h *MessageHandler) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")
	
	ctx := context.Background()
	
	// Delete all messages
	h.db.Collection("messages").DeleteMany(ctx, bson.M{"conversation_id": conversationID})
	
	// Delete conversation
	objID, _ := primitive.ObjectIDFromHex(conversationID)
	_, err := h.db.Collection("conversations").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
