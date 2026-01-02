package handlers

import (
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.AssociationService
}

func NewHandler(service *services.AssociationService) *Handler {
	return &Handler{service: service}
}

// === Association Endpoints ===

func (h *Handler) CreateAssociation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	
	var req models.CreateAssociationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assoc, err := h.service.CreateAssociation(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, assoc)
}

func (h *Handler) GetAssociations(c *gin.Context) {
	userID := middleware.GetUserID(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	associations, err := h.service.GetMyAssociations(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, associations)
}

func (h *Handler) GetAssociation(c *gin.Context) {
	id := c.Param("id")

	assoc, err := h.service.GetAssociation(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Association not found"})
		return
	}

	c.JSON(http.StatusOK, assoc)
}

func (h *Handler) UpdateAssociation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req models.UpdateAssociationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assoc, err := h.service.UpdateAssociation(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assoc)
}

func (h *Handler) DeleteAssociation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	if err := h.service.DeleteAssociation(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Association deleted"})
}

// === Member Endpoints ===

func (h *Handler) JoinAssociation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req models.JoinAssociationRequest
	c.ShouldBindJSON(&req)

	member, err := h.service.JoinAssociation(id, userID, req.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

func (h *Handler) LeaveAssociation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	if err := h.service.LeaveAssociation(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left association"})
}

func (h *Handler) GetMembers(c *gin.Context) {
	id := c.Param("id")

	members, err := h.service.GetMembers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *Handler) UpdateMemberRole(c *gin.Context) {
	userID := middleware.GetUserID(c)
	associationID := c.Param("id")
	targetUserID := c.Param("uid")

	var req models.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMemberRole(associationID, targetUserID, userID, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

// === Contribution Endpoints ===

func (h *Handler) PayContribution(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req struct {
		WalletID    string  `json:"wallet_id" binding:"required"`
		Pin         string  `json:"pin" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Period      string  `json:"period" binding:"required"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.service.PayContribution(id, userID, req.WalletID, req.Pin, req.Amount, req.Period, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *Handler) GetTreasury(c *gin.Context) {
	id := c.Param("id")

	report, err := h.service.GetTreasuryReport(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *Handler) DistributeFunds(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req struct {
		RecipientMemberID string  `json:"recipient_member_id" binding:"required"`
		WalletID          string  `json:"wallet_id" binding:"required"`
		Amount            float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.service.DistributeFunds(id, userID, req.RecipientMemberID, req.WalletID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

// === Loan Endpoints ===

func (h *Handler) RequestLoan(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req models.LoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.service.RequestLoan(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

func (h *Handler) ApproveLoan(c *gin.Context) {
	userID := middleware.GetUserID(c)
	loanID := c.Param("loanId")

	var req struct {
		Approve  bool   `json:"approve"`
		Reason   string `json:"reason"`
		WalletID string `json:"wallet_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ApproveLoan(loanID, userID, req.WalletID, req.Approve, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan processed"})
}

func (h *Handler) RepayLoan(c *gin.Context) {
	userID := middleware.GetUserID(c)
	loanID := c.Param("loanId")

	var req struct {
		WalletID string  `json:"wallet_id" binding:"required"`
		Pin      string  `json:"pin" binding:"required"`
		Amount   float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RepayLoan(loanID, userID, req.WalletID, req.Pin, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repayment recorded"})
}

func (h *Handler) GetLoans(c *gin.Context) {
	id := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	loans, err := h.service.GetLoans(id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loans)
}

// === Meeting Endpoints ===

func (h *Handler) CreateMeeting(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req models.CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting, err := h.service.CreateMeeting(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, meeting)
}

func (h *Handler) GetMeetings(c *gin.Context) {
	id := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	meetings, err := h.service.GetMeetings(id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meetings)
}

func (h *Handler) RecordAttendance(c *gin.Context) {
	userID := middleware.GetUserID(c)
	meetingID := c.Param("meetingId")

	var req models.RecordAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RecordAttendance(meetingID, userID, req.Attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance recorded"})
}

func (h *Handler) UpdateMinutes(c *gin.Context) {
	userID := middleware.GetUserID(c)
	meetingID := c.Param("meetingId")

	var req models.UpdateMinutesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMinutes(meetingID, userID, req.Minutes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Minutes updated"})
}

// === Admin Endpoints ===

func (h *Handler) GetAllAssociations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	associations, err := h.service.GetAllAssociations(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, associations)
}
