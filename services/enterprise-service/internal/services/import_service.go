package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"strconv"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
)

// ImportInvoiceRequest (Point 8)
type ImportInvoiceRequest struct {
	EnterpriseID string            `json:"enterprise_id"`
	ServiceID    string            `json:"service_id"`
	FileContent  []byte            `json:"-"` // Handled via Multipart
	ColumnMap    map[string]string `json:"column_map"` // "phone" -> "2" (index) or "B" (excel col)
	FileType     string            `json:"file_type"`    // CSV, EXCEL
}

// ParseAndGenerateInvoices (Point 9)
// Dynamically maps columns to Invoice fields
func (s *BillingService) ParseAndGenerateInvoices(ctx context.Context, req ImportInvoiceRequest) ([]models.Invoice, error) {
	if req.FileType == "CSV" {
		return s.parseCSV(ctx, req)
	}
	// TODO: Implement Excel parsing (needs external lib)
	return nil, errors.New("unsupported file type (only CSV implemented for MVP)")
}

func (s *BillingService) parseCSV(ctx context.Context, req ImportInvoiceRequest) ([]models.Invoice, error) {
	reader := csv.NewReader(bufio.NewReader(bytes.NewReader(req.FileContent)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Fetch Enterprise to resolve Pricing
	ent, err := s.entRepo.FindByID(ctx, req.EnterpriseID)
	if err != nil {
		return nil, errors.New("failed to find enterprise: " + err.Error())
	}

	// Resolve Service Definition
	var unitPrice float64 = 0
	billingType := "FIXED" // Default
	serviceName := "Service"

	foundService := false
	for _, group := range ent.ServiceGroups {
		for _, svc := range group.Services {
			if svc.ID == req.ServiceID {
				unitPrice = svc.BasePrice
				billingType = svc.BillingType
				serviceName = svc.Name
				foundService = true
				break
			}
		}
		if foundService {
			break
		}
	}
	
	// If not in custom services, check transport/school configs (optional)
	// For now, assuming imported files align with Generic Services or manually mapped amounts
	if !foundService {
		// Just log? Or assume 0 price?
	}

	var invoices []models.Invoice

	// Map expected fields
	colClientStr := req.ColumnMap["client_identifier"]
	colAmountStr := req.ColumnMap["amount"]
	colConsumpStr := req.ColumnMap["consumption"]
	
	colClient, _ := strconv.Atoi(colClientStr)
	
	hasAmount := colAmountStr != ""
	hasConsump := colConsumpStr != ""
	
	colAmount, _ := strconv.Atoi(colAmountStr)
	colConsump, _ := strconv.Atoi(colConsumpStr)

	// Iterate rows
	for i, row := range records {
		if i == 0 { continue } // Assume header

		if len(row) <= colClient { continue } // Skip invalid rows

		var amount float64
		var description string

		// Priority 1: Direct Amount in CSV
		if hasAmount && len(row) > colAmount && row[colAmount] != "" {
			parsedAmt, err := strconv.ParseFloat(row[colAmount], 64)
			if err == nil {
				amount = parsedAmt
				description = serviceName + " (Facture Directe)"
			}
		} 
		
		// Priority 2: Usage Calculation
		if amount == 0 && hasConsump && len(row) > colConsump && row[colConsump] != "" {
			consumption, err := strconv.ParseFloat(row[colConsump], 64)
			if err == nil && billingType == "USAGE" {
				amount = consumption * unitPrice
				description = serviceName + " : " + row[colConsump] + " unités"
			}
		}

		if amount <= 0 {
			continue // Skip if no valid amount calculated
		}

		inv := models.Invoice{
			ClientID:      row[colClient], 
			ClientName:    "Client Importé", // In future: lookup user profile
			ServiceID:     req.ServiceID,
			ServiceName:   serviceName,
			Amount:        amount,
			Currency:      "XOF", // Should come from Enterprise Settings
			Status:        models.InvoiceStatusDraft,
			Description:   description,
			EnterpriseID:  ent.ID,
			CreatedAt:     time.Now(), // Ensure created_at is set
		}
		invoices = append(invoices, inv)
	}

	return invoices, nil
}
