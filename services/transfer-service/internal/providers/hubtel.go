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

// HubtelConfig holds configuration for Hubtel API
type HubtelConfig struct {
	ClientID        string
	ClientSecret    string
	MerchantAccount string
	BaseURL         string // https://api.hubtel.com
}

// HubtelProvider implements CollectionProvider for Hubtel
type HubtelProvider struct {
	config     HubtelConfig
	httpClient *http.Client
}

// NewHubtelCollectionProvider creates a new Hubtel collection provider
func NewHubtelCollectionProvider(config HubtelConfig) *HubtelProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.hubtel.com"
	}
	return &HubtelProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *HubtelProvider) GetName() string {
	return "hubtel"
}

func (h *HubtelProvider) GetSupportedCountries() []string {
	// Hubtel primarily operates in Ghana
	return []string{"GH"}
}

// ======================= API TYPES =======================

// HubtelReceiveMoneyRequest represents a receive money request
type HubtelReceiveMoneyRequest struct {
	Amount          float64 `json:"amount"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	ClientReference string  `json:"clientReference"`
	CallbackURL     string  `json:"callbackUrl"`
	ReturnURL       string  `json:"returnUrl,omitempty"`
	CancellationURL string  `json:"cancellationUrl,omitempty"`
	// Mobile Money specific
	Channel        string `json:"channel,omitempty"` // mtn-gh, vodafone-gh, tigo-gh, airtel-gh
	CustomerMsisdn string `json:"customerMsisdn,omitempty"`
}

// HubtelPaymentResponse represents the API response
type HubtelPaymentResponse struct {
	ResponseCode string             `json:"responseCode"`
	Status       string             `json:"status"`
	Data         *HubtelPaymentData `json:"data"`
	Message      string             `json:"message,omitempty"`
}

// HubtelPaymentData contains the payment details
type HubtelPaymentData struct {
	CheckoutURL           string  `json:"checkoutUrl"`
	CheckoutID            string  `json:"checkoutId"`
	CheckoutDirectURL     string  `json:"checkoutDirectUrl,omitempty"`
	ClientReference       string  `json:"clientReference"`
	Amount                float64 `json:"amount"`
	Currency              string  `json:"currency,omitempty"`
	TransactionID         string  `json:"transactionId,omitempty"`
	ExternalTransactionID string  `json:"externalTransactionId,omitempty"`
	Description           string  `json:"description,omitempty"`
}

// HubtelStatusResponse represents the status check response
type HubtelStatusResponse struct {
	ResponseCode string            `json:"responseCode"`
	Status       string            `json:"status"`
	Data         *HubtelStatusData `json:"data"`
	Message      string            `json:"message,omitempty"`
}

// HubtelStatusData contains transaction status details
type HubtelStatusData struct {
	TransactionID         string  `json:"transactionId"`
	ExternalTransactionID string  `json:"externalTransactionId,omitempty"`
	ClientReference       string  `json:"clientReference"`
	Amount                float64 `json:"amount"`
	Charges               float64 `json:"charges"`
	AmountAfterCharges    float64 `json:"amountAfterCharges"`
	Currency              string  `json:"currency"`
	Description           string  `json:"description"`
	TransactionStatus     string  `json:"transactionStatus"` // Pending, Success, Failed
	PaymentMethod         string  `json:"paymentMethod"`
	CustomerPhoneNumber   string  `json:"customerPhoneNumber,omitempty"`
	CustomerName          string  `json:"customerName,omitempty"`
}

// HubtelMoMoRequest represents a direct mobile money request
type HubtelMoMoRequest struct {
	CustomerName         string  `json:"CustomerName"`
	CustomerMsisdn       string  `json:"CustomerMsisdn"`
	CustomerEmail        string  `json:"CustomerEmail,omitempty"`
	Channel              string  `json:"Channel"` // mtn-gh, vodafone-gh, tigo-gh
	Amount               float64 `json:"Amount"`
	PrimaryCallbackUrl   string  `json:"PrimaryCallbackUrl"`
	SecondaryCallbackUrl string  `json:"SecondaryCallbackUrl,omitempty"`
	ClientReference      string  `json:"ClientReference"`
	Description          string  `json:"Description,omitempty"`
}

// HubtelMoMoResponse represents direct MoMo API response
type HubtelMoMoResponse struct {
	ResponseCode string          `json:"ResponseCode"`
	Data         *HubtelMoMoData `json:"Data"`
	Message      string          `json:"Message,omitempty"`
}

// HubtelMoMoData contains MoMo transaction data
type HubtelMoMoData struct {
	TransactionId         string  `json:"TransactionId"`
	ClientReference       string  `json:"ClientReference"`
	Description           string  `json:"Description"`
	Amount                float64 `json:"Amount"`
	Charges               float64 `json:"Charges"`
	AmountAfterCharges    float64 `json:"AmountAfterCharges"`
	ExternalTransactionId string  `json:"ExternalTransactionId,omitempty"`
}

// ======================= CHANNEL MAPPING =======================

// GetChannelForNetwork returns the Hubtel channel code for a network
func (h *HubtelProvider) GetChannelForNetwork(network string) string {
	channels := map[string]string{
		"MTN":        "mtn-gh",
		"VODAFONE":   "vodafone-gh",
		"TIGO":       "tigo-gh",
		"AIRTEL":     "airtel-gh",
		"AIRTELTIGO": "tigo-gh",
	}

	if channel, ok := channels[network]; ok {
		return channel
	}
	return "mtn-gh" // Default to MTN
}

// ======================= AUTHENTICATION =======================

// getAuthHeader returns the Basic Auth header for Hubtel API
func (h *HubtelProvider) getAuthHeader() string {
	auth := h.config.ClientID + ":" + h.config.ClientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// ======================= COLLECTION PROVIDER INTERFACE =======================

// InitiateCollection initiates a mobile money collection via Hubtel
func (h *HubtelProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	log.Printf("[Hubtel] 🚀 Initiating collection for Ghana GHS %.2f", req.Amount)

	// Validate credentials
	if h.config.ClientID == "" {
		log.Printf("[Hubtel] ❌ Client ID not configured!")
		log.Printf("[Hubtel] 💡 Configure 'client_id' in admin panel for Hubtel")
		return nil, fmt.Errorf("Hubtel Client ID not configured")
	}
	if h.config.ClientSecret == "" {
		log.Printf("[Hubtel] ❌ Client Secret not configured!")
		log.Printf("[Hubtel] 💡 Configure 'client_secret' in admin panel for Hubtel")
		return nil, fmt.Errorf("Hubtel Client Secret not configured")
	}
	if h.config.MerchantAccount == "" {
		log.Printf("[Hubtel] ❌ Merchant Account not configured!")
		log.Printf("[Hubtel] 💡 Configure 'merchant_key' in admin panel for Hubtel")
		return nil, fmt.Errorf("Hubtel Merchant Account not configured")
	}

	// Generate client reference if not provided
	clientRef := req.ReferenceID
	if clientRef == "" {
		clientRef = "HUB-" + uuid.New().String()[:8]
	}

	log.Printf("[Hubtel]    Client Ref: %s", clientRef)
	log.Printf("[Hubtel]    Phone: %s", req.PhoneNumber)

	// Determine callback URL
	callbackURL := req.RedirectURL
	if callbackURL == "" {
		callbackURL = "https://api.zekora.com/webhooks/hubtel"
	}

	// Use direct MoMo API if phone number is provided
	if req.PhoneNumber != "" {
		return h.initiateDirectMoMo(ctx, req, clientRef, callbackURL)
	}

	// Otherwise use checkout API
	return h.initiateCheckout(ctx, req, clientRef, callbackURL)
}

// initiateDirectMoMo uses the direct MoMo API for USSD push
func (h *HubtelProvider) initiateDirectMoMo(ctx context.Context, req *CollectionRequest, clientRef, callbackURL string) (*CollectionResponse, error) {
	log.Printf("[Hubtel] 📱 Using Direct MoMo API (USSD Push)")

	// Detect network from phone number or use default
	channel := h.GetChannelForNetwork("MTN") // Default to MTN
	if network, ok := req.Metadata["network"]; ok {
		channel = h.GetChannelForNetwork(network)
	}

	log.Printf("[Hubtel]    Channel: %s", channel)

	momoReq := HubtelMoMoRequest{
		CustomerName:       req.Email, // Use email as name if available
		CustomerMsisdn:     req.PhoneNumber,
		CustomerEmail:      req.Email,
		Channel:            channel,
		Amount:             req.Amount,
		PrimaryCallbackUrl: callbackURL,
		ClientReference:    clientRef,
		Description:        fmt.Sprintf("Deposit %s", clientRef),
	}

	if momoReq.CustomerName == "" {
		momoReq.CustomerName = "Customer"
	}

	// Serialize request
	jsonBody, err := json.Marshal(momoReq)
	if err != nil {
		log.Printf("[Hubtel] ❌ Failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	// Build URL
	url := fmt.Sprintf("%s/v1/merchantaccount/merchants/%s/receive/mobilemoney",
		h.config.BaseURL, h.config.MerchantAccount)

	log.Printf("[Hubtel] 📡 Sending request to: %s", url)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", h.getAuthHeader())
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("[Hubtel] ❌ HTTP request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	log.Printf("[Hubtel]    Response status: %d", resp.StatusCode)
	log.Printf("[Hubtel]    Response: %s", string(body))

	// Handle non-success responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[Hubtel] ❌ API error: %s", string(body))
		return nil, fmt.Errorf("Hubtel API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var momoResp HubtelMoMoResponse
	if err := json.Unmarshal(body, &momoResp); err != nil {
		log.Printf("[Hubtel] ❌ Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check response code
	if momoResp.ResponseCode != "0000" && momoResp.ResponseCode != "0001" {
		errMsg := momoResp.Message
		if errMsg == "" {
			errMsg = fmt.Sprintf("Error code: %s", momoResp.ResponseCode)
		}
		log.Printf("[Hubtel] ❌ Transaction failed: %s", errMsg)
		return nil, fmt.Errorf("Hubtel error: %s", errMsg)
	}

	transactionID := ""
	if momoResp.Data != nil {
		transactionID = momoResp.Data.TransactionId
	}

	log.Printf("[Hubtel] ✅ MoMo request initiated: TransactionID=%s", transactionID)

	return &CollectionResponse{
		ReferenceID:       clientRef,
		ProviderReference: transactionID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          "GHS",
		Message:           "USSD prompt sent to phone. Approve to complete payment.",
	}, nil
}

// initiateCheckout uses the checkout API for web payments
func (h *HubtelProvider) initiateCheckout(ctx context.Context, req *CollectionRequest, clientRef, callbackURL string) (*CollectionResponse, error) {
	log.Printf("[Hubtel] 🌐 Using Checkout API")

	checkoutReq := HubtelReceiveMoneyRequest{
		Amount:          req.Amount,
		Title:           "Deposit",
		Description:     fmt.Sprintf("Deposit %s", clientRef),
		ClientReference: clientRef,
		CallbackURL:     callbackURL,
		ReturnURL:       req.RedirectURL,
	}

	// Serialize request
	jsonBody, err := json.Marshal(checkoutReq)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	// Build URL
	url := fmt.Sprintf("%s/v1/merchantaccount/merchants/%s/receive/onlinecheckout",
		h.config.BaseURL, h.config.MerchantAccount)

	log.Printf("[Hubtel] 📡 Sending request to: %s", url)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", h.getAuthHeader())
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("[Hubtel] ❌ HTTP request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	log.Printf("[Hubtel]    Response status: %d", resp.StatusCode)

	// Handle non-success responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[Hubtel] ❌ API error: %s", string(body))
		return nil, fmt.Errorf("Hubtel API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var checkoutResp HubtelPaymentResponse
	if err := json.Unmarshal(body, &checkoutResp); err != nil {
		log.Printf("[Hubtel] ❌ Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if checkoutResp.ResponseCode != "0000" {
		errMsg := checkoutResp.Message
		if errMsg == "" {
			errMsg = fmt.Sprintf("Error code: %s", checkoutResp.ResponseCode)
		}
		return nil, fmt.Errorf("Hubtel error: %s", errMsg)
	}

	checkoutURL := ""
	checkoutID := ""
	if checkoutResp.Data != nil {
		checkoutURL = checkoutResp.Data.CheckoutURL
		checkoutID = checkoutResp.Data.CheckoutID
	}

	log.Printf("[Hubtel] ✅ Checkout created: ID=%s, URL=%s", checkoutID, checkoutURL)

	return &CollectionResponse{
		ReferenceID:       clientRef,
		ProviderReference: checkoutID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          "GHS",
		PaymentLink:       checkoutURL,
		Message:           "Redirect user to payment page",
	}, nil
}

// GetAvailableMethods returns available payment methods for a country
func (h *HubtelProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{CollectionMethodMobileMoney}, nil
}

// VerifyCollection checks the status of a collection
func (h *HubtelProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	log.Printf("[Hubtel] 🔍 Checking status for: %s", referenceID)

	// Build URL
	url := fmt.Sprintf("%s/v1/merchantaccount/merchants/%s/transactions/status?clientReference=%s",
		h.config.BaseURL, h.config.MerchantAccount, referenceID)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", h.getAuthHeader())

	// Execute request
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("status check failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status check failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var statusResp HubtelStatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}

	// Map status
	status := CollectionStatusPending
	var failureReason string

	if statusResp.Data != nil {
		switch statusResp.Data.TransactionStatus {
		case "Success", "Successful":
			status = CollectionStatusSuccessful
		case "Failed", "Declined", "Cancelled":
			status = CollectionStatusFailed
			failureReason = statusResp.Message
		case "Pending":
			status = CollectionStatusPending
		}
	}

	providerRef := ""
	if statusResp.Data != nil {
		providerRef = statusResp.Data.TransactionID
	}

	log.Printf("[Hubtel] ✅ Status for %s: %s", referenceID, status)

	return &CollectionResponse{
		ReferenceID:       referenceID,
		ProviderReference: providerRef,
		Status:            status,
		Message:           failureReason,
	}, nil
}

// HubtelWebhookPayload represents the webhook callback payload from Hubtel
type HubtelWebhookPayload struct {
	ResponseCode          string  `json:"ResponseCode"`
	Status                string  `json:"Status"`
	TransactionId         string  `json:"TransactionId"`
	ClientReference       string  `json:"ClientReference"`
	Amount                float64 `json:"Amount"`
	Charges               float64 `json:"Charges"`
	AmountAfterCharges    float64 `json:"AmountAfterCharges"`
	Currency              string  `json:"Currency"`
	Description           string  `json:"Description"`
	ExternalTransactionId string  `json:"ExternalTransactionId"`
	CustomerPhoneNumber   string  `json:"CustomerPhoneNumber"`
	PaymentMethod         string  `json:"PaymentMethod"`
}

// ParseWebhook parses a Hubtel webhook payload
func (h *HubtelProvider) ParseWebhook(payload []byte) (*HubtelWebhookPayload, error) {
	var webhook HubtelWebhookPayload
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}
	return &webhook, nil
}

// VerifyWebhook verifies the webhook (Hubtel doesn't use signatures by default)
func (h *HubtelProvider) VerifyWebhook(signature string, payload []byte) bool {
	// Hubtel doesn't require signature verification
	// But you can implement IP whitelisting or other validation
	return true
}
