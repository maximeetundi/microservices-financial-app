package handlers

import (
	"net/http"

	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionHandler struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

// CreateSubscription POST /api/v1/enterprises/:id/subscriptions
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	enterpriseID := c.Param("id")
	var req models.Subscription
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Enterprise ID
	entOID, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Enterprise ID"})
		return
	}
	req.EnterpriseID = entOID
	req.Status = "ACTIVE"
	
	// Set NextBilling based on Frequency if not provided?
	// For now assume client sends it or we default to Tomorrow/Next Month.
	// We'll trust the input for MVP or default to Now.
	if req.NextBillingAt.IsZero() {
		req.NextBillingAt = time.Now().AddDate(0, 1, 0) // Default 1 month
	}

	if err := h.repo.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// ListSubscriptions GET /api/v1/enterprises/:id/subscriptions
// Optional Query Param: service_id
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	enterpriseID := c.Param("id")
	serviceID := c.Query("service_id")

	subs, err := h.repo.FindByEnterprise(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions", "details": err.Error()})
		return
	}

	// Filter by ServiceID if provided
	if serviceID != "" {
		var filtered []interface{}
		for _, sub := range subs {
			if sub.ServiceID == serviceID {
				filtered = append(filtered, sub)
			}
		}
		c.JSON(http.StatusOK, filtered)
		return
	}

	c.JSON(http.StatusOK, subs)
}
