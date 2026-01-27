package messaging

import (
	"time"
)

// EventEnvelope wraps all events with standard metadata
// Implements the Event Envelope design pattern for traceability and versioning
type EventEnvelope struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Source        string                 `json:"source"`
	Timestamp     time.Time              `json:"timestamp"`
	Version       string                 `json:"version"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	Data          interface{}            `json:"data"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// NewEventEnvelope creates a new event envelope with the given type and data
func NewEventEnvelope(eventType, source string, data interface{}) *EventEnvelope {
	return &EventEnvelope{
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
		Data:      data,
		Metadata:  make(map[string]interface{}),
	}
}

// WithCorrelationID sets the correlation ID for event tracing
func (e *EventEnvelope) WithCorrelationID(id string) *EventEnvelope {
	e.CorrelationID = id
	return e
}

// WithMetadata adds metadata to the event
func (e *EventEnvelope) WithMetadata(key string, value interface{}) *EventEnvelope {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// === User Events ===

type UserRegisteredEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Currency  string `json:"currency"`
	Country   string `json:"country"`
}

// === Wallet Events ===

type WalletCreatedEvent struct {
	WalletID   string  `json:"wallet_id"`
	UserID     string  `json:"user_id"`
	Currency   string  `json:"currency"`
	WalletType string  `json:"wallet_type"`
	Balance    float64 `json:"balance"`
}

type WalletBalanceUpdatedEvent struct {
	WalletID   string  `json:"wallet_id"`
	UserID     string  `json:"user_id"`
	OldBalance float64 `json:"old_balance"`
	NewBalance float64 `json:"new_balance"`
	Amount     float64 `json:"amount"`
	Operation  string  `json:"operation"` // credit, debit
	Reference  string  `json:"reference"`
}

// === Transfer Events ===

type TransferInitiatedEvent struct {
	TransferID   string  `json:"transfer_id"`
	FromUserID   string  `json:"from_user_id"`
	ToUserID     string  `json:"to_user_id"`
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Fee          float64 `json:"fee"`
}

type TransferCompletedEvent struct {
	TransferID   string  `json:"transfer_id"`
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
	Fee          float64 `json:"fee"`
	Status       string  `json:"status"`
}

// === Exchange Events ===

type ExchangeCompletedEvent struct {
	ExchangeID   string  `json:"exchange_id"`
	UserID       string  `json:"user_id"`
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	FromAmount   float64 `json:"from_amount"`
	ToAmount     float64 `json:"to_amount"`
	ExchangeRate float64 `json:"exchange_rate"`
	Fee          float64 `json:"fee"`
}

type FiatExchangeRequestEvent struct {
	ExchangeID   string  `json:"exchange_id"`
	UserID       string  `json:"user_id"`
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	FromAmount   float64 `json:"from_amount"`
	ToAmount     float64 `json:"to_amount"`
	Fee          float64 `json:"fee"`
}

// === Payment Events ===

type PaymentRequestEvent struct {
	RequestID         string                 `json:"request_id"`
	ReferenceID       string                 `json:"reference_id"`
	Type              string                 `json:"type"` // exchange, transfer
	UserID            string                 `json:"user_id"`
	FromWalletID      string                 `json:"from_wallet_id"`
	DestinationUserID string                 `json:"destination_user_id,omitempty"` // For auto-resolution
	ToWalletID        string                 `json:"to_wallet_id,omitempty"`
	DebitAmount       float64                `json:"debit_amount"`
	CreditAmount      float64                `json:"credit_amount,omitempty"`
	Currency          string                 `json:"currency"`
	Description       string                 `json:"description,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentStatusEvent struct {
	RequestID   string `json:"request_id"`
	ReferenceID string `json:"reference_id"`
	Type        string `json:"type"`
	Status      string `json:"status"` // success, failed
	Error       string `json:"error,omitempty"`
}

// === Card Events ===

type CardLoadedEvent struct {
	CardID         string  `json:"card_id"`
	UserID         string  `json:"user_id"`
	SourceWalletID string  `json:"source_wallet_id"`
	Amount         float64 `json:"amount"`
	Fee            float64 `json:"fee"`
}

type CardTransactionEvent struct {
	CardID        string  `json:"card_id"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Merchant      string  `json:"merchant"`
	Status        string  `json:"status"`
}

// === Notification Events ===

type NotificationEvent struct {
	UserID    string                 `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	ActionUrl string                 `json:"action_url,omitempty"`
}
