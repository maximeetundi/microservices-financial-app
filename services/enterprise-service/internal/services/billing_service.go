package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BillingService struct {
	subRepo     *repository.SubscriptionRepository
	invRepo     *repository.InvoiceRepository
	batchRepo   *repository.BatchRepository
	entRepo     *repository.EnterpriseRepository
	notifClient *NotificationClient
}

func NewBillingService(
	subRepo *repository.SubscriptionRepository, 
	invRepo *repository.InvoiceRepository, 
	batchRepo *repository.BatchRepository, 
	entRepo *repository.EnterpriseRepository,
	notifClient *NotificationClient,
) *BillingService {
	return &BillingService{
		subRepo:     subRepo,
		invRepo:     invRepo,
		batchRepo:   batchRepo,
		entRepo:     entRepo,
		notifClient: notifClient,
	}
}

func (s *BillingService) ListInvoicesForClient(ctx context.Context, enterpriseID, clientID string) ([]models.Invoice, error) {
	return s.invRepo.FindByEnterpriseAndClientID(ctx, enterpriseID, clientID)
}

// CreateInvoice (Single manual or system generated)
func (s *BillingService) CreateInvoice(ctx context.Context, inv *models.Invoice) error {
	inv.Status = models.InvoiceStatusDraft
	if err := s.invRepo.Create(ctx, inv); err != nil {
		return err
	}
	
	// Send notification to client
	if inv.ClientID != "" && s.notifClient != nil {
		ent, err := s.entRepo.FindByID(ctx, inv.EnterpriseID.Hex())
		if err == nil {
			dueDate := inv.DueDate.Format("02/01/2006")
			if err := s.notifClient.NotifyInvoice(ctx, inv.ClientID, ent.Name, inv.ServiceName, inv.Amount, dueDate); err != nil {
				log.Printf("Failed to send invoice notification: %v", err)
			}
		}
	}
	
	return nil
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

// CalculatePenalty calculates the penalty amount based on the penalty config
// Supports grace period with unit (DAYS, WEEKS, MONTHS)
// Supports frequencies: DAILY, WEEKLY, MONTHLY, QUARTERLY, SEMIANNUAL, ANNUAL
// Respects max_penalty_months and max_penalty_amount limits
func CalculatePenalty(config *models.PenaltyConfig, periodAmount float64, daysOverdue int, dueDateMonthsAgo int) float64 {
	if config == nil {
		return 0
	}

	// Convert grace period to days based on unit
	graceDays := config.GracePeriod
	switch config.GraceUnit {
	case "WEEKS":
		graceDays = config.GracePeriod * 7
	case "MONTHS":
		graceDays = config.GracePeriod * 30
	}

	if daysOverdue <= graceDays {
		return 0
	}

	// Check if we've exceeded the max penalty months limit
	if config.MaxPenaltyMonths > 0 && dueDateMonthsAgo > config.MaxPenaltyMonths {
		// Penalty period has ended, cap at max months
		dueDateMonthsAgo = config.MaxPenaltyMonths
	}

	var penaltyAmount float64
	effectiveDays := daysOverdue - graceDays

	// Calculate base penalty by type
	switch config.Type {
	case "FIXED":
		penaltyAmount = config.Value
	case "PERCENTAGE":
		penaltyAmount = periodAmount * (config.Percentage / 100)
	case "HYBRID":
		// Hybrid = fixed amount + percentage of period amount
		penaltyAmount = config.Value + (periodAmount * (config.Percentage / 100))
	default:
		return 0
	}

	// Apply frequency multiplier
	switch config.Frequency {
	case "DAILY":
		penaltyAmount = penaltyAmount * float64(effectiveDays)
	case "WEEKLY":
		weeks := effectiveDays / 7
		if weeks < 1 {
			weeks = 1
		}
		penaltyAmount = penaltyAmount * float64(weeks)
	case "MONTHLY":
		months := effectiveDays / 30
		if months < 1 {
			months = 1
		}
		penaltyAmount = penaltyAmount * float64(months)
	case "QUARTERLY":
		quarters := effectiveDays / 90
		if quarters < 1 {
			quarters = 1
		}
		penaltyAmount = penaltyAmount * float64(quarters)
	case "SEMIANNUAL":
		semesters := effectiveDays / 180
		if semesters < 1 {
			semesters = 1
		}
		penaltyAmount = penaltyAmount * float64(semesters)
	case "ANNUAL":
		years := effectiveDays / 365
		if years < 1 {
			years = 1
		}
		penaltyAmount = penaltyAmount * float64(years)
	}

	// Apply max penalty amount cap if set
	if config.MaxPenaltyAmount > 0 && penaltyAmount > config.MaxPenaltyAmount {
		penaltyAmount = config.MaxPenaltyAmount
	}

	return penaltyAmount
}

// CalculateUsagePrice calculates the total price based on consumption and pricing tiers
// Supports: FIXED (simple price/unit), TIERED (different price per tier), THRESHOLD (bonuses when reaching thresholds)
func CalculateUsagePrice(serviceDef *models.ServiceDefinition, consumption float64) (totalPrice float64, breakdown []map[string]interface{}) {
	if serviceDef == nil || consumption <= 0 {
		return 0, nil
	}

	// Default to FIXED mode if not specified
	pricingMode := serviceDef.PricingMode
	if pricingMode == "" {
		pricingMode = "FIXED"
	}

	switch pricingMode {
	case "FIXED":
		// Simple calculation: consumption * base_price
		totalPrice = consumption * serviceDef.BasePrice
		breakdown = append(breakdown, map[string]interface{}{
			"label":       "Prix unitaire fixe",
			"consumption": consumption,
			"price_unit":  serviceDef.BasePrice,
			"subtotal":    totalPrice,
		})

	case "TIERED":
		// Tiered pricing: each tier has its own price
		// E.g., 0-100 kWh at 50 XOF, 101-500 at 75 XOF, 500+ at 100 XOF
		remainingConsumption := consumption
		
		for _, tier := range serviceDef.PricingTiers {
			if remainingConsumption <= 0 {
				break
			}
			
			tierMin := tier.MinConsumption
			tierMax := tier.MaxConsumption
			if tierMax < 0 {
				tierMax = remainingConsumption + tierMin // Unlimited
			}
			
			// Calculate consumption in this tier
			tierRange := tierMax - tierMin
			tierConsumption := remainingConsumption
			if tierConsumption > tierRange {
				tierConsumption = tierRange
			}
			
			if tierConsumption > 0 {
				tierTotal := tierConsumption * tier.PricePerUnit
				
				// Apply percentage bonus if any
				if tier.PercentBonus > 0 {
					tierTotal += tierTotal * (tier.PercentBonus / 100)
				}
				
				// Apply fixed bonus if any
				tierTotal += tier.FixedBonus
				
				totalPrice += tierTotal
				
				breakdown = append(breakdown, map[string]interface{}{
					"label":        tier.Label,
					"min":          tierMin,
					"max":          tierMax,
					"consumption":  tierConsumption,
					"price_unit":   tier.PricePerUnit,
					"fixed_bonus":  tier.FixedBonus,
					"percent_bonus": tier.PercentBonus,
					"subtotal":     tierTotal,
				})
				
				remainingConsumption -= tierConsumption
			}
		}

	case "THRESHOLD":
		// Threshold pricing: base price + bonuses when thresholds are reached
		// The last matching tier applies
		basePrice := serviceDef.BasePrice
		var matchedTier *models.PricingTier
		
		for i := range serviceDef.PricingTiers {
			tier := &serviceDef.PricingTiers[i]
			if consumption >= tier.MinConsumption {
				matchedTier = tier
			}
		}
		
		// Calculate base amount
		baseAmount := consumption * basePrice
		totalPrice = baseAmount
		
		breakdown = append(breakdown, map[string]interface{}{
			"label":       "Prix de base",
			"consumption": consumption,
			"price_unit":  basePrice,
			"subtotal":    baseAmount,
		})
		
		// Apply threshold bonuses if a tier matched
		if matchedTier != nil {
			// Override price per unit if tier has one
			if matchedTier.PricePerUnit > 0 {
				priceChange := (matchedTier.PricePerUnit - basePrice) * consumption
				totalPrice += priceChange
				breakdown = append(breakdown, map[string]interface{}{
					"label":       matchedTier.Label + " - Ajustement prix",
					"price_unit":  matchedTier.PricePerUnit,
					"subtotal":    priceChange,
				})
			}
			
			// Apply fixed bonus
			if matchedTier.FixedBonus > 0 {
				totalPrice += matchedTier.FixedBonus
				breakdown = append(breakdown, map[string]interface{}{
					"label":    matchedTier.Label + " - Bonus fixe",
					"subtotal": matchedTier.FixedBonus,
				})
			}
			
			// Apply percentage bonus
			if matchedTier.PercentBonus > 0 {
				percentAmount := baseAmount * (matchedTier.PercentBonus / 100)
				totalPrice += percentAmount
				breakdown = append(breakdown, map[string]interface{}{
					"label":    matchedTier.Label + " - Bonus " + fmt.Sprintf("%.1f%%", matchedTier.PercentBonus),
					"subtotal": percentAmount,
				})
			}
		}

	default:
		// Fallback to fixed
		totalPrice = consumption * serviceDef.BasePrice
	}

	return totalPrice, breakdown
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
			for _, group := range ent.ServiceGroups {
				for _, svc := range group.Services {
					if svc.ID == sub.ServiceID {
						unitPrice = svc.BasePrice
						break
					}
				}
			}
			finalAmount = c * unitPrice
		}
		
		// Fallback to Base Price if still 0 (e.g. Fixed service manual trigger)
		if finalAmount == 0 {
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
