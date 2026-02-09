package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *services.ProductService
	orderRepo      *repository.OrderRepository
	storageService *services.StorageService
}

func NewProductHandler(productService *services.ProductService, orderRepo *repository.OrderRepository, storageService *services.StorageService) *ProductHandler {
	return &ProductHandler{productService: productService, orderRepo: orderRepo, storageService: storageService}
}

// ListByShop returns products for a shop
func (h *ProductHandler) ListByShop(c *gin.Context) {
	shopSlug := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categorySlug := c.Query("category")
	status := c.DefaultQuery("status", "active")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.productService.ListByShop(c.Request.Context(), shopSlug, page, pageSize, categorySlug, status, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get returns a product by slug
func (h *ProductHandler) Get(c *gin.Context) {
	shopSlug := c.Param("id")
	productSlug := c.Param("productSlug")

	product, err := h.productService.GetBySlug(c.Request.Context(), shopSlug, productSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetByID returns a product by ID
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Create creates a new product
func (h *ProductHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// Update updates a product
func (h *ProductHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	productID := c.Param("id")

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Update(c.Request.Context(), productID, &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Delete deletes a product
func (h *ProductHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	productID := c.Param("id")

	if err := h.productService.Delete(c.Request.Context(), productID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// GetDigitalDownload returns a temporary download URL for a digital product.
// Authorization: user must have a completed order containing this product.
func (h *ProductHandler) GetDigitalDownload(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID := c.Param("id")
	product, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if !product.IsDigital {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not digital"})
		return
	}
	if product.DigitalFileURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Digital file not available"})
		return
	}

	if h.orderRepo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Order repository not available"})
		return
	}

	hasAccess, err := h.orderRepo.BuyerHasCompletedOrderWithProduct(c.Request.Context(), userID, product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate purchase"})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	expiry := 15 * time.Minute
	if h.storageService != nil {
		url, err := h.storageService.PresignGet(c.Request.Context(), product.DigitalFileURL, expiry)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"url": url, "expires_in_seconds": int(expiry.Seconds())})
			return
		}
	}

	// Fallback to stored URL (still protected by this endpoint)
	c.JSON(http.StatusOK, gin.H{"url": product.DigitalFileURL})
}
