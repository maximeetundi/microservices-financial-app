package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

type PaymentRequestConsumer struct {
	mqClient     *database.RabbitMQClient
	walletClient *WalletClient
	walletRepo   *repository.WalletRepository
}

func NewPaymentRequestConsumer(mqClient *database.RabbitMQClient, walletClient *WalletClient, walletRepo *repository.WalletRepository) *PaymentRequestConsumer {
	return &PaymentRequestConsumer{
		mqClient:     mqClient,
		walletClient: walletClient,
		walletRepo:   walletRepo,
	}
}

func (c *PaymentRequestConsumer) Start() error {
	msgs, err := c.mqClient.Consume("payments")
	if err != nil {
		return fmt.Errorf("failed to start consuming payments queue: %w", err)
	}

	go func() {
		for d := range msgs {
			var req models.PaymentRequestEvent
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Error unmarshalling payment request: %v", err)
				d.Ack(false) // Ack to remove bad message
				continue
			}

			log.Printf("Processing payment request: %s from service: %s", req.TransactionID, req.OriginService)

			err := c.processPayment(&req)
			
			status := "completed"
			errorMsg := ""
			if err != nil {
				log.Printf("Payment processing failed for %s: %v", req.TransactionID, err)
				status = "failed"
				errorMsg = err.Error()
			}

			// Publish status event
			statusEvent := models.PaymentStatusEvent{
				TransactionID: req.TransactionID,
				Status:        status,
				Timestamp:     time.Now(),
				Error:         errorMsg,
			}

			eventBytes, _ := json.Marshal(statusEvent)
			
			// Routing key: payment.status.{completed|failed}
			routingKey := fmt.Sprintf("payment.status.%s", status)
			c.mqClient.PublishToExchange("payment.events", routingKey, eventBytes)

			d.Ack(false)
		}
	}()

	log.Println("Payment Request Consumer started")
	return nil
}

func (c *PaymentRequestConsumer) processPayment(req *models.PaymentRequestEvent) error {
	opsPerformed := 0

	// 1. Debit Source Wallet (if provided)
	if req.SourceWalletID != "" {
		debitReq := &WalletTransactionRequest{
			UserID:    req.UserID,
			WalletID:  req.SourceWalletID,
			Amount:    req.Amount,
			Type:      "debit",
			Currency:  req.Currency,
			Reference: req.Reference,
		}

		if err := c.walletClient.ProcessTransaction(debitReq); err != nil {
			return fmt.Errorf("failed to debit source wallet: %w", err)
		}
		opsPerformed++
	}

	// 2. Resolve Destination Wallet if missing
	if req.DestinationWalletID == "" && req.DestinationUserID != "" {
		walletID, err := c.walletRepo.GetWalletIDByUserAndCurrency(req.DestinationUserID, req.Currency)
		if err == nil && walletID != "" {
			req.DestinationWalletID = walletID
		} else {
			// If wallet not found, fail.
			return fmt.Errorf("destination wallet not found for user %s and currency %s", req.DestinationUserID, req.Currency)
		}
	}

	// 3. Credit Destination Wallet (if provided or resolved)
	if req.DestinationWalletID != "" {
		if req.DestinationUserID == "" {
			return fmt.Errorf("missing destination user id for credit")
		}

		creditReq := &WalletTransactionRequest{
			UserID:    req.DestinationUserID,
			WalletID:  req.DestinationWalletID,
			Amount:    req.Amount,
			Type:      "credit",
			Currency:  req.Currency,
			Reference: req.Reference,
		}

		if err := c.walletClient.ProcessTransaction(creditReq); err != nil {
			// Compensation logic: credit back the source if it was debited in this same transaction
			if req.SourceWalletID != "" {
				compReq := &WalletTransactionRequest{
					UserID:    req.UserID,
					WalletID:  req.SourceWalletID,
					Amount:    req.Amount,
					Type:      "credit",
					Currency:  req.Currency,
					Reference: "REFUND_" + req.Reference,
				}
				_ = c.walletClient.ProcessTransaction(compReq)
			}

			return fmt.Errorf("failed to credit destination wallet: %w", err)
		}
		opsPerformed++
	}

	if opsPerformed == 0 {
		return fmt.Errorf("invalid payment request: neither source nor destination wallet provided")
	}

	return nil
}
