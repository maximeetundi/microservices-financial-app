package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ExchangeClient provides HTTP client to communicate with exchange-service for currency conversion
type ExchangeClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewExchangeClient creates a new exchange client
func NewExchangeClient() *ExchangeClient {
	exchangeURL := os.Getenv("EXCHANGE_SERVICE_URL")
	if exchangeURL == "" {
		exchangeURL = "http://exchange-service:8086"
	}
	return &ExchangeClient{
		baseURL:    exchangeURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// QuoteResponse represents the response from exchange-service quote endpoint
type QuoteResponse struct {
	Quote struct {
		FromCurrency string  `json:"from_currency"`
		ToCurrency   string  `json:"to_currency"`
		FromAmount   float64 `json:"from_amount"`
		ToAmount     float64 `json:"to_amount"`
		ExchangeRate float64 `json:"exchange_rate"`
		Fee          float64 `json:"fee"`
	} `json:"quote"`
	Error string `json:"error"`
}

// GetConversionQuote gets a quote for converting currency
func (c *ExchangeClient) GetConversionQuote(userID, fromCurrency, toCurrency string, amount float64) (*QuoteResponse, error) {
	payload := map[string]interface{}{
		"from_currency": strings.ToUpper(fromCurrency),
		"to_currency":   strings.ToUpper(toCurrency),
		"amount":        amount,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/fiat/quote", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchange service unavailable: %w", err)
	}
	defer resp.Body.Close()

	var quoteResp QuoteResponse
	json.NewDecoder(resp.Body).Decode(&quoteResp)

	if resp.StatusCode != http.StatusOK {
		if quoteResp.Error != "" {
			return nil, fmt.Errorf(quoteResp.Error)
		}
		return nil, fmt.Errorf("quote failed with status %d", resp.StatusCode)
	}

	return &quoteResp, nil
}

// ConvertAmount converts an amount from one currency to another using fallback rates
// This is used when exchange-service is unavailable
func (c *ExchangeClient) ConvertAmount(amount float64, fromCurrency, toCurrency string) (float64, error) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	// Same currency, no conversion needed
	if fromCurrency == toCurrency {
		return amount, nil
	}

	// Fallback exchange rates (to XOF as base)
	ratesInXOF := map[string]float64{
		"XOF": 1.0,
		"XAF": 1.0,
		"EUR": 655.957,
		"USD": 600.0,
		"GBP": 765.0,
		"CAD": 442.0,
		"NGN": 0.4,
		"GHS": 48.0,
		"KES": 3.92,
		"ZAR": 32.4,
		"MAD": 60.0,
	}

	fromRate, fromOk := ratesInXOF[fromCurrency]
	toRate, toOk := ratesInXOF[toCurrency]

	if !fromOk || !toOk {
		// If currency not in fallback, try using exchange service
		quote, err := c.GetConversionQuote("system", fromCurrency, toCurrency, amount)
		if err != nil {
			return 0, fmt.Errorf("unsupported currency pair: %s to %s", fromCurrency, toCurrency)
		}
		return quote.Quote.ToAmount, nil
	}

	// Convert: from -> XOF -> to
	// amount in 'from' * (fromRate in XOF) / (toRate in XOF) = amount in 'to'
	convertedAmount := amount * fromRate / toRate
	return convertedAmount, nil
}
