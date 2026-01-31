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

// WaveConfig holds Wave API configuration
type WaveConfig struct {
	Name        string // Internal name (e.g., wave, wave_ci)
	APIKey      string
	MerchantID  string
	WebhookKey  string
	BaseURL     string
	Environment string
}

// WaveProvider implements PayoutProvider for Wave
type WaveProvider struct {
	config     WaveConfig
	httpClient *http.Client
}

// NewWaveProvider creates a new Wave provider
func NewWaveProvider(config WaveConfig) *WaveProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.wave.com/v1"
	}
	if config.Name == "" {
		config.Name = "wave"
	}
	return &WaveProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *WaveProvider) GetName() string {
	return p.config.Name
}

func (p *WaveProvider) GetSupportedCountries() []string {
	return []string{
		"SN", // Sénégal
		"CI", // Côte d'Ivoire
		"ML", // Mali
		"BF", // Burkina Faso
		"GM", // Gambia
		"UG", // Uganda
	}
}

func (p *WaveProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return []AvailableMethod{
		{
			Code:           "wave",
			Name:           "Wave",
			Type:           "mobile",
			Countries:      p.GetSupportedCountries(),
			Currencies:     []string{"XOF", "XAF", "UGX"},
			RequiredFields: []string{"phone"},
		},
	}, nil
}

func (p *WaveProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	// Wave is mobile money only, no bank support
	return []Bank{}, nil
}

func (p *WaveProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	switch country {
	case "SN":
		return []MobileOperator{
			{Code: "WAVE", Name: "Wave", Country: "SN", NumberPrefix: []string{"70"}},
		}, nil
	case "CI":
		return []MobileOperator{
			{Code: "WAVE", Name: "Wave", Country: "CI", NumberPrefix: []string{"07"}},
		}, nil
	default:
		return []MobileOperator{
			{Code: "WAVE", Name: "Wave", Country: country, NumberPrefix: []string{}},
		}, nil
	}
}

func (p *WaveProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.RecipientPhone == "" {
		return fmt.Errorf("recipient phone required for Wave")
	}
	return nil
}

func (p *WaveProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	// Wave has very low fees (free for most transfers)
	fee := 0.0
	return &PayoutResponse{
		ProviderName:      "wave",
		ProviderReference: "",
		Amount:            req.Amount,
		Currency:          req.Currency,
		Fee:               fee,
		TotalAmount:       req.Amount + fee,
		Status:            "quote",
	}, nil
}

func (p *WaveProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	payload := map[string]interface{}{
		"recipient_mobile": req.RecipientPhone,
		"amount":           fmt.Sprintf("%.2f", req.Amount),
		"currency":         req.Currency,
		"client_reference": req.ExternalReference,
		"name":             req.RecipientName,
	}

	body, _ := json.Marshal(payload)
	resp, err := p.makeRequest(ctx, "POST", p.config.BaseURL+"/payout", body)
	if err != nil {
		return nil, fmt.Errorf("Wave payout failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Wave API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ID             string `json:"id"`
		Status         string `json:"status"`
		Amount         string `json:"amount"`
		Currency       string `json:"currency"`
		ReceiverMobile string `json:"receiver_mobile"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Wave response: %w", err)
	}

	return &PayoutResponse{
		ProviderName:      "wave",
		ProviderReference: result.ID,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Fee:               0,
		TotalAmount:       req.Amount,
		Status:            result.Status,
	}, nil
}

func (p *WaveProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	resp, err := p.makeRequest(ctx, "GET", fmt.Sprintf("%s/payout/%s", p.config.BaseURL, referenceID), nil)
	if err != nil {
		return nil, fmt.Errorf("Wave status check failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Wave status response: %w", err)
	}

	status := "pending"
	switch result.Status {
	case "succeeded", "complete":
		status = "completed"
	case "failed", "cancelled":
		status = "failed"
	case "pending", "processing":
		status = "pending"
	}

	return &PayoutStatusResponse{
		Status:    status,
		UpdatedAt: time.Now(),
	}, nil
}

func (p *WaveProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("Wave does not support payout cancellation")
}

func (p *WaveProvider) makeRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
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
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	return p.httpClient.Do(req)
}
