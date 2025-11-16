package models

import (
	"time"
)

type ExchangeRate struct {
	ID           string    `json:"id" db:"id"`
	FromCurrency string    `json:"from_currency" db:"from_currency"`
	ToCurrency   string    `json:"to_currency" db:"to_currency"`
	Rate         float64   `json:"rate" db:"rate"`
	BidPrice     float64   `json:"bid_price" db:"bid_price"`     // Prix d'achat
	AskPrice     float64   `json:"ask_price" db:"ask_price"`     // Prix de vente
	Spread       float64   `json:"spread" db:"spread"`           // Différence bid/ask
	Source       string    `json:"source" db:"source"`           // coinbase, binance, etc.
	Volume24h    float64   `json:"volume_24h" db:"volume_24h"`   // Volume 24h
	Change24h    float64   `json:"change_24h" db:"change_24h"`   // Changement % 24h
	LastUpdated  time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Exchange struct {
	ID                   string     `json:"id" db:"id"`
	UserID               string     `json:"user_id" db:"user_id"`
	FromWalletID         string     `json:"from_wallet_id" db:"from_wallet_id"`
	ToWalletID           string     `json:"to_wallet_id" db:"to_wallet_id"`
	FromCurrency         string     `json:"from_currency" db:"from_currency"`
	ToCurrency           string     `json:"to_currency" db:"to_currency"`
	FromAmount           float64    `json:"from_amount" db:"from_amount"`
	ToAmount             float64    `json:"to_amount" db:"to_amount"`
	ExchangeRate         float64    `json:"exchange_rate" db:"exchange_rate"`
	Fee                  float64    `json:"fee" db:"fee"`
	FeePercentage        float64    `json:"fee_percentage" db:"fee_percentage"`
	Status               string     `json:"status" db:"status"` // pending, completed, failed, cancelled
	QuoteID              string     `json:"quote_id" db:"quote_id"`
	ExternalTransactionID *string   `json:"external_transaction_id,omitempty" db:"external_transaction_id"`
	CompletedAt          *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

type Quote struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	FromCurrency      string    `json:"from_currency" db:"from_currency"`
	ToCurrency        string    `json:"to_currency" db:"to_currency"`
	FromAmount        float64   `json:"from_amount" db:"from_amount"`
	ToAmount          float64   `json:"to_amount" db:"to_amount"`
	ExchangeRate      float64   `json:"exchange_rate" db:"exchange_rate"`
	Fee               float64   `json:"fee" db:"fee"`
	FeePercentage     float64   `json:"fee_percentage" db:"fee_percentage"`
	ValidUntil        time.Time `json:"valid_until" db:"valid_until"`
	EstimatedDelivery string    `json:"estimated_delivery" db:"estimated_delivery"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type TradingOrder struct {
	ID            string     `json:"id" db:"id"`
	UserID        string     `json:"user_id" db:"user_id"`
	WalletID      string     `json:"wallet_id" db:"wallet_id"`
	OrderType     string     `json:"order_type" db:"order_type"`     // market, limit, stop_loss, take_profit
	Side          string     `json:"side" db:"side"`                 // buy, sell
	FromCurrency  string     `json:"from_currency" db:"from_currency"`
	ToCurrency    string     `json:"to_currency" db:"to_currency"`
	Amount        float64    `json:"amount" db:"amount"`
	Price         *float64   `json:"price,omitempty" db:"price"`     // Pour les ordres limit
	StopPrice     *float64   `json:"stop_price,omitempty" db:"stop_price"` // Pour stop loss
	FilledAmount  float64    `json:"filled_amount" db:"filled_amount"`
	RemainingAmount float64  `json:"remaining_amount" db:"remaining_amount"`
	Status        string     `json:"status" db:"status"` // open, filled, cancelled, partial
	Fee           float64    `json:"fee" db:"fee"`
	ExecutedAt    *time.Time `json:"executed_at,omitempty" db:"executed_at"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type Trade struct {
	ID           string    `json:"id" db:"id"`
	BuyOrderID   string    `json:"buy_order_id" db:"buy_order_id"`
	SellOrderID  string    `json:"sell_order_id" db:"sell_order_id"`
	FromCurrency string    `json:"from_currency" db:"from_currency"`
	ToCurrency   string    `json:"to_currency" db:"to_currency"`
	Amount       float64   `json:"amount" db:"amount"`
	Price        float64   `json:"price" db:"price"`
	Fee          float64   `json:"fee" db:"fee"`
	ExecutedAt   time.Time `json:"executed_at" db:"executed_at"`
}

type P2POffer struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	Type            string     `json:"type" db:"type"`            // buy, sell
	FromCurrency    string     `json:"from_currency" db:"from_currency"`
	ToCurrency      string     `json:"to_currency" db:"to_currency"`
	Amount          float64    `json:"amount" db:"amount"`
	Price           float64    `json:"price" db:"price"`
	MinAmount       float64    `json:"min_amount" db:"min_amount"`
	MaxAmount       float64    `json:"max_amount" db:"max_amount"`
	PaymentMethods  string     `json:"payment_methods" db:"payment_methods"` // JSON array
	Terms           string     `json:"terms" db:"terms"`
	AutoReply       bool       `json:"auto_reply" db:"auto_reply"`
	Status          string     `json:"status" db:"status"` // active, paused, completed
	CompletedTrades int        `json:"completed_trades" db:"completed_trades"`
	SuccessRate     float64    `json:"success_rate" db:"success_rate"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type P2PTrade struct {
	ID          string     `json:"id" db:"id"`
	OfferID     string     `json:"offer_id" db:"offer_id"`
	BuyerID     string     `json:"buyer_id" db:"buyer_id"`
	SellerID    string     `json:"seller_id" db:"seller_id"`
	Amount      float64    `json:"amount" db:"amount"`
	Price       float64    `json:"price" db:"price"`
	TotalValue  float64    `json:"total_value" db:"total_value"`
	Status      string     `json:"status" db:"status"` // pending, paid, disputed, completed, cancelled
	ChatID      string     `json:"chat_id" db:"chat_id"`
	EscrowID    *string    `json:"escrow_id,omitempty" db:"escrow_id"`
	PaidAt      *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Request/Response models
type QuoteRequest struct {
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	FromAmount   *float64 `json:"from_amount,omitempty"`
	ToAmount     *float64 `json:"to_amount,omitempty"`
}

type ExchangeRequest struct {
	FromWalletID string `json:"from_wallet_id" binding:"required"`
	ToWalletID   string `json:"to_wallet_id" binding:"required"`
	QuoteID      string `json:"quote_id" binding:"required"`
}

type BuyOrderRequest struct {
	Currency     string  `json:"currency" binding:"required"`     // Crypto à acheter
	PayCurrency  string  `json:"pay_currency" binding:"required"` // Devise de paiement
	Amount       float64 `json:"amount" binding:"required,gt=0"`  // Montant à acheter
	OrderType    string  `json:"order_type" binding:"required,oneof=market limit"`
	LimitPrice   *float64 `json:"limit_price,omitempty"`
}

type SellOrderRequest struct {
	Currency     string  `json:"currency" binding:"required"`     // Crypto à vendre
	ReceiveCurrency string `json:"receive_currency" binding:"required"` // Devise à recevoir
	Amount       float64 `json:"amount" binding:"required,gt=0"`  // Montant à vendre
	OrderType    string  `json:"order_type" binding:"required,oneof=market limit"`
	LimitPrice   *float64 `json:"limit_price,omitempty"`
}

type LimitOrderRequest struct {
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	LimitPrice   float64 `json:"limit_price" binding:"required,gt=0"`
	OrderType    string  `json:"order_type" binding:"required,oneof=buy sell"`
}

type StopLossRequest struct {
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	StopPrice    float64 `json:"stop_price" binding:"required,gt=0"`
	LimitPrice   *float64 `json:"limit_price,omitempty"`
}

type P2POfferRequest struct {
	Type           string   `json:"type" binding:"required,oneof=buy sell"`
	FromCurrency   string   `json:"from_currency" binding:"required"`
	ToCurrency     string   `json:"to_currency" binding:"required"`
	Amount         float64  `json:"amount" binding:"required,gt=0"`
	Price          float64  `json:"price" binding:"required,gt=0"`
	MinAmount      float64  `json:"min_amount" binding:"required,gt=0"`
	MaxAmount      float64  `json:"max_amount" binding:"required,gt=0"`
	PaymentMethods []string `json:"payment_methods" binding:"required"`
	Terms          string   `json:"terms" binding:"required"`
	AutoReply      bool     `json:"auto_reply"`
}

type Market struct {
	Symbol       string  `json:"symbol"`        // BTC/USD
	BaseAsset    string  `json:"base_asset"`    // BTC
	QuoteAsset   string  `json:"quote_asset"`   // USD
	Price        float64 `json:"price"`
	Change24h    float64 `json:"change_24h"`
	Volume24h    float64 `json:"volume_24h"`
	High24h      float64 `json:"high_24h"`
	Low24h       float64 `json:"low_24h"`
	BidPrice     float64 `json:"bid_price"`
	AskPrice     float64 `json:"ask_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

type OrderBook struct {
	Symbol    string      `json:"symbol"`
	Bids      []OrderLevel `json:"bids"` // Prix d'achat
	Asks      []OrderLevel `json:"asks"` // Prix de vente
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderLevel struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"` // Nombre d'ordres à ce niveau
}

type Portfolio struct {
	TotalValue    float64            `json:"total_value"`
	TotalPnL      float64            `json:"total_pnl"`
	TotalPnLPerc  float64            `json:"total_pnl_percentage"`
	Holdings      []Holding          `json:"holdings"`
	Performance   PerformanceMetrics `json:"performance"`
}

type Holding struct {
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	Value         float64 `json:"value"`
	AvgBuyPrice   float64 `json:"avg_buy_price"`
	CurrentPrice  float64 `json:"current_price"`
	PnL           float64 `json:"pnl"`
	PnLPercentage float64 `json:"pnl_percentage"`
}

type PerformanceMetrics struct {
	TotalTrades      int     `json:"total_trades"`
	WinningTrades    int     `json:"winning_trades"`
	LosingTrades     int     `json:"losing_trades"`
	WinRate          float64 `json:"win_rate"`
	TotalVolume      float64 `json:"total_volume"`
	TotalFees        float64 `json:"total_fees"`
	BestTrade        float64 `json:"best_trade"`
	WorstTrade       float64 `json:"worst_trade"`
	ProfitFactor     float64 `json:"profit_factor"`
}