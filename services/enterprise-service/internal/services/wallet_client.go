package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WalletClient handles communication with wallet-service
type WalletClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewWalletClient(baseURL string) *WalletClient {
	return &WalletClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WalletInfo represents wallet data from wallet-service
type WalletInfo struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Name     string  `json:"name"`
	Type     string  `json:"wallet_type"`
}

// CreateWalletRequest for internal wallet creation
type CreateWalletRequest struct {
	UserID   string `json:"user_id"`
	Currency string `json:"currency"`
}

// CreateWallet creates a new wallet via wallet-service internal API
func (c *WalletClient) CreateWallet(ctx context.Context, ownerUserID, currency string) (*WalletInfo, error) {
	reqBody := CreateWalletRequest{
		UserID:   ownerUserID,
		Currency: currency,
	}
	
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	url := c.baseURL + "/api/v1/internal/wallets"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call wallet service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Wallet   WalletInfo `json:"wallet"`
		Existing bool       `json:"existing"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &result.Wallet, nil
}

// GetWallets retrieves wallets for a user from wallet-service
func (c *WalletClient) GetWallets(ctx context.Context, userID string) ([]WalletInfo, error) {
	url := fmt.Sprintf("%s/api/v1/internal/wallets?user_id=%s", c.baseURL, userID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call wallet service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Wallets []WalletInfo `json:"wallets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Wallets, nil
}

// GetWallet retrieves a specific wallet by ID
func (c *WalletClient) GetWallet(ctx context.Context, walletID string) (*WalletInfo, error) {
	url := fmt.Sprintf("%s/api/v1/internal/wallets/%s", c.baseURL, walletID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call wallet service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}
	
	var wallet WalletInfo
	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &wallet, nil
}
