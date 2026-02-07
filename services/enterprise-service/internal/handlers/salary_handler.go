package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/services"
	"github.com/gin-gonic/gin"
)

type SalaryHandler struct {
	payrollService *services.PayrollService
	employeeService *services.EmployeeService
}

func NewSalaryHandler(payrollService *services.PayrollService, employeeService *services.EmployeeService) *SalaryHandler {
	return &SalaryHandler{payrollService: payrollService, employeeService: employeeService}
}

// GetMySalary returns payroll payment history for the authenticated user in a given enterprise.
// GET /api/v1/employees/me/salary?enterprise_id=...
func (h *SalaryHandler) GetMySalary(c *gin.Context) {
	enterpriseID := c.Query("enterprise_id")
	if enterpriseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "enterprise_id query parameter required"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	employee, err := h.employeeService.GetEmployeeByUserAndEnterprise(c.Request.Context(), userID, enterpriseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found in this enterprise"})
		return
	}

	run, err := h.payrollService.PreparePayroll(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute salary"})
		return
	}

	// History: derived from completed payroll runs
	runs, err := h.payrollService.ListAllPayrollRuns(c.Request.Context(), enterpriseID)
	if err != nil {
		// If no runs yet, still return salary config + empty payments
		runs = nil
	}

	payments := make([]gin.H, 0)
	for _, r := range runs {
		for _, d := range r.Details {
			if d.EmployeeID == employee.ID {
				payments = append(payments, gin.H{
					"period_month": r.PeriodMonth,
					"period_year":  r.PeriodYear,
					"net_pay":      d.NetPay,
					"executed_at":   r.ExecutedAt,
					"status":       d.Status,
					"currency":     "XOF",
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"salary_config": gin.H{
			"base_amount":  employee.SalaryConfig.BaseAmount,
			"net_payable":  employee.SalaryConfig.NetPayable,
			"frequency":    employee.SalaryConfig.Frequency,
			"bonuses":      employee.SalaryConfig.Bonuses,
			"deductions":   employee.SalaryConfig.Deductions,
		},
		"payments": payments,
		"preview": gin.H{
			"period_month": run.PeriodMonth,
			"period_year":  run.PeriodYear,
			"net_pay":      func() float64 {
				for _, d := range run.Details {
					if d.EmployeeID == employee.ID {
						return d.NetPay
					}
				}
				return 0
			}(),
			"currency": "XOF",
		},
	})
}
