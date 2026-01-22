package models

import "time"

// PaymentRequestEvent is sent to wallet service for payment processing
type PaymentRequestEvent struct {
	TransactionID       string                 `json:"transaction_id"`
	SourceWalletID      string                 `json:"source_wallet_id"`
	UserID              string                 `json:"user_id"`
	Amount              float64                `json:"amount"`
	Currency            string                 `json:"currency"`
	DestinationWalletID string                 `json:"destination_wallet_id"`
	DestinationUserID   string                 `json:"destination_user_id"`
	Reference           string                 `json:"reference"`
	OriginService       string                 `json:"origin_service"`
	MetaData            map[string]interface{} `json:"meta_data"`
}

// PaymentStatusEvent is received from wallet service after payment processing
type PaymentStatusEvent struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"` // completed, failed
	TxHash        string    `json:"tx_hash,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	Error         string    `json:"error,omitempty"`
}

// NotificationEvent is sent to notification service
type NotificationEvent struct {
	UserID    string                 `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}
