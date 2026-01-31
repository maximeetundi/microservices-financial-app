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

// CinetPayConfig holds CinetPay API configuration
type CinetPayConfig struct {
	APIKey    string
	SiteID    string
	SecretKey string
	BaseURL   string
}

// CinetPayProvider implements PayoutProvider for CinetPay
type CinetPayProvider struct {
	config     CinetPayConfig
	httpClient *http.Client
}

// NewCinetPayProvider creates a new CinetPay provider
func NewCinetPayProvider(config CinetPayConfig) *CinetPayProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api-checkout.cinetpay.com/v2"
	}
	return &CinetPayProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *CinetPayProvider) GetName() string {
	return "cinetpay"
}

func (p *CinetPayProvider) GetSupportedCountries() []string {
	return []string{
		"CI", // Côte d'Ivoire
		"SN", // Sénégal
		"CM", // Cameroun
		"BF", // Burkina Faso
		"TG", // Togo
		"BJ", // Bénin
		"ML", // Mali
		"GN", // Guinée
	}
}

func (p *CinetPayProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	methods := []AvailableMethod{
		{
			Code:           "mobile_money",
			Name:           "Mobile Money",
			Type:           "mobile",
			Countries:      []string{"CI", "SN", "CM", "BF", "TG", "BJ", "ML", "GN"},
			Currencies:     []string{"XOF", "XAF"},
			RequiredFields: []string{"phone"},
		},
		{
			Code:           "bank_transfer",
			Name:           "Bank Transfer",
			Type:           "bank",
			Countries:      []string{"CI", "SN", "CM"},
			Currencies:     []string{"XOF", "XAF"},
			RequiredFields: []string{"bank_code", "account_number"},
		},
	}
	return methods, nil
}

func (p *CinetPayProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	banks := []Bank{
		{Code: "BICICI", Name: "BICICI", Country: "CI"},
		{Code: "SGBCI", Name: "Société Générale CI", Country: "CI"},
		{Code: "BOAD", Name: "BOAD", Country: "CI"},
		{Code: "ECOBANK", Name: "Ecobank", Country: "CI"},
	}
	return banks, nil
}

func (p *CinetPayProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	switch country {
	case "CI":
		return []MobileOperator{
			{Code: "ORANGE_CI", Name: "Orange Money", Country: "CI", NumberPrefix: []string{"07", "08", "09"}},
			{Code: "MTN_CI", Name: "MTN Mobile Money", Country: "CI", NumberPrefix: []string{"05", "04"}},
			{Code: "MOOV_CI", Name: "Moov Money", Country: "CI", NumberPrefix: []string{"01"}},
		}, nil
	case "SN":
		return []MobileOperator{
			{Code: "ORANGE_SN", Name: "Orange Money", Country: "SN", NumberPrefix: []string{"77", "78"}},
			{Code: "WAVE_SN", Name: "Wave", Country: "SN", NumberPrefix: []string{"70"}},
			{Code: "FREE_SN", Name: "Free Money", Country: "SN", NumberPrefix: []string{"76"}},
		}, nil
	default:
		return []MobileOperator{}, nil
	}
}

func (p *CinetPayProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.RecipientPhone == "" && req.BankAccountNumber == "" {
		return fmt.Errorf("recipient phone or bank account required")
	}
	return nil
}

func (p *CinetPayProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	fee := req.Amount * 0.015 // 1.5% fee estimate
	return &PayoutResponse{
		ProviderName:      "cinetpay",
		ProviderReference: "",
		Amount:            req.Amount,
		Currency:          req.Currency,
		Fee:               fee,
		TotalAmount:       req.Amount + fee,
		Status:            "quote",
	}, nil
}

func (p *CinetPayProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	payload := map[string]interface{}{
		"apikey":         p.config.APIKey,
		"site_id":        p.config.SiteID,
		"transaction_id": req.ExternalReference,
		"amount":         req.Amount,
		"currency":       req.Currency,
		"description":    req.Description,
		"receiver":       req.RecipientPhone,
		"notify_url":     req.CallbackURL,
	}

	body, _ := json.Marshal(payload)
	resp, err := p.makeRequest(ctx, "POST", p.config.BaseURL+"/payment/deposit", body)
	if err != nil {
		return nil, fmt.Errorf("CinetPay payout failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TransactionID string `json:"transaction_id"`
			Status        string `json:"status"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse CinetPay response: %w", err)
	}

	if result.Code != "00" {
		return nil, fmt.Errorf("CinetPay error: %s", result.Message)
	}

	fee := req.Amount * 0.015
	return &PayoutResponse{
		ProviderName:      "cinetpay",
		ProviderReference: result.Data.TransactionID,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Fee:               fee,
		TotalAmount:       req.Amount + fee,
		Status:            result.Data.Status,
	}, nil
}

func (p *CinetPayProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	payload := map[string]interface{}{
		"apikey":         p.config.APIKey,
		"site_id":        p.config.SiteID,
		"transaction_id": referenceID,
	}

	body, _ := json.Marshal(payload)
	resp, err := p.makeRequest(ctx, "POST", p.config.BaseURL+"/payment/check", body)
	if err != nil {
		return nil, fmt.Errorf("CinetPay status check failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Status string `json:"status"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse CinetPay status response: %w", err)
	}

	status := "pending"
	switch result.Data.Status {
	case "ACCEPTED":
		status = "completed"
	case "REFUSED", "ERROR":
		status = "failed"
	case "PENDING":
		status = "pending"
	}

	return &PayoutStatusResponse{
		Status:    status,
		UpdatedAt: time.Now(),
	}, nil
}

func (p *CinetPayProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("CinetPay does not support payout cancellation")
}

func (p *CinetPayProvider) makeRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return p.httpClient.Do(req)
}
