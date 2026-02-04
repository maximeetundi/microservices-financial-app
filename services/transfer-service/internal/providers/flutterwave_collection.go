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
	PublicKey     string
	SecretKey     string
	EncryptionKey string
	WebhookSecret string
	BaseURL       string
}

// FlutterwaveCollectionProvider implements Flutterwave payment collection
type FlutterwaveCollectionProvider struct {
	config     FlutterwaveConfig
	httpClient *http.Client
}

// NewFlutterwaveCollectionProvider creates a new Flutterwave provider
func NewFlutterwaveCollectionProvider(config FlutterwaveConfig) *FlutterwaveCollectionProvider {
	return &FlutterwaveCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (f *FlutterwaveCollectionProvider) GetName() string {
	return "flutterwave"
}

func (f *FlutterwaveCollectionProvider) GetSupportedCountries() []string {
	return []string{"NG", "GH", "KE", "UG", "TZ", "ZA", "RW", "CI", "CM", "SN"}
}

func (f *FlutterwaveCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	// Flutterwave supports card, bank transfer, and mobile money
	methods := []CollectionMethod{
		CollectionMethodCard,
		CollectionMethodBankTransfer,
		CollectionMethodMobileMoney,
	}
	return methods, nil
}

// FlutterwaveCollectionRequest represents the request body for Flutterwave
type flutterwavePaymentRequest struct {
	TxRef          string                   `json:"tx_ref"`
	Amount         float64                  `json:"amount"`
	Currency       string                   `json:"currency"`
	RedirectURL    string                   `json:"redirect_url"`
	Customer       flutterwaveCustomer      `json:"customer"`
	Customizations flutterwaveCustomization `json:"customizations,omitempty"`
	PaymentOptions string                   `json:"payment_options,omitempty"`
}

type flutterwaveCustomer struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber,omitempty"`
	Name        string `json:"name,omitempty"`
}

type flutterwaveCustomization struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Logo        string `json:"logo,omitempty"`
}

type flutterwavePaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

func (f *FlutterwaveCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// Prepare Flutterwave request
	fwReq := flutterwavePaymentRequest{
		TxRef:       req.ReferenceID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		RedirectURL: req.RedirectURL,
		Customer: flutterwaveCustomer{
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Name:        fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		},
		Customizations: flutterwaveCustomization{
			Title:       "Payment",
			Description: fmt.Sprintf("Deposit %.2f %s", req.Amount, req.Currency),
		},
		PaymentOptions: "card,banktransfer,mobilemoney",
	}

	// Marshal request
	jsonData, err := json.Marshal(fwReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", f.config.BaseURL+"/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+f.config.SecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := f.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flutterwave API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var fwResp flutterwavePaymentResponse
	if err := json.Unmarshal(body, &fwResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if fwResp.Status != "success" {
		return nil, fmt.Errorf("flutterwave error: %s", fwResp.Message)
	}

	// Return collection response
	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: req.ReferenceID, // Flutterwave uses tx_ref
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       fwResp.Data.Link,
		Message:           "Click the link to complete payment",
	}, nil
}

func (f *FlutterwaveCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	// Verify transaction with Flutterwave
	url := fmt.Sprintf("%s/transactions/verify_by_reference?tx_ref=%s", f.config.BaseURL, referenceID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+f.config.SecretKey)

	resp, err := f.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("verification error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			Status        string  `json:"status"`
			Amount        float64 `json:"amount"`
			Currency      string  `json:"currency"`
			ID            int     `json:"id"`
			TxRef         string  `json:"tx_ref"`
			FlwRef        string  `json:"flw_ref"`
			ChargedAmount float64 `json:"charged_amount"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.Data.Status == "successful" {
		status = CollectionStatusSuccessful
	} else if result.Data.Status == "failed" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       result.Data.TxRef,
		ProviderReference: result.Data.FlwRef,
		Status:            status,
		Amount:            result.Data.Amount,
		Currency:          result.Data.Currency,
	}, nil
}
