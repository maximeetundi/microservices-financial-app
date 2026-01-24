package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WalletClient struct {
	baseURL    string
	httpClient *http.Client
}

type WalletBalance struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type TransactionRequest struct {
	UserID   string  `json:"user_id"`
	WalletID string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"` // debit, credit
	Currency string  `json:"currency"`
	Reference string `json:"reference"`
}

func NewWalletClient(baseURL string) *WalletClient {
	return &WalletClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (w *WalletClient) GetWalletBalance(userID, currency string) (*WalletBalance, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/balance?user_id=%s&currency=%s", w.baseURL, userID, currency)
	
	resp, err := w.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet balance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}

	var balance WalletBalance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, fmt.Errorf("failed to decode wallet balance response: %w", err)
	}

	return &balance, nil
}

func (w *WalletClient) GetWalletBalanceByID(walletID, token string) (float64, string, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/%s/balance", w.baseURL, walletID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}
	
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get wallet balance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}

	var response struct {
		Balance float64 `json:"balance"`
		// API might not return currency in simple balance endpoint, implied from context or we fetch wallet details
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, "", fmt.Errorf("failed to decode wallet balance: %w", err)
	}

	return response.Balance, "", nil 
}

func (w *WalletClient) ProcessTransaction(req *TransactionRequest) error {
	url := fmt.Sprintf("%s/api/v1/wallets/transaction", w.baseURL)
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction request: %w", err)
	}

	resp, err := w.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to process wallet transaction: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}

	return nil
}

func (w *WalletClient) CheckWalletExists(userID, currency string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/exists?user_id=%s&currency=%s", w.baseURL, userID, currency)
	
	resp, err := w.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to check wallet existence: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}