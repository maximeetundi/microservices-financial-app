package models

import "time"

type PaymentRequestEvent struct {
	TransactionID       string                 `json:"transaction_id"`
	SourceWalletID      string                 `json:"source_wallet_id"`
	UserID              string                 `json:"user_id"` // Owner of source wallet
	Amount              float64                `json:"amount"`
	Currency            string                 `json:"currency"`
	DestinationWalletID string                 `json:"destination_wallet_id"` // Optional
	DestinationUserID   string                 `json:"destination_user_id"`   // Required if DestinationWalletID is present
	Reference           string                 `json:"reference"`
	OriginService       string                 `json:"origin_service"`
	MetaData            map[string]interface{} `json:"meta_data"`
}

type PaymentStatusEvent struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"` // completed, failed
	TxHash        string    `json:"tx_hash,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	Error         string    `json:"error,omitempty"`
}
