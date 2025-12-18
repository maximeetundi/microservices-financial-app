package services

import (
	"encoding/json"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/streadway/amqp"
)

// Consumer handles RabbitMQ message consumption for wallet-service
type Consumer struct {
	channel       *amqp.Channel
	walletService *WalletService
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(channel *amqp.Channel, walletService *WalletService) *Consumer {
	return &Consumer{
		channel:       channel,
		walletService: walletService,
	}
}

// Start begins consuming messages from all subscribed queues
func (c *Consumer) Start() error {
	// Start consumers for each queue
	go c.consumeTransferCompleted()
	go c.consumeExchangeCompleted()
	go c.consumeCardLoaded()
	go c.consumeUserRegistered()

	log.Println("Wallet service consumers started")
	return nil
}

// consumeTransferCompleted handles transfer.completed events
func (c *Consumer) consumeTransferCompleted() {
	msgs, err := c.channel.Consume(
		"wallet.transfer_completed", // queue
		"",                          // consumer
		false,                       // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)
	if err != nil {
		log.Printf("Failed to consume transfer_completed: %v", err)
		return
	}

	for msg := range msgs {
		c.handleTransferCompleted(msg)
	}
}

func (c *Consumer) handleTransferCompleted(msg amqp.Delivery) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal transfer event: %v", err)
		msg.Nack(false, true) // Requeue
		return
	}

	log.Printf("Processing transfer_completed event: %v", event)

	// Extract transfer details
	fromWalletID, _ := event["from_wallet_id"].(string)
	toWalletID, _ := event["to_wallet_id"].(string)
	amount, _ := event["amount"].(float64)
	fee, _ := event["fee"].(float64)

	// Debit source wallet
	if fromWalletID != "" {
		totalDebit := amount + fee
		err := c.walletService.balanceService.UpdateBalance(fromWalletID, totalDebit, "debit")
		if err != nil {
			log.Printf("Failed to debit wallet %s: %v", fromWalletID, err)
			msg.Nack(false, true)
			return
		}
	}

	// Credit destination wallet
	if toWalletID != "" {
		err := c.walletService.balanceService.UpdateBalance(toWalletID, amount, "credit")
		if err != nil {
			log.Printf("Failed to credit wallet %s: %v", toWalletID, err)
			msg.Nack(false, true)
			return
		}
	}

	msg.Ack(false)
	log.Printf("Successfully processed transfer_completed event")
}

// consumeExchangeCompleted handles exchange.completed events
func (c *Consumer) consumeExchangeCompleted() {
	msgs, err := c.channel.Consume(
		"wallet.exchange_completed", // queue
		"",                          // consumer
		false,                       // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)
	if err != nil {
		log.Printf("Failed to consume exchange_completed: %v", err)
		return
	}

	for msg := range msgs {
		c.handleExchangeCompleted(msg)
	}
}

func (c *Consumer) handleExchangeCompleted(msg amqp.Delivery) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal exchange event: %v", err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Processing exchange_completed event: %v", event)

	// Extract exchange details
	fromWalletID, _ := event["from_wallet_id"].(string)
	toWalletID, _ := event["to_wallet_id"].(string)
	fromAmount, _ := event["from_amount"].(float64)
	toAmount, _ := event["to_amount"].(float64)
	fee, _ := event["fee"].(float64)

	// Debit source wallet (from_amount + fee)
	if fromWalletID != "" {
		totalDebit := fromAmount + fee
		err := c.walletService.balanceService.UpdateBalance(fromWalletID, totalDebit, "debit")
		if err != nil {
			log.Printf("Failed to debit wallet %s: %v", fromWalletID, err)
			msg.Nack(false, true)
			return
		}
	}

	// Credit destination wallet with exchanged amount
	if toWalletID != "" {
		err := c.walletService.balanceService.UpdateBalance(toWalletID, toAmount, "credit")
		if err != nil {
			log.Printf("Failed to credit wallet %s: %v", toWalletID, err)
			msg.Nack(false, true)
			return
		}
	}

	msg.Ack(false)
	log.Printf("Successfully processed exchange_completed event")
}

// consumeCardLoaded handles card.loaded events
func (c *Consumer) consumeCardLoaded() {
	msgs, err := c.channel.Consume(
		"wallet.card_loaded", // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Printf("Failed to consume card_loaded: %v", err)
		return
	}

	for msg := range msgs {
		c.handleCardLoaded(msg)
	}
}

func (c *Consumer) handleCardLoaded(msg amqp.Delivery) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal card loaded event: %v", err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Processing card_loaded event: %v", event)

	// Extract card load details
	sourceWalletID, _ := event["source_wallet_id"].(string)
	amount, _ := event["amount"].(float64)
	fee, _ := event["fee"].(float64)

	// Debit the source wallet
	if sourceWalletID != "" {
		totalDebit := amount + fee
		err := c.walletService.balanceService.UpdateBalance(sourceWalletID, totalDebit, "debit")
		if err != nil {
			log.Printf("Failed to debit wallet %s for card load: %v", sourceWalletID, err)
			msg.Nack(false, true)
			return
		}
	}

	msg.Ack(false)
	log.Printf("Successfully processed card_loaded event")
}

// consumeUserRegistered handles user.registered events
func (c *Consumer) consumeUserRegistered() {
	_, err := c.channel.QueueDeclare(
		"user.registered", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Printf("Failed to declare user.registered queue: %v", err)
		return
	}

	msgs, err := c.channel.Consume(
		"user.registered", // queue
		"",                // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Printf("Failed to consume user.registered: %v", err)
		return
	}

	for msg := range msgs {
		c.handleUserRegistered(msg)
	}
}

func (c *Consumer) handleUserRegistered(msg amqp.Delivery) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal user registered event: %v", err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Processing user.registered event: %v", event)

	userID, _ := event["user_id"].(string)
	currency, _ := event["currency"].(string)

	if userID != "" && currency != "" {
		// Create default wallet
		name := "Main Wallet"
		desc := "Default wallet created on registration"
		
		req := &models.CreateWalletRequest{
			Currency:    currency,
			WalletType:  "fiat",
			Name:        &name,
			Description: &desc,
		}
		
		// Internal call, no auth token needed
		_, err := c.walletService.CreateWallet(userID, req)
		if err != nil {
			log.Printf("Failed to create default wallet for user %s: %v", userID, err)
			msg.Nack(false, true)
			return
		}
		log.Printf("Created default wallet for user %s with currency %s", userID, currency)
	}

	msg.Ack(false)
}
