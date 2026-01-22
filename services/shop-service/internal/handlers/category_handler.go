package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/services"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// ListByShop returns categories for a shop
func (h *CategoryHandler) ListByShop(c *gin.Context) {
	shopSlug := c.Param("slug")

	categories, err := h.categoryService.ListByShop(c.Request.Context(), shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// ListWithHierarchy returns categories with nested children
func (h *CategoryHandler) ListWithHierarchy(c *gin.Context) {
	shopSlug := c.Param("slug")

	categories, err := h.categoryService.ListWithHierarchy(c.Request.Context(), shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// Get returns a category by ID
func (h *CategoryHandler) Get(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Create creates a new category
func (h *CategoryHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Update updates a category
func (h *CategoryHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	categoryID := c.Param("id")

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), categoryID, &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete deletes a category
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	categoryID := c.Param("id")

	if err := h.categoryService.Delete(c.Request.Context(), categoryID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
