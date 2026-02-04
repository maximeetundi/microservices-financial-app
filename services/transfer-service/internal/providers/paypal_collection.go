package providers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PayPalConfig holds PayPal API configuration
type PayPalConfig struct {
	ClientID     string
	ClientSecret string
	Mode         string // "sandbox" or "live"
	BaseURL      string
}

// PayPalCollectionProvider implements PayPal payment collection
type PayPalCollectionProvider struct {
	config      PayPalConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// NewPayPalCollectionProvider creates a new PayPal provider
func NewPayPalCollectionProvider(config PayPalConfig) *PayPalCollectionProvider {
	return &PayPalCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *PayPalCollectionProvider) GetName() string {
	return "paypal"
}

func (p *PayPalCollectionProvider) GetSupportedCountries() []string {
	return []string{"*"} // PayPal is global
}

func (p *PayPalCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{
		CollectionMethodCard,
	}, nil
}

// getAccessToken retrieves or refreshes PayPal OAuth token
func (p *PayPalCollectionProvider) getAccessToken(ctx context.Context) error {
	// Check if token is still valid
	if time.Now().Before(p.tokenExpiry) && p.accessToken != "" {
		return nil
	}

	// Request new token
	auth := base64.StdEncoding.EncodeToString([]byte(p.config.ClientID + ":" + p.config.ClientSecret))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/v1/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(httpReq)
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
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	p.accessToken = tokenResp.AccessToken
	p.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

type paypalOrderRequest struct {
	Intent             string                   `json:"intent"`
	PurchaseUnits      []paypalPurchaseUnit     `json:"purchase_units"`
	ApplicationContext paypalApplicationContext `json:"application_context"`
}

type paypalPurchaseUnit struct {
	Amount paypalAmount `json:"amount"`
}

type paypalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type paypalApplicationContext struct {
	ReturnURL string `json:"return_url"`
	CancelURL string `json:"cancel_url"`
}

type paypalOrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Links  []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
}

func (p *PayPalCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// Always get fresh token
	if err := p.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	orderReq := paypalOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []paypalPurchaseUnit{
			{
				Amount: paypalAmount{
					CurrencyCode: req.Currency,
					Value:        fmt.Sprintf("%.2f", req.Amount),
				},
			},
		},
		ApplicationContext: paypalApplicationContext{
			ReturnURL: req.RedirectURL + "?success=true",
			CancelURL: req.RedirectURL + "?cancelled=true",
		},
	}

	jsonData, err := json.Marshal(orderReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/v2/checkout/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.accessToken)
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("paypal API error (status %d): %s", resp.StatusCode, string(body))
	}

	var orderResp paypalOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// Find approval link
	var approvalURL string
	for _, link := range orderResp.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: orderResp.ID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       approvalURL,
		Message:           "Redirecting to PayPal",
	}, nil
}

func (p *PayPalCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	if err := p.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s", p.config.BaseURL, referenceID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.accessToken)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result paypalOrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.Status == "COMPLETED" || result.Status == "APPROVED" {
		status = CollectionStatusSuccessful
	} else if result.Status == "VOIDED" || result.Status == "CANCELLED" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       referenceID,
		ProviderReference: result.ID,
		Status:            status,
	}, nil
}
