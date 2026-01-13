package services

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
)

type BillingService struct {
	subRepo   *repository.SubscriptionRepository
	invRepo   *repository.InvoiceRepository
	batchRepo *repository.BatchRepository
}

func NewBillingService(subRepo *repository.SubscriptionRepository, invRepo *repository.InvoiceRepository, batchRepo *repository.BatchRepository) *BillingService {
	return &BillingService{
		subRepo:   subRepo,
		invRepo:   invRepo,
		batchRepo: batchRepo,
	}
}

// CreateInvoice (Single manual or system generated)
func (s *BillingService) CreateInvoice(ctx context.Context, inv *models.Invoice) error {
	inv.Status = models.InvoiceStatusDraft
	return s.invRepo.Create(ctx, inv)
}

// ProcessRecurringBilling (Cron Job)
// Scans active subscriptions and generates Invoices
func (s *BillingService) ProcessRecurringBilling(ctx context.Context) error {
	subs, err := s.subRepo.FindDueSubscriptions(ctx)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		// Generate Invoice
		inv := models.Invoice{
			EnterpriseID:  sub.EnterpriseID,
			ClientID:      sub.ClientID,
			ClientName:    "Subscriber", // Ideally fetch User details
			ServiceID:     sub.ServiceID,
			Amount:        sub.Amount,
			DueDate:       time.Now().AddDate(0, 0, 7), // Due in 7 days
			Status:        models.InvoiceStatusSent, // Auto-sent
			SentAt:        time.Now(),
		}
		
		if err := s.invRepo.Create(ctx, &inv); err != nil {
			// Log error but continue
			continue
		}

		// Update Subscription NextBilling
		now := time.Now()
		sub.LastBillingAt = &now
		
		switch sub.BillingFrequency {
		case "WEEKLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 0, 7)
		case "MONTHLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 1, 0)
		case "QUARTERLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 3, 0)
		case "ANNUALLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(1, 0, 0)
		default:
			// Default to Monthly if unknown
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 1, 0)
		}
		
		// Update sub in DB (Missing Update method in Repo for now)
		// s.subRepo.Update(ctx, &sub)
	}
	return nil
}
