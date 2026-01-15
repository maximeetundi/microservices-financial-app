package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// ExportHandler handles enterprise data export
type ExportHandler struct {
	entRepo  *repository.EnterpriseRepository
	empRepo  *repository.EmployeeRepository
	subRepo  *repository.SubscriptionRepository
	invRepo  *repository.InvoiceRepository
	payRepo  *repository.PayrollRepository
}

func NewExportHandler(
	entRepo *repository.EnterpriseRepository,
	empRepo *repository.EmployeeRepository,
	subRepo *repository.SubscriptionRepository,
	invRepo *repository.InvoiceRepository,
	payRepo *repository.PayrollRepository,
) *ExportHandler {
	return &ExportHandler{
		entRepo: entRepo,
		empRepo: empRepo,
		subRepo: subRepo,
		invRepo: invRepo,
		payRepo: payRepo,
	}
}

// EnterpriseExport represents the complete enterprise data export
type EnterpriseExport struct {
	ExportedAt     time.Time                `json:"exported_at"`
	ExportVersion  string                   `json:"export_version"`
	Enterprise     *models.Enterprise       `json:"enterprise"`
	Employees      []models.Employee        `json:"employees"`
	Subscriptions  []models.Subscription    `json:"subscriptions"`
	Invoices       []models.Invoice         `json:"invoices"`
	PayrollRuns    []models.PayrollRun      `json:"payroll_runs"`
	SecurityPolicies []models.SecurityPolicy `json:"security_policies"`
	ServiceGroups  []models.ServiceGroup    `json:"service_groups"`
	Metadata       ExportMetadata           `json:"metadata"`
}

type ExportMetadata struct {
	TotalEmployees    int     `json:"total_employees"`
	ActiveEmployees   int     `json:"active_employees"`
	TotalSubscriptions int    `json:"total_subscriptions"`
	TotalInvoices     int     `json:"total_invoices"`
	TotalRevenue      float64 `json:"total_revenue"`
	TotalPayrollRuns  int     `json:"total_payroll_runs"`
}

// ExportEnterpriseData exports all enterprise data as JSON
// GET /enterprises/:id/export
func (h *ExportHandler) ExportEnterpriseData(c *gin.Context) {
	enterpriseID := c.Param("id")
	userID := c.GetString("user_id")

	// Get enterprise
	ent, err := h.entRepo.FindByID(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}

	// Verify ownership (only owner can export)
	if ent.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the owner can export enterprise data"})
		return
	}

	// Gather all data
	employees, _ := h.empRepo.FindByEnterprise(c.Request.Context(), enterpriseID)
	subscriptions, _ := h.subRepo.FindByEnterprise(c.Request.Context(), enterpriseID)
	invoices, _ := h.invRepo.FindByEnterprise(c.Request.Context(), enterpriseID)
	payrollRuns, _ := h.payRepo.FindByEnterprise(c.Request.Context(), enterpriseID)

	// Calculate metadata
	activeEmps := 0
	for _, emp := range employees {
		if emp.Status == models.EmployeeStatusActive {
			activeEmps++
		}
	}

	totalRevenue := 0.0
	for _, inv := range invoices {
		if inv.Status == "PAID" {
			totalRevenue += inv.Amount
		}
	}

	export := EnterpriseExport{
		ExportedAt:       time.Now(),
		ExportVersion:    "1.0",
		Enterprise:       ent,
		Employees:        employees,
		Subscriptions:    subscriptions,
		Invoices:         invoices,
		PayrollRuns:      payrollRuns,
		SecurityPolicies: ent.SecurityPolicies,
		ServiceGroups:    ent.ServiceGroups,
		Metadata: ExportMetadata{
			TotalEmployees:     len(employees),
			ActiveEmployees:    activeEmps,
			TotalSubscriptions: len(subscriptions),
			TotalInvoices:      len(invoices),
			TotalRevenue:       totalRevenue,
			TotalPayrollRuns:   len(payrollRuns),
		},
	}

	// Mark enterprise as exported (add timestamp)
	ent.LastExportedAt = time.Now()
	h.entRepo.Update(c.Request.Context(), ent)

	// Return as JSON file download
	exportJSON, _ := json.MarshalIndent(export, "", "  ")
	
	filename := ent.Name + "_export_" + time.Now().Format("2006-01-02") + ".json"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", exportJSON)
}

// CheckExportStatus returns whether the enterprise has been exported recently
// GET /enterprises/:id/export/status
func (h *ExportHandler) CheckExportStatus(c *gin.Context) {
	enterpriseID := c.Param("id")

	ent, err := h.entRepo.FindByID(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}

	hasRecentExport := false
	if !ent.LastExportedAt.IsZero() {
		// Export is valid for 30 days
		hasRecentExport = time.Since(ent.LastExportedAt) < 30*24*time.Hour
	}

	c.JSON(http.StatusOK, gin.H{
		"has_export":       !ent.LastExportedAt.IsZero(),
		"has_recent_export": hasRecentExport,
		"last_exported_at":  ent.LastExportedAt,
		"can_delete":        hasRecentExport,
	})
}
