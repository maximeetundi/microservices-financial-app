package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
)

type PaymentStatusConsumer struct {
	mqClient            *database.RabbitMQClient
	exchangeService     *ExchangeService
	fiatExchangeService *FiatExchangeService
}

func NewPaymentStatusConsumer(mqClient *database.RabbitMQClient, exchangeService *ExchangeService) *PaymentStatusConsumer {
	return &PaymentStatusConsumer{
		mqClient:        mqClient,
		exchangeService: exchangeService,
	}
}

// SetFiatExchangeService allows setting the fiat exchange service after construction
func (c *PaymentStatusConsumer) SetFiatExchangeService(fiatService *FiatExchangeService) {
	c.fiatExchangeService = fiatService
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

			txID := event.TransactionID
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
			// Check for fiat exchange transactions
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
				d.Ack(false)
				continue
			}

			log.Printf("Received payment status for %s Exchange %s Step %s: %s",
				map[bool]string{true: "Fiat", false: "Crypto"}[isFiat], exchangeID, step, event.Status)

			if event.Status == "completed" {
				if isFiat && c.fiatExchangeService != nil {
					// Handle fiat exchange
					if step == "1_debit" {
						c.fiatExchangeService.CompleteFiatExchangeCredit(exchangeID)
					} else if step == "2_credit" {
						c.fiatExchangeService.FinalizeFiatExchange(exchangeID)
					}
				} else {
					// Handle crypto exchange
					if step == "1_debit" {
						c.exchangeService.CompleteExchangeCredit(exchangeID)
					} else if step == "2_credit" {
						c.exchangeService.FinalizeExchange(exchangeID)
					}
				}
			} else {
				// Failed
				reason := event.Error
				if reason == "" {
					reason = "Payment failed"
				}
				if isFiat && c.fiatExchangeService != nil {
					c.fiatExchangeService.FailFiatExchange(exchangeID, reason)
				} else {
					c.exchangeService.FailExchange(exchangeID, reason)
				}
			}

			d.Ack(false)
		}
	}()

	log.Println("Payment Status Consumer started")
	return nil
}

