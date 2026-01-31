package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LygosConfig holds Lygos API configuration
type LygosConfig struct {
	APIKey     string
	ShopName   string
	BaseURL    string
	WebhookURL string
}

// LygosProvider implements PayoutProvider for Lygos
// Note: Lygos is primarily for payment collection (deposits), not payouts
type LygosProvider struct {
	config     LygosConfig
	httpClient *http.Client
}

// NewLygosProvider creates a new Lygos provider
func NewLygosProvider(config LygosConfig) *LygosProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.lygosapp.com/v1"
	}
	return &LygosProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *LygosProvider) GetName() string {
	return "lygos"
}

func (p *LygosProvider) GetSupportedCountries() []string {
	return []string{
		"CI", // Côte d'Ivoire
		"SN", // Sénégal
	}
}

func (p *LygosProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return []AvailableMethod{
		{
			Code:           "mobile_money",
			Name:           "Mobile Money",
			Type:           "mobile",
			Countries:      p.GetSupportedCountries(),
			Currencies:     []string{"XOF"},
			RequiredFields: []string{"phone"},
		},
	}, nil
}

func (p *LygosProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	return []Bank{}, nil
}

func (p *LygosProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	switch country {
	case "CI":
		return []MobileOperator{
			{Code: "ORANGE_CI", Name: "Orange Money", Country: "CI", NumberPrefix: []string{"07"}},
			{Code: "MTN_CI", Name: "MTN Mobile Money", Country: "CI", NumberPrefix: []string{"05"}},
			{Code: "MOOV_CI", Name: "Moov Money", Country: "CI", NumberPrefix: []string{"01"}},
		}, nil
	case "SN":
		return []MobileOperator{
			{Code: "ORANGE_SN", Name: "Orange Money", Country: "SN", NumberPrefix: []string{"77", "78"}},
			{Code: "WAVE_SN", Name: "Wave", Country: "SN", NumberPrefix: []string{"70"}},
		}, nil
	default:
		return []MobileOperator{}, nil
	}
}

func (p *LygosProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.Amount < 100 {
		return fmt.Errorf("minimum amount is 100 XOF")
	}
	return nil
}

func (p *LygosProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	fee := req.Amount * 0.02 // 2% fee estimate
	return &PayoutResponse{
		ProviderName:      "lygos",
		ProviderReference: "",
		Amount:            req.Amount,
		Currency:          req.Currency,
		Fee:               fee,
		TotalAmount:       req.Amount + fee,
		Status:            "quote",
	}, nil
}

// CreatePaymentGateway creates a payment collection link (deposit)
func (p *LygosProvider) CreatePaymentGateway(ctx context.Context, amount float64, orderID, message, successURL, failureURL string) (string, string, error) {
	payload := map[string]interface{}{
		"amount":      int(amount),
		"shop_name":   p.config.ShopName,
		"message":     message,
		"order_id":    orderID,
		"success_url": successURL,
		"failure_url": failureURL,
	}

	body, _ := json.Marshal(payload)
	resp, err := p.makeRequest(ctx, "POST", p.config.BaseURL+"/gateway", body)
	if err != nil {
		return "", "", fmt.Errorf("Lygos gateway creation failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("Lygos API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ID   string `json:"id"`
		Link string `json:"link"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse Lygos response: %w", err)
	}

	return result.ID, result.Link, nil
}

func (p *LygosProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	// Lygos is primarily for collection, payouts may not be fully supported
	return nil, fmt.Errorf("Lygos does not support direct payouts, use CreatePaymentGateway for deposits")
}

func (p *LygosProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	resp, err := p.makeRequest(ctx, "GET", fmt.Sprintf("%s/gateway/%s", p.config.BaseURL, referenceID), nil)
	if err != nil {
		return nil, fmt.Errorf("Lygos status check failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Lygos status response: %w", err)
	}

	status := "pending"
	switch result.Status {
	case "paid", "completed":
		status = "completed"
	case "failed", "expired":
		status = "failed"
	case "pending", "created":
		status = "pending"
	}

	return &PayoutStatusResponse{
		Status:    status,
		UpdatedAt: time.Now(),
	}, nil
}

func (p *LygosProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("Lygos does not support cancellation")
}

func (p *LygosProvider) makeRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", p.config.APIKey)

	return p.httpClient.Do(req)
}
