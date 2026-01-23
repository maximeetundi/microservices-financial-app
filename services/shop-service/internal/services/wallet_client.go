package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type WalletClient struct {
	baseURL    string
	httpClient *http.Client
}

type Wallet struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewWalletClient() *WalletClient {
	baseURL := os.Getenv("WALLET_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://wallet-service:8083"
	}
	return &WalletClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *WalletClient) GetWallet(walletID, token string) (*Wallet, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/%s", c.baseURL, walletID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wallet service returned %d: %s", resp.StatusCode, string(body))
	}
	
	var response struct {
		Wallet *Wallet `json:"wallet"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode wallet: %w", err)
	}
	if response.Wallet == nil {
		return nil, fmt.Errorf("wallet not found in response")
	}
	return response.Wallet, nil
}

func (c *WalletClient) GetUserWallets(userID, token string) ([]Wallet, error) {
	url := fmt.Sprintf("%s/api/v1/wallets?user_id=%s", c.baseURL, userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wallet service returned %d: %s", resp.StatusCode, string(body))
	}
	
	var response struct {
		Wallets []Wallet `json:"wallets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode wallets: %w", err)
	}
	return response.Wallets, nil
}

func (c *WalletClient) ValidateWalletOwnership(walletID, userID, token string) (bool, error) {
	wallet, err := c.GetWallet(walletID, token)
	if err != nil {
		return false, err
	}
	return wallet.UserID == userID, nil
}

func (c *WalletClient) GetWalletBalance(walletID, token string) (float64, string, error) {
	wallet, err := c.GetWallet(walletID, token)
	if err != nil {
		return 0, "", err
	}
	return wallet.Balance, wallet.Currency, nil
}
