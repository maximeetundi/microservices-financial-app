package services

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"strconv"

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
func (s *BillingService) ParseAndGenerateInvoices(req ImportInvoiceRequest) ([]models.Invoice, error) {
	if req.FileType == "CSV" {
		return s.parseCSV(req)
	}
	// TODO: Implement Excel parsing (needs external lib)
	return nil, errors.New("unsupported file type (only CSV implemented for MVP)")
}

func (s *BillingService) parseCSV(req ImportInvoiceRequest) ([]models.Invoice, error) {
	reader := csv.NewReader(bufio.NewReader(bytes.NewReader(req.FileContent)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var invoices []models.Invoice

	// Map expected fields: "client_identifier", "amount", "consumption"
	// User provides: "client_identifier" -> "0" (Column 0), "amount" -> "2" (Column 2) OR "consumption" -> "3" (Column 3)
	
	colClientStr := req.ColumnMap["client_identifier"]
	colAmountStr := req.ColumnMap["amount"]
	colConsumpStr := req.ColumnMap["consumption"]
	
	colClient, _ := strconv.Atoi(colClientStr)
	
	// Check if we have explicit amount or consumption
	hasAmount := colAmountStr != ""
	hasConsump := colConsumpStr != ""
	
	colAmount, _ := strconv.Atoi(colAmountStr)
	colConsump, _ := strconv.Atoi(colConsumpStr)

	// Fetch Service Definition if consumption based (need UnitPrice)
	var unitPrice float64 = 0
	if hasConsump && !hasAmount {
		// Mock fetching service definition. In real world: s.entRepo.FindServiceDefinition(req.EnterpriseID, req.ServiceID)
		// For MVP, we assume UnitPrice 1 or provided in request? 
		// Ideally, we'd lookup the service definition from the Enterprise model here.
		// For now, let's assume a passed-in UnitPrice or default.
		unitPrice = 1.0 // Placeholder
	}

	// Iterate rows
	for i, row := range records {
		if i == 0 { continue } // Assume header

		if len(row) <= colClient { continue }

		var amount float64
		var description string

		if hasAmount && len(row) > colAmount {
			amount, _ = strconv.ParseFloat(row[colAmount], 64)
			description = "Direct Amount Bill"
		} else if hasConsump && len(row) > colConsump {
			consumption, _ := strconv.ParseFloat(row[colConsump], 64)
			amount = consumption * unitPrice
			description = "Consumption Bill: " + row[colConsump] + " units"
		} else {
			continue // Cannot calculate
		}

		inv := models.Invoice{
			ClientID:      row[colClient], 
			ClientName:    "Imported Client",
			ServiceID:     req.ServiceID,
			Amount:        amount,
			Status:        models.InvoiceStatusDraft,
			Description:   description,
		}
		invoices = append(invoices, inv)
	}

	return invoices, nil
}
