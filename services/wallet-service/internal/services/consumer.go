package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

	// Subscribe to payment events (requests from other services)
	if err := c.kafkaClient.Subscribe(messaging.TopicPaymentEvents, c.handlePaymentEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to payment events: %v", err)
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

	if userEvent.UserID != "" {
		// Determine currency
		currency := "USD" // Default
		
		// Prioritize explicit currency from event if available
		if userEvent.Currency != "" {
			currency = userEvent.Currency
		} else if userEvent.Country != "" {
			currency = getCurrencyByCountry(userEvent.Country)
		}

		// Create default wallet
		name := "Main Wallet"
		desc := fmt.Sprintf("Default %s wallet created on registration", currency)

		req := &models.CreateWalletRequest{
			Currency:    currency,
			WalletType:  "fiat",
			Name:        &name,
			Description: &desc,
		}

		// Internal call, no auth token needed
		if _, err := c.walletService.CreateWallet(userEvent.UserID, req); err != nil {
			log.Printf("Failed to create default wallet (Currency: %s): %v", currency, err)
		}
	}
	
	return nil
}

// getCurrencyByCountry returns the currency code for a given country code (ISO 2 chars)
func getCurrencyByCountry(countryCode string) string {
	switch strings.ToUpper(countryCode) {
	// === EUROPE ===
	case "AT", "BE", "CY", "EE", "FI", "FR", "DE", "GR", "IE", "IT", "LV", "LT", "LU", "MT", "NL", "PT", "SK", "SI", "ES":
		return "EUR"
	case "GB", "GBR": return "GBP"
	case "CH", "CHE": return "CHF"
	case "SE", "SWE": return "SEK"
	case "NO", "NOR": return "NOK"
	case "DK", "DNK": return "DKK"
	case "PL", "POL": return "PLN"
	case "CZ", "CZE": return "CZK"
	case "HU", "HUN": return "HUF"
	case "RO", "ROU": return "RON"
	case "BG", "BGR": return "BGN"
	case "HR", "HRV": return "HRK"
	case "RU", "RUS": return "RUB" // Russia
	case "UA", "UKR": return "UAH"
	case "BY", "BLR": return "BYN"

	// === NORTH AMERICA ===
	case "US", "USA": return "USD"
	case "CA", "CAN": return "CAD"
	case "MX", "MEX": return "MXN"

	// === SOUTH AMERICA ===
	case "BR", "BRA": return "BRL"
	case "AR", "ARG": return "ARS"
	case "CL", "CHL": return "CLP"
	case "CO", "COL": return "COP"
	case "PE", "PER": return "PEN"
	case "UY", "URY": return "UYU"
	case "VE", "VEN": return "VES"
	case "BO", "BOL": return "BOB"
	case "PY", "PRY": return "PYG"

	// === ASIA ===
	case "CN", "CHN": return "CNY" // China
	case "JP", "JPN": return "JPY" // Japan
	case "KR", "KOR": return "KRW" // South Korea
	case "KP", "PRK": return "KPW" // North Korea
	case "IN", "IND": return "INR"
	case "ID", "IDN": return "IDR"
	case "MY", "MYS": return "MYR"
	case "SG", "SGP": return "SGD"
	case "TH", "THA": return "THB"
	case "VN", "VNM": return "VND"
	case "PH", "PHL": return "PHP"
	case "PK", "PAK": return "PKR"
	case "BD", "BGD": return "BDT"
	case "HK", "HKD": return "HKD"
	case "TW", "TWN": return "TWD"
	case "AF", "AFG": return "AFN" // Afghanistan
	case "IR", "IRN": return "IRR" // Iran
	case "IQ", "IRQ": return "IQD"
	case "LB", "LBN": return "LBP"
	case "IL", "ISR": return "ILS" // Israel
	case "SA", "SAU": return "SAR"
	case "AE", "ARE": return "AED"
	case "QA", "QAT": return "QAR"
	case "TR", "TUR": return "TRY"

	// === AFRICA ===
	// CFA Franc BEAC (CEMAC)
	case "CM", "CMR", "CF", "CAF", "TD", "TCD", "CG", "COG", "GA", "GAB", "GQ", "GNQ":
		return "XAF"
	// CFA Franc BCEAO (UEMOA)
	case "CI", "CIV", "BJ", "BEN", "BF", "BFA", "GW", "GNB", "ML", "MLI", "NE", "NER", "SN", "SEN", "TG", "TGO":
		return "XOF"
	case "NG", "NGA": return "NGN" // Nigeria
	case "ZA", "ZAF": return "ZAR" // South Africa
	case "EG", "EGY": return "EGP" // Egypt
	case "MA", "MAR": return "MAD" // Morocco
	case "KE", "KEN": return "KES"
	case "GH", "GHA": return "GHS"
	case "DZ", "DZA": return "DZD"
	case "TN", "TUN": return "TND"
	case "ET", "ETH": return "ETB"
	case "RW", "RWF": return "RWF"
	case "UG", "UGA": return "UGX"
	case "TZ", "TZA": return "TZS"
	case "AO", "AGO": return "AOA"
	case "MZ", "MOZ": return "MZN"
	case "ZW", "ZWE": return "ZWL" // Zimbabwe
	case "CD", "COD": return "CDF"

	// === OCEANIA ===
	case "AU", "AUS": return "AUD"
	case "NZ", "NZL": return "NZD"

	default:
		return "USD"
	}
}

// handlePaymentEvent processes payment.request events from other services (e.g. exchange-service)
func (c *Consumer) handlePaymentEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	// We only care about payment requests here
	if event.Type != messaging.EventPaymentRequest {
		return nil
	}

	log.Printf("[Kafka] Processing payment request event: %s from %s", event.ID, event.Source)

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var req messaging.PaymentRequestEvent
	if err := json.Unmarshal(dataBytes, &req); err != nil {
		log.Printf("Failed to unmarshal payment request event: %v", err)
		return err
	}

	// Helper to publish status
	publishStatus := func(status, errorMsg string) {
		statusEvent := messaging.PaymentStatusEvent{
			RequestID:   req.RequestID,
			ReferenceID: req.ReferenceID,
			Type:        req.Type,
			Status:      status,
			Error:       errorMsg,
		}

		eventType := messaging.EventPaymentSuccess
		if status == "failed" {
			eventType = messaging.EventPaymentFailed
		}

		envelope := messaging.NewEventEnvelope(eventType, "wallet-service", statusEvent)
		// Propagate correlation ID
		if event.CorrelationID != "" {
			envelope.WithCorrelationID(event.CorrelationID)
		}

		c.kafkaClient.Publish(ctx, messaging.TopicPaymentEvents, envelope)
	}

	var updateErr error

	// Handle Debit
	if req.DebitAmount > 0 && req.FromWalletID != "" {
		updateErr = c.walletService.balanceService.UpdateBalance(req.FromWalletID, req.DebitAmount, "debit")
		if updateErr != nil {
			log.Printf("Failed to process payment debit for %s: %v", req.RequestID, updateErr)
			publishStatus("failed", fmt.Sprintf("Debit failed: %v", updateErr))
			return nil
		}
	}

	// Handle Credit
	if req.CreditAmount > 0 && req.ToWalletID != "" {
		updateErr = c.walletService.balanceService.UpdateBalance(req.ToWalletID, req.CreditAmount, "credit")
		if updateErr != nil {
			log.Printf("Failed to process payment credit for %s: %v", req.RequestID, updateErr)
			
			// If we also did a debit, we should ideally rollback, but for MVP we log critical error.
			// In a real system, we'd use a saga or 2PC, or reverse the debit.
			if req.DebitAmount > 0 {
				log.Printf("CRITICAL: Debit succeeded but credit failed for %s. Funds may be lost/stuck.", req.RequestID)
				// Attempt rollback
				c.walletService.balanceService.UpdateBalance(req.FromWalletID, req.DebitAmount, "credit")
			}
			
			publishStatus("failed", fmt.Sprintf("Credit failed: %v", updateErr))
			return nil
		}
	}

	// If successful
	log.Printf("[Kafka] Successfully processed payment request %s", req.RequestID)
	publishStatus("success", "")

	return nil
}
