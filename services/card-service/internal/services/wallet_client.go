package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WalletClient provides HTTP client to communicate with wallet-service
type WalletClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewWalletClient creates a new wallet service client
func NewWalletClient(baseURL string) *WalletClient {
	return &WalletClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WalletBalance represents wallet balance response
type WalletBalance struct {
	WalletID  string  `json:"wallet_id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Available float64 `json:"available"`
}

// DebitRequest represents a debit request to wallet
type DebitRequest struct {
	WalletID    string  `json:"wallet_id"`
	Amount      float64 `json:"amount"`
	Fee         float64 `json:"fee"`
	Description string  `json:"description"`
	Reference   string  `json:"reference"`
}

// DebitResponse represents the response from a debit operation
type DebitResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id"`
	NewBalance    float64 `json:"new_balance"`
}

// GetBalance gets the balance of a wallet
func (c *WalletClient) GetBalance(walletID, userToken string) (*WalletBalance, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/%s/balance", c.baseURL, walletID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet balance: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var balance WalletBalance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &balance, nil
}

// DebitWallet debits an amount from a wallet for card loading
func (c *WalletClient) DebitWallet(walletID, userToken string, req *DebitRequest) (*DebitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/%s/debit", c.baseURL, walletID)
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(respBody))
	}
	
	var debitResp DebitResponse
	if err := json.NewDecoder(resp.Body).Decode(&debitResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &debitResp, nil
}

// CreditWallet credits an amount to a wallet (for refunds, gift card redemption)
func (c *WalletClient) CreditWallet(walletID, userToken string, amount float64, description, reference string) (*DebitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/%s/credit", c.baseURL, walletID)
	
	req := map[string]interface{}{
		"amount":      amount,
		"description": description,
		"reference":   reference,
	}
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to credit wallet: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(respBody))
	}
	
	var creditResp DebitResponse
	if err := json.NewDecoder(resp.Body).Decode(&creditResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &creditResp, nil
}

// CheckSufficientBalance checks if a wallet has sufficient balance
func (c *WalletClient) CheckSufficientBalance(walletID, userToken string, amount float64) (bool, error) {
	balance, err := c.GetBalance(walletID, userToken)
	if err != nil {
		return false, err
	}
	
	return balance.Available >= amount, nil
}
