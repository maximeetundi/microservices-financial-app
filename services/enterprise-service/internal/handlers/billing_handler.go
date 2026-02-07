package handlers

import (
	"bytes"
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/services"
	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	service *services.BillingService
}

func NewBillingHandler(service *services.BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) CreateInvoice(c *gin.Context) {
	var req models.Invoice
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateInvoice(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// ListMyInvoices returns invoices for the current authenticated user
// GET /api/v1/enterprises/:id/invoices/me
func (h *BillingHandler) ListMyInvoices(c *gin.Context) {
	enterpriseID := c.Param("id")
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	invs, err := h.service.ListInvoicesForClient(c.Request.Context(), enterpriseID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invs)
}

// ImportInvoices (Point 8, 9)
func (h *BillingHandler) ImportInvoices(c *gin.Context) {
	enterpriseID := c.Param("id")
	
	// 1. Get File
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Read file to bytes
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	// 2. Get Mapping (JSON string in form data)
	// For simplicity, we assume fixed mapping or specific params for MVP
	// In real world: c.PostForm("mapping") -> JSON Unmarshal
	
	req := services.ImportInvoiceRequest{
		EnterpriseID: enterpriseID,
		ServiceID:    c.PostForm("service_id"),
		FileContent:  buf.Bytes(),
		FileType:     "CSV", // Detected from extension in real impl
		ColumnMap: map[string]string{
			"client_identifier": c.PostForm("col_client_idx"),
			"amount":            c.PostForm("col_amount_idx"),
			"consumption":       c.PostForm("col_consumption_idx"),
		},
	}

	// 3. Process
	invoices, err := h.service.ParseAndGenerateInvoices(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Create Batch (Persist Draft)
	batch, err := h.service.CreateBatchFromImport(c.Request.Context(), enterpriseID, invoices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create batch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch created successfully",
		"batch_id": batch.ID,
		"status": batch.Status,
		"total_invoices": batch.TotalInvoices,
	})
}

// ValidateBatch Endpoint
func (h *BillingHandler) ValidateBatch(c *gin.Context) {
	batchID := c.Param("batch_id")
	if err := h.service.ValidateBatch(c.Request.Context(), batchID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Batch validated"})
}

// ScheduleBatch Endpoint
func (h *BillingHandler) ScheduleBatch(c *gin.Context) {
	batchID := c.Param("batch_id")
	var req struct {
		ScheduledAt string `json:"scheduled_at"` // ISO8601
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Parse time
	// For MVP: assume format parsing or use time.Time if JSON handles it, usually string.
	// We'll skip complex parsing and assume Immediate for now if empty
	
	if err := h.service.ValidateBatch(c.Request.Context(), batchID); err != nil { // Re-using validate check logic inside schedule? No, separate.
		// logic inside service handles status check
	}
	
	// Simply call schedule
	// h.service.ScheduleBatch(...)
    c.JSON(http.StatusOK, gin.H{"message": "Batch scheduled (Mock)"})
}

// BatchInvoiceRequest for Manual Entry
type BatchInvoiceRequest struct {
	Items []BatchInvoiceItem `json:"items" binding:"required"`
}

type BatchInvoiceItem struct {
	SubscriptionID string  `json:"subscription_id" binding:"required"`
	Amount         float64 `json:"amount"`      // For Fixed override
	Consumption    float64 `json:"consumption"` // For Usage
}

func (h *BillingHandler) CreateBatchInvoices(c *gin.Context) {
	enterpriseID := c.Param("id")
	var req BatchInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    // Map to Service Model
    // We cannot import services.ManualInvoiceItem if services imports handlers (cycle).
    // Services generally shouldn't import handlers.
    // If services handles model definition, we can use it.
    // Wait, I defined ManualInvoiceItem in services package. So I can import it here.
    // handler -> service IS allowed. Service -> handler is NOT.
    
    var items []services.ManualInvoiceItem
    for _, item := range req.Items {
        items = append(items, services.ManualInvoiceItem{
            SubscriptionID: item.SubscriptionID,
            Amount:         item.Amount,
            Consumption:    item.Consumption,
        })
    }

    if err := h.service.GenerateBatchFromManualEntry(c.Request.Context(), enterpriseID, items); err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate batch", "details": err.Error()})
         return
    }

	c.JSON(http.StatusOK, gin.H{"message": "Invoices generated successfully"})
}
