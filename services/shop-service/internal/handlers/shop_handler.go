package handlers

import (
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/services"
	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	shopService *services.ShopService
}

func NewShopHandler(shopService *services.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

// List returns all public shops
func (h *ShopHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.shopService.ListPublic(c.Request.Context(), page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get returns a shop by slug
func (h *ShopHandler) Get(c *gin.Context) {
	slug := c.Param("id")

	shop, err := h.shopService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// GetByID returns a shop by ID
func (h *ShopHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	shop, err := h.shopService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// GetByWalletID returns a shop by Wallet ID
func (h *ShopHandler) GetByWalletID(c *gin.Context) {
	walletID := c.Param("wallet_id")

	shop, err := h.shopService.GetByWalletID(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// Create creates a new shop
func (h *ShopHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := c.GetHeader("Authorization")
	shop, err := h.shopService.Create(c.Request.Context(), &req, userID, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shop)
}

// Update updates a shop
func (h *ShopHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	shopID := c.Param("id")

	var req models.UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := c.GetHeader("Authorization")
	shop, err := h.shopService.Update(c.Request.Context(), shopID, &req, userID, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// Delete deletes a shop
func (h *ShopHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	shopID := c.Param("id")

	if err := h.shopService.Delete(c.Request.Context(), shopID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shop deleted"})
}

// GetMyShops returns shops owned or managed by the user
func (h *ShopHandler) GetMyShops(c *gin.Context) {
	userID := c.GetString("user_id")

	shops, err := h.shopService.GetMyShops(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shops": shops})
}

// InviteManager invites a new manager to the shop
func (h *ShopHandler) InviteManager(c *gin.Context) {
	userID := c.GetString("user_id")
	shopID := c.Param("id")

	var req models.InviteManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.shopService.InviteManager(c.Request.Context(), shopID, &req, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager invited"})
}

// RemoveManager removes a manager from the shop
func (h *ShopHandler) RemoveManager(c *gin.Context) {
	userID := c.GetString("user_id")
	shopID := c.Param("id")
	targetUserID := c.Param("userId")

	if err := h.shopService.RemoveManager(c.Request.Context(), shopID, targetUserID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed"})
}
