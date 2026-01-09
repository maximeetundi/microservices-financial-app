package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
)

type PaymentStatusConsumer struct {
	mqClient        *database.RabbitMQClient
	exchangeService *ExchangeService
}

func NewPaymentStatusConsumer(mqClient *database.RabbitMQClient, exchangeService *ExchangeService) *PaymentStatusConsumer {
	return &PaymentStatusConsumer{
		mqClient:        mqClient,
		exchangeService: exchangeService,
	}
}

func (c *PaymentStatusConsumer) Start() error {
	msgs, err := c.mqClient.Consume("exchange.payment.updates")
	if err != nil {
		return fmt.Errorf("failed to start consuming exchange payment updates queue: %w", err)
	}

	go func() {
		for d := range msgs {
			var event models.PaymentStatusEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshalling payment status event: %v", err)
				d.Ack(false)
				continue
			}

			// Logic to check if this transaction belongs to exchange service
			// We check if we can retrieve meta_data or simple ID parsing if metadata isn't on status.
			// The Status Event currently only has TransactionID.
			// But `processExchange` creating IDs like `TX-EXC-DEBIT-...` 
			
			// We need to parse TransactionID or better, rely on Metadata passing through?
			// `PaymentStatusEvent` struct I defined in transfer-service has NO metadata.
			// This is a limitation.
			// But I used `TransactionID` in `processExchange`:
			// `fmt.Sprintf("TX-EXC-DEBIT-%s", exchange.ID)`
			// `fmt.Sprintf("TX-EXC-CREDIT-%s", exchange.ID)`
			
			// So I can parse string.
			
			txID := event.TransactionID
			var exchangeID, step string
			
			// Simple string check
			if len(txID) > 13 && txID[:13] == "TX-EXC-DEBIT-" {
				exchangeID = txID[13:]
				step = "1_debit"
			} else if len(txID) > 14 && txID[:14] == "TX-EXC-CREDIT-" {
				exchangeID = txID[14:]
				step = "2_credit"
			} else {
				// Not an exchange transaction
				d.Ack(false)
				continue
			}

			log.Printf("Received payment status for Exchange %s Step %s: %s", exchangeID, step, event.Status)

			if event.Status == "completed" {
				if step == "1_debit" {
					c.exchangeService.CompleteExchangeCredit(exchangeID)
				} else if step == "2_credit" {
					c.exchangeService.FinalizeExchange(exchangeID)
				}
			} else {
				// Failed
				reason := event.Error
				if reason == "" {
					reason = "Payment failed"
				}
				c.exchangeService.FailExchange(exchangeID, reason)
			}

			d.Ack(false)
		}
	}()

	log.Println("Payment Status Consumer started")
	return nil
}
