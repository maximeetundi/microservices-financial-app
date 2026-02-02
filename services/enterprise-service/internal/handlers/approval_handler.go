package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/services"
	"github.com/gin-gonic/gin"
)

type ApprovalHandler struct {
	approvalService *services.ApprovalService
	employeeService *services.EmployeeService
	authClient      *services.AuthClient
}

func NewApprovalHandler(approvalService *services.ApprovalService, employeeService *services.EmployeeService, authClient *services.AuthClient) *ApprovalHandler {
	return &ApprovalHandler{
		approvalService: approvalService,
		employeeService: employeeService,
		authClient:      authClient,
	}
}

// GetPendingApprovals returns all pending approvals for an enterprise
// GET /enterprises/:id/approvals
func (h *ApprovalHandler) GetPendingApprovals(c *gin.Context) {
	enterpriseID := c.Param("id")

	approvals, err := h.approvalService.GetPendingApprovals(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, approvals)
}

// GetApproval returns a single approval by ID
// GET /approvals/:id
func (h *ApprovalHandler) GetApproval(c *gin.Context) {
	approvalID := c.Param("id")

	approval, err := h.approvalService.GetApprovalByID(c.Request.Context(), approvalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approval not found"})
		return
	}

	c.JSON(http.StatusOK, approval)
}

// InitiateAction creates a new action that requires multi-approval
// POST /enterprises/:id/actions
func (h *ApprovalHandler) InitiateAction(c *gin.Context) {
	enterpriseID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		ActionType  string                 `json:"action_type" binding:"required"`
		ActionName  string                 `json:"action_name" binding:"required"`
		Description string                 `json:"description"`
		Payload     map[string]interface{} `json:"payload"`
		Amount      float64                `json:"amount"`
		Currency    string                 `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get employee who is initiating
	employee, err := h.employeeService.GetEmployeeByUserAndEnterprise(c.Request.Context(), userID, enterpriseID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized for this enterprise"})
		return
	}

	// Check if action requires approval
	requiresApproval, err := h.approvalService.RequiresApproval(c.Request.Context(), enterpriseID, req.ActionType, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !requiresApproval {
		// Action can be executed directly
		c.JSON(http.StatusOK, gin.H{
			"requires_approval": false,
			"message":           "Action does not require approval, proceed directly",
		})
		return
	}

	// Create approval request
	approval, err := h.approvalService.InitiateAction(
		c.Request.Context(),
		enterpriseID,
		employee,
		req.ActionType,
		req.ActionName,
		req.Description,
		req.Payload,
		req.Amount,
		req.Currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"requires_approval": true,
		"approval":          approval,
		"message":           "Action requires approval from other administrators",
	})
}

// ApproveAction adds current admin's approval to an action
// POST /approvals/:id/approve
func (h *ApprovalHandler) ApproveAction(c *gin.Context) {
	approvalID := c.Param("id")
	userID := c.GetString("user_id")
	authToken := c.GetHeader("Authorization")
	if len(authToken) > 7 && authToken[:7] == "Bearer " {
		authToken = authToken[7:]
	}

	var req struct {
		PIN string `json:"pin" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify PIN via AuthClient
	if h.authClient != nil {
		valid, err := h.authClient.VerifyPin(userID, req.PIN, authToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "PIN verification failed: " + err.Error()})
			return
		}
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
			return
		}
	}

	// Get the approval to find the enterprise ID
	approval, err := h.approvalService.GetApprovalByID(c.Request.Context(), approvalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approval not found"})
		return
	}

	// Get employee
	employee, err := h.employeeService.GetEmployeeByUserAndEnterprise(c.Request.Context(), userID, approval.EnterpriseID.Hex())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized for this enterprise"})
		return
	}

	// Check if employee is admin
	if !employee.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only administrators can approve actions"})
		return
	}

	// Add approval
	err = h.approvalService.ApproveAction(c.Request.Context(), approvalID, employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if action is now fully approved
	updatedApproval, _ := h.approvalService.GetApprovalByID(c.Request.Context(), approvalID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Approval recorded",
		"approval": updatedApproval,
		"status":   updatedApproval.Status,
	})
}

// RejectAction adds current admin's rejection to an action
// POST /approvals/:id/reject
func (h *ApprovalHandler) RejectAction(c *gin.Context) {
	approvalID := c.Param("id")
	userID := c.GetString("user_id")
	authToken := c.GetHeader("Authorization")
	if len(authToken) > 7 && authToken[:7] == "Bearer " {
		authToken = authToken[7:]
	}

	var req struct {
		PIN    string `json:"pin" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify PIN via AuthClient
	if h.authClient != nil {
		valid, err := h.authClient.VerifyPin(userID, req.PIN, authToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "PIN verification failed: " + err.Error()})
			return
		}
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
			return
		}
	}

	// Get the approval to find the enterprise ID
	approval, err := h.approvalService.GetApprovalByID(c.Request.Context(), approvalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approval not found"})
		return
	}

	// Get employee
	employee, err := h.employeeService.GetEmployeeByUserAndEnterprise(c.Request.Context(), userID, approval.EnterpriseID.Hex())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized for this enterprise"})
		return
	}

	// Check if employee is admin
	if !employee.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only administrators can reject actions"})
		return
	}

	// Add rejection
	err = h.approvalService.RejectAction(c.Request.Context(), approvalID, employee, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action rejected"})
}
