package services

import (
	"bytes"
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

// DeductPayment deducts amount from user's wallet for ticket purchase
func (c *WalletClient) DeductPayment(userID, walletID, pin string, amount float64, currency, description string) (string, error) {
	payload := map[string]interface{}{
		"amount":      amount,
		"description": description,
		"pin":         pin,
		"currency":    currency, // Send event currency for conversion
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/wallets/%s/withdraw", c.baseURL, walletID), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("wallet service unavailable: %w", err)
	}
	defer resp.Body.Close()

	var walletResp WalletResponse
	json.NewDecoder(resp.Body).Decode(&walletResp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if walletResp.Error != "" {
			return "", fmt.Errorf(walletResp.Error)
		}
		return "", fmt.Errorf("payment failed with status %d", resp.StatusCode)
	}

	return walletResp.TransactionID, nil
}

// CreditPayment credits amount to organizer's wallet (optional, can be added later)
func (c *WalletClient) CreditPayment(userID, walletID string, amount float64, currency, description string) error {
	payload := map[string]interface{}{
		"amount":      amount,
		"method":      "ticket_sale",
		"description": description,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/wallets/%s/deposit", c.baseURL, walletID), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("wallet service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("credit failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetUserWallets fetches all wallets for a user (used for refunds)
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
