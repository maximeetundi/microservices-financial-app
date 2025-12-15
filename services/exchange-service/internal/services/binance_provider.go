package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BinanceConfig holds Binance API configuration
type BinanceConfig struct {
	APIKey    string // BINANCE_API_KEY
	APISecret string // BINANCE_API_SECRET
	BaseURL   string // https://api.binance.com (prod) or https://testnet.binance.vision (testnet)
	TestMode  bool
}

// BinanceProvider provides real crypto trading via Binance
type BinanceProvider struct {
	config     BinanceConfig
	httpClient *http.Client
}

// NewBinanceProvider creates a new Binance provider
func NewBinanceProvider(config BinanceConfig) *BinanceProvider {
	if config.BaseURL == "" {
		if config.TestMode {
			config.BaseURL = "https://testnet.binance.vision"
		} else {
			config.BaseURL = "https://api.binance.com"
		}
	}
	
	return &BinanceProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ========================
// MARKET DATA
// ========================

// CryptoPrice represents a cryptocurrency price
type CryptoPrice struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Change24h float64 `json:"change_24h"`
	Volume24h float64 `json:"volume_24h"`
	High24h   float64 `json:"high_24h"`
	Low24h    float64 `json:"low_24h"`
}

// GetPrice gets the current price of a symbol
func (b *BinanceProvider) GetPrice(ctx context.Context, symbol string) (*CryptoPrice, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/24hr?symbol=%s", b.config.BaseURL, symbol)
	
	resp, err := b.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Symbol             string `json:"symbol"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
		LastPrice          string `json:"lastPrice"`
		HighPrice          string `json:"highPrice"`
		LowPrice           string `json:"lowPrice"`
		Volume             string `json:"volume"`
		QuoteVolume        string `json:"quoteVolume"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	price, _ := strconv.ParseFloat(result.LastPrice, 64)
	change, _ := strconv.ParseFloat(result.PriceChangePercent, 64)
	high, _ := strconv.ParseFloat(result.HighPrice, 64)
	low, _ := strconv.ParseFloat(result.LowPrice, 64)
	volume, _ := strconv.ParseFloat(result.QuoteVolume, 64)
	
	return &CryptoPrice{
		Symbol:    symbol,
		Price:     price,
		Change24h: change,
		High24h:   high,
		Low24h:    low,
		Volume24h: volume,
	}, nil
}

// GetAllPrices gets prices for all supported trading pairs
func (b *BinanceProvider) GetAllPrices(ctx context.Context) ([]CryptoPrice, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/price", b.config.BaseURL)
	
	resp, err := b.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var results []struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	
	prices := make([]CryptoPrice, 0)
	// Filter for common pairs
	relevantPairs := map[string]bool{
		"BTCUSDT": true, "ETHUSDT": true, "BNBUSDT": true, "SOLUSDT": true,
		"XRPUSDT": true, "ADAUSDT": true, "DOGEUSDT": true, "DOTUSDT": true,
		"BTCEUR": true, "ETHEUR": true,
	}
	
	for _, r := range results {
		if relevantPairs[r.Symbol] {
			price, _ := strconv.ParseFloat(r.Price, 64)
			prices = append(prices, CryptoPrice{
				Symbol: r.Symbol,
				Price:  price,
			})
		}
	}
	
	return prices, nil
}

// ========================
// TRADING
// ========================

// TradeRequest represents a buy/sell crypto request
type TradeRequest struct {
	UserID     string  `json:"user_id"`
	Symbol     string  `json:"symbol"`     // BTCUSDT, ETHUSDT
	Side       string  `json:"side"`       // BUY or SELL
	Type       string  `json:"type"`       // MARKET or LIMIT
	Quantity   float64 `json:"quantity"`   // Amount of crypto
	QuoteQty   float64 `json:"quote_qty"`  // Amount in quote currency (for MARKET orders)
	Price      float64 `json:"price"`      // For LIMIT orders
	ReferenceID string `json:"reference_id"`
}

// TradeResponse represents the result of a trade
type TradeResponse struct {
	OrderID           string  `json:"order_id"`
	Symbol            string  `json:"symbol"`
	Side              string  `json:"side"`
	Type              string  `json:"type"`
	Status            string  `json:"status"`
	ExecutedQty       float64 `json:"executed_qty"`
	CummulativeQuoteQty float64 `json:"cummulative_quote_qty"`
	Price             float64 `json:"price"`
	Fills             []TradeFill `json:"fills"`
}

// TradeFill represents a single fill in a trade
type TradeFill struct {
	Price   float64 `json:"price"`
	Qty     float64 `json:"qty"`
	Commission float64 `json:"commission"`
	CommissionAsset string `json:"commission_asset"`
}

// ExecuteTrade executes a buy or sell order on Binance
func (b *BinanceProvider) ExecuteTrade(ctx context.Context, req *TradeRequest) (*TradeResponse, error) {
	if b.config.APIKey == "" || b.config.APISecret == "" {
		return nil, fmt.Errorf("Binance API credentials not configured")
	}
	
	params := url.Values{}
	params.Set("symbol", req.Symbol)
	params.Set("side", strings.ToUpper(req.Side))
	params.Set("type", strings.ToUpper(req.Type))
	params.Set("newClientOrderId", req.ReferenceID)
	
	if req.Type == "MARKET" {
		if req.QuoteQty > 0 {
			// Buy with specific quote amount (e.g., buy $100 worth of BTC)
			params.Set("quoteOrderQty", fmt.Sprintf("%.8f", req.QuoteQty))
		} else {
			params.Set("quantity", fmt.Sprintf("%.8f", req.Quantity))
		}
	} else {
		// LIMIT order
		params.Set("quantity", fmt.Sprintf("%.8f", req.Quantity))
		params.Set("price", fmt.Sprintf("%.8f", req.Price))
		params.Set("timeInForce", "GTC") // Good Till Cancelled
	}
	
	resp, err := b.signedRequest(ctx, "POST", "/api/v3/order", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result struct {
		OrderId             int64  `json:"orderId"`
		Symbol              string `json:"symbol"`
		Status              string `json:"status"`
		Side                string `json:"side"`
		Type                string `json:"type"`
		ExecutedQty         string `json:"executedQty"`
		CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
		Price               string `json:"price"`
		Fills               []struct {
			Price           string `json:"price"`
			Qty             string `json:"qty"`
			Commission      string `json:"commission"`
			CommissionAsset string `json:"commissionAsset"`
		} `json:"fills"`
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if result.Code != 0 && result.Code != 200 {
		return nil, fmt.Errorf("Binance error %d: %s", result.Code, result.Msg)
	}
	
	execQty, _ := strconv.ParseFloat(result.ExecutedQty, 64)
	cumQuoteQty, _ := strconv.ParseFloat(result.CummulativeQuoteQty, 64)
	price, _ := strconv.ParseFloat(result.Price, 64)
	
	fills := make([]TradeFill, len(result.Fills))
	for i, f := range result.Fills {
		p, _ := strconv.ParseFloat(f.Price, 64)
		q, _ := strconv.ParseFloat(f.Qty, 64)
		c, _ := strconv.ParseFloat(f.Commission, 64)
		fills[i] = TradeFill{
			Price:           p,
			Qty:             q,
			Commission:      c,
			CommissionAsset: f.CommissionAsset,
		}
	}
	
	return &TradeResponse{
		OrderID:             fmt.Sprintf("%d", result.OrderId),
		Symbol:              result.Symbol,
		Side:                result.Side,
		Type:                result.Type,
		Status:              result.Status,
		ExecutedQty:         execQty,
		CummulativeQuoteQty: cumQuoteQty,
		Price:               price,
		Fills:               fills,
	}, nil
}

// GetOrderStatus gets the status of an order
func (b *BinanceProvider) GetOrderStatus(ctx context.Context, symbol, orderID string) (*TradeResponse, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", orderID)
	
	resp, err := b.signedRequest(ctx, "GET", "/api/v3/order", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		OrderId     int64  `json:"orderId"`
		Symbol      string `json:"symbol"`
		Status      string `json:"status"`
		Side        string `json:"side"`
		Type        string `json:"type"`
		ExecutedQty string `json:"executedQty"`
		CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	execQty, _ := strconv.ParseFloat(result.ExecutedQty, 64)
	cumQuoteQty, _ := strconv.ParseFloat(result.CummulativeQuoteQty, 64)
	
	return &TradeResponse{
		OrderID:             fmt.Sprintf("%d", result.OrderId),
		Symbol:              result.Symbol,
		Side:                result.Side,
		Type:                result.Type,
		Status:              result.Status,
		ExecutedQty:         execQty,
		CummulativeQuoteQty: cumQuoteQty,
	}, nil
}

// CancelOrder cancels an open order
func (b *BinanceProvider) CancelOrder(ctx context.Context, symbol, orderID string) error {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", orderID)
	
	resp, err := b.signedRequest(ctx, "DELETE", "/api/v3/order", params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to cancel order: %s", string(body))
	}
	
	return nil
}

// ========================
// ACCOUNT
// ========================

// AccountBalance represents an asset balance
type AccountBalance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
}

// GetAccountBalances gets all balances from Binance account
func (b *BinanceProvider) GetAccountBalances(ctx context.Context) ([]AccountBalance, error) {
	params := url.Values{}
	
	resp, err := b.signedRequest(ctx, "GET", "/api/v3/account", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Balances []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	balances := make([]AccountBalance, 0)
	for _, b := range result.Balances {
		free, _ := strconv.ParseFloat(b.Free, 64)
		locked, _ := strconv.ParseFloat(b.Locked, 64)
		if free > 0 || locked > 0 {
			balances = append(balances, AccountBalance{
				Asset:  b.Asset,
				Free:   free,
				Locked: locked,
			})
		}
	}
	
	return balances, nil
}

// ========================
// CONVERT (Simple Buy/Sell)
// ========================

// ConvertQuoteRequest for getting a quote for conversion
type ConvertQuoteRequest struct {
	FromAsset string  `json:"from_asset"` // e.g., USDT, EUR
	ToAsset   string  `json:"to_asset"`   // e.g., BTC, ETH
	Amount    float64 `json:"amount"`
}

// ConvertQuoteResponse contains the conversion quote
type ConvertQuoteResponse struct {
	FromAsset     string  `json:"from_asset"`
	ToAsset       string  `json:"to_asset"`
	FromAmount    float64 `json:"from_amount"`
	ToAmount      float64 `json:"to_amount"`
	ExchangeRate  float64 `json:"exchange_rate"`
	Fee           float64 `json:"fee"`
	ValidUntil    time.Time `json:"valid_until"`
}

// GetConvertQuote gets a quote for simple crypto conversion
func (b *BinanceProvider) GetConvertQuote(ctx context.Context, req *ConvertQuoteRequest) (*ConvertQuoteResponse, error) {
	// Determine the trading pair
	symbol := req.FromAsset + req.ToAsset
	reverseSymbol := req.ToAsset + req.FromAsset
	
	// Try to get price for the pair
	price, err := b.GetPrice(ctx, symbol)
	if err != nil {
		// Try reverse pair
		price, err = b.GetPrice(ctx, reverseSymbol)
		if err != nil {
			return nil, fmt.Errorf("trading pair not found: %s or %s", symbol, reverseSymbol)
		}
		// If reverse, invert the calculation
		toAmount := req.Amount * price.Price
		fee := toAmount * 0.001 // 0.1% fee
		
		return &ConvertQuoteResponse{
			FromAsset:    req.FromAsset,
			ToAsset:      req.ToAsset,
			FromAmount:   req.Amount,
			ToAmount:     toAmount - fee,
			ExchangeRate: price.Price,
			Fee:          fee,
			ValidUntil:   time.Now().Add(30 * time.Second),
		}, nil
	}
	
	// Normal direction
	toAmount := req.Amount / price.Price
	fee := toAmount * 0.001 // 0.1% fee
	
	return &ConvertQuoteResponse{
		FromAsset:    req.FromAsset,
		ToAsset:      req.ToAsset,
		FromAmount:   req.Amount,
		ToAmount:     toAmount - fee,
		ExchangeRate: 1 / price.Price,
		Fee:          fee,
		ValidUntil:   time.Now().Add(30 * time.Second),
	}, nil
}

// ========================
// HELPERS
// ========================

func (b *BinanceProvider) signedRequest(ctx context.Context, method, path string, params url.Values) (*http.Response, error) {
	// Add timestamp
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	
	// Sort and create query string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var queryParts []string
	for _, k := range keys {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}
	queryString := strings.Join(queryParts, "&")
	
	// Create signature
	h := hmac.New(sha256.New, []byte(b.config.APISecret))
	h.Write([]byte(queryString))
	signature := hex.EncodeToString(h.Sum(nil))
	
	// Add signature to query
	queryString += "&signature=" + signature
	
	fullURL := fmt.Sprintf("%s%s?%s", b.config.BaseURL, path, queryString)
	
	var reqBody io.Reader
	if method == "POST" {
		reqBody = bytes.NewBufferString(queryString)
		fullURL = b.config.BaseURL + path
	}
	
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("X-MBX-APIKEY", b.config.APIKey)
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	
	return b.httpClient.Do(req)
}
