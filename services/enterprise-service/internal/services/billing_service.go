package services

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BillingService struct {
	subRepo   *repository.SubscriptionRepository
	invRepo   *repository.InvoiceRepository
	batchRepo *repository.BatchRepository
	entRepo   *repository.EnterpriseRepository
}

func NewBillingService(subRepo *repository.SubscriptionRepository, invRepo *repository.InvoiceRepository, batchRepo *repository.BatchRepository, entRepo *repository.EnterpriseRepository) *BillingService {
	return &BillingService{
		subRepo:   subRepo,
		invRepo:   invRepo,
		batchRepo: batchRepo,
		entRepo:   entRepo,
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
		case "DAILY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 0, 1)
		case "WEEKLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 0, 7)
		case "MONTHLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 1, 0)
		case "QUARTERLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 3, 0)
		case "ANNUALLY":
			sub.NextBillingAt = sub.NextBillingAt.AddDate(1, 0, 0)
		case "CUSTOM":
			// Ensure CustomInterval is valid, default to 30 days if 0
			days := sub.CustomInterval
			if days <= 0 {
				days = 30
			}
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 0, days)
		default:
			// Default to Monthly if unknown
			sub.NextBillingAt = sub.NextBillingAt.AddDate(0, 1, 0)
		}
		
		// Update sub in DB
		s.subRepo.Update(ctx, &sub)
	}
	return nil
}
// ManualInvoiceItem for Service Layer
type ManualInvoiceItem struct {
	SubscriptionID string
	Amount         float64
	Consumption    float64
}

func (s *BillingService) GenerateBatchFromManualEntry(ctx context.Context, enterpriseID string, items []ManualInvoiceItem) error {
	ent, err := s.entRepo.FindByID(ctx, enterpriseID)
	if err != nil {
		return err
	}

	var invoices []models.Invoice

	for _, item := range items {
		sub, err := s.subRepo.FindByID(ctx, item.SubscriptionID)
		if err != nil {
			continue // Skip invalid subs
		}

		finalAmount := item.Amount
		var consumption *float64

		// Calculate based on Usage if Consumtpion provided and Amount is 0
		if finalAmount == 0 && item.Consumption > 0 {
			c := item.Consumption
			consumption = &c
			
			// Find Service Def
			var unitPrice float64
			for _, svc := range ent.CustomServices {
				if svc.ID == sub.ServiceID {
					unitPrice = svc.BasePrice
					break
				}
			}
			finalAmount = c * unitPrice
		}
		
		// Fallback to Base Price if still 0 (e.g. Fixed service manual trigger)
			finalAmount = sub.Amount
        }
        
		entID, _ := primitive.ObjectIDFromHex(enterpriseID)
		
		inv := models.Invoice{
			EnterpriseID:  entID,
			ClientID:      sub.ClientID,
			ClientName:    "Subscriber", // Fetch name if possible, or leave generic
			ServiceID:     sub.ServiceID,
			Amount:        finalAmount,
			Consumption:   consumption,
			DueDate:       time.Now().AddDate(0, 0, 7),
			Status:        models.InvoiceStatusDraft, 
			SentAt:        time.Now(),
		}
		invoices = append(invoices, inv)
	}
	
	// Create Batch
	if len(invoices) > 0 {
		_, err := s.CreateBatchFromImport(ctx, enterpriseID, invoices) // Reuse CreateBatchFromImport logic which saves invoices and batch
		return err
	}

	return nil
}
