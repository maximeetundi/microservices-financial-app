package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
)

type WalletClient struct {
	baseURL    string
	httpClient *http.Client
}

type WalletTransactionRequest struct {
	UserID    string  `json:"user_id"`
	WalletID  string  `json:"wallet_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"` // debit, credit
	Currency  string  `json:"currency"`
	Reference string  `json:"reference"`
}

func NewWalletClient(cfg *config.Config) *WalletClient {
	return &WalletClient{
		baseURL: cfg.WalletServiceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *WalletClient) ProcessTransaction(req *WalletTransactionRequest) error {
	url := fmt.Sprintf("%s/api/v1/wallets/transaction", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	log.Printf("[WALLET-CLIENT DEBUG] Sending %s request to %s: UserID=%s, WalletID=%s, Amount=%f, Currency=%s",
		req.Type, url, req.UserID, req.WalletID, req.Amount, req.Currency)

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[WALLET-CLIENT ERROR] HTTP request failed: %v", err)
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read response body for error details
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[WALLET-CLIENT ERROR] Wallet service returned status %d: %s", resp.StatusCode, string(body))
		return fmt.Errorf("wallet service returned status: %d - %s", resp.StatusCode, string(body))
	}

	log.Printf("[WALLET-CLIENT DEBUG] Transaction successful: %s %s", req.Type, req.WalletID)
	return nil
}
// GetUserWallets fetches wallets for a specific user via internal API
func (c *WalletClient) GetUserWallets(userID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/internal/wallets?user_id=%s", c.baseURL, userID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet service returned status: %d", resp.StatusCode)
	}

	var result struct {
		Wallets []map[string]interface{} `json:"wallets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Wallets, nil
}

// CreateUserWallet creates a wallet for a user via internal API
func (c *WalletClient) CreateUserWallet(userID, currency string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/internal/wallets", c.baseURL)
	
	reqBody := map[string]string{
		"user_id":  userID,
		"currency": currency,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Wallet map[string]interface{} `json:"wallet"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Safely get ID
	if id, ok := result.Wallet["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("wallet id not found in response")
}
