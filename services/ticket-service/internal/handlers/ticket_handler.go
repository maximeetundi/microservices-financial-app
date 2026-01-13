package handlers

import (
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/services"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service *services.TicketService
}

func NewTicketHandler(service *services.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

// === Event Handlers ===

// CreateEvent creates a new event
func (h *TicketHandler) CreateEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

// GetEvent returns an event by ID
func (h *TicketHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")

	event, err := h.service.GetEvent(eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

// GetEventByCode returns an event by its code
func (h *TicketHandler) GetEventByCode(c *gin.Context) {
	code := c.Param("code")

	event, err := h.service.GetEventByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

// GetMyEvents returns events created by the current user
func (h *TicketHandler) GetMyEvents(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.service.GetMyEvents(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetActiveEvents returns all active events
func (h *TicketHandler) GetActiveEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.service.GetActiveEvents(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// UpdateEvent updates an event
func (h *TicketHandler) UpdateEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	var req models.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.service.UpdateEvent(eventID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"event":   event,
	})
}

// PublishEvent publishes an event
func (h *TicketHandler) PublishEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	if err := h.service.PublishEvent(eventID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event published successfully"})
}

// DeleteEvent deletes an event
func (h *TicketHandler) DeleteEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	if err := h.service.DeleteEvent(eventID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// GetEventStats returns statistics for an event
func (h *TicketHandler) GetEventStats(c *gin.Context) {
	eventID := c.Param("id")

	stats, err := h.service.GetEventStats(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetEventTickets returns sold tickets for an event
func (h *TicketHandler) GetEventTickets(c *gin.Context) {
	eventID := c.Param("id")
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tickets, err := h.service.GetEventTickets(eventID, search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}

// === Ticket Handlers ===

// PurchaseTicket purchases a ticket
func (h *TicketHandler) PurchaseTicket(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.PurchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := c.GetHeader("Authorization")
	ticket, err := h.service.PurchaseTicket(userID, &req, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ticket purchased successfully",
		"ticket":  ticket,
	})
}

// GetMyTickets returns tickets owned by the current user
func (h *TicketHandler) GetMyTickets(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tickets, err := h.service.GetMyTickets(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}

// GetTicket returns a ticket by ID
func (h *TicketHandler) GetTicket(c *gin.Context) {
	ticketID := c.Param("id")

	ticket, err := h.service.GetTicket(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

// VerifyTicket verifies a ticket by code or QR data
func (h *TicketHandler) VerifyTicket(c *gin.Context) {
	var req models.VerifyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := req.TicketCode
	if code == "" {
		code = req.QRData
	}

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_code or qr_data is required"})
		return
	}

	result, err := h.service.VerifyTicket(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UseTicket marks a ticket as used
func (h *TicketHandler) UseTicket(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ticketID := c.Param("id")

	if err := h.service.UseTicket(ticketID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket marked as used"})
}

// RefundTicket refunds a ticket
// RefundTicket refunds a ticket
func (h *TicketHandler) RefundTicket(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ticketID := c.Param("id")

	var req models.RefundTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason is required: " + err.Error()})
		return
	}

	if err := h.service.RefundTicket(ticketID, userID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket refunded successfully"})
}

// CancelEvent cancels an event and refunds all tickets
func (h *TicketHandler) CancelEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eventID := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}

	if err := h.service.CancelAndRefundEvent(eventID, userID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event cancelled and refunds initiated"})
}

// GetAvailableIcons returns available icons for ticket tiers
func (h *TicketHandler) GetAvailableIcons(c *gin.Context) {
	icons := h.service.GetAvailableIcons()
	c.JSON(http.StatusOK, gin.H{"icons": icons})
}

// HealthCheck returns service health status
func (h *TicketHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "ticket-service",
	})
}
