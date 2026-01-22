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

func (c *WalletClient) GetWallet(walletID string) (*Wallet, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/%s", c.baseURL, walletID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wallet service returned %d: %s", resp.StatusCode, string(body))
	}
	
	var wallet Wallet
	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		return nil, fmt.Errorf("failed to decode wallet: %w", err)
	}
	return &wallet, nil
}

func (c *WalletClient) GetUserWallets(userID string) ([]Wallet, error) {
	url := fmt.Sprintf("%s/api/v1/wallets?user_id=%s", c.baseURL, userID)
	
	resp, err := c.httpClient.Get(url)
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

func (c *WalletClient) ValidateWalletOwnership(walletID, userID string) (bool, error) {
	wallet, err := c.GetWallet(walletID)
	if err != nil {
		return false, err
	}
	return wallet.UserID == userID, nil
}

func (c *WalletClient) GetWalletBalance(walletID string) (float64, string, error) {
	wallet, err := c.GetWallet(walletID)
	if err != nil {
		return 0, "", err
	}
	return wallet.Balance, wallet.Currency, nil
}
