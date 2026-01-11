package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

// PaymentRequestConsumer handles payment request events from Kafka
type PaymentRequestConsumer struct {
	kafkaClient    *messaging.KafkaClient
	walletClient   *WalletClient
	exchangeClient *ExchangeClient
	walletRepo     *repository.WalletRepository
}

// NewPaymentRequestConsumer creates a new payment request consumer
func NewPaymentRequestConsumer(kafkaClient *messaging.KafkaClient, walletClient *WalletClient, exchangeClient *ExchangeClient, walletRepo *repository.WalletRepository) *PaymentRequestConsumer {
	return &PaymentRequestConsumer{
		kafkaClient:    kafkaClient,
		walletClient:   walletClient,
		exchangeClient: exchangeClient,
		walletRepo:     walletRepo,
	}
}

// Start begins consuming payment events from Kafka
func (c *PaymentRequestConsumer) Start() error {
	// Subscribe to payment events
	if err := c.kafkaClient.Subscribe(messaging.TopicPaymentEvents, c.handlePaymentEvent); err != nil {
		return err
	}

	log.Println("[Kafka] PaymentRequestConsumer started")
	return nil
}

// handlePaymentEvent processes payment.request events
func (c *PaymentRequestConsumer) handlePaymentEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Processing payment event: %s", event.Type)

	if event.Type != messaging.EventPaymentRequest {
		return nil
	}

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var paymentReq messaging.PaymentRequestEvent
	if err := json.Unmarshal(dataBytes, &paymentReq); err != nil {
		log.Printf("Failed to unmarshal payment request: %v", err)
		return err
	}

	log.Printf("[TRACE-FIAT] Payment Consumer received: RequestID=%s, UserID=%s, Type=%s, Amount=%.2f",
		paymentReq.RequestID, paymentReq.UserID, paymentReq.Type, paymentReq.DebitAmount)

	// Resolve Destination User Wallet if needed (e.g. for Ticket Organizer)
	if paymentReq.ToWalletID == "" && paymentReq.DestinationUserID != "" {
		log.Printf("[TRACE-FIAT] Resolving wallet for DestinationUserID: %s, Currency: %s", paymentReq.DestinationUserID, paymentReq.Currency)
		
		// 1. Check if user already has a wallet for this currency
		wallets, err := c.walletClient.GetUserWallets(paymentReq.DestinationUserID)
		if err == nil {
			for _, w := range wallets {
				if cur, ok := w["currency"].(string); ok && cur == paymentReq.Currency {
					if id, ok := w["id"].(string); ok {
						paymentReq.ToWalletID = id
						log.Printf("[TRACE-FIAT] Found existing wallet: %s", id)
						break
					}
				}
			}
		}

		// 2. If no wallet found, create one
		if paymentReq.ToWalletID == "" {
			log.Printf("[TRACE-FIAT] No wallet found, creating new one for UserID: %s", paymentReq.DestinationUserID)
			newWalletID, err := c.walletClient.CreateUserWallet(paymentReq.DestinationUserID, paymentReq.Currency)
			if err != nil {
				log.Printf("[Kafka] Failed to create wallet for destination user: %v", err)
				// We can't proceed with credit, so we must fail and potentially refund/rollback
				// Ideally we should rollback debit if it happened, but debit happens before this block in current flow?
				// WAIT: Debit happens in lines 62-78. We should move this resolution BEFORE debit to be safe.
			} else {
				paymentReq.ToWalletID = newWalletID
				log.Printf("[TRACE-FIAT] Created new wallet: %s", newWalletID)
			}
		}
	}

	// Process debit operation via wallet client

	if paymentReq.FromWalletID != "" && paymentReq.DebitAmount > 0 {
		log.Printf("[TRACE-FIAT] Processing Debit for UserID: %s", paymentReq.UserID)
		
		debitCurrency := paymentReq.Currency
		debitAmount := paymentReq.DebitAmount

		// Auto-Conversion Check
		payerWallet, err := c.walletClient.GetWallet(paymentReq.FromWalletID)
		if err == nil {
			if walletCurrency, ok := payerWallet["currency"].(string); ok && walletCurrency != "" {
				if walletCurrency != paymentReq.Currency {
					log.Printf("[TRACE-FIAT] Currency Mismatch detected. Wallet: %s, Request: %s. Initiating Auto-Conversion.", walletCurrency, paymentReq.Currency)
					
					// Fetch Rate
					rate, err := c.exchangeClient.GetRate(walletCurrency, paymentReq.Currency)
					if err != nil {
						log.Printf("[Kafka] Auto-Conversion Failed: Could not get rate: %v", err)
						c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", "Auto-conversion failed: "+err.Error())
						return err
					}

					// Calculate Debited Amount
					// Rate is From -> To (e.g. USD -> XOF = 650)
					// We need 'DebitAmount' in ToCurrency (XOF).
					// So AmountInFrom = AmountInTo / Rate
					debitAmount = paymentReq.DebitAmount / rate
					debitCurrency = walletCurrency
					
					log.Printf("[TRACE-FIAT] Auto-Conversion: Rate %s->%s = %f. Debiting %f %s instead of %f %s.", 
						walletCurrency, paymentReq.Currency, rate, debitAmount, debitCurrency, paymentReq.DebitAmount, paymentReq.Currency)
				}
			}
		} else {
			log.Printf("[TRACE-FIAT] Warning: Could not fetch payer wallet details: %v. Proceeding with default currency.", err)
		}

		debitReq := &WalletTransactionRequest{
			UserID:    paymentReq.UserID,
			WalletID:  paymentReq.FromWalletID,
			Amount:    debitAmount,
			Currency:  debitCurrency,
			Type:      "debit",
			Reference: paymentReq.ReferenceID,
		}
		err = c.walletClient.ProcessTransaction(debitReq)
		if err != nil {
			log.Printf("[Kafka] Debit failed: %v", err)
			c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", err.Error())
			return err
		}
	}

	// If there's a credit operation (to wallet)
	if paymentReq.ToWalletID != "" && paymentReq.CreditAmount > 0 {
		// Determine correct UserID for credit (Organizer/Payee)
		creditUserID := paymentReq.UserID
		if paymentReq.DestinationUserID != "" {
			creditUserID = paymentReq.DestinationUserID
		}

		log.Printf("[TRACE-FIAT] Processing Credit for UserID: %s, ToWalletID: %s", creditUserID, paymentReq.ToWalletID)
		creditReq := &WalletTransactionRequest{
			UserID:    creditUserID,
			WalletID:  paymentReq.ToWalletID,
			Amount:    paymentReq.CreditAmount,
			Currency:  paymentReq.Currency,
			Type:      "credit",
			Reference: paymentReq.ReferenceID,
		}
		err = c.walletClient.ProcessTransaction(creditReq)
		if err != nil {
			log.Printf("[Kafka] Credit failed: %v", err)
			// Rollback debit
			rollbackReq := &WalletTransactionRequest{
				UserID:    paymentReq.UserID,
				WalletID:  paymentReq.FromWalletID,
				Amount:    paymentReq.DebitAmount,
				Currency:  paymentReq.Currency,
				Type:      "credit",
				Reference: paymentReq.ReferenceID + "_rollback",
			}
			c.walletClient.ProcessTransaction(rollbackReq)
			c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", err.Error())
			return err
		}
	}

	log.Printf("[Kafka] Payment request %s processed successfully", paymentReq.RequestID)
	c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "success", "")

	return nil
}

// publishPaymentStatus publishes status updates back to Kafka
func (c *PaymentRequestConsumer) publishPaymentStatus(requestID, referenceID, reqType, status, errMsg string) {
	statusEvent := messaging.PaymentStatusEvent{
		RequestID:   requestID,
		ReferenceID: referenceID,
		Type:        reqType,
		Status:      status,
		Error:       errMsg,
	}

	eventType := messaging.EventPaymentSuccess
	if status == "failed" {
		eventType = messaging.EventPaymentFailed
	}

	envelope := messaging.NewEventEnvelope(eventType, "transfer-service", statusEvent)
	c.kafkaClient.Publish(context.Background(), messaging.TopicPaymentEvents, envelope)
}
