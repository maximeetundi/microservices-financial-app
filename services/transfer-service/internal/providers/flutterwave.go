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

// FlutterwaveConfig holds Flutterwave API configuration
type FlutterwaveConfig struct {
	SecretKey   string // FLUTTERWAVE_SECRET_KEY
	PublicKey   string // FLUTTERWAVE_PUBLIC_KEY
	EncryptKey  string // FLUTTERWAVE_ENCRYPT_KEY
	BaseURL     string // Default: https://api.flutterwave.com/v3
	CallbackURL string // Webhook callback URL
}

// FlutterwaveProvider implements PayoutProvider for Flutterwave
type FlutterwaveProvider struct {
	config     FlutterwaveConfig
	httpClient *http.Client
}

// NewFlutterwaveProvider creates a new Flutterwave provider
func NewFlutterwaveProvider(config FlutterwaveConfig) *FlutterwaveProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.flutterwave.com/v3"
	}
	
	return &FlutterwaveProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *FlutterwaveProvider) GetName() string {
	return "flutterwave"
}

func (f *FlutterwaveProvider) GetSupportedCountries() []string {
	return []string{
		// West Africa
		"NG", "GH", "CI", "SN", "CM", "BJ", "TG", "BF", "ML", "NE",
		// East Africa
		"KE", "UG", "TZ", "RW",
		// Central Africa
		"CD", "GA", "CG",
		// Southern Africa
		"ZA", "ZM", "MW",
	}
}

func (f *FlutterwaveProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	methods := []AvailableMethod{}
	
	// Mobile Money available in most African countries
	mobileMoneyCountries := map[string]bool{
		"GH": true, "KE": true, "UG": true, "TZ": true, "RW": true,
		"CI": true, "SN": true, "CM": true, "BJ": true, "TG": true,
		"BF": true, "ML": true, "NE": true, "ZM": true, "MW": true,
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
	
	// Bank Transfer available everywhere
	methods = append(methods, AvailableMethod{
		Method:           PayoutMethodBankTransfer,
		Name:             "Bank Transfer",
		EstimatedMinutes: 60 * 24, // 1 day
		Fee:              2.0,
		FeeType:          "percentage",
		MinAmount:        10,
		MaxAmount:        50000,
	})
	
	return methods, nil
}

func (f *FlutterwaveProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	url := fmt.Sprintf("%s/banks/%s", f.config.BaseURL, country)
	
	resp, err := f.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	
	banks := make([]Bank, len(response.Data))
	for i, b := range response.Data {
		banks[i] = Bank{
			Code:    b.Code,
			Name:    b.Name,
			Country: country,
		}
	}
	
	return banks, nil
}

func (f *FlutterwaveProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	// Flutterwave supported mobile operators by country
	operators := map[string][]MobileOperator{
		"GH": {
			{Code: "MTN", Name: "MTN Mobile Money", Countries: []string{"GH"}},
			{Code: "VODAFONE", Name: "Vodafone Cash", Countries: []string{"GH"}},
			{Code: "TIGO", Name: "AirtelTigo Money", Countries: []string{"GH"}},
		},
		"KE": {
			{Code: "MPESA", Name: "M-Pesa", Countries: []string{"KE"}},
		},
		"UG": {
			{Code: "MTN", Name: "MTN Mobile Money", Countries: []string{"UG"}},
			{Code: "AIRTEL", Name: "Airtel Money", Countries: []string{"UG"}},
		},
		"CI": {
			{Code: "MTN", Name: "MTN Mobile Money", Countries: []string{"CI"}},
			{Code: "ORANGE", Name: "Orange Money", Countries: []string{"CI"}},
			{Code: "MOOV", Name: "Moov Money", Countries: []string{"CI"}},
		},
		"SN": {
			{Code: "ORANGE", Name: "Orange Money", Countries: []string{"SN"}},
			{Code: "FREE", Name: "Free Money", Countries: []string{"SN"}},
		},
		"CM": {
			{Code: "MTN", Name: "MTN Mobile Money", Countries: []string{"CM"}},
			{Code: "ORANGE", Name: "Orange Money", Countries: []string{"CM"}},
		},
	}
	
	if ops, ok := operators[country]; ok {
		return ops, nil
	}
	return []MobileOperator{}, nil
}

func (f *FlutterwaveProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	// Validate based on payout method
	switch req.PayoutMethod {
	case PayoutMethodMobileMoney:
		if req.MobileNumber == "" || req.MobileOperator == "" {
			return fmt.Errorf("mobile number and operator required for mobile money")
		}
	case PayoutMethodBankTransfer:
		if req.AccountNumber == "" || req.BankCode == "" {
			return fmt.Errorf("account number and bank code required for bank transfer")
		}
	}
	return nil
}

func (f *FlutterwaveProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	// Calculate fees based on method
	var fee float64
	switch req.PayoutMethod {
	case PayoutMethodMobileMoney:
		fee = req.Amount * 0.01 // 1%
	case PayoutMethodBankTransfer:
		fee = req.Amount * 0.02 // 2%
	default:
		fee = req.Amount * 0.015
	}
	
	return &PayoutResponse{
		ProviderName:   f.GetName(),
		ReferenceID:    req.ReferenceID,
		Status:         PayoutStatusPending,
		Fee:            fee,
		AmountReceived: req.Amount,
		ExchangeRate:   1.0, // Same currency assumed
	}, nil
}

func (f *FlutterwaveProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	var endpoint string
	var payload map[string]interface{}
	
	switch req.PayoutMethod {
	case PayoutMethodMobileMoney:
		endpoint = f.config.BaseURL + "/transfers"
		payload = map[string]interface{}{
			"account_bank":   req.MobileOperator,
			"account_number": req.MobileNumber,
			"amount":         req.Amount,
			"currency":       req.Currency,
			"narration":      req.Narration,
			"reference":      req.ReferenceID,
			"beneficiary_name": req.RecipientName,
			"callback_url":   f.config.CallbackURL,
			"meta": map[string]string{
				"mobile_number": req.MobileNumber,
				"sender_country": "internal",
			},
		}
		
	case PayoutMethodBankTransfer:
		endpoint = f.config.BaseURL + "/transfers"
		payload = map[string]interface{}{
			"account_bank":   req.BankCode,
			"account_number": req.AccountNumber,
			"amount":         req.Amount,
			"currency":       req.Currency,
			"narration":      req.Narration,
			"reference":      req.ReferenceID,
			"beneficiary_name": req.RecipientName,
			"callback_url":   f.config.CallbackURL,
		}
		
	default:
		return nil, fmt.Errorf("unsupported payout method: %s", req.PayoutMethod)
	}
	
	body, _ := json.Marshal(payload)
	resp, err := f.makeRequest(ctx, "POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, _ := io.ReadAll(resp.Body)
	
	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ID        int     `json:"id"`
			Reference string  `json:"reference"`
			Status    string  `json:"status"`
			Amount    float64 `json:"amount"`
			Fee       float64 `json:"fee"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	
	if response.Status != "success" {
		return nil, fmt.Errorf("flutterwave error: %s", response.Message)
	}
	
	return &PayoutResponse{
		ProviderName:      f.GetName(),
		ProviderReference: fmt.Sprintf("%d", response.Data.ID),
		ReferenceID:       req.ReferenceID,
		Status:            PayoutStatusProcessing,
		Fee:               response.Data.Fee,
		AmountReceived:    response.Data.Amount,
	}, nil
}

func (f *FlutterwaveProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	url := fmt.Sprintf("%s/transfers?reference=%s", f.config.BaseURL, referenceID)
	
	resp, err := f.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var response struct {
		Status string `json:"status"`
		Data   []struct {
			ID        int    `json:"id"`
			Reference string `json:"reference"`
			Status    string `json:"status"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("transfer not found")
	}
	
	transfer := response.Data[0]
	
	var status PayoutStatus
	switch transfer.Status {
	case "NEW", "PENDING":
		status = PayoutStatusPending
	case "PROCESSING":
		status = PayoutStatusProcessing
	case "SUCCESSFUL":
		status = PayoutStatusCompleted
	case "FAILED":
		status = PayoutStatusFailed
	default:
		status = PayoutStatusPending
	}
	
	return &PayoutStatusResponse{
		ReferenceID:       referenceID,
		ProviderReference: fmt.Sprintf("%d", transfer.ID),
		Status:            status,
	}, nil
}

func (f *FlutterwaveProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("cancel not supported for Flutterwave transfers")
}

func (f *FlutterwaveProvider) makeRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+f.config.SecretKey)
	req.Header.Set("Content-Type", "application/json")
	
	return f.httpClient.Do(req)
}
