package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// PayPalConfig holds PayPal API configuration
type PayPalConfig struct {
	ClientID     string
	ClientSecret string
	WebhookID    string
	BaseURL      string // https://api-m.sandbox.paypal.com or https://api-m.paypal.com
	Mode         string // sandbox or live
}

// PayPalProvider implements payout and collection for PayPal
type PayPalProvider struct {
	config      PayPalConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

func normalizePayPalBaseURL(raw string, mode string) string {
	base := strings.TrimRight(strings.TrimSpace(raw), "/")
	// Strip version suffixes if user provided /v1 or /v2
	for _, suffix := range []string{"/v1", "/v2"} {
		if strings.HasSuffix(base, suffix) {
			base = strings.TrimSuffix(base, suffix)
		}
	}

	// Some users configure https://api.paypal.com which is not the correct REST host.
	// Use the documented REST hosts.
	lower := strings.ToLower(base)
	if strings.Contains(lower, "api.paypal.com") || strings.Contains(lower, "api.sandbox.paypal.com") {
		if strings.ToLower(mode) == "live" {
			return "https://api-m.paypal.com"
		}
		return "https://api-m.sandbox.paypal.com"
	}

	return base
}

func (p *PayPalProvider) apiBaseURL() string {
	return normalizePayPalBaseURL(p.config.BaseURL, p.config.Mode)
}

// NewPayPalProvider creates a new PayPal provider
func NewPayPalProvider(config PayPalConfig) *PayPalProvider {
	if config.BaseURL == "" {
		if config.Mode == "live" {
			config.BaseURL = "https://api-m.paypal.com"
		} else {
			config.BaseURL = "https://api-m.sandbox.paypal.com"
		}
	}
	config.BaseURL = normalizePayPalBaseURL(config.BaseURL, config.Mode)
	return &PayPalProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *PayPalProvider) GetName() string {
	return "paypal"
}

func (p *PayPalProvider) GetSupportedCountries() []string {
	return []string{
		// Americas
		"US", "CA", "MX", "BR", "AR",
		// Europe
		"GB", "DE", "FR", "ES", "IT", "NL", "BE", "AT", "CH", "SE", "DK", "NO", "FI", "PL", "PT", "IE",
		// Asia Pacific
		"AU", "NZ", "JP", "SG", "HK", "KR", "IN", "TH", "MY", "PH", "ID",
		// Middle East
		"AE", "IL", "SA",
	}
}

// ==================== AUTHENTICATION ====================

type PayPalTokenResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken retrieves or refreshes the OAuth access token
func (p *PayPalProvider) GetAccessToken(ctx context.Context) (string, error) {
	log.Printf("[PayPal] 🔐 GetAccessToken called")

	// Return cached token if still valid
	if p.accessToken != "" && time.Now().Before(p.tokenExpiry) {
		log.Printf("[PayPal] ✅ Using cached token (expires: %v)", p.tokenExpiry)
		return p.accessToken, nil
	}

	log.Printf("[PayPal] 🔄 Requesting new OAuth token...")
	log.Printf("[PayPal]    Mode: %s", p.config.Mode)
	log.Printf("[PayPal]    Base URL: %s", p.config.BaseURL)

	// Validate credentials before making request
	if p.config.ClientID == "" {
		log.Printf("[PayPal] ❌ ERROR: ClientID is EMPTY!")
		log.Printf("[PayPal] 💡 Configure 'client_id' in admin panel for PayPal")
		return "", fmt.Errorf("PayPal ClientID not configured - set in admin panel")
	}
	if p.config.ClientSecret == "" {
		log.Printf("[PayPal] ❌ ERROR: ClientSecret is EMPTY!")
		log.Printf("[PayPal] 💡 Configure 'client_secret' in admin panel for PayPal")
		return "", fmt.Errorf("PayPal ClientSecret not configured - set in admin panel")
	}

	// Log credential info (masked for security)
	clientIDMasked := p.config.ClientID
	if len(clientIDMasked) > 12 {
		clientIDMasked = clientIDMasked[:6] + "****" + clientIDMasked[len(clientIDMasked)-4:]
	}
	log.Printf("[PayPal]    Client ID: %s", clientIDMasked)
	log.Printf("[PayPal]    Client Secret: %d characters", len(p.config.ClientSecret))

	baseURL := p.apiBaseURL()
	tokenURL := baseURL + "/v1/oauth2/token"
	log.Printf("[PayPal] 📡 Sending token request to: %s", tokenURL)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL,
		bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		log.Printf("[PayPal] ❌ Failed to create request: %v", err)
		return "", err
	}

	req.SetBasicAuth(p.config.ClientID, p.config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Printf("[PayPal] ❌ Token request failed (network error): %v", err)
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[PayPal] ❌ Token request failed!")
		log.Printf("[PayPal]    Status: %d", resp.StatusCode)
		log.Printf("[PayPal]    Response: %s", string(body))
		log.Printf("[PayPal] ═══════════════════════════════════════════════════════")
		log.Printf("[PayPal] 💡 POSSIBLE CAUSES:")
		log.Printf("[PayPal]    - Invalid client_id or client_secret")
		log.Printf("[PayPal]    - Credentials are for wrong environment (sandbox vs live)")
		log.Printf("[PayPal]    - App not approved on PayPal Developer Portal")
		log.Printf("[PayPal] 💡 TO FIX:")
		log.Printf("[PayPal]    1. Go to https://developer.paypal.com/dashboard/applications/")
		log.Printf("[PayPal]    2. Check your app credentials for %s mode", p.config.Mode)
		log.Printf("[PayPal]    3. Update in Admin Panel -> Agrégateurs -> PayPal")
		log.Printf("[PayPal] ═══════════════════════════════════════════════════════")
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp PayPalTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		log.Printf("[PayPal] ❌ Failed to decode token response: %v", err)
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	p.accessToken = tokenResp.AccessToken
	p.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // Refresh 60s before expiry

	log.Printf("[PayPal] ✅ Token obtained successfully!")
	log.Printf("[PayPal]    Token type: %s", tokenResp.TokenType)
	log.Printf("[PayPal]    App ID: %s", tokenResp.AppID)
	log.Printf("[PayPal]    Expires in: %d seconds", tokenResp.ExpiresIn)

	return p.accessToken, nil
}

// ==================== ORDERS (COLLECTIONS) ====================

type PayPalOrder struct {
	ID            string               `json:"id"`
	Status        string               `json:"status"`
	Intent        string               `json:"intent"`
	PurchaseUnits []PayPalPurchaseUnit `json:"purchase_units"`
	CreateTime    string               `json:"create_time"`
	Links         []PayPalLink         `json:"links"`
}

type PayPalPurchaseUnit struct {
	ReferenceID string          `json:"reference_id,omitempty"`
	Amount      PayPalAmount    `json:"amount"`
	Description string          `json:"description,omitempty"`
	CustomID    string          `json:"custom_id,omitempty"` // Use for wallet_id
	InvoiceID   string          `json:"invoice_id,omitempty"`
	Payee       *PayPalPayee    `json:"payee,omitempty"`
	Payments    *PayPalPayments `json:"payments,omitempty"`
}

type PayPalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PayPalPayee struct {
	EmailAddress string `json:"email_address,omitempty"`
	MerchantID   string `json:"merchant_id,omitempty"`
}

type PayPalPayments struct {
	Captures []PayPalCapture `json:"captures"`
}

type PayPalCapture struct {
	ID     string       `json:"id"`
	Status string       `json:"status"`
	Amount PayPalAmount `json:"amount"`
}

type PayPalLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type CreateOrderRequest struct {
	Intent             string                    `json:"intent"` // CAPTURE or AUTHORIZE
	PurchaseUnits      []PayPalPurchaseUnit      `json:"purchase_units"`
	ApplicationContext *PayPalApplicationContext `json:"application_context,omitempty"`
}

type PayPalApplicationContext struct {
	BrandName          string `json:"brand_name,omitempty"`
	Locale             string `json:"locale,omitempty"`
	LandingPage        string `json:"landing_page,omitempty"`
	UserAction         string `json:"user_action,omitempty"`
	ReturnURL          string `json:"return_url,omitempty"`
	CancelURL          string `json:"cancel_url,omitempty"`
	ShippingPreference string `json:"shipping_preference,omitempty"`
}

// CreateOrder initiates a PayPal payment (for deposits)
func (p *PayPalProvider) CreateOrder(ctx context.Context, amount float64, currency, walletID, description string) (*PayPalOrder, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	orderReq := CreateOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []PayPalPurchaseUnit{
			{
				Amount: PayPalAmount{
					CurrencyCode: currency,
					Value:        fmt.Sprintf("%.2f", amount),
				},
				Description: description,
				CustomID:    walletID, // Store wallet_id for webhook processing
			},
		},
	}

	body, _ := json.Marshal(orderReq)
	baseURL := p.apiBaseURL()
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v2/checkout/orders", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create order request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create order failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var order PayPalOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	log.Printf("[PayPal] Created order %s for %.2f %s (wallet: %s)", order.ID, amount, currency, walletID)
	return &order, nil
}

// CaptureOrder captures an authorized payment
func (p *PayPalProvider) CaptureOrder(ctx context.Context, orderID string) (*PayPalOrder, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	baseURL := p.apiBaseURL()
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/v2/checkout/orders/%s/capture", baseURL, orderID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("capture order request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("capture order failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var order PayPalOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode capture response: %w", err)
	}

	log.Printf("[PayPal] Captured order %s, status: %s", order.ID, order.Status)
	return &order, nil
}

// GetOrder retrieves order details
func (p *PayPalProvider) GetOrder(ctx context.Context, orderID string) (*PayPalOrder, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	baseURL := p.apiBaseURL()
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/v2/checkout/orders/%s", baseURL, orderID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get order request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get order failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var order PayPalOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &order, nil
}

// ==================== PAYOUTS ====================

type PayPalPayout struct {
	SenderBatchHeader PayPalSenderBatchHeader `json:"sender_batch_header"`
	Items             []PayPalPayoutItem      `json:"items"`
}

type PayPalSenderBatchHeader struct {
	SenderBatchID string `json:"sender_batch_id"`
	EmailSubject  string `json:"email_subject,omitempty"`
	EmailMessage  string `json:"email_message,omitempty"`
}

type PayPalPayoutItem struct {
	RecipientType   string             `json:"recipient_type"` // EMAIL, PHONE, PAYPAL_ID
	Amount          PayPalPayoutAmount `json:"amount"`
	Receiver        string             `json:"receiver"`
	Note            string             `json:"note,omitempty"`
	SenderItemID    string             `json:"sender_item_id"`
	RecipientWallet string             `json:"recipient_wallet,omitempty"` // PAYPAL, VENMO
}

type PayPalPayoutAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type PayPalPayoutResponse struct {
	BatchHeader PayPalBatchHeader `json:"batch_header"`
	Links       []PayPalLink      `json:"links"`
}

type PayPalBatchHeader struct {
	PayoutBatchID     string                  `json:"payout_batch_id"`
	BatchStatus       string                  `json:"batch_status"`
	SenderBatchHeader PayPalSenderBatchHeader `json:"sender_batch_header"`
}

type PayPalWebhookVerification struct {
	AuthAlgo         string `json:"auth_algo"`
	CertURL          string `json:"cert_url"`
	TransmissionID   string `json:"transmission_id"`
	TransmissionSig  string `json:"transmission_sig"`
	TransmissionTime string `json:"transmission_time"`
	WebhookID        string `json:"webhook_id"`
	WebhookEvent     []byte `json:"webhook_event"`
}

type PayPalVerificationResponse struct {
	VerificationStatus string `json:"verification_status"`
}

// CreatePayout sends money to a PayPal account (for withdrawals)
func (p *PayPalProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	// Call internal payout logic
	resp, err := p.createPayoutInternal(ctx, req.RecipientEmail, req.Amount, req.Currency, req.ReferenceID, req.Narration)
	if err != nil {
		return nil, err
	}

	return &PayoutResponse{
		ProviderName:      "paypal",
		ProviderReference: resp.BatchHeader.PayoutBatchID,
		ReferenceID:       req.ReferenceID,
		Status:            PayoutStatusAccordingTo(resp.BatchHeader.BatchStatus),
		AmountReceived:    req.Amount,
		ReceivedCurrency:  req.Currency,
	}, nil
}

// createPayoutInternal sends money to a PayPal account (for withdrawals)
func (p *PayPalProvider) createPayoutInternal(ctx context.Context, recipientEmail string, amount float64, currency, reference, note string) (*PayPalPayoutResponse, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	payout := PayPalPayout{
		SenderBatchHeader: PayPalSenderBatchHeader{
			SenderBatchID: reference,
			EmailSubject:  "You have a payment",
			EmailMessage:  note,
		},
		Items: []PayPalPayoutItem{
			{
				RecipientType: "EMAIL",
				Amount: PayPalPayoutAmount{
					Value:    fmt.Sprintf("%.2f", amount),
					Currency: currency,
				},
				Receiver:     recipientEmail,
				Note:         note,
				SenderItemID: reference + "_1",
			},
		},
	}

	body, _ := json.Marshal(payout)
	baseURL := p.apiBaseURL()
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v1/payments/payouts", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create payout request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create payout failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var payoutResp PayPalPayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&payoutResp); err != nil {
		return nil, fmt.Errorf("failed to decode payout response: %w", err)
	}

	log.Printf("[PayPal] Created payout batch %s for %.2f %s to %s",
		payoutResp.BatchHeader.PayoutBatchID, amount, currency, recipientEmail)
	return &payoutResp, nil
}

// getPayoutStatusInternal retrieves payout batch status (Internal)
func (p *PayPalProvider) getPayoutStatusInternal(ctx context.Context, payoutBatchID string) (*PayPalBatchHeader, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	baseURL := p.apiBaseURL()
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/v1/payments/payouts/%s", baseURL, payoutBatchID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get payout status request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get payout status failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var batchResp struct {
		BatchHeader PayPalBatchHeader `json:"batch_header"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&batchResp); err != nil {
		return nil, fmt.Errorf("failed to decode payout status response: %w", err)
	}

	return &batchResp.BatchHeader, nil
}

// VerifyWebhookSignature validates a PayPal webhook event
func (p *PayPalProvider) VerifyWebhookSignature(ctx context.Context, headers http.Header, body []byte) (bool, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get access token: %w", err)
	}

	verification := PayPalWebhookVerification{
		AuthAlgo:         headers.Get("PAYPAL-AUTH-ALGO"),
		CertURL:          headers.Get("PAYPAL-CERT-URL"),
		TransmissionID:   headers.Get("PAYPAL-TRANSMISSION-ID"),
		TransmissionSig:  headers.Get("PAYPAL-TRANSMISSION-SIG"),
		TransmissionTime: headers.Get("PAYPAL-TRANSMISSION-TIME"),
		WebhookID:        p.config.WebhookID,
		WebhookEvent:     body,
	}

	reqBody, _ := json.Marshal(verification)
	baseURL := p.apiBaseURL()
	req, err := http.NewRequestWithContext(ctx, "POST",
		baseURL+"/v1/notifications/verify-webhook-signature", bytes.NewBuffer(reqBody))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("webhook verification request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("webhook verification failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var verifyResp PayPalVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return false, fmt.Errorf("failed to decode verification response: %w", err)
	}

	return verifyResp.VerificationStatus == "SUCCESS", nil
}

func (p *PayPalProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("PayPal does not support payout cancellation")
}

// ==================== AVAILABLE METHODS ====================

func (p *PayPalProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return []AvailableMethod{
		{Code: "paypal_wallet", Name: "PayPal Wallet", Type: "wallet"},
		{Code: "paypal_card", Name: "Debit or Credit Card", Type: "card"},
		{Code: "paypal_paylater", Name: "Pay Later", Type: "credit"},
	}, nil
}

func (p *PayPalProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	return []Bank{}, nil
}

func (p *PayPalProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	return []MobileOperator{}, nil
}

func (p *PayPalProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.RecipientEmail == "" {
		return fmt.Errorf("recipient email is required for PayPal")
	}
	return nil
}

func (p *PayPalProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	// Simple quote simulation
	return &PayoutResponse{
		ProviderName:     "paypal",
		Status:           PayoutStatusPending,
		Fee:              0,
		AmountReceived:   req.Amount,
		ReceivedCurrency: req.Currency,
	}, nil
}
