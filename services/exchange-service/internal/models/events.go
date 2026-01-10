package models

import "time"

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

type PaymentStatusEvent struct {
	TransactionID string    `json:"request_id"` // Matches messaging.PaymentStatusEvent.RequestID
	Status        string    `json:"status"` // completed, failed
	TxHash        string    `json:"tx_hash,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	Error         string    `json:"error,omitempty"`
}
