package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OrangeMoneyConfig holds Orange Money API configuration
type OrangeMoneyConfig struct {
	ClientID     string
	ClientSecret string
	MerchantKey  string
	BaseURL      string
	Country      string // CI, SN, ML, etc.
}

// OrangeMoneyCollectionProvider implements Orange Money payment collection
type OrangeMoneyCollectionProvider struct {
	config      OrangeMoneyConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// NewOrangeMoneyCollectionProvider creates a new Orange Money provider
func NewOrangeMoneyCollectionProvider(config OrangeMoneyConfig) *OrangeMoneyCollectionProvider {
	return &OrangeMoneyCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (o *OrangeMoneyCollectionProvider) GetName() string {
	return "orange_money"
}

func (o *OrangeMoneyCollectionProvider) GetSupportedCountries() []string {
	return []string{"CI", "SN", "ML", "BF", "CM", "GN", "NE", "MG", "CD"}
}

func (o *OrangeMoneyCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{
		CollectionMethodMobileMoney,
	}, nil
}

// getAccessToken retrieves OAuth 2.0 token from Orange
func (o *OrangeMoneyCollectionProvider) getAccessToken(ctx context.Context) error {
	if time.Now().Before(o.tokenExpiry) && o.accessToken != "" {
		return nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.config.BaseURL+"/oauth/v2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.SetBasicAuth(o.config.ClientID, o.config.ClientSecret)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	o.accessToken = tokenResp.AccessToken
	o.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

type orangeMoneyPaymentRequest struct {
	MerchantKey string `json:"merchant_key"`
	Currency    string `json:"currency"`
	OrderID     string `json:"order_id"`
	Amount      int    `json:"amount"`
	ReturnURL   string `json:"return_url"`
	CancelURL   string `json:"cancel_url"`
	NotifURL    string `json:"notif_url"`
	Lang        string `json:"lang"`
	Reference   string `json:"reference"`
}

type orangeMoneyPaymentResponse struct {
	PayToken   string `json:"pay_token"`
	PaymentURL string `json:"payment_url"`
	NotifToken string `json:"notif_token"`
}

func (o *OrangeMoneyCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	if err := o.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	omReq := orangeMoneyPaymentRequest{
		MerchantKey: o.config.MerchantKey,
		Currency:    req.Currency,
		OrderID:     req.ReferenceID,
		Amount:      int(req.Amount),
		ReturnURL:   req.RedirectURL + "?success=true",
		CancelURL:   req.RedirectURL + "?cancelled=true",
		NotifURL:    req.RedirectURL,
		Lang:        "fr",
		Reference:   req.ReferenceID,
	}

	jsonData, err := json.Marshal(omReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Country-specific URL
	paymentURL := fmt.Sprintf("%s/%s/v1/webpayment", o.config.BaseURL, strings.ToLower(req.Country))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", paymentURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+o.accessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("orange money API error (status %d): %s", resp.StatusCode, string(body))
	}

	var omResp orangeMoneyPaymentResponse
	if err := json.Unmarshal(body, &omResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: omResp.PayToken,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       omResp.PaymentURL,
		Message:           "Redirecting to Orange Money",
		Metadata: map[string]string{
			"notif_token": omResp.NotifToken,
		},
	}, nil
}

func (o *OrangeMoneyCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	if err := o.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	// Orange Money verification endpoint
	url := fmt.Sprintf("%s/%s/v1/transactionstatus/%s", o.config.BaseURL, strings.ToLower(o.config.Country), referenceID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+o.accessToken)

	resp, err := o.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		Status   string `json:"status"`
		OrderID  string `json:"order_id"`
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.Status == "SUCCESS" {
		status = CollectionStatusSuccessful
	} else if result.Status == "FAILED" || result.Status == "EXPIRED" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       result.OrderID,
		ProviderReference: result.OrderID,
		Status:            status,
		Amount:            float64(result.Amount),
		Currency:          result.Currency,
	}, nil
}
