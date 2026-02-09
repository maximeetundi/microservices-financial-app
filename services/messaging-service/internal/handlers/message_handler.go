package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
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
	ReadAt         *time.Time         `json:"read_at,omitempty" bson:"read_at,omitempty"`
	Status         string             `json:"status" bson:"status"` // sending, sent, delivered, read
}

type Participant struct {
	UserID string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Email  string `json:"email,omitempty" bson:"email,omitempty"`
	Phone  string `json:"phone,omitempty" bson:"phone,omitempty"`
	Online bool   `json:"online,omitempty" bson:"-"` // Not stored in DB, computed at runtime
}

type ConversationContext struct {
	Type        string `json:"type,omitempty" bson:"type,omitempty"` // e.g. "shop"
	ShopID      string `json:"shop_id,omitempty" bson:"shop_id,omitempty"`
	ShopName    string `json:"shop_name,omitempty" bson:"shop_name,omitempty"`
	ProductID   string `json:"product_id,omitempty" bson:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty" bson:"product_name,omitempty"`
}

type Conversation struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Participants []Participant      `json:"participants" bson:"participants"`
	LastMessage  *Message           `json:"last_message,omitempty" bson:"last_message,omitempty"`
	Context      *ConversationContext `json:"context,omitempty" bson:"context,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	// Helper fields for display (not stored in DB)
	Name        string `json:"name,omitempty" bson:"-"`
	UnreadCount int    `json:"unread_count" bson:"-"`
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

// getMultiplePresence fetches online status for multiple users from auth-service
func (h *MessageHandler) getMultiplePresence(userIDs []string) map[string]string {
	result := make(map[string]string)

	if len(userIDs) == 0 {
		return result
	}

	client := &http.Client{Timeout: 3 * time.Second}

	payload := map[string]interface{}{
		"user_ids": userIDs,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return result
	}

	req, err := http.NewRequest("POST", "http://auth-service:8081/api/v1/users/presence/batch", bytes.NewBuffer(jsonData))
	if err != nil {
		return result
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[PRESENCE] Error calling auth-service: %v", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return result
	}

	var response struct {
		Presences []struct {
			UserID   string `json:"user_id"`
			Status   string `json:"status"`
			LastSeen string `json:"last_seen"`
		} `json:"presences"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return result
	}

	for _, p := range response.Presences {
		result[p.UserID] = p.Status
	}

	return result
}

// GetConversations gets all conversations for the current user with unread counts and participant status
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

	// Collect other participant IDs for batch presence check
	otherUserIDs := make([]string, 0)

	// Enrich each conversation
	for i := range conversations {
		convID := conversations[i].ID.Hex()

		// Set conversation name based on context type
		if conversations[i].Context != nil && conversations[i].Context.Type == "shop" && conversations[i].Context.ShopName != "" {
			// For shop conversations, use the shop name
			conversations[i].Name = conversations[i].Context.ShopName
		} else {
			// For regular conversations, find the other participant's name
			for _, p := range conversations[i].Participants {
				if p.UserID != userID {
					conversations[i].Name = p.Name
					break
				}
			}
		}

		// Find the other participant and get online status
		for j, p := range conversations[i].Participants {
			if p.UserID != userID {
				otherUserIDs = append(otherUserIDs, p.UserID)

				// Check online status (will be updated after batch call)
				conversations[i].Participants[j].Online = false
				break
			}
		}

		// Count unread messages for this conversation (messages not sent by user and not read)
		unreadCount, _ := h.db.Collection("messages").CountDocuments(ctx, bson.M{
			"conversation_id": convID,
			"sender_id":       bson.M{"$ne": userID},
			"read":            false,
		})
		conversations[i].UnreadCount = int(unreadCount)

		// Ensure last_message is populated (if missing, fetch the latest)
		if conversations[i].LastMessage == nil {
			var lastMsg Message
			err := h.db.Collection("messages").FindOne(ctx,
				bson.M{"conversation_id": convID},
				options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}),
			).Decode(&lastMsg)
			if err == nil {
				conversations[i].LastMessage = &lastMsg
			}
		}
	}

	// Batch check online status from auth-service
	if len(otherUserIDs) > 0 {
		onlineStatus := h.getMultiplePresence(otherUserIDs)

		// Update participants with online status
		for i := range conversations {
			for j, p := range conversations[i].Participants {
				if status, ok := onlineStatus[p.UserID]; ok {
					conversations[i].Participants[j].Online = (status == "online")
				}
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

	// Send notification to other participants if they're not in chat
	go h.sendMessageNotifications(conversationID, userID, message.Content, req.SenderName)

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
		Context          *ConversationContext `json:"context"`
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
		// Backfill context if the conversation already exists but doesn't have it.
		if existingConv.Context == nil && req.Context != nil {
			objID := existingConv.ID
			_, _ = h.db.Collection("conversations").UpdateOne(ctx,
				bson.M{"_id": objID},
				bson.M{"$set": bson.M{"context": req.Context, "updated_at": time.Now()}},
			)
			existingConv.Context = req.Context
		}

		// Set conversation name based on context type
		if existingConv.Context != nil && existingConv.Context.Type == "shop" && existingConv.Context.ShopName != "" {
			existingConv.Name = existingConv.Context.ShopName
		} else {
			existingConv.Name = req.ParticipantName
		}

		// Conversation already exists, return it
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
		Context:   req.Context,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.db.Collection("conversations").InsertOne(ctx, conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	// Set conversation name based on context type
	if conversation.Context != nil && conversation.Context.Type == "shop" && conversation.Context.ShopName != "" {
		conversation.Name = conversation.Context.ShopName
	} else {
		conversation.Name = req.ParticipantName
	}

	c.JSON(http.StatusCreated, conversation)
}

// MarkAsRead marks a message as read
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	messageID := c.Param("id")
	userID, _ := getUserID(c)

	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	now := time.Now()
	_, err = h.db.Collection("messages").UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"read":    true,
			"read_at": now,
			"status":  "read",
		}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message_id": messageID,
		"read_at":    now.Format(time.RFC3339),
		"read_by":    userID,
	})
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

// sendMessageNotifications sends push notifications to conversation participants who are not currently in chat
func (h *MessageHandler) sendMessageNotifications(conversationID, senderID, content, senderName string) {
	log.Printf("[NOTIF] Starting sendMessageNotifications for conversation %s, sender %s", conversationID, senderID)
	ctx := context.Background()

	// Get conversation to find participants
	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		log.Printf("[NOTIF] Error parsing conversation ID: %v", err)
		return
	}

	var conv Conversation
	err = h.db.Collection("conversations").FindOne(ctx, bson.M{"_id": objID}).Decode(&conv)
	if err != nil {
		log.Printf("[NOTIF] Error fetching conversation: %v", err)
		return
	}

	log.Printf("[NOTIF] Found %d participants in conversation", len(conv.Participants))

	// For each participant except sender, check if they're in chat and send notification
	for _, p := range conv.Participants {
		if p.UserID == senderID {
			log.Printf("[NOTIF] Skipping sender %s", p.UserID)
			continue // Skip sender
		}

		log.Printf("[NOTIF] Checking if user %s is in chat", p.UserID)
		// Check if user is currently in chat
		if h.isUserInChat(p.UserID) {
			log.Printf("[NOTIF] User %s IS in chat, skipping notification", p.UserID)
			continue // User is viewing messages, no notification needed
		}

		log.Printf("[NOTIF] User %s is NOT in chat, sending notification", p.UserID)
		// User is not in chat, send notification
		h.sendNotification(p.UserID, senderName, content)
	}
}

// isUserInChat checks if a user is currently viewing the messages page
func (h *MessageHandler) isUserInChat(userID string) bool {
	// Call auth-service to check chat activity
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://auth-service:8081/api/v1/users/chat-activity/" + userID)
	if err != nil {
		return false // If we can't check, assume not in chat
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}

	var result struct {
		InChat bool `json:"in_chat"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	return result.InChat
}

// sendNotification sends a push notification via notification-service
func (h *MessageHandler) sendNotification(userID, senderName, content string) {
	log.Printf("[NOTIF] Sending notification to user %s from %s: %s", userID, senderName, content)
	client := &http.Client{Timeout: 5 * time.Second}

	// Truncate content for notification
	notifContent := content
	if len(notifContent) > 100 {
		notifContent = notifContent[:97] + "..."
	}

	// Use a default sender name if empty
	if senderName == "" {
		senderName = "Nouveau message"
	}

	payload := map[string]interface{}{
		"user_id": userID,
		"title":   senderName,
		"message": notifContent,
		"type":    "message",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[NOTIF] Error marshaling payload: %v", err)
		return
	}

	log.Printf("[NOTIF] Calling notification-service with payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", "http://notification-service:8087/api/v1/notifications", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[NOTIF] Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[NOTIF] Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("[NOTIF] Notification response status: %d", resp.StatusCode)
}
