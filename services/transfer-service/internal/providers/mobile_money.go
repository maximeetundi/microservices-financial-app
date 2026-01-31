package providers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ==================== MTN MOBILE MONEY PROVIDER ====================

// MTNMomoConfig holds MTN Mobile Money API configuration
type MTNMomoConfig struct {
	Name            string // Internal name (e.g., mtn_momo, mtn_ci)
	SubscriptionKey string // Ocp-Apim-Subscription-Key
	APIUser         string
	APIKey          string
	BaseURL         string // https://sandbox.momodeveloper.mtn.com or https://momodeveloper.mtn.com
	Environment     string // sandbox or production
	CallbackURL     string
}

// MTNMomoProvider implements MTN Mobile Money collections and disbursements
type MTNMomoProvider struct {
	config      MTNMomoConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// NewMTNMomoProvider creates a new MTN MoMo provider
func NewMTNMomoProvider(config MTNMomoConfig) *MTNMomoProvider {
	if config.BaseURL == "" {
		if config.Environment == "production" {
			config.BaseURL = "https://momodeveloper.mtn.com"
		} else {
			config.BaseURL = "https://sandbox.momodeveloper.mtn.com"
		}
	}
	return &MTNMomoProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *MTNMomoProvider) GetName() string {
	if m.config.Name != "" {
		return m.config.Name
	}
	return "mtn_momo"
}

func (m *MTNMomoProvider) GetSupportedCountries() []string {
	return []string{
		"BJ", "CI", "CM", "CD", "GH", "GN", "LR", "RW", "UG", "ZA", "ZM",
	}
}

// MTN MoMo token response
type MTNTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken retrieves OAuth token for MTN MoMo
func (m *MTNMomoProvider) GetAccessToken(ctx context.Context, product string) (string, error) {
	if m.accessToken != "" && time.Now().Before(m.tokenExpiry) {
		return m.accessToken, nil
	}

	url := fmt.Sprintf("%s/%s/token/", m.config.BaseURL, product)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", err
	}

	// Basic auth with API User and API Key
	auth := base64.StdEncoding.EncodeToString([]byte(m.config.APIUser + ":" + m.config.APIKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.SubscriptionKey)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp MTNTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	m.accessToken = tokenResp.AccessToken
	m.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	return m.accessToken, nil
}

// MTN MoMo collection request
type MTNCollectionRequest struct {
	Amount       string   `json:"amount"`
	Currency     string   `json:"currency"`
	ExternalID   string   `json:"externalId"`
	Payer        MTNPayer `json:"payer"`
	PayerMessage string   `json:"payerMessage"`
	PayeeNote    string   `json:"payeeNote"`
}

type MTNPayer struct {
	PartyIDType string `json:"partyIdType"` // MSISDN
	PartyID     string `json:"partyId"`     // Phone number
}

// RequestToPay initiates a collection request (user deposits money)
func (m *MTNMomoProvider) RequestToPay(ctx context.Context, phoneNumber string, amount float64, currency, externalID, message string) (string, error) {
	token, err := m.GetAccessToken(ctx, "collection")
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	referenceID := uuid.New().String()

	collReq := MTNCollectionRequest{
		Amount:     fmt.Sprintf("%.0f", amount),
		Currency:   currency,
		ExternalID: externalID,
		Payer: MTNPayer{
			PartyIDType: "MSISDN",
			PartyID:     phoneNumber,
		},
		PayerMessage: message,
		PayeeNote:    "Deposit to wallet",
	}

	body, _ := json.Marshal(collReq)
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", m.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("X-Target-Environment", m.config.Environment)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.SubscriptionKey)
	req.Header.Set("Content-Type", "application/json")
	if m.config.CallbackURL != "" {
		req.Header.Set("X-Callback-Url", m.config.CallbackURL)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to pay failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request to pay failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[MTN-MoMo] Request to pay initiated: %s for %.2f %s from %s", referenceID, amount, currency, phoneNumber)
	return referenceID, nil
}

// MTN MoMo transaction status
type MTNTransactionStatus struct {
	Amount                 string   `json:"amount"`
	Currency               string   `json:"currency"`
	ExternalID             string   `json:"externalId"`
	FinancialTransactionID string   `json:"financialTransactionId"`
	Payer                  MTNPayer `json:"payer"`
	Status                 string   `json:"status"` // SUCCESSFUL, FAILED, PENDING
	Reason                 string   `json:"reason,omitempty"`
}

// GetRequestToPayStatus checks the status of a collection request
func (m *MTNMomoProvider) GetRequestToPayStatus(ctx context.Context, referenceID string) (*MTNTransactionStatus, error) {
	token, err := m.GetAccessToken(ctx, "collection")
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/collection/v1_0/requesttopay/%s", m.config.BaseURL, referenceID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Target-Environment", m.config.Environment)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.SubscriptionKey)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get status request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get status failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var status MTNTransactionStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return &status, nil
}

// MTN Transfer request for disbursements
type MTNTransferRequest struct {
	Amount       string   `json:"amount"`
	Currency     string   `json:"currency"`
	ExternalID   string   `json:"externalId"`
	Payee        MTNPayer `json:"payee"`
	PayerMessage string   `json:"payerMessage"`
	PayeeNote    string   `json:"payeeNote"`
}

// Transfer sends money to a phone number (for withdrawals)
func (m *MTNMomoProvider) Transfer(ctx context.Context, phoneNumber string, amount float64, currency, externalID, message string) (string, error) {
	token, err := m.GetAccessToken(ctx, "disbursement")
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	referenceID := uuid.New().String()

	transferReq := MTNTransferRequest{
		Amount:     fmt.Sprintf("%.0f", amount),
		Currency:   currency,
		ExternalID: externalID,
		Payee: MTNPayer{
			PartyIDType: "MSISDN",
			PartyID:     phoneNumber,
		},
		PayerMessage: message,
		PayeeNote:    "Withdrawal from wallet",
	}

	body, _ := json.Marshal(transferReq)
	url := fmt.Sprintf("%s/disbursement/v1_0/transfer", m.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("X-Target-Environment", m.config.Environment)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.SubscriptionKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("transfer request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("transfer failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[MTN-MoMo] Transfer initiated: %s for %.2f %s to %s", referenceID, amount, currency, phoneNumber)
	return referenceID, nil
}

// CreatePayout implements PayoutProvider interface
func (m *MTNMomoProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	refID, err := m.Transfer(ctx, req.RecipientPhone, req.Amount, req.Currency, req.ReferenceID, req.Narration)
	if err != nil {
		return nil, err
	}

	return &PayoutResponse{
		ProviderName:      m.GetName(),
		ProviderReference: refID,
		ReferenceID:       req.ReferenceID,
		Status:            PayoutStatusPending, // Async
		AmountReceived:    req.Amount,
		ReceivedCurrency:  req.Currency,
		Fee:               0,
	}, nil
}

// CancelPayout cancels a pending payout
func (m *MTNMomoProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("MTN MoMo cancellation not supported")
}

func (m *MTNMomoProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return []AvailableMethod{
		{Code: "mtn_momo", Name: "MTN Mobile Money", Type: "mobile"},
	}, nil
}

func (m *MTNMomoProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	return []Bank{}, nil
}

func (m *MTNMomoProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	return []MobileOperator{
		{Code: "mtn", Name: "MTN", Countries: m.GetSupportedCountries()},
	}, nil
}

func (m *MTNMomoProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.RecipientPhone == "" {
		return fmt.Errorf("recipient phone is required for MTN MoMo")
	}
	return nil
}

func (m *MTNMomoProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	return &PayoutResponse{
		ProviderName: m.GetName(),
		Status:       PayoutStatusPending,
		Fee:          0, // Calculate based on amount if needed
		TotalAmount:  req.Amount,
	}, nil
}

// ==================== ORANGE MONEY PROVIDER ====================

// OrangeMoneyConfig holds Orange Money API configuration
type OrangeMoneyConfig struct {
	Name         string // Internal name (e.g., orange_money, orange_money_ci)
	ClientID     string
	ClientSecret string
	MerchantKey  string
	BaseURL      string // https://api.orange.com/orange-money-webpay/dev/v1
	Country      string // CI, SN, CM, etc.
}

// OrangeMoneyProvider implements Orange Money payments
type OrangeMoneyProvider struct {
	config      OrangeMoneyConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// NewOrangeMoneyProvider creates a new Orange Money provider
func NewOrangeMoneyProvider(config OrangeMoneyConfig) *OrangeMoneyProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.orange.com/orange-money-webpay/dev/v1"
	}
	return &OrangeMoneyProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (o *OrangeMoneyProvider) GetName() string {
	if o.config.Name != "" {
		return o.config.Name
	}
	return "orange_money"
}

func (o *OrangeMoneyProvider) GetSupportedCountries() []string {
	return []string{
		"CI", "SN", "CM", "ML", "GN", "BF", "NE", "MG", "CD", "CG", "GA", "BI", "LR", "CF",
	}
}

// Orange OAuth token response
type OrangeTokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // seconds, usually 3600
}

// GetAccessToken retrieves OAuth token for Orange Money
func (o *OrangeMoneyProvider) GetAccessToken(ctx context.Context) (string, error) {
	if o.accessToken != "" && time.Now().Before(o.tokenExpiry) {
		return o.accessToken, nil
	}

	url := "https://api.orange.com/oauth/v3/token"
	body := "grant_type=client_credentials"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(o.config.ClientID + ":" + o.config.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var tokenResp OrangeTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	o.accessToken = tokenResp.AccessToken
	o.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	return o.accessToken, nil
}

// Orange Money web payment request
type OrangeWebPaymentRequest struct {
	MerchantKey string `json:"merchant_key"`
	Currency    string `json:"currency"` // XOF, XAF
	OrderID     string `json:"order_id"`
	Amount      int    `json:"amount"` // Amount in integer (no decimals for XOF)
	ReturnURL   string `json:"return_url"`
	CancelURL   string `json:"cancel_url"`
	NotifURL    string `json:"notif_url"`
	Lang        string `json:"lang,omitempty"` // fr or en
}

// Orange Money payment response
type OrangePaymentResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	PayToken   string `json:"pay_token"`
	PaymentURL string `json:"payment_url"`
	NotifToken string `json:"notif_token"`
}

// InitiateWebPayment creates a web payment session for deposits
func (o *OrangeMoneyProvider) InitiateWebPayment(ctx context.Context, amount float64, currency, orderID, returnURL, notifURL string) (*OrangePaymentResponse, error) {
	token, err := o.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	payReq := OrangeWebPaymentRequest{
		MerchantKey: o.config.MerchantKey,
		Currency:    currency, // XOF or XAF
		OrderID:     orderID,
		Amount:      int(amount), // No decimals for XOF/XAF
		ReturnURL:   returnURL,
		CancelURL:   returnURL,
		NotifURL:    notifURL,
		Lang:        "fr",
	}

	body, _ := json.Marshal(payReq)
	url := o.config.BaseURL + "/webpayment"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("payment request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("payment request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var payResp OrangePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&payResp); err != nil {
		return nil, fmt.Errorf("failed to decode payment response: %w", err)
	}

	log.Printf("[OrangeMoney] Payment initiated: %s for %d %s, PayToken: %s", orderID, int(amount), currency, payResp.PayToken)
	return &payResp, nil
}

// CreatePayout implements PayoutProvider interface
func (o *OrangeMoneyProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	// Payouts are not directly supported via WebPayment API (only collections)
	// We return an error or implement if another API endpoint exists
	return nil, fmt.Errorf("Orange Money payouts not supported via this API")
}

func (o *OrangeMoneyProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	// Orange Money Payouts are not supported, so status check is also not relevant for payouts
	return nil, fmt.Errorf("Orange Money payouts not supported")
}

// Orange Money transaction status
type OrangeTransactionStatus struct {
	Status   string `json:"status"`
	OrderID  string `json:"order_id"`
	Amount   int    `json:"amount"`
	PayToken string `json:"pay_token"`
	TxnID    string `json:"txnid"`
	Message  string `json:"message"`
}

// GetTransactionStatus checks payment status
func (o *OrangeMoneyProvider) GetTransactionStatus(ctx context.Context, orderID, amount, payToken string) (*OrangeTransactionStatus, error) {
	token, err := o.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/transactionstatus?order_id=%s&amount=%s&pay_token=%s", o.config.BaseURL, orderID, amount, payToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("status request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var status OrangeTransactionStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return &status, nil
}

// CancelPayout cancels a pending payout
func (o *OrangeMoneyProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("Orange Money cancellation not supported")
}

func (o *OrangeMoneyProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return []AvailableMethod{
		{Code: "orange_money", Name: "Orange Money", Type: "mobile"},
	}, nil
}

func (o *OrangeMoneyProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	return []Bank{}, nil
}

func (o *OrangeMoneyProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	return []MobileOperator{
		{Code: "orange", Name: "Orange", Countries: o.GetSupportedCountries()},
	}, nil
}

func (o *OrangeMoneyProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.RecipientPhone == "" {
		return fmt.Errorf("recipient phone is required for Orange Money")
	}
	return nil
}

func (o *OrangeMoneyProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	return &PayoutResponse{
		ProviderName: o.GetName(),
		Status:       PayoutStatusPending,
		Fee:          0,
		TotalAmount:  req.Amount,
	}, nil
}

// ==================== PESAPAL PROVIDER ====================

// PesapalConfig holds Pesapal API configuration
type PesapalConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	BaseURL        string // https://pay.pesapal.com/v3 or https://cybqa.pesapal.com/pesapalv3
}

// PesapalProvider implements Pesapal payments for East Africa
type PesapalProvider struct {
	config      PesapalConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// NewPesapalProvider creates a new Pesapal provider
func NewPesapalProvider(config PesapalConfig) *PesapalProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://cybqa.pesapal.com/pesapalv3" // Sandbox by default
	}
	return &PesapalProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *PesapalProvider) GetName() string {
	return "pesapal"
}

func (p *PesapalProvider) GetSupportedCountries() []string {
	return []string{
		"KE", "UG", "TZ", "RW", "ZM", "MW", "ZW",
	}
}

// Pesapal auth response
type PesapalAuthResponse struct {
	Token      string `json:"token"`
	ExpiryDate string `json:"expiryDate"`
	Error      string `json:"error,omitempty"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

// GetAccessToken retrieves auth token for Pesapal
func (p *PesapalProvider) GetAccessToken(ctx context.Context) (string, error) {
	if p.accessToken != "" && time.Now().Before(p.tokenExpiry) {
		return p.accessToken, nil
	}

	authReq := map[string]string{
		"consumer_key":    p.config.ConsumerKey,
		"consumer_secret": p.config.ConsumerSecret,
	}

	body, _ := json.Marshal(authReq)
	url := p.config.BaseURL + "/api/Auth/RequestToken"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var authResp PesapalAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	if authResp.Error != "" {
		return "", fmt.Errorf("auth error: %s", authResp.Error)
	}

	p.accessToken = authResp.Token
	// Parse expiry date and set tokenExpiry
	if expiry, err := time.Parse("2006-01-02T15:04:05.999999999", authResp.ExpiryDate); err == nil {
		p.tokenExpiry = expiry.Add(-60 * time.Second)
	} else {
		p.tokenExpiry = time.Now().Add(4 * time.Minute) // Default 5 minutes
	}

	return p.accessToken, nil
}

// Pesapal order request
type PesapalOrderRequest struct {
	ID             string                `json:"id"`
	Currency       string                `json:"currency"`
	Amount         float64               `json:"amount"`
	Description    string                `json:"description"`
	CallbackURL    string                `json:"callback_url"`
	NotificationID string                `json:"notification_id"`
	BillingAddress PesapalBillingAddress `json:"billing_address"`
}

type PesapalBillingAddress struct {
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

// Pesapal order response
type PesapalOrderResponse struct {
	OrderTrackingID string `json:"order_tracking_id"`
	MerchantRef     string `json:"merchant_reference"`
	RedirectURL     string `json:"redirect_url"`
	Error           string `json:"error,omitempty"`
	Status          string `json:"status"`
}

// SubmitOrderRequest creates a payment order
func (p *PesapalProvider) SubmitOrderRequest(ctx context.Context, orderID string, amount float64, currency, description, callbackURL, notificationID, email string) (*PesapalOrderResponse, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	orderReq := PesapalOrderRequest{
		ID:             orderID,
		Currency:       currency,
		Amount:         amount,
		Description:    description,
		CallbackURL:    callbackURL,
		NotificationID: notificationID,
		BillingAddress: PesapalBillingAddress{
			EmailAddress: email,
		},
	}

	body, _ := json.Marshal(orderReq)
	url := p.config.BaseURL + "/api/Transactions/SubmitOrderRequest"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("submit order request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("submit order failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var orderResp PesapalOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	log.Printf("[Pesapal] Order created: %s, TrackingID: %s", orderID, orderResp.OrderTrackingID)
	return &orderResp, nil
}

// Pesapal transaction status
type PesapalTransactionStatus struct {
	PaymentMethod            string  `json:"payment_method"`
	Amount                   float64 `json:"amount"`
	CreatedDate              string  `json:"created_date"`
	ConfirmationCode         string  `json:"confirmation_code"`
	PaymentStatusDescription string  `json:"payment_status_description"`
	Description              string  `json:"description"`
	Message                  string  `json:"message"`
	PaymentAccount           string  `json:"payment_account"`
	Currency                 string  `json:"currency"`
	Status                   int     `json:"status_code"`
	MerchantReference        string  `json:"merchant_reference"`
}

// GetTransactionStatus retrieves transaction status
func (p *PesapalProvider) GetTransactionStatus(ctx context.Context, orderTrackingID string) (*PesapalTransactionStatus, error) {
	token, err := p.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/api/Transactions/GetTransactionStatus?orderTrackingId=%s", p.config.BaseURL, orderTrackingID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get status request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get status failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var status PesapalTransactionStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return &status, nil
}

// CreatePayout implements PayoutProvider interface
func (p *PesapalProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	return nil, fmt.Errorf("Pesapal payouts not supported")
}

func (p *PesapalProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	return nil, fmt.Errorf("Pesapal payouts not supported")
}

// CancelPayout cancels a pending payout
func (p *PesapalProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("Pesapal cancellation not supported")
}

func (p *PesapalProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return []AvailableMethod{
		{Code: "pesapal", Name: "Pesapal", Type: "mobile"},
	}, nil
}

func (p *PesapalProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	return []Bank{}, nil
}

func (p *PesapalProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	return []MobileOperator{}, nil
}

func (p *PesapalProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	if req.RecipientEmail == "" && req.RecipientPhone == "" {
		return fmt.Errorf("recipient email or phone is required for Pesapal")
	}
	return nil
}

func (p *PesapalProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	return &PayoutResponse{
		ProviderName: p.GetName(),
		Status:       PayoutStatusPending,
		Fee:          0,
		TotalAmount:  req.Amount,
	}, nil
}

// ==================== CHIPPER CASH PROVIDER ====================

// ChipperConfig holds Chipper Cash API configuration
type ChipperConfig struct {
	APIKey    string
	APISecret string
	BaseURL   string
}

// ChipperProvider implements Chipper Cash for African remittances
type ChipperProvider struct {
	config     ChipperConfig
	httpClient *http.Client
}

// NewChipperProvider creates a new Chipper Cash provider
func NewChipperProvider(config ChipperConfig) *ChipperProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.sandbox.chipper.cash/v1"
	}
	return &ChipperProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ChipperProvider) GetName() string {
	return "chipper"
}

func (c *ChipperProvider) GetSupportedCountries() []string {
	return []string{
		"GH", "KE", "NG", "RW", "TZ", "UG", "ZA", "US", "GB",
	}
}

func (c *ChipperProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return []AvailableMethod{
		{Code: "chipper_wallet", Name: "Chipper Network", Type: "mobile"},
	}, nil
}

func (c *ChipperProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	return []Bank{}, nil
}

func (c *ChipperProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	return []MobileOperator{}, nil
}

func (c *ChipperProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	return nil
}

func (c *ChipperProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	return &PayoutResponse{
		ProviderName: "chipper",
		Status:       "quote",
		Fee:          0,
		TotalAmount:  req.Amount,
	}, nil
}

func (c *ChipperProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	return nil, fmt.Errorf("Chipper Cash integration not yet implemented")
}

func (c *ChipperProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	return nil, fmt.Errorf("Chipper Cash integration not yet implemented")
}

func (c *ChipperProvider) CancelPayout(ctx context.Context, referenceID string) error {
	return fmt.Errorf("Chipper Cash cancellation not supported")
}
