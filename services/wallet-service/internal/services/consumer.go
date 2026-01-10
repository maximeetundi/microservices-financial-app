package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
)

// Consumer handles Kafka message consumption for wallet-service
type Consumer struct {
	kafkaClient   *messaging.KafkaClient
	walletService *WalletService
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(kafkaClient *messaging.KafkaClient, walletService *WalletService) *Consumer {
	return &Consumer{
		kafkaClient:   kafkaClient,
		walletService: walletService,
	}
}

// Start begins consuming messages from all subscribed topics
func (c *Consumer) Start() error {
	// Subscribe to transfer events
	if err := c.kafkaClient.Subscribe(messaging.TopicTransferEvents, c.handleTransferEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to transfer events: %v", err)
	}

	// Subscribe to exchange events
	if err := c.kafkaClient.Subscribe(messaging.TopicExchangeEvents, c.handleExchangeEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to exchange events: %v", err)
	}

	// Subscribe to card events
	if err := c.kafkaClient.Subscribe(messaging.TopicCardEvents, c.handleCardEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to card events: %v", err)
	}

	// Subscribe to user events
	if err := c.kafkaClient.Subscribe(messaging.TopicUserEvents, c.handleUserEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to user events: %v", err)
	}

	log.Println("[Kafka] Wallet service consumers started")
	return nil
}

// handleTransferEvent processes transfer.completed events
func (c *Consumer) handleTransferEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Processing transfer event: %s", event.Type)

	if event.Type != messaging.EventTransferCompleted {
		return nil // Not a transfer completed event
	}

	// Parse the event data
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var transferEvent messaging.TransferCompletedEvent
	if err := json.Unmarshal(dataBytes, &transferEvent); err != nil {
		log.Printf("Failed to unmarshal transfer event: %v", err)
		return err
	}

	// Debit source wallet
	if transferEvent.FromWalletID != "" {
		totalDebit := transferEvent.Amount + transferEvent.Fee
		err := c.walletService.balanceService.UpdateBalance(transferEvent.FromWalletID, totalDebit, "debit")
		if err != nil {
			log.Printf("Failed to debit wallet %s: %v", transferEvent.FromWalletID, err)
			return err
		}
	}

	// Credit destination wallet
	if transferEvent.ToWalletID != "" {
		err := c.walletService.balanceService.UpdateBalance(transferEvent.ToWalletID, transferEvent.Amount, "credit")
		if err != nil {
			log.Printf("Failed to credit wallet %s: %v", transferEvent.ToWalletID, err)
			return err
		}
	}

	log.Printf("[Kafka] Successfully processed transfer_completed event")
	return nil
}

// handleExchangeEvent processes exchange.completed events
func (c *Consumer) handleExchangeEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Processing exchange event: %s", event.Type)

	if event.Type != messaging.EventExchangeCompleted {
		return nil
	}

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var exchangeEvent messaging.ExchangeCompletedEvent
	if err := json.Unmarshal(dataBytes, &exchangeEvent); err != nil {
		log.Printf("Failed to unmarshal exchange event: %v", err)
		return err
	}

	// Debit source wallet (from_amount + fee)
	if exchangeEvent.FromWalletID != "" {
		totalDebit := exchangeEvent.FromAmount + exchangeEvent.Fee
		err := c.walletService.balanceService.UpdateBalance(exchangeEvent.FromWalletID, totalDebit, "debit")
		if err != nil {
			log.Printf("Failed to debit wallet %s: %v", exchangeEvent.FromWalletID, err)
			return err
		}
	}

	// Credit destination wallet with exchanged amount
	if exchangeEvent.ToWalletID != "" {
		err := c.walletService.balanceService.UpdateBalance(exchangeEvent.ToWalletID, exchangeEvent.ToAmount, "credit")
		if err != nil {
			log.Printf("Failed to credit wallet %s: %v", exchangeEvent.ToWalletID, err)
			return err
		}
	}

	log.Printf("[Kafka] Successfully processed exchange_completed event")
	return nil
}

// handleCardEvent processes card.loaded events
func (c *Consumer) handleCardEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Processing card event: %s", event.Type)

	if event.Type != messaging.EventCardLoaded {
		return nil
	}

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var cardEvent messaging.CardLoadedEvent
	if err := json.Unmarshal(dataBytes, &cardEvent); err != nil {
		log.Printf("Failed to unmarshal card loaded event: %v", err)
		return err
	}

	// Debit the source wallet
	if cardEvent.SourceWalletID != "" {
		totalDebit := cardEvent.Amount + cardEvent.Fee
		err := c.walletService.balanceService.UpdateBalance(cardEvent.SourceWalletID, totalDebit, "debit")
		if err != nil {
			log.Printf("Failed to debit wallet %s for card load: %v", cardEvent.SourceWalletID, err)
			return err
		}
	}

	log.Printf("[Kafka] Successfully processed card_loaded event")
	return nil
}

// handleUserEvent processes user.registered events
func (c *Consumer) handleUserEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Processing user event: %s", event.Type)

	if event.Type != messaging.EventUserRegistered {
		return nil
	}

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var userEvent messaging.UserRegisteredEvent
	if err := json.Unmarshal(dataBytes, &userEvent); err != nil {
		log.Printf("Failed to unmarshal user registered event: %v", err)
		return err
	}

	if userEvent.UserID != "" && userEvent.Currency != "" {
		// Create default wallet
		name := "Main Wallet"
		desc := "Default wallet created on registration"

		req := &models.CreateWalletRequest{
			Currency:    userEvent.Currency,
			WalletType:  "fiat",
			Name:        &name,
			Description: &desc,
		}

		// Internal call, no auth token needed
		_, err := c.walletService.CreateWallet(userEvent.UserID, req)
		if err != nil {
			log.Printf("Failed to create default wallet for user %s: %v", userEvent.UserID, err)
			return err
		}
		log.Printf("[Kafka] Created default wallet for user %s with currency %s", userEvent.UserID, userEvent.Currency)
	}

	return nil
}
