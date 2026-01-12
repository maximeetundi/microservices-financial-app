package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/repository"
)

type PaymentConsumer struct {
	kafkaClient  *messaging.KafkaClient
	donationRepo *repository.DonationRepository
	campaignRepo *repository.CampaignRepository
}

func NewPaymentConsumer(
	kafkaClient *messaging.KafkaClient,
	donationRepo *repository.DonationRepository, 
	campaignRepo *repository.CampaignRepository,
) *PaymentConsumer {
	return &PaymentConsumer{
		kafkaClient:  kafkaClient,
		donationRepo: donationRepo,
		campaignRepo: campaignRepo,
	}
}

func (c *PaymentConsumer) Start() {
	// Subscribe to payment status events
	// Note: We need to filter for events relevant to donations. 
	
	err := c.kafkaClient.Subscribe(messaging.TopicPaymentEvents, c.handlePaymentEvent)
	if err != nil {
		log.Fatalf("Failed to start payment consumer: %v", err)
	}
}

func (c *PaymentConsumer) handlePaymentEvent(ctx context.Context, msg *messaging.EventEnvelope) error {
	if msg.Type != messaging.EventPaymentSuccess && msg.Type != messaging.EventPaymentFailed {
		return nil // Ignore other events
	}

	var event messaging.PaymentStatusEvent
	dataBytes, _ := json.Marshal(msg.Data)
	if err := json.Unmarshal(dataBytes, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payment status: %w", err)
	}

	// Check if this payment is for a donation
	// We use Type check if available, or try to find donation by ID.
	if event.Type != "donation" {
		return nil
	}

	// RequestID is the Donation ID (set in InitiateDonation)
	donationID := event.RequestID
	
	donation, err := c.donationRepo.GetByID(donationID)
	if err != nil {
		// Log but don't error loudly if not found (maybe it's for another service if types overlap?)
		// But type="donation" should be unique to us.
		log.Printf("Donation not found for payment ID %s: %v", donationID, err)
		return nil
	}

	if event.Status == "success" {
		if err := c.donationRepo.UpdateStatus(donationID, models.DonationStatusPaid, event.ReferenceID); err != nil {
			log.Printf("Failed to update donation status: %v", err)
			return err
		}
		
		// Update Campaign Collected Amount
		if err := c.campaignRepo.IncrementCollectedAmount(donation.CampaignID.Hex(), donation.Amount); err != nil {
			log.Printf("Failed to increment campaign amount: %v", err)
			// Don't fail the consumer, but this is an issue.
		}

		// Send Notification to Creator?
		// "You received a donation of X from Y"
		// "New donation on campaign Z"
		c.sendNotification(donation, "success")
		
	} else {
		c.donationRepo.UpdateStatus(donationID, models.DonationStatusFailed, "")
		c.sendNotification(donation, "failed")
	}

	return nil
}

func (c *PaymentConsumer) sendNotification(donation *models.Donation, status string) {
	// 1. Get Campaign to know Creator and Title
	campaign, err := c.campaignRepo.GetByID(donation.CampaignID.Hex())
	if err != nil {
		log.Printf("Failed to fetch campaign for notification: %v", err)
		return
	}

	var title, message string
	if status == "success" {
		title = "Nouveau Don Reçu! 🎉"
		donorName := "Un donateur anonyme"
		if !donation.IsAnonymous {
			// We only have DonorID. We might need to fetch name? 
			// Or just say "A user". Frontend can resolve name if needed.
			// Ideally we fetch name from user-service or auth-service.
			// For MVP, "Un utilisateur" or just UserID? 
			// Let's rely on Frontend or enrich message later.
			donorName = "Un donateur" 
		}
		message = fmt.Sprintf("%s a fait un don de %.2f %s sur votre campagne '%s'.", 
			donorName, donation.Amount, donation.Currency, campaign.Title)
		
		// Notification to Campaign Creator
		notif := messaging.NotificationEvent{
			UserID:  campaign.CreatorID,
			Type:    "donation_received",
			Title:   title,
			Message: message,
			Data: map[string]interface{}{
				"campaign_id": campaign.ID.Hex(),
				"donation_id": donation.ID.Hex(),
				"amount":      donation.Amount,
			},
		}
		env := messaging.NewEventEnvelope(messaging.EventNotificationCreated, "donation-service", notif)
		c.kafkaClient.Publish(context.Background(), messaging.TopicNotificationEvents, env)
		
		// Notification to Donor (Thank you)
		if donation.DonorID != "" {
			thankYouMsg := fmt.Sprintf("Merci pour votre don de %.2f %s à la campagne '%s'!", 
				donation.Amount, donation.Currency, campaign.Title)
			donorNotif := messaging.NotificationEvent{
				UserID:  donation.DonorID,
				Type:    "donation_sent",
				Title:   "Don Confirmé ✅",
				Message: thankYouMsg,
			}
			donorEnv := messaging.NewEventEnvelope(messaging.EventNotificationCreated, "donation-service", donorNotif)
			c.kafkaClient.Publish(context.Background(), messaging.TopicNotificationEvents, donorEnv)
		}

	}
}
