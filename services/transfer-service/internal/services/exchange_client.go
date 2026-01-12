package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ExchangeClient struct {
	baseURL    string
	httpClient *http.Client
}

type RateResponse struct {
	Rate float64 `json:"rate"`
}

func NewExchangeClient() *ExchangeClient {
	exchangeURL := os.Getenv("EXCHANGE_SERVICE_URL")
	if exchangeURL == "" {
		exchangeURL = "http://exchange-service:8085" // Default port for exchange service
	}
	return &ExchangeClient{
		baseURL: exchangeURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ExchangeClient) GetRate(from, to string) (float64, error) {
	if from == to {
		return 1.0, nil
	}

	url := fmt.Sprintf("%s/api/v1/rates/%s/%s", c.baseURL, from, to)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get rate, status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Assuming response format of GetSpecificRate: {"rate": 655.95}
	if rate, ok := result["rate"].(float64); ok {
		return rate, nil
	}
	
	return 0, fmt.Errorf("invalid rate response")
}
