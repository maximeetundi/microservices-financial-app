package handlers

import (
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/services"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService *services.ClientService
}

func NewClientHandler(clientService *services.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

// InviteClient invites a client to a private shop
func (h *ClientHandler) InviteClient(c *gin.Context) {
	userID := c.GetString("user_id")
	shopID := c.Param("id")

	var req models.InviteClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.clientService.InviteClient(c.Request.Context(), shopID, &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Invitation envoyée",
		"client":  client,
	})
}

// GetMyInvitations returns pending invitations for the authenticated user
func (h *ClientHandler) GetMyInvitations(c *gin.Context) {
	email := c.GetString("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found in token"})
		return
	}

	invitations, err := h.clientService.GetPendingInvitations(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invitations": invitations})
}

// AcceptInvitation accepts a shop invitation with PIN verification
func (h *ClientHandler) AcceptInvitation(c *gin.Context) {
	userID := c.GetString("user_id")
	email := c.GetString("email")

	var req models.AcceptClientInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.clientService.AcceptInvitation(c.Request.Context(), req.InvitationID, userID, email, req.PIN); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation acceptée"})
}

// DeclineInvitation declines a shop invitation
func (h *ClientHandler) DeclineInvitation(c *gin.Context) {
	email := c.GetString("email")
	invitationID := c.Param("id")

	if err := h.clientService.DeclineInvitation(c.Request.Context(), invitationID, email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation refusée"})
}

// RevokeClientAccess revokes a client's access to a private shop
func (h *ClientHandler) RevokeClientAccess(c *gin.Context) {
	userID := c.GetString("user_id")
	shopID := c.Param("id")
	clientID := c.Param("clientId")

	if err := h.clientService.RevokeClientAccess(c.Request.Context(), shopID, clientID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accès révoqué"})
}

// ListShopClients returns all clients for a shop
func (h *ClientHandler) ListShopClients(c *gin.Context) {
	userID := c.GetString("user_id")
	shopID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.clientService.ListShopClients(c.Request.Context(), shopID, userID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetMyPrivateShops returns private shops where user has access
func (h *ClientHandler) GetMyPrivateShops(c *gin.Context) {
	userID := c.GetString("user_id")

	shops, err := h.clientService.GetMyPrivateShops(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shops": shops})
}
