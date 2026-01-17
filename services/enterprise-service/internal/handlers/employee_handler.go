package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/services"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

// InviteEmployee (Point 2)
func (h *EmployeeHandler) InviteEmployee(c *gin.Context) {
	var req models.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.InviteEmployee(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite employee", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Invitation sent", "employee_id": req.ID})
}

// AcceptInvitation (Point 2: PIN Validation)
func (h *EmployeeHandler) AcceptInvitation(c *gin.Context) {
	var req struct {
		EmployeeID string `json:"employee_id"`
		Pin        string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ConfirmEmployee(c.Request.Context(), req.EmployeeID, req.Pin); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted"})
}

// PromoteEmployee (Point 4)
func (h *EmployeeHandler) PromoteEmployee(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NewRole    string                  `json:"new_role"`
		NewSalary  *models.SalaryConfig    `json:"new_salary"`
		Permissions *models.AdminPermission `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PromoteEmployee(c.Request.Context(), id, req.NewRole, req.NewSalary, req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee promoted successfully"})
}

// ListEmployees (For Dashboard)
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	id := c.Param("id")
	employees, err := h.service.ListByEnterprise(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list employees"})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// TerminateEmployee (Point 4)
func (h *EmployeeHandler) TerminateEmployee(c *gin.Context) {
	enterpriseID := c.Param("id")
	employeeID := c.Param("employeeId")
	
	if err := h.service.TerminateEmployee(c.Request.Context(), enterpriseID, employeeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to terminate employee", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee terminated successfully"})
}

// GetMyEmployee returns the current user's employee record for a specific enterprise
// This is used for permission-based UI filtering
func (h *EmployeeHandler) GetMyEmployee(c *gin.Context) {
	enterpriseID := c.Query("enterprise_id")
	if enterpriseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "enterprise_id query parameter required"})
		return
	}
	
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	employee, err := h.service.GetEmployeeByUserAndEnterprise(c.Request.Context(), userID.(string), enterpriseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found in this enterprise"})
		return
	}
	
	c.JSON(http.StatusOK, employee)
}

