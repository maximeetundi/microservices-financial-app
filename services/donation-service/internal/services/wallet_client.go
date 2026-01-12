package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// WalletClient provides HTTP client to communicate with wallet-service
type WalletClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewWalletClient creates a new wallet client
func NewWalletClient() *WalletClient {
	walletURL := os.Getenv("WALLET_SERVICE_URL")
	if walletURL == "" {
		walletURL = "http://wallet-service:8083"
	}
	return &WalletClient{
		baseURL:    walletURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// WalletResponse represents the response from wallet-service
type WalletResponse struct {
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
	Error         string `json:"error"`
}

// CheckBalance checks if user has enough funds (optional pre-check)
func (c *WalletClient) CheckBalance(userID, walletID string) (float64, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/wallets/%s/balance", c.baseURL, walletID), nil)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return 0, "", fmt.Errorf("failed to get balance")
	}
	
	var result struct {
		Balance float64 `json:"balance"`
		Currency string `json:"currency"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Balance, result.Currency, nil
}

// GetUserWallets fetches all wallets for a user
func (c *WalletClient) GetUserWallets(userID string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/internal/wallets?user_id=%s", c.baseURL, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wallet service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch wallets: status %d", resp.StatusCode)
	}

	var result struct {
		Wallets []map[string]interface{} `json:"wallets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Wallets, nil
}
