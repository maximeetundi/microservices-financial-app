package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/repository"
)

type DonationService struct {
	donationRepo *repository.DonationRepository
	campaignRepo *repository.CampaignRepository
	walletClient *WalletClient
	userClient   *UserClient
	kafkaClient  *messaging.KafkaClient
}

func NewDonationService(
	donationRepo *repository.DonationRepository, 
	campaignRepo *repository.CampaignRepository, 
	walletClient *WalletClient,
	userClient   *UserClient,
	kafkaClient  *messaging.KafkaClient,
) *DonationService {
	return &DonationService{
		donationRepo: donationRepo,
		campaignRepo: campaignRepo,
		walletClient: walletClient,
		userClient:   userClient,
		kafkaClient:  kafkaClient,
	}
}

type InitiateDonationRequest struct {
	CampaignID string  `json:"campaign_id"`
	DonorID    string  `json:"donor_id"` // User ID of donor
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	WalletID   string  `json:"wallet_id"` // Donor's wallet
	PIN        string  `json:"pin"`
	Message    string  `json:"message"`
	IsAnonymous bool   `json:"is_anonymous"`
	Frequency  models.DonationFrequency `json:"frequency"`
	FormData   map[string]interface{}   `json:"form_data"`
}

func (s *DonationService) InitiateDonation(req *InitiateDonationRequest, token string) (*models.Donation, error) {
	// 0. Verify PIN securely
	if err := s.userClient.VerifyPin(req.DonorID, req.PIN, token); err != nil {
		return nil, fmt.Errorf("security check failed: %w", err)
	}

	// 1. Get Campaign
	campaign, err := s.campaignRepo.GetByID(req.CampaignID)
	if err != nil {
		return nil, errors.New("campaign not found")
	}
	

	if campaign.Status != models.CampaignStatusActive {
		return nil, errors.New("campaign is not active")
	}

	// 1b. Validate Amount Limits
	if campaign.MinAmount > 0 && req.Amount < campaign.MinAmount {
		return nil, fmt.Errorf("amount below minimum allowed (%f %s)", campaign.MinAmount, campaign.Currency)
	}
	if campaign.MaxAmount > 0 && req.Amount > campaign.MaxAmount {
		return nil, fmt.Errorf("amount exceeds maximum allowed (%f %s)", campaign.MaxAmount, campaign.Currency)
	}

	// 2. Validate Frequency
	if req.Frequency != "" && req.Frequency != models.FrequencyOneTime {
		if !campaign.AllowRecurring {
			return nil, errors.New("recurring donations not allowed for this campaign")
		}
	} else {
		req.Frequency = models.FrequencyOneTime
	}

	// 3. Create Donation Pending Record
	donation := &models.Donation{
		CampaignID:  campaign.ID,
		DonorID:     req.DonorID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Message:     req.Message,
		IsAnonymous: req.IsAnonymous,
		Frequency:   req.Frequency,
		Status:      models.DonationStatusPending,
		PaymentWalletID: req.WalletID,
		FormData:    req.FormData,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if req.Frequency != models.FrequencyOneTime {
		// Calculate next payment date
		// Logic: If user pays NOW, next date is +1 Period.
		// If manual validation is required, we still set next date.
		now := time.Now()
		var nextDate time.Time
		switch req.Frequency {
		case models.FrequencyMonthly:
			nextDate = now.AddDate(0, 1, 0)
		case models.FrequencyQuarterly:
			nextDate = now.AddDate(0, 3, 0)
		case models.FrequencyAnnually:
			nextDate = now.AddDate(1, 0, 0)
		}
		donation.NextPaymentDate = &nextDate
	}
	
	if err := s.donationRepo.Create(donation); err != nil {
		return nil, err
	}

	// 4. Create Payment Request Event
	// We need the Organizer's (Creator's) Wallet to deposit to.
	// We need to fetch Organizer's wallet for the currency.
	
	wallets, err := s.walletClient.GetUserWallets(campaign.CreatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch organizer wallets: %w", err)
	}

	var destWalletID string
	for _, w := range wallets {
		if cur, ok := w["currency"].(string); ok && cur == req.Currency {
			if id, ok := w["id"].(string); ok {
				destWalletID = id
				break
			}
		}
	}
	
	if destWalletID == "" {
		// Fallback: If no wallet found for created currency, maybe default wallet? 
		// Or error out. Better to error or create a wallet?
		// For MVP, error out.
		return nil, fmt.Errorf("organizer does not have a wallet for currency %s", req.Currency)
	}

	paymentReq := messaging.PaymentRequestEvent{
		RequestID:         donation.ID.Hex(), // Use Donation ID as RequestID
		UserID:            req.DonorID,       // Payer
		FromWalletID:      req.WalletID,
		ToWalletID:        destWalletID,
		DebitAmount:       req.Amount,
		CreditAmount:      req.Amount, // Or less fees?
		Currency:          req.Currency,
		Type:              "donation",
		ReferenceID:       fmt.Sprintf("DON-%s", donation.ID.Hex()),
	}

	envelope := messaging.NewEventEnvelope(messaging.EventPaymentRequest, "donation-service", paymentReq)
	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicPaymentEvents, envelope); err != nil {
		return nil, fmt.Errorf("failed to publish payment request: %w", err)
	}

	return donation, nil
}

func (s *DonationService) ListDonations(campaignID string, limit, offset int64) ([]*models.Donation, error) {
	donations, err := s.donationRepo.ListByCampaign(campaignID, limit, offset)
	if err != nil {
		return nil, err
	}
	// Filter anonymity? 
	// The repo returns full objects. The Handler should filter sensitive data (DonorID if anonymous).
	return donations, nil
}
