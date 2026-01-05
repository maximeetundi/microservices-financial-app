package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAssociationChat gets messages for an association chat
func (h *MessageHandler) GetAssociationChat(c *gin.Context) {
	associationID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	ctx := context.Background()
	messages := []Message{}

	cursor, err := h.db.Collection("messages").Find(ctx, bson.M{
		"conversation_id": "assoc_" + associationID,
	}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)))
	
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

// SendAssociationMessage sends a message in an association chat
func (h *MessageHandler) SendAssociationMessage(c *gin.Context) {
	associationID := c.Param("id")
	userID := c.GetHeader("X-User-ID")
	
	var req struct {
		Content     string      `json:"content"`
		MessageType string      `json:"message_type"`
		Attachment  *Attachment `json:"attachment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Default to text
	if req.MessageType == "" {
		req.MessageType = "text"
	}

	message := Message{
		ID:             primitive.NewObjectID(),
		ConversationID: "assoc_" + associationID, // Prefix for association chats
		SenderID:       userID,
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

	c.JSON(http.StatusCreated, message)
}

// DeleteAssociationMessage deletes a message from association chat
func (h *MessageHandler) DeleteAssociationMessage(c *gin.Context) {
	associationID := c.Param("id")
	messageID := c.Param("messageId")
	
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	
	// Verify message belongs to this association
	_, err = h.db.Collection("messages").DeleteOne(ctx, bson.M{
		"_id":             objID,
		"conversation_id": "assoc_" + associationID,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
