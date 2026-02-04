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

// PaystackConfig holds Paystack API configuration
type PaystackConfig struct {
	PublicKey string
	SecretKey string
	BaseURL   string
}

// PaystackCollectionProvider implements Paystack payment collection
type PaystackCollectionProvider struct {
	config     PaystackConfig
	httpClient *http.Client
}

// NewPaystackCollectionProvider creates a new Paystack provider
func NewPaystackCollectionProvider(config PaystackConfig) *PaystackCollectionProvider {
	return &PaystackCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *PaystackCollectionProvider) GetName() string {
	return "paystack"
}

func (p *PaystackCollectionProvider) GetSupportedCountries() []string {
	return []string{"NG", "GH", "ZA", "KE"}
}

func (p *PaystackCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{
		CollectionMethodCard,
		CollectionMethodBankTransfer,
		CollectionMethodMobileMoney,
	}, nil
}

type paystackInitializeRequest struct {
	Email       string `json:"email"`
	Amount      int    `json:"amount"` // In kobo/pesewas (smallest currency unit)
	Currency    string `json:"currency,omitempty"`
	Reference   string `json:"reference"`
	CallbackURL string `json:"callback_url,omitempty"`
}

type paystackInitializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

func (p *PaystackCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// Convert amount to smallest unit (kobo for NGN, pesewas for GHS)
	amountInKobo := int(req.Amount * 100)

	psReq := paystackInitializeRequest{
		Email:       req.Email,
		Amount:      amountInKobo,
		Currency:    req.Currency,
		Reference:   req.ReferenceID,
		CallbackURL: req.RedirectURL,
	}

	jsonData, err := json.Marshal(psReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.config.SecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paystack API error (status %d): %s", resp.StatusCode, string(body))
	}

	var psResp paystackInitializeResponse
	if err := json.Unmarshal(body, &psResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if !psResp.Status {
		return nil, fmt.Errorf("paystack error: %s", psResp.Message)
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: psResp.Data.AccessCode,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       psResp.Data.AuthorizationURL,
		Message:           "Redirecting to Paystack",
	}, nil
}

func (p *PaystackCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", p.config.BaseURL, referenceID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.config.SecretKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Status    string `json:"status"`
			Reference string `json:"reference"`
			Amount    int    `json:"amount"`
			Currency  string `json:"currency"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.Data.Status == "success" {
		status = CollectionStatusSuccessful
	} else if result.Data.Status == "failed" || result.Data.Status == "abandoned" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       result.Data.Reference,
		ProviderReference: result.Data.Reference,
		Status:            status,
		Amount:            float64(result.Data.Amount) / 100, // Convert back from kobo
		Currency:          result.Data.Currency,
	}, nil
}
