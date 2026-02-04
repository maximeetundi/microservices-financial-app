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

// CollectionMethod represents how user deposits money
type CollectionMethod string

const (
	CollectionMethodCard         CollectionMethod = "card"
	CollectionMethodBankTransfer CollectionMethod = "bank_transfer"
	CollectionMethodMobileMoney  CollectionMethod = "mobile_money"
	CollectionMethodUSSD         CollectionMethod = "ussd"
)

// CollectionStatus represents the status of a deposit
type CollectionStatus string

const (
	CollectionStatusPending    CollectionStatus = "pending"
	CollectionStatusSuccessful CollectionStatus = "successful"
	CollectionStatusFailed     CollectionStatus = "failed"
	CollectionStatusCancelled  CollectionStatus = "cancelled"
)

// CollectionRequest represents a deposit/collection request
type CollectionRequest struct {
	ReferenceID string  `json:"reference_id"`
	UserID      string  `json:"user_id"`
	WalletID    string  `json:"wallet_id,omitempty"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`

	// Collection method
	Method CollectionMethod `json:"method"`

	// For card payments
	CardNumber     string `json:"card_number,omitempty"`
	CardExpiry     string `json:"card_expiry,omitempty"`
	CardCVV        string `json:"card_cvv,omitempty"`
	CardHolderName string `json:"card_holder_name,omitempty"`

	// For mobile money
	MobileNumber   string `json:"mobile_number,omitempty"`
	MobileOperator string `json:"mobile_operator,omitempty"`

	// For bank transfer
	BankCode string `json:"bank_code,omitempty"`

	// User info
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Country     string `json:"country"`

	// Redirect URL after payment
	RedirectURL string `json:"redirect_url,omitempty"`

	// Metadata
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CollectionResponse represents the result of a collection
type CollectionResponse struct {
	ReferenceID       string           `json:"reference_id"`
	ProviderReference string           `json:"provider_reference"`
	Status            CollectionStatus `json:"status"`
	Amount            float64          `json:"amount"`
	Currency          string           `json:"currency"`
	Fee               float64          `json:"fee"`
	NetAmount         float64          `json:"net_amount"`

	// For redirect-based payments (card, bank)
	PaymentLink string `json:"payment_link,omitempty"`

	// For USSD/Mobile Money
	USSDCode string `json:"ussd_code,omitempty"`

	Message string `json:"message,omitempty"`

	// Metadata from provider (includes wallet_id, user_id)
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CollectionProvider interface for deposit/collection providers
type CollectionProvider interface {
	GetName() string
	GetSupportedCountries() []string
	GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error)
	InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error)
	VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error)
}

// ========================
// Flutterwave Collection
// ========================

type FlutterwaveCollectionProvider struct {
	config     FlutterwaveConfig
	httpClient *http.Client
}

func NewFlutterwaveCollectionProvider(config FlutterwaveConfig) *FlutterwaveCollectionProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.flutterwave.com/v3"
	}
	return &FlutterwaveCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (f *FlutterwaveCollectionProvider) GetName() string {
	return "flutterwave"
}

func (f *FlutterwaveCollectionProvider) GetSupportedCountries() []string {
	return []string{"NG", "GH", "KE", "UG", "TZ", "ZA", "CI", "SN", "CM"}
}

func (f *FlutterwaveCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	methods := []CollectionMethod{CollectionMethodCard}

	mobileMoneyCountries := map[string]bool{
		"GH": true, "KE": true, "UG": true, "TZ": true, "RW": true,
		"CI": true, "SN": true, "CM": true,
	}

	if mobileMoneyCountries[country] {
		methods = append(methods, CollectionMethodMobileMoney)
	}

	if country == "NG" {
		methods = append(methods, CollectionMethodBankTransfer, CollectionMethodUSSD)
	}

	return methods, nil
}

func (f *FlutterwaveCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	if req.Metadata == nil {
		req.Metadata = make(map[string]string)
	}
	if req.WalletID != "" {
		req.Metadata["wallet_id"] = req.WalletID
	}

	payload := map[string]interface{}{
		"tx_ref":       req.ReferenceID,
		"amount":       req.Amount,
		"currency":     req.Currency,
		"redirect_url": req.RedirectURL,
		"customer": map[string]string{
			"email":       req.Email,
			"phonenumber": req.PhoneNumber,
			"name":        req.CardHolderName,
		},
		"customizations": map[string]string{
			"title":       "Zekora Deposit",
			"description": "Deposit funds to your wallet",
		},
		"meta": req.Metadata,
	}

	// Add payment type based on method
	switch req.Method {
	case CollectionMethodMobileMoney:
		payload["payment_options"] = "mobilemoney"
	case CollectionMethodBankTransfer:
		payload["payment_options"] = "banktransfer"
	case CollectionMethodUSSD:
		payload["payment_options"] = "ussd"
	default:
		payload["payment_options"] = "card"
	}

	body, _ := json.Marshal(payload)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", f.config.BaseURL+"/payments", bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+f.config.SecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Link string `json:"link"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("flutterwave error: %s", result.Message)
	}

	return &CollectionResponse{
		ReferenceID: req.ReferenceID,
		Status:      CollectionStatusPending,
		Amount:      req.Amount,
		Currency:    req.Currency,
		PaymentLink: result.Data.Link,
	}, nil
}

func (f *FlutterwaveCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	url := fmt.Sprintf("%s/transactions/verify_by_reference?tx_ref=%s", f.config.BaseURL, referenceID)

	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	httpReq.Header.Set("Authorization", "Bearer "+f.config.SecretKey)

	resp, err := f.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
		Data   struct {
			ID            int         `json:"id"`
			TxRef         string      `json:"tx_ref"`
			Amount        float64     `json:"amount"`
			Currency      string      `json:"currency"`
			ChargedAmount float64     `json:"charged_amount"`
			AppFee        float64     `json:"app_fee"`
			Status        string      `json:"status"`
			Meta          interface{} `json:"meta"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var status CollectionStatus
	switch result.Data.Status {
	case "successful":
		status = CollectionStatusSuccessful
	case "failed":
		status = CollectionStatusFailed
	default:
		status = CollectionStatusPending
	}

	var metaMap map[string]interface{}
	if m, ok := result.Data.Meta.(map[string]interface{}); ok {
		metaMap = m
	}

	return &CollectionResponse{
		ReferenceID:       referenceID,
		ProviderReference: fmt.Sprintf("%d", result.Data.ID),
		Status:            status,
		Amount:            result.Data.Amount,
		Currency:          result.Data.Currency,
		Fee:               result.Data.AppFee,
		NetAmount:         result.Data.ChargedAmount - result.Data.AppFee,
		Metadata:          convertMeta(metaMap),
	}, nil
}

// ========================
// Stripe Collection (Cards, SEPA)
// ========================

type StripeCollectionProvider struct {
	config     StripeConfig
	httpClient *http.Client
}

func NewStripeCollectionProvider(config StripeConfig) *StripeCollectionProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.stripe.com/v1"
	}
	return &StripeCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *StripeCollectionProvider) GetName() string {
	return "stripe"
}

func (s *StripeCollectionProvider) GetSupportedCountries() []string {
	return []string{"US", "CA", "GB", "DE", "FR", "IT", "ES", "NL", "BE", "AT", "CH"}
}

func (s *StripeCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	methods := []CollectionMethod{CollectionMethodCard}

	// SEPA countries
	sepaCountries := map[string]bool{
		"DE": true, "FR": true, "IT": true, "ES": true, "NL": true,
		"BE": true, "AT": true, "PT": true, "IE": true, "FI": true,
	}

	if sepaCountries[country] {
		methods = append(methods, CollectionMethodBankTransfer)
	}

	return methods, nil
}

func (s *StripeCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// Create Stripe Checkout Session
	payload := fmt.Sprintf(
		"mode=payment&success_url=%s&cancel_url=%s&line_items[0][price_data][currency]=%s&line_items[0][price_data][product_data][name]=Wallet+Deposit&line_items[0][price_data][unit_amount]=%d&line_items[0][quantity]=1&metadata[reference_id]=%s&metadata[user_id]=%s&metadata[wallet_id]=%s",
		req.RedirectURL,
		req.RedirectURL,
		req.Currency,
		int(req.Amount*100), // Stripe uses cents
		req.ReferenceID,
		req.UserID,
	)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL+"/checkout/sessions", bytes.NewBufferString(payload))
	httpReq.SetBasicAuth(s.config.SecretKey, "")
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: result.ID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       result.URL,
	}, nil
}

func (s *StripeCollectionProvider) VerifyCollection(ctx context.Context, sessionID string) (*CollectionResponse, error) {
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", s.config.BaseURL+"/checkout/sessions/"+sessionID, nil)
	httpReq.SetBasicAuth(s.config.SecretKey, "")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		ID            string            `json:"id"`
		PaymentStatus string            `json:"payment_status"`
		AmountTotal   int               `json:"amount_total"`
		Currency      string            `json:"currency"`
		Metadata      map[string]string `json:"metadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var status CollectionStatus
	switch result.PaymentStatus {
	case "paid":
		status = CollectionStatusSuccessful
	case "unpaid", "no_payment_required":
		status = CollectionStatusPending
	default:
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       result.Metadata["reference_id"],
		ProviderReference: result.ID,
		Status:            status,
		Amount:            float64(result.AmountTotal) / 100,
		Currency:          result.Currency,
		Metadata:          result.Metadata,
	}, nil
}

// convertMeta converts map[string]interface{} to map[string]string
func convertMeta(input map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range input {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}

// ========================
// CinetPay Collection Wrapper
// ========================

type CinetPayCollectionProvider struct {
	provider *CinetPayProvider
}

func NewCinetPayCollectionProvider(config CinetPayConfig) *CinetPayCollectionProvider {
	return &CinetPayCollectionProvider{
		provider: NewCinetPayProvider(config),
	}
}

func (c *CinetPayCollectionProvider) GetName() string {
	return c.provider.GetName()
}

func (c *CinetPayCollectionProvider) GetSupportedCountries() []string {
	return c.provider.GetSupportedCountries()
}

func (c *CinetPayCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{CollectionMethodMobileMoney, CollectionMethodCard}, nil
}

func (c *CinetPayCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// CinetPay collection not fully implemented - placeholder
	return &CollectionResponse{
		ReferenceID: req.ReferenceID,
		Status:      CollectionStatusPending,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Message:     "CinetPay collection initiated",
	}, nil
}

func (c *CinetPayCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	return &CollectionResponse{
		ReferenceID: referenceID,
		Status:      CollectionStatusPending,
	}, nil
}

// ========================
// Wave Collection Wrapper
// ========================

type WaveCollectionProvider struct {
	provider *WaveProvider
}

func NewWaveCollectionProvider(config WaveConfig) *WaveCollectionProvider {
	return &WaveCollectionProvider{
		provider: NewWaveProvider(config),
	}
}

func (w *WaveCollectionProvider) GetName() string {
	return w.provider.GetName()
}

func (w *WaveCollectionProvider) GetSupportedCountries() []string {
	return w.provider.GetSupportedCountries()
}

func (w *WaveCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{CollectionMethodMobileMoney}, nil
}

func (w *WaveCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	return &CollectionResponse{
		ReferenceID: req.ReferenceID,
		Status:      CollectionStatusPending,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Message:     "Wave collection initiated",
	}, nil
}

func (w *WaveCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	return &CollectionResponse{
		ReferenceID: referenceID,
		Status:      CollectionStatusPending,
	}, nil
}

// ========================
// PayPal Collection Wrapper
// ========================

type PayPalCollectionProvider struct {
	provider *PayPalProvider
}

func NewPayPalCollectionProvider(config PayPalConfig) *PayPalCollectionProvider {
	return &PayPalCollectionProvider{
		provider: NewPayPalProvider(config),
	}
}

func (p *PayPalCollectionProvider) GetName() string {
	return p.provider.GetName()
}

func (p *PayPalCollectionProvider) GetSupportedCountries() []string {
	return p.provider.GetSupportedCountries()
}

func (p *PayPalCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{CollectionMethodCard}, nil
}

func (p *PayPalCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	order, err := p.provider.CreateOrder(ctx, req.Amount, req.Currency, req.WalletID, "Wallet Deposit")
	if err != nil {
		return nil, err
	}

	// Find approval URL
	var paymentLink string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			paymentLink = link.Href
			break
		}
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: order.ID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       paymentLink,
		Message:           "Redirecting to PayPal",
	}, nil
}

func (p *PayPalCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	order, err := p.provider.GetOrder(ctx, referenceID)
	if err != nil {
		return nil, err
	}

	status := CollectionStatusPending
	if order.Status == "COMPLETED" {
		status = CollectionStatusSuccessful
	} else if order.Status == "VOIDED" || order.Status == "CANCELLED" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       referenceID,
		ProviderReference: order.ID,
		Status:            status,
	}, nil
}

// ========================
// Orange Money Collection Wrapper
// ========================

type OrangeMoneyCollectionProvider struct {
	provider *OrangeMoneyProvider
}

func NewOrangeMoneyCollectionProvider(config OrangeMoneyConfig) *OrangeMoneyCollectionProvider {
	return &OrangeMoneyCollectionProvider{
		provider: NewOrangeMoneyProvider(config),
	}
}

func (o *OrangeMoneyCollectionProvider) GetName() string {
	return o.provider.GetName()
}

func (o *OrangeMoneyCollectionProvider) GetSupportedCountries() []string {
	return o.provider.GetSupportedCountries()
}

func (o *OrangeMoneyCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{CollectionMethodMobileMoney}, nil
}

func (o *OrangeMoneyCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	resp, err := o.provider.InitiateWebPayment(ctx, req.Amount, req.Currency, req.ReferenceID, req.RedirectURL, req.RedirectURL)
	if err != nil {
		return nil, err
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: resp.PayToken,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       resp.PaymentURL,
		Message:           "Redirecting to Orange Money",
	}, nil
}

func (o *OrangeMoneyCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	return &CollectionResponse{
		ReferenceID: referenceID,
		Status:      CollectionStatusPending,
	}, nil
}

// ========================
// MTN MoMo Collection Wrapper
// ========================

type MTNMoMoCollectionProvider struct {
	provider *MTNMomoProvider
}

func NewMTNMoMoCollectionProvider(config MTNMomoConfig) *MTNMoMoCollectionProvider {
	return &MTNMoMoCollectionProvider{
		provider: NewMTNMomoProvider(config),
	}
}

func (m *MTNMoMoCollectionProvider) GetName() string {
	return m.provider.GetName()
}

func (m *MTNMoMoCollectionProvider) GetSupportedCountries() []string {
	return m.provider.GetSupportedCountries()
}

func (m *MTNMoMoCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{CollectionMethodMobileMoney}, nil
}

func (m *MTNMoMoCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	refID, err := m.provider.RequestToPay(ctx, req.MobileNumber, req.Amount, req.Currency, req.ReferenceID, "Wallet Deposit")
	if err != nil {
		return nil, err
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: refID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Message:           "MTN MoMo payment request sent",
	}, nil
}

func (m *MTNMoMoCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	status, err := m.provider.GetRequestToPayStatus(ctx, referenceID)
	if err != nil {
		return nil, err
	}

	collStatus := CollectionStatusPending
	switch status.Status {
	case "SUCCESSFUL":
		collStatus = CollectionStatusSuccessful
	case "FAILED":
		collStatus = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       referenceID,
		ProviderReference: status.FinancialTransactionID,
		Status:            collStatus,
	}, nil
}
