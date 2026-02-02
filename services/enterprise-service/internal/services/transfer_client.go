package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// TransferClient provides HTTP client to communicate with transfer-service
type TransferClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTransferClient creates a new transfer client
func NewTransferClient() *TransferClient {
	transferURL := os.Getenv("TRANSFER_SERVICE_URL")
	if transferURL == "" {
		transferURL = "http://transfer-service:8084"
	}
	return &TransferClient{
		baseURL:    transferURL,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// BulkTransferRequest represents a request to execute multiple transfers
type BulkTransferRequest struct {
	Transfers []TransferItem `json:"transfers"`
	BatchID   string         `json:"batch_id"`
	Source    string         `json:"source"` // "PAYROLL", "BILLING", etc.
}

// TransferItem represents a single transfer in a bulk request
type TransferItem struct {
	RecipientUserID string  `json:"recipient_user_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Description     string  `json:"description"`
	Reference       string  `json:"reference"`
}

// BulkTransferResponse represents the response from bulk transfer
type BulkTransferResponse struct {
	Success       bool                 `json:"success"`
	BatchID       string               `json:"batch_id"`
	TransactionID string               `json:"transaction_id"`
	TotalAmount   float64              `json:"total_amount"`
	TotalCount    int                  `json:"total_count"`
	SuccessCount  int                  `json:"success_count"`
	FailedCount   int                  `json:"failed_count"`
	Results       []TransferResultItem `json:"results"`
}

// TransferResultItem represents result for a single transfer
type TransferResultItem struct {
	Reference     string `json:"reference"`
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id,omitempty"`
	Error         string `json:"error,omitempty"`
}

// ExecuteBulkTransfer executes multiple transfers (for payroll/billing)
func (c *TransferClient) ExecuteBulkTransfer(sourceWalletID string, token string, req BulkTransferRequest) (*BulkTransferResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/internal/bulk-transfer", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Source-Wallet-ID", sourceWalletID)
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("transfer service unavailable: %w", err)
	}
	defer resp.Body.Close()

	var transferResp BulkTransferResponse
	if err := json.NewDecoder(resp.Body).Decode(&transferResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return &transferResp, fmt.Errorf("transfer failed with status %d", resp.StatusCode)
	}

	return &transferResp, nil
}

// SingleTransferRequest for simple P2P transfers
type SingleTransferRequest struct {
	RecipientUserID string  `json:"recipient_user_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Description     string  `json:"description"`
}

// SingleTransferResponse represents simple transfer response
type SingleTransferResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message,omitempty"`
	Error         string `json:"error,omitempty"`
}

// ExecuteSingleTransfer executes a single P2P transfer
func (c *TransferClient) ExecuteSingleTransfer(sourceWalletID string, token string, req SingleTransferRequest) (*SingleTransferResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/internal/transfer", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Source-Wallet-ID", sourceWalletID)
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("transfer service unavailable: %w", err)
	}
	defer resp.Body.Close()

	var transferResp SingleTransferResponse
	if err := json.NewDecoder(resp.Body).Decode(&transferResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return &transferResp, fmt.Errorf("transfer failed: %s", transferResp.Error)
	}

	return &transferResp, nil
}
