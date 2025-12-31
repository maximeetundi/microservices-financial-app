package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/services"
)

type SupportHandler struct {
	supportService *services.SupportService
}

func NewSupportHandler(supportService *services.SupportService) *SupportHandler {
	return &SupportHandler{
		supportService: supportService,
	}
}

// UploadFile uploads an attachment
func (h *SupportHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size (e.g. 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(400, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Save file to temporary location
	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, fmt.Sprintf("upload_%d_%s", time.Now().UnixNano(), file.Filename))
	
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save uploaded file"})
		return
	}
	defer os.Remove(tempPath) // Clean up temp file

	// Upload using service
	url, err := h.supportService.UploadAttachment(c, tempPath, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	c.JSON(200, gin.H{"url": url})
}

// CreateConversation creates a new support conversation
func (h *SupportHandler) CreateConversation(c *gin.Context) {
	userID := c.GetString("user_id")
	userName := c.GetString("user_name")
	userEmail := c.GetString("user_email")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, msg, err := h.supportService.CreateConversation(userID, userName, userEmail, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Conversation créée avec succès",
		"conversation": conv,
		"first_response": msg,
	})
}

// GetConversations returns user's conversations
func (h *SupportHandler) GetConversations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	conversations, err := h.supportService.GetUserConversations(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
		"total":         len(conversations),
	})
}

// GetConversation returns a single conversation with messages
func (h *SupportHandler) GetConversation(c *gin.Context) {
	conversationID := c.Param("id")

	conv, err := h.supportService.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Mark messages as read
	userID := c.GetString("user_id")
	h.supportService.MarkMessagesAsRead(conversationID, userID)

	c.JSON(http.StatusOK, gin.H{
		"conversation": conv,
	})
}

// SendMessage sends a message in a conversation
func (h *SupportHandler) SendMessage(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetString("user_id")
	userName := c.GetString("user_name")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userMsg, aiResponse, err := h.supportService.SendMessage(
		conversationID,
		userID,
		userName,
		models.MessageTypeUser,
		&req,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message": userMsg,
	}

	if aiResponse != nil {
		response["ai_response"] = aiResponse
	}

	c.JSON(http.StatusOK, response)
}

// GetMessages returns messages for a conversation
func (h *SupportHandler) GetMessages(c *gin.Context) {
	conversationID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.supportService.GetMessages(conversationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    len(messages),
	})
}

// EscalateConversation escalates to human agent
func (h *SupportHandler) EscalateConversation(c *gin.Context) {
	conversationID := c.Param("id")

	var req models.EscalateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.supportService.EscalateConversation(conversationID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Conversation transférée à un conseiller humain",
		"status":  "escalated",
	})
}

// CloseConversation closes a conversation
func (h *SupportHandler) CloseConversation(c *gin.Context) {
	conversationID := c.Param("id")

	var req models.CloseConversationRequest
	c.ShouldBindJSON(&req)

	if err := h.supportService.CloseConversation(conversationID, req.Rating, req.Feedback); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Conversation fermée avec succès",
	})
}

// === Admin Handlers ===

// AdminGetConversations returns all conversations for admin
func (h *SupportHandler) AdminGetConversations(c *gin.Context) {
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	conversations, err := h.supportService.GetAllConversations(status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
		"total":         len(conversations),
	})
}

// AdminAssignAgent assigns an agent to a conversation
func (h *SupportHandler) AdminAssignAgent(c *gin.Context) {
	conversationID := c.Param("id")

	var req models.AssignAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.supportService.AssignAgent(conversationID, req.AgentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Agent assigné avec succès",
	})
}

// AdminSendMessage allows admin/agent to send a message
func (h *SupportHandler) AdminSendMessage(c *gin.Context) {
	conversationID := c.Param("id")
	agentID := c.GetString("user_id")
	agentName := c.GetString("user_name")

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, _, err := h.supportService.SendMessage(
		conversationID,
		agentID,
		agentName,
		models.MessageTypeAgent,
		&req,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

// AdminGetAgents returns all agents
func (h *SupportHandler) AdminGetAgents(c *gin.Context) {
	agents, err := h.supportService.GetAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
		"total":  len(agents),
	})
}

// AdminGetStats returns support statistics
func (h *SupportHandler) AdminGetStats(c *gin.Context) {
	stats, err := h.supportService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// AdminUpdateAgentAvailability updates agent availability
func (h *SupportHandler) AdminUpdateAgentAvailability(c *gin.Context) {
	agentID := c.Param("id")

	var req struct {
		Available bool `json:"available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.supportService.UpdateAgentAvailability(agentID, req.Available); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Disponibilité mise à jour",
	})
}
