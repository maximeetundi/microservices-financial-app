package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ExchangeClient struct {
	baseURL    string
	httpClient *http.Client
}

type ExchangeRate struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Rate         float64 `json:"rate"`
	Timestamp    string  `json:"timestamp"`
}

func NewExchangeClient() *ExchangeClient {
	baseURL := os.Getenv("EXCHANGE_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://exchange-service:8085"
	}
	return &ExchangeClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ExchangeClient) GetExchangeRate(fromCurrency, toCurrency string) (*ExchangeRate, error) {
	if fromCurrency == toCurrency {
		return &ExchangeRate{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         1.0,
			Timestamp:    time.Now().Format(time.RFC3339),
		}, nil
	}
	
	url := fmt.Sprintf("%s/api/v1/rates?from=%s&to=%s", c.baseURL, fromCurrency, toCurrency)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("exchange service returned %d: %s", resp.StatusCode, string(body))
	}
	
	var rate ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&rate); err != nil {
		return nil, fmt.Errorf("failed to decode exchange rate: %w", err)
	}
	return &rate, nil
}

func (c *ExchangeClient) ConvertAmount(amount float64, fromCurrency, toCurrency string) (float64, float64, error) {
	rate, err := c.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return 0, 0, err
	}
	convertedAmount := amount * rate.Rate
	return convertedAmount, rate.Rate, nil
}
