package services

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
)

// PaymentStatusConsumer consumes payment status events from Kafka
type PaymentStatusConsumer struct {
	kafkaClient         *messaging.KafkaClient
	exchangeService     *ExchangeService
	fiatExchangeService *FiatExchangeService
}

// NewPaymentStatusConsumer creates a new payment status consumer
func NewPaymentStatusConsumer(kafkaClient *messaging.KafkaClient, exchangeService *ExchangeService) *PaymentStatusConsumer {
	return &PaymentStatusConsumer{
		kafkaClient:     kafkaClient,
		exchangeService: exchangeService,
	}
}

// SetFiatExchangeService allows setting the fiat exchange service after construction
func (c *PaymentStatusConsumer) SetFiatExchangeService(fiatService *FiatExchangeService) {
	c.fiatExchangeService = fiatService
}

// Start begins consuming payment status events
func (c *PaymentStatusConsumer) Start() error {
	log.Println("[Kafka] Starting PaymentStatusConsumer for exchange-service...")

	if err := c.kafkaClient.Subscribe(messaging.TopicPaymentEvents, c.handlePaymentEvent); err != nil {
		log.Printf("[Kafka] Failed to subscribe to payment events: %v", err)
		return err
	}

	log.Println("[Kafka] PaymentStatusConsumer started")
	return nil
}

func (c *PaymentStatusConsumer) handlePaymentEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	// Only process payment success/failed events
	if event.Type != messaging.EventPaymentSuccess && event.Type != messaging.EventPaymentFailed {
		return nil
	}

	log.Printf("[Kafka] Exchange-service received payment event: %s", event.Type)

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var statusEvent models.PaymentStatusEvent
	if err := json.Unmarshal(dataBytes, &statusEvent); err != nil {
		log.Printf("Failed to unmarshal payment status event: %v", err)
		return err
	}
	log.Printf("[TRACE-FIAT] Decoded StatusEvent: %+v", statusEvent)

	txID := statusEvent.TransactionID
	var exchangeID, step string
	var isFiat bool

	// Check for crypto exchange transactions
	if strings.HasPrefix(txID, "TX-EXC-DEBIT-") {
		exchangeID = txID[13:]
		step = "1_debit"
		isFiat = false
	} else if strings.HasPrefix(txID, "TX-EXC-CREDIT-") {
		exchangeID = txID[14:]
		step = "2_credit"
		isFiat = false
	} else if strings.HasPrefix(txID, "TX-FIAT-DEBIT-") {
		exchangeID = txID[14:]
		step = "1_debit"
		isFiat = true
	} else if strings.HasPrefix(txID, "TX-FIAT-CREDIT-") {
		exchangeID = txID[15:]
		step = "2_credit"
		isFiat = true
	} else {
		// Not an exchange transaction
		return nil
	}

	log.Printf("[Kafka] Processing %s exchange %s step %s: %s",
		map[bool]string{true: "Fiat", false: "Crypto"}[isFiat], exchangeID, step, statusEvent.Status)

	if statusEvent.Status == "completed" || statusEvent.Status == "success" {
		if isFiat && c.fiatExchangeService != nil {
			if step == "1_debit" {
				log.Printf("[TRACE-FIAT] Payment Consumer: Step 1 debit success, calling CompleteFiatExchangeCredit for %s", exchangeID)
				c.fiatExchangeService.CompleteFiatExchangeCredit(exchangeID)
			} else if step == "2_credit" {
				log.Printf("[TRACE-FIAT] Payment Consumer: Step 2 credit success, calling FinalizeFiatExchange for %s", exchangeID)
				c.fiatExchangeService.FinalizeFiatExchange(exchangeID)
			}
		} else {
			if step == "1_debit" {
				c.exchangeService.CompleteExchangeCredit(exchangeID)
			} else if step == "2_credit" {
				c.exchangeService.FinalizeExchange(exchangeID)
			}
		}
	} else {
		reason := statusEvent.Error
		if reason == "" {
			reason = "Payment failed"
		}
		if isFiat && c.fiatExchangeService != nil {
			c.fiatExchangeService.FailFiatExchange(exchangeID, reason)
		} else {
			c.exchangeService.FailExchange(exchangeID, reason)
		}
	}

	return nil
}
