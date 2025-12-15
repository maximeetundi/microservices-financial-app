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

// ThunesConfig holds Thunes API configuration
type ThunesConfig struct {
	APIKey      string // THUNES_API_KEY
	APISecret   string // THUNES_API_SECRET
	BaseURL     string // Default: https://api.thunes.com
	CallbackURL string // Webhook callback URL
}

// ThunesProvider implements PayoutProvider for Thunes (global coverage)
type ThunesProvider struct {
	config     ThunesConfig
	httpClient *http.Client
}

// NewThunesProvider creates a new Thunes provider
func NewThunesProvider(config ThunesConfig) *ThunesProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.thunes.com/v2"
	}
	
	return &ThunesProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (t *ThunesProvider) GetName() string {
	return "thunes"
}

func (t *ThunesProvider) GetSupportedCountries() []string {
	// Thunes supports 130+ countries - listing main ones
	return []string{
		// Asia
		"PH", "ID", "VN", "TH", "MY", "SG", "IN", "BD", "PK", "NP", "LK",
		// Middle East
		"AE", "SA", "KW", "QA", "BH", "OM",
		// Africa
		"NG", "GH", "KE", "UG", "TZ", "ZA", "EG", "MA", "TN",
		// Latin America
		"MX", "BR", "CO", "PE", "CL", "AR",
		// Europe
		"GB", "DE", "FR", "ES", "IT", "NL", "BE", "PL",
	}
}

func (t *ThunesProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	methods := []AvailableMethod{}
	
	// E-wallets (Asia primarily)
	eWalletCountries := map[string]bool{
		"PH": true, "ID": true, "TH": true, "MY": true, "VN": true,
	}
	if eWalletCountries[country] {
		methods = append(methods, AvailableMethod{
			Method:           PayoutMethodWallet,
			Name:             "E-Wallet",
			EstimatedMinutes: 5,
			Fee:              1.5,
			FeeType:          "percentage",
			MinAmount:        1,
			MaxAmount:        10000,
		})
	}
	
	// Bank transfer everywhere
	methods = append(methods, AvailableMethod{
		Method:           PayoutMethodBankTransfer,
		Name:             "Bank Transfer",
		EstimatedMinutes: 60 * 24, // 1 day
		Fee:              2.5,
		FeeType:          "percentage",
		MinAmount:        10,
		MaxAmount:        100000,
	})
	
	// Cash pickup in select countries
	cashPickupCountries := map[string]bool{
		"PH": true, "MX": true, "CO": true, "PE": true, "NG": true, "GH": true,
	}
	if cashPickupCountries[country] {
		methods = append(methods, AvailableMethod{
			Method:           PayoutMethodCashPickup,
			Name:             "Cash Pickup",
			EstimatedMinutes: 30,
			Fee:              3.0,
			FeeType:          "percentage",
			MinAmount:        10,
			MaxAmount:        3000,
		})
	}
	
	// Mobile money in Africa
	mobileMoneyCountries := map[string]bool{
		"NG": true, "GH": true, "KE": true, "UG": true, "TZ": true,
	}
	if mobileMoneyCountries[country] {
		methods = append(methods, AvailableMethod{
			Method:           PayoutMethodMobileMoney,
			Name:             "Mobile Money",
			EstimatedMinutes: 5,
			Fee:              1.0,
			FeeType:          "percentage",
			MinAmount:        1,
			MaxAmount:        5000,
		})
	}
	
	return methods, nil
}

func (t *ThunesProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	// Thunes provides bank lists via API
	url := fmt.Sprintf("%s/payers?country=%s&service_type=bank", t.config.BaseURL, country)
	
	resp, err := t.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var payers []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&payers); err != nil {
		// Return some default banks on error
		return t.getDefaultBanks(country), nil
	}
	
	banks := make([]Bank, len(payers))
	for i, p := range payers {
		banks[i] = Bank{
			Code:    p.ID,
			Name:    p.Name,
			Country: country,
		}
	}
	
	return banks, nil
}

func (t *ThunesProvider) getDefaultBanks(country string) []Bank {
	// Some default banks by country
	defaults := map[string][]Bank{
		"PH": {
			{Code: "BDO", Name: "BDO Unibank", Country: "PH"},
			{Code: "BPI", Name: "Bank of the Philippine Islands", Country: "PH"},
			{Code: "GCASH", Name: "GCash", Country: "PH"},
		},
		"ID": {
			{Code: "BCA", Name: "Bank Central Asia", Country: "ID"},
			{Code: "BRI", Name: "Bank Rakyat Indonesia", Country: "ID"},
			{Code: "MANDIRI", Name: "Bank Mandiri", Country: "ID"},
		},
		"IN": {
			{Code: "SBI", Name: "State Bank of India", Country: "IN"},
			{Code: "HDFC", Name: "HDFC Bank", Country: "IN"},
			{Code: "ICICI", Name: "ICICI Bank", Country: "IN"},
		},
	}
	
	if banks, ok := defaults[country]; ok {
		return banks
	}
	return []Bank{}
}

func (t *ThunesProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	operators := map[string][]MobileOperator{
		"PH": {
			{Code: "GCASH", Name: "GCash", Countries: []string{"PH"}},
			{Code: "PAYMAYA", Name: "PayMaya", Countries: []string{"PH"}},
		},
		"ID": {
			{Code: "OVO", Name: "OVO", Countries: []string{"ID"}},
			{Code: "GOPAY", Name: "GoPay", Countries: []string{"ID"}},
			{Code: "DANA", Name: "DANA", Countries: []string{"ID"}},
		},
		"TH": {
			{Code: "TRUEMONEY", Name: "TrueMoney", Countries: []string{"TH"}},
			{Code: "PROMPTPAY", Name: "PromptPay", Countries: []string{"TH"}},
		},
		"VN": {
			{Code: "MOMO", Name: "MoMo", Countries: []string{"VN"}},
			{Code: "ZALOPAY", Name: "ZaloPay", Countries: []string{"VN"}},
		},
	}
	
	if ops, ok := operators[country]; ok {
		return ops, nil
	}
	return []MobileOperator{}, nil
}

func (t *ThunesProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	switch req.PayoutMethod {
	case PayoutMethodMobileMoney, PayoutMethodWallet:
		if req.MobileNumber == "" {
			return fmt.Errorf("mobile/wallet number required")
		}
	case PayoutMethodBankTransfer:
		if req.AccountNumber == "" {
			return fmt.Errorf("account number required")
		}
	case PayoutMethodCashPickup:
		if req.RecipientName == "" || req.RecipientPhone == "" {
			return fmt.Errorf("recipient name and phone required for cash pickup")
		}
	}
	return nil
}

func (t *ThunesProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	var feePercent float64
	switch req.PayoutMethod {
	case PayoutMethodMobileMoney, PayoutMethodWallet:
		feePercent = 0.015
	case PayoutMethodBankTransfer:
		feePercent = 0.025
	case PayoutMethodCashPickup:
		feePercent = 0.03
	default:
		feePercent = 0.02
	}
	
	fee := req.Amount * feePercent
	
	return &PayoutResponse{
		ProviderName:   t.GetName(),
		ReferenceID:    req.ReferenceID,
		Status:         PayoutStatusPending,
		Fee:            fee,
		AmountReceived: req.Amount,
		ExchangeRate:   1.0,
	}, nil
}

func (t *ThunesProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	payload := map[string]interface{}{
		"external_id":     req.ReferenceID,
		"payer_id":        t.getPayerID(req),
		"mode":            "DESTINATION_AMOUNT",
		"destination": map[string]interface{}{
			"amount":   req.Amount,
			"currency": req.Currency,
		},
		"beneficiary": map[string]interface{}{
			"firstname":      req.RecipientName,
			"mobile_number":  req.MobileNumber,
			"account_number": req.AccountNumber,
		},
	}
	
	body, _ := json.Marshal(payload)
	
	resp, err := t.makeRequest(ctx, "POST", t.config.BaseURL+"/transactions", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, _ := io.ReadAll(resp.Body)
	
	var response struct {
		ID              string  `json:"id"`
		Status          string  `json:"status"`
		DestinationAmt  float64 `json:"destination_amount"`
		Fee             float64 `json:"fee"`
		WholesaleFxRate float64 `json:"wholesale_fx_rate"`
	}
	
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	
	return &PayoutResponse{
		ProviderName:      t.GetName(),
		ProviderReference: response.ID,
		ReferenceID:       req.ReferenceID,
		Status:            PayoutStatusProcessing,
		Fee:               response.Fee,
		ExchangeRate:      response.WholesaleFxRate,
		AmountReceived:    response.DestinationAmt,
	}, nil
}

func (t *ThunesProvider) getPayerID(req *PayoutRequest) string {
	// Map to Thunes payer IDs based on method and country
	// In production, this would be configured per country/method
	return fmt.Sprintf("%s_%s_%s", req.RecipientCountry, req.PayoutMethod, req.MobileOperator)
}

func (t *ThunesProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	url := fmt.Sprintf("%s/transactions?external_id=%s", t.config.BaseURL, referenceID)
	
	resp, err := t.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var transactions []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}
	
	if len(transactions) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}
	
	tx := transactions[0]
	var status PayoutStatus
	switch tx.Status {
	case "CREATED", "SUBMITTED":
		status = PayoutStatusPending
	case "CONFIRMED":
		status = PayoutStatusProcessing
	case "COMPLETED":
		status = PayoutStatusCompleted
	case "REJECTED", "CANCELLED":
		status = PayoutStatusFailed
	default:
		status = PayoutStatusPending
	}
	
	return &PayoutStatusResponse{
		ReferenceID:       referenceID,
		ProviderReference: tx.ID,
		Status:            status,
	}, nil
}

func (t *ThunesProvider) CancelPayout(ctx context.Context, referenceID string) error {
	url := fmt.Sprintf("%s/transactions/%s/cancel", t.config.BaseURL, referenceID)
	
	resp, err := t.makeRequest(ctx, "POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to cancel transaction")
	}
	
	return nil
}

func (t *ThunesProvider) makeRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.SetBasicAuth(t.config.APIKey, t.config.APISecret)
	req.Header.Set("Content-Type", "application/json")
	
	return t.httpClient.Do(req)
}
