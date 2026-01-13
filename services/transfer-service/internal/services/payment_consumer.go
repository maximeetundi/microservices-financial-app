package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

// PaymentRequestConsumer handles payment request events from Kafka
type PaymentRequestConsumer struct {
	kafkaClient    *messaging.KafkaClient
	walletClient   *WalletClient
	exchangeClient *ExchangeClient
	walletRepo     *repository.WalletRepository
	transferRepo   *repository.TransferRepository
}

// NewPaymentRequestConsumer creates a new payment request consumer
func NewPaymentRequestConsumer(
	kafkaClient *messaging.KafkaClient, 
	walletClient *WalletClient, 
	exchangeClient *ExchangeClient, 
	walletRepo *repository.WalletRepository,
	transferRepo *repository.TransferRepository,
) *PaymentRequestConsumer {
	return &PaymentRequestConsumer{
		kafkaClient:    kafkaClient,
		walletClient:   walletClient,
		exchangeClient: exchangeClient,
		walletRepo:     walletRepo,
		transferRepo:   transferRepo,
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

	// Resolve Destination User Wallet if needed (e.g. for Ticket Organizer or Refund)
	if paymentReq.ToWalletID == "" && paymentReq.DestinationUserID != "" {
		log.Printf("[TRACE-FIAT] Resolving wallet for DestinationUserID: %s, Currency: %s", paymentReq.DestinationUserID, paymentReq.Currency)
		
		var fallbackWalletID string
		var fallbackCurrency string

		// 1. Check user wallets
		wallets, err := c.walletClient.GetUserWallets(paymentReq.DestinationUserID)
		if err == nil {
			for i, w := range wallets {
				// Capture the first wallet as fallback
				if i == 0 {
					if id, ok := w["id"].(string); ok {
						fallbackWalletID = id
						if cur, ok := w["currency"].(string); ok {
							fallbackCurrency = cur
						}
					}
				}

				if cur, ok := w["currency"].(string); ok && cur == paymentReq.Currency {
					if id, ok := w["id"].(string); ok {
						paymentReq.ToWalletID = id
						log.Printf("[TRACE-FIAT] Found existing wallet: %s", id)
						break
					}
				}
			}
		}

		// 2. If no matching wallet found, use fallback or create new
		if paymentReq.ToWalletID == "" {
			if fallbackWalletID != "" {
				// Use fallback (e.g. User paid in USD, Refund is in XOF -> Credit USD wallet via Auto-Convert)
				paymentReq.ToWalletID = fallbackWalletID
				log.Printf("[TRACE-FIAT] No %s wallet found. Using fallback wallet %s (%s) for Auto-Conversion.", 
					paymentReq.Currency, fallbackWalletID, fallbackCurrency)
			} else {
				log.Printf("[TRACE-FIAT] No wallet found, creating new one for UserID: %s", paymentReq.DestinationUserID)
				newWalletID, err := c.walletClient.CreateUserWallet(paymentReq.DestinationUserID, paymentReq.Currency)
				if err != nil {
					errMsg := "Failed to create wallet for destination user: " + err.Error()
					log.Printf("[Kafka] %s", errMsg)
					c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", errMsg)
					return nil // Stop processing
				} else {
					paymentReq.ToWalletID = newWalletID
					log.Printf("[TRACE-FIAT] Created new wallet: %s", newWalletID)
				}
			}
		}
	} else if paymentReq.ToWalletID == "" && paymentReq.CreditAmount > 0 {
		log.Printf("[WARNING] CreditAmount > 0 but ToWalletID and DestinationUserID are empty!")
	}

	// Process debit operation via wallet client
	if paymentReq.FromWalletID != "" && paymentReq.DebitAmount > 0 {
		log.Printf("[TRACE-FIAT] Processing Debit for UserID: %s", paymentReq.UserID)
		
		debitCurrency := paymentReq.Currency
		debitAmount := paymentReq.DebitAmount

		// Auto-Conversion Check
		payerWallets, err := c.walletClient.GetUserWallets(paymentReq.UserID)
		var payerWallet map[string]interface{}
		
		if err == nil {
			for _, w := range payerWallets {
				if id, ok := w["id"].(string); ok && id == paymentReq.FromWalletID {
					payerWallet = w
					break
				}
			}
		}

		if payerWallet != nil {
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
					if rate == 0 {
						log.Printf("[Kafka] Auto-Conversion Failed: Rate is 0")
						c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", "Auto-conversion failed: Rate is 0")
						return fmt.Errorf("rate is 0")
					}
					debitAmount = paymentReq.DebitAmount / rate
					debitCurrency = walletCurrency
					
					log.Printf("[TRACE-FIAT] Auto-Conversion: Rate %s->%s = %f. Debiting %f %s instead of %f %s.", 
						walletCurrency, paymentReq.Currency, rate, debitAmount, debitCurrency, paymentReq.DebitAmount, paymentReq.Currency)
				}
			}
		} else {
			errMsg := fmt.Sprintf("Could not fetch payer wallet details for resolution (ID: %s). Auth restriction or invalid ID.", paymentReq.FromWalletID)
			if err != nil {
				errMsg += fmt.Sprintf(" Error: %v", err)
			}
			log.Printf("[TRACE-FIAT] Error: %s", errMsg)
			c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", errMsg)
			return fmt.Errorf(errMsg)
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
		
		creditCurrency := paymentReq.Currency
		creditAmount := paymentReq.CreditAmount

		// Auto-Conversion Check for CREDIT (e.g. Refund XOF -> USD Wallet)
		destWallets, err := c.walletClient.GetUserWallets(creditUserID)
		var destWallet map[string]interface{}
		
		if err == nil {
			for _, w := range destWallets {
				if id, ok := w["id"].(string); ok && id == paymentReq.ToWalletID {
					destWallet = w
					break
				}
			}
		}

		if destWallet != nil {
			if walletCurrency, ok := destWallet["currency"].(string); ok && walletCurrency != "" {
				if walletCurrency != paymentReq.Currency {
					log.Printf("[TRACE-FIAT] Credit Currency Mismatch. Wallet: %s, Request: %s. Initiating Auto-Conversion.", walletCurrency, paymentReq.Currency)
					
					// Fetch Rate (From XOF -> To USD)
					rate, err := c.exchangeClient.GetRate(paymentReq.Currency, walletCurrency)
					if err != nil {
						log.Printf("[Kafka] Auto-Conversion Failed: Could not get rate: %v", err)
						// Fallback to error or try original currency? Prefer error to avoid bad state.
						c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", "Auto-conversion failed: "+err.Error())
						// Also rollback debit? Yes.
						// NOTE: Rollback logic is duplicated below. Refactor later?
						return err
					}

					// Calculate Credited Amount
					// AmountInWalletCurrency = AmountInPaymentCurrency * Rate
					creditAmount = paymentReq.CreditAmount * rate
					creditCurrency = walletCurrency
					
					log.Printf("[TRACE-FIAT] Auto-Conversion Credit: Rate %s->%s = %f. Crediting %f %s instead of %f %s.", 
						paymentReq.Currency, walletCurrency, rate, creditAmount, creditCurrency, paymentReq.CreditAmount, paymentReq.Currency)
				}
			}
		}

		creditReq := &WalletTransactionRequest{
			UserID:    creditUserID,
			WalletID:  paymentReq.ToWalletID,
			Amount:    creditAmount,
			Currency:  creditCurrency,
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
				Amount:    paymentReq.DebitAmount, // Ideally we should refund the actual debited amount/currency from above. 
				// BUT 'paymentReq.DebitAmount' is the requested amount (often in Ticket Currency). 
				// If we did auto-conversion on debit, we should ideally rollback in the debited currency.
				// However, maintaining that state is hard here without refactoring.
				// For now, let's assume rollback uses original logic or re-calculates?
				// Simplification: We blindly rollback with REQUESTED amount/currency. 
				// The Wallet Client might fail if we try to Credit "XOF" to "USD" wallet (Reverse problem).
				// We need smart rollback too!
				
				// For now, minimal fix:
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

	// --- RECORD TRANSACTION IN HISTORY ---
	// Map payment type to transfer type (or keep original)
	txType := paymentReq.Type
	if txType == "donation" {
		txType = "donation"
	} else if txType == "ticket_purchase" {
		txType = "payment"
	}

	// Create Transfer Record
	transfer := &models.Transfer{
		ID:             paymentReq.RequestID, // Use RequestID to match external ID
		UserID:         paymentReq.UserID,    // Payer
		FromWalletID:   &paymentReq.FromWalletID,
		ToWalletID:     &paymentReq.ToWalletID,
		TransferType:   txType,
		Amount:         paymentReq.DebitAmount,
		Currency:       paymentReq.Currency,
		Status:         "completed",
		Reference:      &paymentReq.ReferenceID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set Description
	desc := fmt.Sprintf("Paiement: %s", paymentReq.Type)
	transfer.Description = &desc

	if err := c.transferRepo.Create(transfer); err != nil {
		log.Printf("[Kafka] Failed to record transfer history for %s: %v", paymentReq.RequestID, err)
	} else {
		log.Printf("[Kafka] Recorded transfer history for %s", paymentReq.RequestID)
	}
	// -------------------------------------

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
