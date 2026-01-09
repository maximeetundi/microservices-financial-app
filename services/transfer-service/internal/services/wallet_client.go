package services

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wallet service returned status: %d", resp.StatusCode)
	}

	return nil
}
