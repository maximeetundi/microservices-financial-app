package handlers

import (
	"net/http"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

// === Emergency Fund Endpoints ===

// CreateEmergencyFund creates a new emergency fund for an association
func (h *ExtendedHandler) CreateEmergencyFund(c *gin.Context) {
	userID := middleware.GetUserID(c)
	associationID := c.Param("id")

	var req struct {
		MonthlyContribution float64 `json:"monthly_contribution" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Verify user is admin of association

	query := `
		INSERT INTO emergency_funds (association_id, monthly_contribution, balance)
		VALUES ($1, $2, 0)
		ON CONFLICT (association_id) DO UPDATE 
		SET monthly_contribution = $2
		RETURNING id, association_id, balance, monthly_contribution, created_at
	`

	var fund struct {
		ID                  string  `json:"id"`
		AssociationID       string  `json:"association_id"`
		Balance             float64 `json:"balance"`
		MonthlyContribution float64 `json:"monthly_contribution"`
		CreatedAt           string  `json:"created_at"`
	}

	err := h.db.QueryRow(query, associationID, req.MonthlyContribution).Scan(
		&fund.ID, &fund.AssociationID, &fund.Balance, &fund.MonthlyContribution, &fund.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, fund)
}

// GetEmergencyFund gets emergency fund details
func (h *ExtendedHandler) GetEmergencyFund(c *gin.Context) {
	associationID := c.Param("id")

	query := `
		SELECT id, association_id, balance, monthly_contribution, created_at
		FROM emergency_funds WHERE association_id = $1
	`

	var fund struct {
		ID                  string  `json:"id"`
		AssociationID       string  `json:"association_id"`
		Balance             float64 `json:"balance"`
		MonthlyContribution float64 `json:"monthly_contribution"`
		CreatedAt           string  `json:"created_at"`
	}

	err := h.db.QueryRow(query, associationID).Scan(
		&fund.ID, &fund.AssociationID, &fund.Balance, &fund.MonthlyContribution, &fund.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emergency fund not found"})
		return
	}

	c.JSON(http.StatusOK, fund)
}

// ContributeToEmergencyFund - member pays monthly contribution
func (h *ExtendedHandler) ContributeToEmergencyFund(c *gin.Context) {
	userID := middleware.GetUserID(c)
	associationID := c.Param("id")

	var req struct {
		Period   string  `json:"period" binding:"required"`
		WalletID string  `json:"wallet_id" binding:"required"`
		Pin      string  `json:"pin" binding:"required"`
		Amount   float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get emergency fund
	var fundID string
	var monthlyAmount float64
	err := h.db.QueryRow(`SELECT id, monthly_contribution FROM emergency_funds WHERE association_id = $1`, associationID).
		Scan(&fundID, &monthlyAmount)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emergency fund not found"})
		return
	}

	// Verify amount matches
	if req.Amount != monthlyAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must match monthly contribution"})
		return
	}

	// Get member ID
	var memberID string
	err = h.db.QueryRow(`SELECT id FROM members WHERE association_id = $1 AND user_id = $2`, associationID, userID).Scan(&memberID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	// TODO: Call wallet service to verify PIN and deduct amount

	// Record contribution
	_, err = h.db.Exec(`
		INSERT INTO emergency_contributions (emergency_fund_id, member_id, amount, period, paid, paid_at)
		VALUES ($1, $2, $3, $4, true, NOW())
		ON CONFLICT (emergency_fund_id, member_id, period) DO UPDATE
		SET paid = true, paid_at = NOW()
	`, fundID, memberID, req.Amount, req.Period)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update fund balance
	_, err = h.db.Exec(`UPDATE emergency_funds SET balance = balance + $1 WHERE id = $2`, req.Amount, fundID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contribution recorded", "amount": req.Amount})
}

// RequestEmergencyWithdrawal - request withdrawal from emergency fund (creates approval request)
func (h *ExtendedHandler) RequestEmergencyWithdrawal(c *gin.Context) {
	userID := middleware.GetUserID(c)
	associationID := c.Param("id")

	var req struct {
		EventType     string  `json:"event_type" binding:"required"`
		BeneficiaryID string  `json:"beneficiary_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		Reason        string  `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get emergency fund
	var fundID string
	var balance float64
	err := h.db.QueryRow(`SELECT id, balance FROM emergency_funds WHERE association_id = $1`, associationID).
		Scan(&fundID, &balance)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emergency fund not found"})
		return
	}

	if balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient emergency fund balance"})
		return
	}

	// Create approval request (4/5 required)
	var approvalID string
	err = h.db.QueryRow(`
		INSERT INTO approval_requests (association_id, request_type, amount, metadata, status, required_approvals, current_approvals, created_by)
		VALUES ($1, 'emergency_withdrawal', $2, $3::jsonb, 'pending', 4, 0, $4)
		RETURNING id
	`, associationID, req.Amount, `{"event_type": "`+req.EventType+`", "reason": "`+req.Reason+`"}`, userID).Scan(&approvalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create withdrawal record
	var withdrawalID string
	err = h.db.QueryRow(`
		INSERT INTO emergency_withdrawals (emergency_fund_id, event_type, beneficiary_id, amount, reason, status, approval_request_id, created_by)
		VALUES ($1, $2, $3, $4, $5, 'pending', $6, $7)
		RETURNING id
	`, fundID, req.EventType, req.BeneficiaryID, req.Amount, req.Reason, approvalID, userID).Scan(&withdrawalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"withdrawal_id": withdrawalID,
		"approval_id":   approvalID,
		"message":       "Withdrawal request created, waiting for 4/5 approvals",
	})
}

// GetEmergencyWithdrawals - list all emergency withdrawals
func (h *ExtendedHandler) GetEmergencyWithdrawals(c *gin.Context) {
	associationID := c.Param("id")

	rows, err := h.db.Query(`
		SELECT ew.id, ew.event_type, ew.beneficiary_id, m.user_id as beneficiary_name,
		       ew.amount, ew.reason, ew.status, ew.created_at
		FROM emergency_withdrawals ew
		JOIN emergency_funds ef ON ew.emergency_fund_id = ef.id
		JOIN members m ON ew.beneficiary_id = m.id
		WHERE ef.association_id = $1
		ORDER BY ew.created_at DESC
	`, associationID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	withdrawals := []map[string]interface{}{}
	for rows.Next() {
		var w struct {
			ID              string
			EventType       string
			BeneficiaryID   string
			BeneficiaryName string
			Amount          float64
			Reason          string
			Status          string
			CreatedAt       string
		}

		rows.Scan(&w.ID, &w.EventType, &w.BeneficiaryID, &w.BeneficiaryName, &w.Amount, &w.Reason, &w.Status, &w.CreatedAt)

		withdrawals = append(withdrawals, map[string]interface{}{
			"id":               w.ID,
			"event_type":       w.EventType,
			"beneficiary_id":   w.BeneficiaryID,
			"beneficiary_name": w.BeneficiaryName,
			"amount":           w.Amount,
			"reason":           w.Reason,
			"status":           w.Status,
			"created_at":       w.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, withdrawals)
}
