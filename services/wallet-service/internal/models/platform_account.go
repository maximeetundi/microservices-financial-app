package models

import (
	"time"
)

// PlatformAccount represents a platform-owned fiat account (reserves, fees, operations)
// Multiple accounts can exist per currency - system intelligently selects based on priority and balance
type PlatformAccount struct {
	ID          string    `json:"id"`
	Currency    string    `json:"currency"`     // FCFA, EUR, USD
	AccountType string    `json:"account_type"` // reserve, fees, operations, pending
	Name        string    `json:"name"`         // "Réserve FCFA 1", "Réserve FCFA 2"
	Balance     float64   `json:"balance"`
	MinBalance  float64   `json:"min_balance"` // Minimum balance to maintain
	MaxBalance  float64   `json:"max_balance"` // Maximum balance allowed (0 = unlimited)
	Priority    int       `json:"priority"`    // Higher priority = selected first (1-100)
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlatformCryptoWallet represents a crypto wallet owned by the platform with blockchain address
// Multiple wallets can exist per currency/network - system intelligently selects based on priority and balance
// SECURITY: Private keys are ALWAYS encrypted with AES-256-GCM before storage
type PlatformCryptoWallet struct {
	ID                  string    `json:"id"`
	Currency            string    `json:"currency"`    // BTC, ETH, USDT
	Network             string    `json:"network"`     // ethereum, bsc, tron, bitcoin
	WalletType          string    `json:"wallet_type"` // hot, cold, reserve
	Address             string    `json:"address"`     // Blockchain address (public)
	Label               string    `json:"label"`       // "ETH Hot Wallet 1"
	Balance             float64   `json:"balance"`
	MinBalance          float64   `json:"min_balance"`        // Minimum balance to maintain
	MaxBalance          float64   `json:"max_balance"`        // Maximum balance (0 = unlimited)
	Priority            int       `json:"priority"`           // Higher priority = selected first (1-100)
	EncryptedPrivateKey string    `json:"-"`                  // NEVER exposed in JSON - encrypted with vault
	KeyHash             string    `json:"key_hash,omitempty"` // SHA-256 hash for verification (not the key itself)
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// PlatformTransaction represents a double-entry accounting transaction
type PlatformTransaction struct {
	ID              string    `json:"id"`
	TransactionDate time.Time `json:"transaction_date"`

	// Debit side (source)
	DebitAccountID   string `json:"debit_account_id"`
	DebitAccountType string `json:"debit_account_type"` // user_wallet, platform_fiat, platform_crypto

	// Credit side (destination)
	CreditAccountID   string `json:"credit_account_id"`
	CreditAccountType string `json:"credit_account_type"`

	// Amount
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`

	// Categorization
	OperationType string `json:"operation_type"` // exchange, deposit, withdrawal, fee, admin_credit, admin_debit
	ReferenceType string `json:"reference_type"` // exchange_id, deposit_id, withdrawal_id
	ReferenceID   string `json:"reference_id"`

	// Metadata
	Description      string `json:"description"`
	PerformedBy      string `json:"performed_by"`       // user_id or admin_id
	BlockchainTxHash string `json:"blockchain_tx_hash"` // For crypto transactions

	CreatedAt time.Time `json:"created_at"`
}

// AdminCreditRequest represents a manual credit/debit request from admin
type AdminCreditRequest struct {
	AccountID   string  `json:"account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
	Reference   string  `json:"reference"` // External reference (bank transfer, etc.)
}

// CreatePlatformAccountRequest for creating new platform accounts
type CreatePlatformAccountRequest struct {
	Currency    string  `json:"currency" binding:"required"`
	AccountType string  `json:"account_type" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	MinBalance  float64 `json:"min_balance"`
	MaxBalance  float64 `json:"max_balance"`
	Priority    int     `json:"priority"`
}

// CreatePlatformCryptoWalletRequest for adding crypto wallet addresses
type CreatePlatformCryptoWalletRequest struct {
	Currency   string  `json:"currency" binding:"required"`
	Network    string  `json:"network" binding:"required"`
	WalletType string  `json:"wallet_type" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	Label      string  `json:"label"`
	MinBalance float64 `json:"min_balance"`
	MaxBalance float64 `json:"max_balance"`
	Priority   int     `json:"priority"`
}

// Platform account types constants
const (
	AccountTypeReserve    = "reserve"
	AccountTypeFees       = "fees"
	AccountTypeOperations = "operations"
	AccountTypePending    = "pending"
)

// Crypto wallet types constants
const (
	WalletTypeHot     = "hot"
	WalletTypeCold    = "cold"
	WalletTypeReserve = "reserve"
)

// Crypto networks constants
const (
	NetworkEthereum = "ethereum"
	NetworkBSC      = "bsc"
	NetworkTron     = "tron"
	NetworkBitcoin  = "bitcoin"
	NetworkPolygon  = "polygon"
	NetworkSolana   = "solana"
)

// Operation types constants
const (
	OpTypeExchange      = "exchange"
	OpTypeDeposit       = "deposit"
	OpTypeWithdrawal    = "withdrawal"
	OpTypeFee           = "fee"
	OpTypeAdminCredit   = "admin_credit"
	OpTypeAdminDebit    = "admin_debit"
	OpTypeTransfer      = "transfer"
	OpTypeConsolidation = "consolidation" // Admin moving funds to hot/cold
	OpTypeCryptoSend    = "crypto_send"   // External crypto send via hot wallet
)

// Account type constants for transaction tracking
const (
	AccTypeUserWallet     = "user_wallet"
	AccTypePlatformFiat   = "platform_fiat"
	AccTypePlatformCrypto = "platform_crypto"
	AccTypeExternal       = "external"
)
