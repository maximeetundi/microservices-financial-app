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

// StripeConfig holds Stripe API configuration
type StripeConfig struct {
	SecretKey   string // STRIPE_SECRET_KEY
	PublishableKey string // STRIPE_PUBLISHABLE_KEY (for frontend)
	WebhookSecret string // STRIPE_WEBHOOK_SECRET
	BaseURL     string // Default: https://api.stripe.com/v1
}

// StripeProvider implements PayoutProvider for Stripe (Europe, US, Canada)
type StripeProvider struct {
	config     StripeConfig
	httpClient *http.Client
}

// NewStripeProvider creates a new Stripe provider
func NewStripeProvider(config StripeConfig) *StripeProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.stripe.com/v1"
	}
	
	return &StripeProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *StripeProvider) GetName() string {
	return "stripe"
}

func (s *StripeProvider) GetSupportedCountries() []string {
	return []string{
		// Europe (SEPA zone)
		"AT", "BE", "BG", "HR", "CY", "CZ", "DK", "EE", "FI", "FR",
		"DE", "GR", "HU", "IE", "IT", "LV", "LT", "LU", "MT", "NL",
		"PL", "PT", "RO", "SK", "SI", "ES", "SE", "GB", "CH", "NO",
		// North America
		"US", "CA",
		// Oceania
		"AU", "NZ",
		// Asia
		"SG", "HK", "JP",
	}
}

func (s *StripeProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	methods := []AvailableMethod{}
	
	// SEPA countries
	sepaCountries := map[string]bool{
		"AT": true, "BE": true, "BG": true, "HR": true, "CY": true, "CZ": true,
		"DK": true, "EE": true, "FI": true, "FR": true, "DE": true, "GR": true,
		"HU": true, "IE": true, "IT": true, "LV": true, "LT": true, "LU": true,
		"MT": true, "NL": true, "PL": true, "PT": true, "RO": true, "SK": true,
		"SI": true, "ES": true, "SE": true, "GB": true, "CH": true, "NO": true,
	}
	
	if sepaCountries[country] {
		methods = append(methods, AvailableMethod{
			Method:           PayoutMethodSEPA,
			Name:             "SEPA Bank Transfer",
			EstimatedMinutes: 60 * 24, // 1 day
			Fee:              0.5,
			FeeType:          "flat", // €0.50 flat
			MinAmount:        1,
			MaxAmount:        1000000,
		})
		
		methods = append(methods, AvailableMethod{
			Method:           PayoutMethodSEPAInstant,
			Name:             "SEPA Instant",
			EstimatedMinutes: 1, // < 10 seconds
			Fee:              1.0,
			FeeType:          "flat", // €1.00 flat
			MinAmount:        1,
			MaxAmount:        100000,
		})
	}
	
	// US/Canada ACH/Wire
	if country == "US" || country == "CA" {
		methods = append(methods, AvailableMethod{
			Method:           PayoutMethodACH,
			Name:             "ACH Transfer",
			EstimatedMinutes: 60 * 48, // 2 days
			Fee:              0.25,
			FeeType:          "percentage",
			MinAmount:        1,
			MaxAmount:        1000000,
		})
		
		methods = append(methods, AvailableMethod{
			Method:           PayoutMethodWire,
			Name:             "Wire Transfer",
			EstimatedMinutes: 60 * 4, // Same day
			Fee:              25.0,
			FeeType:          "flat",
			MinAmount:        100,
			MaxAmount:        10000000,
		})
	}
	
	// Card payouts (Visa Direct, Mastercard Send)
	methods = append(methods, AvailableMethod{
		Method:           PayoutMethodCard,
		Name:             "Instant Card Payout",
		EstimatedMinutes: 30,
		Fee:              1.5,
		FeeType:          "percentage",
		MinAmount:        1,
		MaxAmount:        10000,
	})
	
	return methods, nil
}

func (s *StripeProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	// Stripe uses IBAN for Europe, routing+account for US
	// No bank list needed - user provides account details directly
	return []Bank{}, nil
}

func (s *StripeProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	// Stripe doesn't support mobile money
	return []MobileOperator{}, nil
}

func (s *StripeProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	switch req.PayoutMethod {
	case PayoutMethodSEPA, PayoutMethodSEPAInstant:
		if req.IBAN == "" {
			return fmt.Errorf("IBAN required for SEPA transfers")
		}
	case PayoutMethodACH:
		if req.AccountNumber == "" || req.RoutingNumber == "" {
			return fmt.Errorf("account and routing number required for ACH")
		}
	case PayoutMethodWire:
		if req.AccountNumber == "" || (req.SwiftCode == "" && req.RoutingNumber == "") {
			return fmt.Errorf("account number and SWIFT/routing required for wire")
		}
	case PayoutMethodCard:
		// Card validation happens on Stripe side
	}
	return nil
}

func (s *StripeProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	var fee float64
	
	switch req.PayoutMethod {
	case PayoutMethodSEPA:
		fee = 0.50 // €0.50 flat
	case PayoutMethodSEPAInstant:
		fee = 1.00 // €1.00 flat
	case PayoutMethodACH:
		fee = req.Amount * 0.0025 // 0.25%
		if fee < 0.10 {
			fee = 0.10
		}
	case PayoutMethodWire:
		fee = 25.00 // $25 flat
	case PayoutMethodCard:
		fee = req.Amount * 0.015 // 1.5%
	default:
		fee = req.Amount * 0.01
	}
	
	return &PayoutResponse{
		ProviderName:   s.GetName(),
		ReferenceID:    req.ReferenceID,
		Status:         PayoutStatusPending,
		Fee:            fee,
		AmountReceived: req.Amount,
		ExchangeRate:   1.0,
	}, nil
}

func (s *StripeProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	// First, create or get connected account for recipient
	// Then create payout to that account
	
	// For simplicity, using Stripe Payouts API directly
	// In production, you'd use Connect for multi-party payouts
	
	var payload map[string]interface{}
	
	switch req.PayoutMethod {
	case PayoutMethodSEPA, PayoutMethodSEPAInstant:
		// Create bank account token first, then payout
		payload = map[string]interface{}{
			"amount":      int(req.Amount * 100), // Stripe uses cents
			"currency":    "eur",
			"description": req.Narration,
			"metadata": map[string]string{
				"reference_id": req.ReferenceID,
				"recipient":    req.RecipientName,
			},
		}
		
	case PayoutMethodACH:
		payload = map[string]interface{}{
			"amount":      int(req.Amount * 100),
			"currency":    "usd",
			"method":      "standard", // or "instant"
			"description": req.Narration,
			"metadata": map[string]string{
				"reference_id": req.ReferenceID,
			},
		}
		
	default:
		return nil, fmt.Errorf("unsupported payout method for Stripe: %s", req.PayoutMethod)
	}
	
	body, _ := json.Marshal(payload)
	
	resp, err := s.makeRequest(ctx, "POST", s.config.BaseURL+"/payouts", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, _ := io.ReadAll(resp.Body)
	
	var response struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	}
	
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	
	return &PayoutResponse{
		ProviderName:      s.GetName(),
		ProviderReference: response.ID,
		ReferenceID:       req.ReferenceID,
		Status:            s.mapStatus(response.Status),
		AmountReceived:    float64(response.Amount) / 100,
	}, nil
}

func (s *StripeProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	// Search by metadata
	url := fmt.Sprintf("%s/payouts?limit=1", s.config.BaseURL)
	
	resp, err := s.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var response struct {
		Data []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("payout not found")
	}
	
	payout := response.Data[0]
	
	return &PayoutStatusResponse{
		ReferenceID:       referenceID,
		ProviderReference: payout.ID,
		Status:            s.mapStatus(payout.Status),
	}, nil
}

func (s *StripeProvider) mapStatus(stripeStatus string) PayoutStatus {
	switch stripeStatus {
	case "pending":
		return PayoutStatusPending
	case "in_transit":
		return PayoutStatusProcessing
	case "paid":
		return PayoutStatusCompleted
	case "failed", "canceled":
		return PayoutStatusFailed
	default:
		return PayoutStatusPending
	}
}

func (s *StripeProvider) CancelPayout(ctx context.Context, referenceID string) error {
	url := fmt.Sprintf("%s/payouts/%s/cancel", s.config.BaseURL, referenceID)
	
	resp, err := s.makeRequest(ctx, "POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to cancel payout")
	}
	
	return nil
}

func (s *StripeProvider) makeRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.SetBasicAuth(s.config.SecretKey, "")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Stripe-Version", "2023-10-16")
	
	return s.httpClient.Do(req)
}
