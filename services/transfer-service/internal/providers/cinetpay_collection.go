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
	APIKey  string
	SiteID  string
	BaseURL string
}

// CinetPayCollectionProvider implements CinetPay payment collection
type CinetPayCollectionProvider struct {
	config     CinetPayConfig
	httpClient *http.Client
}

// NewCinetPayCollectionProvider creates a new CinetPay provider
func NewCinetPayCollectionProvider(config CinetPayConfig) *CinetPayCollectionProvider {
	return &CinetPayCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *CinetPayCollectionProvider) GetName() string {
	return "cinetpay"
}

func (c *CinetPayCollectionProvider) GetSupportedCountries() []string {
	return []string{"CI", "SN", "BF", "ML", "NE", "TG", "BJ", "GW", "CM", "GA", "CG", "TD", "GN", "CD"}
}

func (c *CinetPayCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{
		CollectionMethodCard,
		CollectionMethodMobileMoney,
	}, nil
}

type cinetpayPaymentRequest struct {
	APIKey              string `json:"apikey"`
	SiteID              string `json:"site_id"`
	TransactionID       string `json:"transaction_id"`
	Amount              int    `json:"amount"` // Must be multiple of 5
	Currency            string `json:"currency"`
	Description         string `json:"description"`
	ReturnURL           string `json:"return_url"`
	NotifyURL           string `json:"notify_url"`
	CustomerName        string `json:"customer_name,omitempty"`
	CustomerSurname     string `json:"customer_surname,omitempty"`
	CustomerEmail       string `json:"customer_email,omitempty"`
	CustomerPhoneNumber string `json:"customer_phone_number,omitempty"`
}

type cinetpayPaymentResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PaymentURL   string `json:"payment_url"`
		PaymentToken string `json:"payment_token"`
	} `json:"data"`
}

func (c *CinetPayCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// Round amount to nearest multiple of 5
	amount := int(req.Amount)
	if amount%5 != 0 {
		amount = ((amount + 4) / 5) * 5
	}

	cpReq := cinetpayPaymentRequest{
		APIKey:              c.config.APIKey,
		SiteID:              c.config.SiteID,
		TransactionID:       req.ReferenceID,
		Amount:              amount,
		Currency:            req.Currency,
		Description:         fmt.Sprintf("Deposit %d %s", amount, req.Currency),
		ReturnURL:           req.RedirectURL,
		NotifyURL:           req.RedirectURL, // Use same for now
		CustomerEmail:       req.Email,
		CustomerPhoneNumber: req.PhoneNumber,
		CustomerName:        req.FirstName,
		CustomerSurname:     req.LastName,
	}

	jsonData, err := json.Marshal(cpReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/payment", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cinetpay API error (status %d): %s", resp.StatusCode, string(body))
	}

	var cpResp cinetpayPaymentResponse
	if err := json.Unmarshal(body, &cpResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if cpResp.Code != "00" {
		return nil, fmt.Errorf("cinetpay error: %s", cpResp.Message)
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: cpResp.Data.PaymentToken,
		Status:            CollectionStatusPending,
		Amount:            float64(amount),
		Currency:          req.Currency,
		PaymentLink:       cpResp.Data.PaymentURL,
		Message:           "Redirecting to CinetPay checkout",
	}, nil
}

func (c *CinetPayCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	// CinetPay verification
	verifyReq := map[string]string{
		"apikey":         c.config.APIKey,
		"site_id":        c.config.SiteID,
		"transaction_id": referenceID,
	}

	jsonData, err := json.Marshal(verifyReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/payment/check", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		Code string `json:"code"`
		Data struct {
			Status        string `json:"status"`
			TransactionID string `json:"transaction_id"`
			Amount        int    `json:"amount"`
			Currency      string `json:"currency"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.Data.Status == "ACCEPTED" || result.Data.Status == "00" {
		status = CollectionStatusSuccessful
	} else if result.Data.Status == "REFUSED" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       result.Data.TransactionID,
		ProviderReference: result.Data.TransactionID,
		Status:            status,
		Amount:            float64(result.Data.Amount),
		Currency:          result.Data.Currency,
	}, nil
}
