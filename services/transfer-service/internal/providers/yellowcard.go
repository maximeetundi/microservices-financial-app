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

// YellowCardConfig holds Yellow Card API configuration
type YellowCardConfig struct {
	APIKey      string
	SecretKey   string
	BusinessID  string
	BaseURL     string
	WebhookURL  string
	Environment string
}

// YellowCardProvider implements PayoutProvider for Yellow Card
type YellowCardProvider struct {
	config     YellowCardConfig
	httpClient *http.Client
}

// NewYellowCardProvider creates a new Yellow Card provider
func NewYellowCardProvider(config YellowCardConfig) *YellowCardProvider {
	if config.BaseURL == "" {
		if config.Environment == "live" || config.Environment == "production" {
			config.BaseURL = "https://api.yellowcard.io/v1"
		} else {
			config.BaseURL = "https://sandbox.api.yellowcard.io/v1"
		}
	}
	return &YellowCardProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *YellowCardProvider) GetName() string {
	return "yellowcard"
}

func (p *YellowCardProvider) GetSupportedCountries() []string {
	return []string{
		"NG", // Nigeria
		"GH", // Ghana
		"KE", // Kenya
		"ZA", // South Africa
		"UG", // Uganda
		"TZ", // Tanzania
		"RW", // Rwanda
		"ZM", // Zambia
		"BW", // Botswana
	}
}

func (p *YellowCardProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	methods := []AvailableMethod{
		{
			Code:           "mobile_money",
			Name:           "Mobile Money",
			Type:           "mobile",
			Countries:      []string{"GH", "KE", "UG", "TZ", "RW", "ZM"},
			Currencies:     []string{"GHS", "KES", "UGX", "TZS", "RWF", "ZMK"},
			RequiredFields: []string{"phone"},
		},
		{
			Code:           "bank_transfer",
			Name:           "Bank Transfer",
			Type:           "bank",
			Countries:      []string{"NG", "GH", "KE", "ZA"},
			Currencies:     []string{"NGN", "GHS", "KES", "ZAR"},
			RequiredFields: []string{"bank_code", "account_number"},
		},
	}
	return methods, nil
}

func (p *YellowCardProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	// Fetch banks from YellowCard API
	resp, err := p.makeRequest(ctx, "GET", fmt.Sprintf("%s/banks?country=%s", p.config.BaseURL, country), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch banks: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Banks []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse banks response: %w", err)
	}

	var banks []Bank
	for _, b := range result.Banks {
		banks = append(banks, Bank{
			Code:    b.Code,
			Name:    b.Name,
			Country: country,
		})
	}

	return banks, nil
}

func (p *YellowCardProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	operators := map[string][]MobileOperator{
		"GH": {
			{Code: "MTN", Name: "MTN Mobile Money", Countries: []string{"GH"}, NumberPrefix: []string{"23", "24", "54", "55"}},
			{Code: "VODAFONE", Name: "Vodafone Cash", Countries: []string{"GH"}, NumberPrefix: []string{"20", "50"}},
			{Code: "AIRTELTIGO", Name: "AirtelTigo Money", Countries: []string{"GH"}, NumberPrefix: []string{"26", "27", "56", "57"}},
		},
		"KE": {
			{Code: "MPESA", Name: "M-Pesa", Countries: []string{"KE"}, NumberPrefix: []string{"07", "01"}},
			{Code: "AIRTEL", Name: "Airtel Money", Countries: []string{"KE"}, NumberPrefix: []string{"073", "078"}},
		},
		"UG": {
			{Code: "MTN", Name: "MTN Mobile Money", Countries: []string{"UG"}, NumberPrefix: []string{"77", "78"}},
			{Code: "AIRTEL", Name: "Airtel Money", Countries: []string{"UG"}, NumberPrefix: []string{"70", "75"}},
		},
	}

	if ops, ok := operators[country]; ok {
		return ops, nil
	}
	return []MobileOperator{}, nil
}

func (p *YellowCardProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.RecipientPhone == "" && req.AccountNumber == "" {
		return fmt.Errorf("recipient phone or bank account number required")
	}
	return nil
}

func (p *YellowCardProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	payload := map[string]interface{}{
		"sourceCurrency":      req.Currency,
		"destinationCurrency": req.Currency,
		"amount":              req.Amount,
	}

	body, _ := json.Marshal(payload)
	resp, err := p.makeRequest(ctx, "POST", p.config.BaseURL+"/rates", body)
	if err != nil {
		return nil, fmt.Errorf("YellowCard quote failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Rate  float64 `json:"rate"`
		Fee   float64 `json:"fee"`
		Total float64 `json:"total"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		// Return estimated fee if API fails
		fee := req.Amount * 0.015
		return &PayoutResponse{
			ProviderName:      "yellowcard",
			ProviderReference: "",
			AmountReceived:    req.Amount,
			ReceivedCurrency:  req.Currency,
			Fee:               fee,
			Status:            PayoutStatusPending,
		}, nil
	}

	return &PayoutResponse{
		ProviderName:      "yellowcard",
		ProviderReference: "",
		Amount:            req.Amount,
		Currency:          req.Currency,
		Fee:               result.Fee,
		TotalAmount:       result.Total,
		Status:            "quote",
	}, nil
}

func (p *YellowCardProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	payload := map[string]interface{}{
		"destination": map[string]interface{}{
			"currency":      req.Currency,
			"country":       req.Country,
			"accountNumber": req.BankAccountNumber,
			"bankCode":      req.BankCode,
			"phone":         req.RecipientPhone,
			"name":          req.RecipientName,
		},
		"amount":      req.Amount,
		"reference":   req.ExternalReference,
		"callbackUrl": req.CallbackURL,
	}

	body, _ := json.Marshal(payload)
	resp, err := p.makeRequest(ctx, "POST", p.config.BaseURL+"/payments", body)
	if err != nil {
		return nil, fmt.Errorf("YellowCard payout failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("YellowCard API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ID     string  `json:"id"`
		Status string  `json:"status"`
		Fee    float64 `json:"fee"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YellowCard response: %w", err)
	}

	return &PayoutResponse{
		ProviderName:      "yellowcard",
		ProviderReference: result.ID,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Fee:               result.Fee,
		TotalAmount:       req.Amount + result.Fee,
		Status:            result.Status,
	}, nil
}

func (p *YellowCardProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	resp, err := p.makeRequest(ctx, "GET", fmt.Sprintf("%s/payments/%s", p.config.BaseURL, referenceID), nil)
	if err != nil {
		return nil, fmt.Errorf("YellowCard status check failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YellowCard status response: %w", err)
	}

	status := "pending"
	switch result.Status {
	case "completed", "success":
		status = "completed"
	case "failed", "cancelled", "rejected":
		status = "failed"
	case "pending", "processing":
		status = "pending"
	}

	return &PayoutStatusResponse{
		Status:    status,
		UpdatedAt: time.Now(),
	}, nil
}

func (p *YellowCardProvider) CancelPayout(ctx context.Context, referenceID string) error {
	resp, err := p.makeRequest(ctx, "POST", fmt.Sprintf("%s/payments/%s/cancel", p.config.BaseURL, referenceID), nil)
	if err != nil {
		return fmt.Errorf("YellowCard cancellation failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("YellowCard cancellation error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (p *YellowCardProvider) makeRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
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
	req.Header.Set("X-YC-API-KEY", p.config.APIKey)
	req.Header.Set("X-YC-SECRET-KEY", p.config.SecretKey)
	if p.config.BusinessID != "" {
		req.Header.Set("X-YC-BUSINESS-ID", p.config.BusinessID)
	}

	return p.httpClient.Do(req)
}
