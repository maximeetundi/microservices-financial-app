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
	// Delegate wallet resolution to transfer-service (like TicketService)
	// We pass the Organizer's (Creator's) UserID as DestinationUserID.
	
	paymentReq := messaging.PaymentRequestEvent{
		RequestID:         donation.ID.Hex(), // Use Donation ID as RequestID
		UserID:            req.DonorID,       // Payer
		FromWalletID:      req.WalletID,
		DestinationUserID: campaign.CreatorID, // Set organizer as destination
		ToWalletID:        "",                 // Will be resolved by transfer-service
		DebitAmount:       req.Amount,
		CreditAmount:      req.Amount, // Or less fees?
		Currency:          req.Currency,
		Type:              "donation",
		ReferenceID:       fmt.Sprintf("DON-%s", donation.ID.Hex()),
		Metadata:          map[string]interface{}{
			"campaign_title": campaign.Title,
			"campaign_id":    campaign.ID.Hex(),
		},
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
	return donations, nil
}

// RefundDonation initiates a refund for a specific donation
func (s *DonationService) RefundDonation(donationID, requesterID, reason string) error {
	// 1. Get Donation
	donation, err := s.donationRepo.GetByID(donationID)
	if err != nil {
		return errors.New("donation not found")
	}

	if donation.Status != models.DonationStatusPaid {
		return fmt.Errorf("donation status is %s, cannot refund", donation.Status)
	}

	// 2. Get Campaign to verify Creator
	campaign, err := s.campaignRepo.GetByID(donation.CampaignID.Hex())
	if err != nil {
		return errors.New("associated campaign not found")
	}

	if campaign.CreatorID != requesterID {
		return errors.New("unauthorized: only campaign creator can refund")
	}

	// 3. Identify Source Wallet (Creator's Wallet)
	// We need to find the wallet that originally received the funds? 
	// Or just a wallet of the creator with same currency.
	// For now, let's find a wallet of creator with matching currency.
	serverWallets, err := s.walletClient.GetUserWallets(campaign.CreatorID)
	if err != nil {
		return fmt.Errorf("failed to fetch creator wallets: %w", err)
	}

	var sourceWalletID string
	for _, w := range serverWallets {
		if cur, ok := w["currency"].(string); ok && cur == donation.Currency {
			if id, ok := w["id"].(string); ok {
				sourceWalletID = id
				break
			}
		}
	}
	if sourceWalletID == "" {
		return errors.New("creator does not have a wallet for refund currency")
	}

	// 4. Identify Destination Wallet (Donor's Wallet)
	// We stored `PaymentWalletID` in Donation model.
	destWalletID := donation.PaymentWalletID
	if destWalletID == "" {
		return errors.New("donor wallet not recorded, cannot refund automatically")
	}

	// 5. Trigger Refund Payment Event
	refundReq := messaging.PaymentRequestEvent{
		RequestID:         fmt.Sprintf("REF-%s", donation.ID.Hex()),
		UserID:            requesterID,        // Initiator (Creator)
		FromWalletID:      sourceWalletID,     // Creator pays back
		ToWalletID:        destWalletID,       // Donor receives
		DebitAmount:       donation.Amount,    // Full amount
		CreditAmount:      donation.Amount,
		Currency:          donation.Currency,
		Type:              "refund",
		ReferenceID:       fmt.Sprintf("REF-DON-%s", donation.ID.Hex()),
		Metadata:          map[string]interface{}{"reason": reason},
	}

	envelope := messaging.NewEventEnvelope(messaging.EventPaymentRequest, "donation-service", refundReq)
	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicPaymentEvents, envelope); err != nil {
		return fmt.Errorf("failed to publish refund request: %w", err)
	}

	// 6. Update Donation Status
	if err := s.donationRepo.UpdateStatus(donationID, models.DonationStatusRefunding, ""); err != nil {
		return fmt.Errorf("failed to update donation status: %w", err)
	}

	// 7. Decrement Campaign Amount
	if err := s.campaignRepo.IncrementCollectedAmount(campaign.ID.Hex(), -donation.Amount); err != nil {
		// Log error but proceed as refund is initiated
		fmt.Printf("Failed to decrement campaign amount for refund %s: %v\n", donationID, err)
	}

	return nil
}

// CancelCampaign cancels a campaign and refunds all paid donations
func (s *DonationService) CancelCampaign(campaignID, requesterID, reason, pin, token string) error {
	// 0. Verify PIN
	if err := s.userClient.VerifyPin(requesterID, pin, token); err != nil {
		return fmt.Errorf("security check failed: %w", err)
	}

	// 1. Get Campaign
	campaign, err := s.campaignRepo.GetByID(campaignID)
	if err != nil {
		return errors.New("campaign not found")
	}

	if campaign.CreatorID != requesterID {
		return errors.New("unauthorized: only creator can cancel campaign")
	}

	if campaign.Status == models.CampaignStatusCancelled {
		return errors.New("campaign is already cancelled")
	}

	// 2. Mark Campaign as Cancelled
	// This prevents new donations.
	updates := map[string]interface{}{
		"status": models.CampaignStatusCancelled,
	}
	if err := s.campaignRepo.Update(campaignID, updates); err != nil {
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	// 3. List Data and Refund
	// Note: Pagination? If thousands of donations, this might timeout.
	// For MVP, limit to 1000 or reasonable number.
	donations, err := s.donationRepo.ListByCampaign(campaignID, 1000, 0)
	if err != nil {
		// Log error but campaign is already cancelled.
		fmt.Printf("Failed to list donations for cancellation: %v\n", err)
		return nil
	}

	for _, d := range donations {
		if d.Status == models.DonationStatusPaid {
			// Trigger refund asynchronously or synchronously?
			// Synchronously for now, but handle details carefully.
			// Re-use Refund logic but bypass some checks if needed?
			// reusing RefundDonation is safest as it checks wallets.
			
			// We pass requesterID (creator) to authorize the refund.
			if err := s.RefundDonation(d.ID.Hex(), requesterID, reason); err != nil {
				fmt.Printf("Failed to refund donation %s: %v\n", d.ID.Hex(), err)
				// Continue to next donation
			}
		}
	}

	return nil
}
