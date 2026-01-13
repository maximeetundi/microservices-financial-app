package services

import (
	"context"
	"errors"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateBatchFromImport (Point 8, 9)
func (s *BillingService) CreateBatchFromImport(ctx context.Context, enterpriseID string, invoices []models.Invoice) (*models.InvoiceBatch, error) {
	entOID, _ := primitive.ObjectIDFromHex(enterpriseID)
	
	totalAmount := 0.0
	for _, inv := range invoices {
		totalAmount += inv.Amount
	}

	batch := &models.InvoiceBatch{
		EnterpriseID:  entOID,
		Name:          "Import " + time.Now().Format("2006-01-02 15:04"),
		Status:        "DRAFT",
		TotalInvoices: len(invoices),
		TotalAmount:   totalAmount,
		Invoices:      invoices,
	}

	if err := s.batchRepo.Create(ctx, batch); err != nil {
		return nil, err
	}
	return batch, nil
}

// ValidateBatch (Point 11: Validation)
func (s *BillingService) ValidateBatch(ctx context.Context, batchID string) error {
	batch, err := s.batchRepo.FindByID(ctx, batchID)
	if err != nil { return err }

	if batch.Status != "DRAFT" {
		return errors.New("only draft batches can be validated")
	}

	batch.Status = "VALIDATED"
	return s.batchRepo.Update(ctx, batch)
}

// ScheduleBatch (Point 10: Planification)
func (s *BillingService) ScheduleBatch(ctx context.Context, batchID string, sendAt time.Time) error {
	batch, err := s.batchRepo.FindByID(ctx, batchID)
	if err != nil { return err }

	if batch.Status != "VALIDATED" {
		return errors.New("batch must be validated before scheduling")
	}

	batch.Status = "SCHEDULED"
	batch.ScheduledAt = sendAt
	return s.batchRepo.Update(ctx, batch)
}

// ProcessScheduledBatches (Cron Job: Point 12)
func (s *BillingService) ProcessScheduledBatches(ctx context.Context) error {
	batches, err := s.batchRepo.FindPendingScheduled(ctx)
	if err != nil { return err }

	for _, batch := range batches {
		// Send Invoices
		for _, inv := range batch.Invoices {
			// Save individual invoice to main table
			inv.Status = models.InvoiceStatusSent
			inv.SentAt = time.Now()
			// Link batch ID?
			s.invRepo.Create(ctx, &inv)
			// TODO: Notify User "You have a new bill"
		}
		
		batch.Status = "PROCESSED"
		s.batchRepo.Update(ctx, &batch)
	}
	return nil
}
