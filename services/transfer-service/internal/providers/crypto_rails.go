package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CryptoRailsConfig holds configuration for the crypto rails system
type CryptoRailsConfig struct {
	// Circle USDC API - https://developers.circle.com
	CircleAPIKey string // CIRCLE_API_KEY
	CircleBaseURL string // Default: https://api.circle.com/v1
	
	// Binance for exchange rates and conversions
	BinanceAPIKey    string // BINANCE_API_KEY
	BinanceAPISecret string // BINANCE_API_SECRET
	BinanceBaseURL   string // Default: https://api.binance.com
	
	// Internal pool settings
	UseInternalPoolThreshold float64 // Amount below which we use internal pool (default: 500)
	InternalPoolEnabled      bool
	
	// Stablecoin settings
	PreferredStablecoin string // "USDC" or "USDT"
	
	// Wallet addresses for different chains
	EthereumWallet string
	TronWallet     string
	PolygonWallet  string
}

// CryptoRailsProvider handles invisible crypto conversion for transfers
type CryptoRailsProvider struct {
	config     CryptoRailsConfig
	httpClient *http.Client
	
	// Internal liquidity pools (simulated balances for demo)
	internalPools map[string]float64
}

// ConversionRequest represents a fiat-to-fiat conversion via crypto
type ConversionRequest struct {
	ReferenceID    string  `json:"reference_id"`
	SourceAmount   float64 `json:"source_amount"`
	SourceCurrency string  `json:"source_currency"` // EUR, USD, GBP
	TargetCurrency string  `json:"target_currency"` // XOF, KES, NGN
	
	// Recipient details for final payout
	RecipientCountry string `json:"recipient_country"`
	PayoutMethod     PayoutMethod `json:"payout_method"`
}

// ConversionResponse represents the result of a crypto rails conversion
type ConversionResponse struct {
	ReferenceID      string  `json:"reference_id"`
	SourceAmount     float64 `json:"source_amount"`
	SourceCurrency   string  `json:"source_currency"`
	StablecoinAmount float64 `json:"stablecoin_amount"`
	Stablecoin       string  `json:"stablecoin"` // USDC or USDT
	TargetAmount     float64 `json:"target_amount"`
	TargetCurrency   string  `json:"target_currency"`
	
	// Rates used
	SourceToUSDRate  float64 `json:"source_to_usd_rate"`
	USDToTargetRate  float64 `json:"usd_to_target_rate"`
	
	// Fees
	ConversionFee    float64 `json:"conversion_fee"`
	NetworkFee       float64 `json:"network_fee"`
	TotalFee         float64 `json:"total_fee"`
	
	// Processing info
	UsedInternalPool bool   `json:"used_internal_pool"`
	BlockchainTxHash string `json:"blockchain_tx_hash,omitempty"`
	Status           string `json:"status"`
	EstimatedTime    int    `json:"estimated_time_seconds"`
}

// NewCryptoRailsProvider creates a new crypto rails provider
func NewCryptoRailsProvider(config CryptoRailsConfig) *CryptoRailsProvider {
	if config.CircleBaseURL == "" {
		config.CircleBaseURL = "https://api.circle.com/v1"
	}
	if config.BinanceBaseURL == "" {
		config.BinanceBaseURL = "https://api.binance.com"
	}
	if config.UseInternalPoolThreshold == 0 {
		config.UseInternalPoolThreshold = 2500 // $2500 default
	}
	if config.PreferredStablecoin == "" {
		config.PreferredStablecoin = "USDC"
	}
	
	return &CryptoRailsProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Initialize internal pools with demo balances
		internalPools: map[string]float64{
			"USDC": 100000, // $100k USDC pool
			"USDT": 100000, // $100k USDT pool
		},
	}
}

// GetExchangeRate gets the current exchange rate between currencies
func (c *CryptoRailsProvider) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	// Fiat to USD rates (simplified - in production use forex API)
	fiatToUSD := map[string]float64{
		"USD": 1.0,
		"EUR": 1.08,
		"GBP": 1.27,
		"CHF": 1.12,
		"CAD": 0.74,
		"AUD": 0.66,
	}
	
	// USD to African currencies (CFA zone uses fixed rate to EUR)
	usdToLocal := map[string]float64{
		"XOF": 605.0,  // West African CFA franc
		"XAF": 605.0,  // Central African CFA franc
		"NGN": 1550.0, // Nigerian Naira
		"KES": 153.0,  // Kenyan Shilling
		"GHS": 15.5,   // Ghanaian Cedi
		"ZAR": 18.5,   // South African Rand
		"UGX": 3750.0, // Ugandan Shilling
		"TZS": 2500.0, // Tanzanian Shilling
		"RWF": 1250.0, // Rwandan Franc
	}
	
	// Calculate rate
	fromRate, ok1 := fiatToUSD[from]
	toRate, ok2 := usdToLocal[to]
	
	if ok1 && ok2 {
		return fromRate * toRate, nil
	}
	
	// If target is also fiat
	if ok1 {
		if toUSD, ok := fiatToUSD[to]; ok {
			return fromRate / toUSD, nil
		}
	}
	
	return 0, fmt.Errorf("exchange rate not available for %s to %s", from, to)
}

// ConvertViaStablecoin performs the invisible crypto conversion
func (c *CryptoRailsProvider) ConvertViaStablecoin(ctx context.Context, req *ConversionRequest) (*ConversionResponse, error) {
	// Step 1: Get exchange rates
	sourceToUSDRate, err := c.GetExchangeRate(ctx, req.SourceCurrency, "USD")
	if err != nil {
		// Default rates if not found
		sourceToUSDRate = 1.0
	}
	
	usdToTargetRate, err := c.GetExchangeRate(ctx, "USD", req.TargetCurrency)
	if err != nil {
		return nil, fmt.Errorf("cannot get rate for %s: %w", req.TargetCurrency, err)
	}
	
	// Step 2: Calculate amounts
	usdAmount := req.SourceAmount * sourceToUSDRate
	
	// Step 3: Determine if we use internal pool or blockchain
	useInternalPool := c.config.InternalPoolEnabled && usdAmount <= c.config.UseInternalPoolThreshold
	
	// Step 4: Calculate fees
	var conversionFee, networkFee float64
	if useInternalPool {
		// Lower fees for internal pool (no blockchain)
		conversionFee = usdAmount * 0.005 // 0.5%
		networkFee = 0
	} else {
		// Higher fees for blockchain transfer
		conversionFee = usdAmount * 0.01 // 1%
		networkFee = 2.0 // $2 network fee (Polygon/Tron is cheap)
	}
	
	totalFee := conversionFee + networkFee
	netUSDAmount := usdAmount - totalFee
	
	// Step 5: Calculate target amount
	targetAmount := netUSDAmount * usdToTargetRate
	
	// Step 6: Process the conversion
	var txHash string
	var estimatedTime int
	
	if useInternalPool {
		// Instant internal transfer
		estimatedTime = 5 // 5 seconds
		
		// Check pool balance
		if c.internalPools[c.config.PreferredStablecoin] < netUSDAmount {
			// Fall back to blockchain
			useInternalPool = false
		}
	}
	
	if !useInternalPool {
		// Blockchain transfer required
		estimatedTime = 120 // 2 minutes for Polygon/Tron
		
		// In production, this would call Circle/Fireblocks API
		// For now, generate a mock transaction hash
		txHash = fmt.Sprintf("0x%s%d", req.ReferenceID[:8], time.Now().Unix())
	}
	
	return &ConversionResponse{
		ReferenceID:      req.ReferenceID,
		SourceAmount:     req.SourceAmount,
		SourceCurrency:   req.SourceCurrency,
		StablecoinAmount: netUSDAmount,
		Stablecoin:       c.config.PreferredStablecoin,
		TargetAmount:     targetAmount,
		TargetCurrency:   req.TargetCurrency,
		SourceToUSDRate:  sourceToUSDRate,
		USDToTargetRate:  usdToTargetRate,
		ConversionFee:    conversionFee,
		NetworkFee:       networkFee,
		TotalFee:         totalFee,
		UsedInternalPool: useInternalPool,
		BlockchainTxHash: txHash,
		Status:           "processing",
		EstimatedTime:    estimatedTime,
	}, nil
}

// GetQuote gets a quote for a crypto rails conversion (no execution)
func (c *CryptoRailsProvider) GetQuote(ctx context.Context, req *ConversionRequest) (*ConversionResponse, error) {
	resp, err := c.ConvertViaStablecoin(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Status = "quote"
	return resp, nil
}

// ========================
// Circle USDC Integration
// ========================

// CreateCirclePayment creates a USDC payment via Circle API
func (c *CryptoRailsProvider) CreateCirclePayment(ctx context.Context, amount float64, destinationAddress string, chain string) (string, error) {
	if c.config.CircleAPIKey == "" {
		return "", fmt.Errorf("Circle API key not configured")
	}
	
	payload := map[string]interface{}{
		"idempotencyKey": fmt.Sprintf("%d", time.Now().UnixNano()),
		"amount": map[string]interface{}{
			"amount":   fmt.Sprintf("%.2f", amount),
			"currency": "USD",
		},
		"destination": map[string]interface{}{
			"type":    "blockchain",
			"address": destinationAddress,
			"chain":   chain, // ETH, MATIC, TRX, SOL
		},
	}
	
	body, _ := json.Marshal(payload)
	
	resp, err := c.makeCircleRequest(ctx, "POST", "/payouts", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	return result.Data.ID, nil
}

// GetCirclePaymentStatus checks the status of a Circle payout
func (c *CryptoRailsProvider) GetCirclePaymentStatus(ctx context.Context, paymentID string) (string, error) {
	resp, err := c.makeCircleRequest(ctx, "GET", "/payouts/"+paymentID, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	return result.Data.Status, nil
}

func (c *CryptoRailsProvider) makeCircleRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, c.config.CircleBaseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+c.config.CircleAPIKey)
	req.Header.Set("Content-Type", "application/json")
	
	return c.httpClient.Do(req)
}

// ========================
// Binance Exchange Integration
// ========================

// GetBinancePrice gets real-time price from Binance
func (c *CryptoRailsProvider) GetBinancePrice(ctx context.Context, symbol string) (float64, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", c.config.BinanceBaseURL, symbol)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Price string `json:"price"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	
	var price float64
	fmt.Sscanf(result.Price, "%f", &price)
	
	return price, nil
}

// ========================
// Pool Management
// ========================

// GetPoolBalance returns the current balance of the internal pool
func (c *CryptoRailsProvider) GetPoolBalance(stablecoin string) float64 {
	return c.internalPools[stablecoin]
}

// UpdatePoolBalance updates the pool balance (for deposits/withdrawals)
func (c *CryptoRailsProvider) UpdatePoolBalance(stablecoin string, delta float64) {
	c.internalPools[stablecoin] += delta
}

// NeedsRebalancing checks if the pool needs to be rebalanced
func (c *CryptoRailsProvider) NeedsRebalancing() bool {
	threshold := 10000.0 // $10k minimum
	for _, balance := range c.internalPools {
		if balance < threshold {
			return true
		}
	}
	return false
}
