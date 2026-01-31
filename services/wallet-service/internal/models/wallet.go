package models

import (
	"time"
)

// Wallet represents a user's wallet (fiat or crypto)
// SECURITY: Private keys for crypto wallets are ALWAYS encrypted with AES-256-GCM
// The plaintext private key is NEVER stored in the database
type Wallet struct {
	ID                  string  `json:"id" db:"id"`
	UserID              string  `json:"user_id" db:"user_id"`
	Currency            string  `json:"currency" db:"currency"`
	WalletType          string  `json:"wallet_type" db:"wallet_type"` // fiat, crypto
	Status              string  `json:"status" db:"status"`           // active, frozen, blocked, closed
	Balance             float64 `json:"balance" db:"balance"`
	FrozenBalance       float64 `json:"frozen_balance" db:"frozen_balance"`
	WalletAddress       *string `json:"wallet_address,omitempty" db:"wallet_address"`
	PrivateKeyEncrypted *string `json:"-" db:"private_key_encrypted"` // NEVER exposed - encrypted with vault
	KeyHash             *string `json:"-" db:"key_hash"`              // SHA-256 hash for verification
	Name                *string `json:"name,omitempty" db:"name"`
	IsActive            bool    `json:"is_active" db:"is_active"`
	IsHidden            bool    `json:"is_hidden" db:"is_hidden"`
	ExternalID          string  `json:"external_id" db:"external_id"` // Blockchain address or internal reference (self-custody)

	// Hybrid deposit system fields
	DepositMemo        *string    `json:"deposit_memo,omitempty" db:"deposit_memo"`       // Unique memo for XRP/XLM/TON deposits
	SweepStatus        string     `json:"sweep_status" db:"sweep_status"`                 // none, pending, completed
	LastSweptAt        *time.Time `json:"last_swept_at,omitempty" db:"last_swept_at"`     // Last sweep timestamp
	PendingSweepAmount float64    `json:"pending_sweep_amount" db:"pending_sweep_amount"` // Amount waiting to be swept

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID               string    `json:"id" db:"id"`
	FromWalletID     *string   `json:"from_wallet_id,omitempty" db:"from_wallet_id"`
	ToWalletID       *string   `json:"to_wallet_id,omitempty" db:"to_wallet_id"`
	TransactionType  string    `json:"transaction_type" db:"transaction_type"` // transfer, exchange, deposit, withdrawal
	Amount           float64   `json:"amount" db:"amount"`
	Fee              float64   `json:"fee" db:"fee"`
	Currency         string    `json:"currency" db:"currency"`
	Status           string    `json:"status" db:"status"` // pending, completed, failed, cancelled
	BlockchainTxHash *string   `json:"blockchain_tx_hash,omitempty" db:"blockchain_tx_hash"`
	ReferenceID      *string   `json:"reference_id,omitempty" db:"reference_id"`
	Description      *string   `json:"description,omitempty" db:"description"`
	Metadata         *string   `json:"metadata,omitempty" db:"metadata"` // JSON string
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	// Virtual fields for API responses
	FromWallet *Wallet `json:"from_wallet,omitempty" db:"-"`
	ToWallet   *Wallet `json:"to_wallet,omitempty" db:"-"`
}

type Balance struct {
	Currency        string  `json:"currency"`
	Available       float64 `json:"available"`
	Frozen          float64 `json:"frozen"`
	Total           float64 `json:"total"`
	PendingDeposits float64 `json:"pending_deposits"`
}

type CreateWalletRequest struct {
	Currency    string  `json:"currency" binding:"required"`
	WalletType  string  `json:"wallet_type" binding:"required,oneof=fiat crypto"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type SendCryptoRequest struct {
	ToAddress string  `json:"to_address" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	GasPrice  *int64  `json:"gas_price,omitempty"`
	Note      *string `json:"note,omitempty"`
}

type CryptoTransactionEstimate struct {
	EstimatedFee   float64 `json:"estimated_fee"`
	EstimatedTotal float64 `json:"estimated_total"`
	GasPrice       *int64  `json:"gas_price,omitempty"`
	GasLimit       *int64  `json:"gas_limit,omitempty"`
	Currency       string  `json:"currency"`
}

type CryptoAddress struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
	Network  string `json:"network"`
	QRCode   string `json:"qr_code,omitempty"`
}

type BlockchainConfirmation struct {
	TxHash        string `json:"tx_hash"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	ToAddress     string `json:"to_address"`
	FromAddress   string `json:"from_address"`
	Confirmations int    `json:"confirmations"`
	BlockHeight   int64  `json:"block_height"`
	Status        string `json:"status"`
}

type WalletStats struct {
	TotalBalance      float64            `json:"total_balance"`
	TotalTransactions int                `json:"total_transactions"`
	BalanceByType     map[string]float64 `json:"balance_by_type"`
	RecentActivity    []Transaction      `json:"recent_activity"`
}

// TransferRequest for internal wallet-to-wallet transfers
type TransferRequest struct {
	FromWalletID string  `json:"from_wallet_id" binding:"required"`
	ToWalletID   string  `json:"to_wallet_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required"`
	Description  string  `json:"description,omitempty"`
}

// TransactionRequest for external service integration (like transfer-service)
type TransactionRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	WalletID  string  `json:"wallet_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Type      string  `json:"type" binding:"required,oneof=debit credit"` // debit or credit
	Currency  string  `json:"currency" binding:"required"`
	Reference string  `json:"reference" binding:"required"`
}
