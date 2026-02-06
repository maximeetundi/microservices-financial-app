package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// PawapayConfig holds configuration for Pawapay API
type PawapayConfig struct {
	APIKey        string
	WebhookSecret string
	BaseURL       string // https://api.sandbox.pawapay.cloud or https://api.pawapay.cloud
}

// PawapayProvider implements CollectionProvider for Pawapay
type PawapayProvider struct {
	config     PawapayConfig
	httpClient *http.Client
}

// NewPawapayCollectionProvider creates a new Pawapay collection provider
func NewPawapayCollectionProvider(config PawapayConfig) *PawapayProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.sandbox.pawapay.cloud"
	}
	return &PawapayProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *PawapayProvider) GetName() string {
	return "pawapay"
}

func (p *PawapayProvider) GetSupportedCountries() []string {
	// Pawapay supports multiple African countries
	return []string{
		"BJ", // Benin
		"BF", // Burkina Faso
		"CM", // Cameroon
		"CI", // Côte d'Ivoire
		"CD", // DRC
		"GH", // Ghana
		"KE", // Kenya
		"MW", // Malawi
		"ML", // Mali
		"MZ", // Mozambique
		"NE", // Niger
		"NG", // Nigeria
		"RW", // Rwanda
		"SN", // Senegal
		"TZ", // Tanzania
		"TG", // Togo
		"UG", // Uganda
		"ZM", // Zambia
	}
}

// ======================= API TYPES =======================

// PawapayDepositRequest represents a deposit initiation request
type PawapayDepositRequest struct {
	DepositID     string                 `json:"depositId"`
	Amount        string                 `json:"amount"`
	Currency      string                 `json:"currency"`
	Correspondent string                 `json:"correspondent"` // MNO code (e.g., MTN_MOMO_BEN)
	Payer         PawapayPayer           `json:"payer"`
	CustomerTimestamp string             `json:"customerTimestamp"`
	StatementDescription string          `json:"statementDescription,omitempty"`
	Metadata      []PawapayMetadataItem  `json:"metadata,omitempty"`
}

// PawapayPayer represents the payer information
type PawapayPayer struct {
	Type    string                 `json:"type"` // MSISDN
	Address PawapayPayerAddress    `json:"address"`
}

// PawapayPayerAddress represents the payer's phone number
type PawapayPayerAddress struct {
	Value string `json:"value"` // Phone number in E.164 format
}

// PawapayMetadataItem represents a metadata key-value pair
type PawapayMetadataItem struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
	IsPII      bool   `json:"isPII,omitempty"`
}

// PawapayDepositResponse represents the API response for a deposit
type PawapayDepositResponse struct {
	DepositID        string    `json:"depositId"`
	Status           string    `json:"status"` // ACCEPTED, SUBMITTED, COMPLETED, FAILED
	Created          time.Time `json:"created"`
	CorrespondentIDs []string  `json:"correspondentIds,omitempty"`
	Amount           string    `json:"amount,omitempty"`
	Currency         string    `json:"currency,omitempty"`
	Country          string    `json:"country,omitempty"`
	Correspondent    string    `json:"correspondent,omitempty"`
	Payer            *PawapayPayer `json:"payer,omitempty"`
	// Error fields
	RejectionReason *PawapayRejectionReason `json:"rejectionReason,omitempty"`
}

// PawapayRejectionReason contains details about why a transaction was rejected
type PawapayRejectionReason struct {
	RejectionCode    string `json:"rejectionCode"`
	RejectionMessage string `json:"rejectionMessage"`
}

// PawapayStatusResponse represents the status check response
type PawapayStatusResponse struct {
	DepositID            string                   `json:"depositId"`
	Status               string                   `json:"status"`
	Amount               string                   `json:"amount"`
	Currency             string                   `json:"currency"`
	Country              string                   `json:"country"`
	Correspondent        string                   `json:"correspondent"`
	Payer                *PawapayPayer            `json:"payer"`
	CustomerTimestamp    string                   `json:"customerTimestamp"`
	StatementDescription string                   `json:"statementDescription"`
	Created              time.Time                `json:"created"`
	ReceivedByRecipient  time.Time                `json:"receivedByRecipient,omitempty"`
	CorrespondentIds     []string                 `json:"correspondentIds"`
	Metadata             []PawapayMetadataItem    `json:"metadata"`
	FailureReason        *PawapayRejectionReason  `json:"failureReason,omitempty"`
}

// ======================= CORRESPONDENT MAPPING =======================

// GetCorrespondentForCountry returns the appropriate MNO code for a country
func (p *PawapayProvider) GetCorrespondentForCountry(country, currency string) string {
	// Map country to primary correspondent
	correspondents := map[string]string{
		"BJ": "MTN_MOMO_BEN",       // Benin - MTN
		"BF": "ORANGE_BFA",         // Burkina Faso - Orange
		"CM": "MTN_MOMO_CMR",       // Cameroon - MTN
		"CI": "MTN_MOMO_CIV",       // Côte d'Ivoire - MTN
		"CD": "VODACOM_MPESA_COD",  // DRC - Vodacom M-Pesa
		"GH": "MTN_MOMO_GHA",       // Ghana - MTN
		"KE": "MPESA_KEN",          // Kenya - M-Pesa
		"MW": "AIRTEL_MWI",         // Malawi - Airtel
		"ML": "ORANGE_MLI",         // Mali - Orange
		"MZ": "VODACOM_MOZ",        // Mozambique - Vodacom
		"NE": "AIRTEL_NER",         // Niger - Airtel
		"NG": "MTN_MOMO_NGA",       // Nigeria - MTN
		"RW": "MTN_MOMO_RWA",       // Rwanda - MTN
		"SN": "ORANGE_SEN",         // Senegal - Orange
		"TZ": "VODACOM_TZA",        // Tanzania - Vodacom
		"TG": "MOOV_TGO",           // Togo - Moov
		"UG": "MTN_MOMO_UGA",       // Uganda - MTN
		"ZM": "MTN_MOMO_ZMB",       // Zambia - MTN
	}

	if correspondent, ok := correspondents[country]; ok {
		return correspondent
	}
	return "MTN_MOMO_" + country // Fallback
}

// ======================= COLLECTION PROVIDER INTERFACE =======================

// InitiateCollection initiates a mobile money collection via Pawapay
func (p *PawapayProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	log.Printf("[Pawapay] 🚀 Initiating collection for %s %.2f %s", req.Country, req.Amount, req.Currency)

	// Validate API key
	if p.config.APIKey == "" {
		log.Printf("[Pawapay] ❌ API Key not configured!")
		return nil, fmt.Errorf("Pawapay API key not configured")
	}

	// Generate deposit ID if not provided
	depositID := req.ReferenceID
	if depositID == "" {
		depositID = uuid.New().String()
	}

	// Get correspondent for country
	correspondent := p.GetCorrespondentForCountry(req.Country, req.Currency)
	log.Printf("[Pawapay]    Correspondent: %s", correspondent)

	// Build request
	depositReq := PawapayDepositRequest{
		DepositID:     depositID,
		Amount:        fmt.Sprintf("%.2f", req.Amount),
		Currency:      req.Currency,
		Correspondent: correspondent,
		Payer: PawapayPayer{
			Type: "MSISDN",
			Address: PawapayPayerAddress{
				Value: req.PhoneNumber,
			},
		},
		CustomerTimestamp:    time.Now().UTC().Format(time.RFC3339),
		StatementDescription: fmt.Sprintf("Deposit %s", depositID[:8]),
	}

	// Add metadata
	if req.UserID != "" {
		depositReq.Metadata = append(depositReq.Metadata, PawapayMetadataItem{
			FieldName:  "user_id",
			FieldValue: req.UserID,
		})
	}
	if req.WalletID != "" {
		depositReq.Metadata = append(depositReq.Metadata, PawapayMetadataItem{
			FieldName:  "wallet_id",
			FieldValue: req.WalletID,
		})
	}

	// Serialize request
	jsonBody, err := json.Marshal(depositReq)
	if err != nil {
		log.Printf("[Pawapay] ❌ Failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	log.Printf("[Pawapay] 📡 Sending request to: %s/deposits", p.config.BaseURL)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/deposits", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("[Pawapay] ❌ HTTP request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := io.ReadAll(resp.Body)

	log.Printf("[Pawapay]    Response status: %d", resp.StatusCode)

	// Handle non-success responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[Pawapay] ❌ API error: %s", string(body))
		return nil, fmt.Errorf("Pawapay API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var depositResp PawapayDepositResponse
	if err := json.Unmarshal(body, &depositResp); err != nil {
		log.Printf("[Pawapay] ❌ Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("[Pawapay] ✅ Deposit initiated: ID=%s, Status=%s", depositResp.DepositID, depositResp.Status)

	// Map status
	status := CollectionStatusPending
	switch depositResp.Status {
	case "ACCEPTED", "SUBMITTED":
		status = CollectionStatusPending
	case "COMPLETED":
		status = CollectionStatusSuccess
	case "FAILED":
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       depositID,
		ProviderReference: depositResp.DepositID,
		Status:            status,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Message:           fmt.Sprintf("Deposit %s: %s", depositResp.DepositID, depositResp.Status),
	}, nil
}

// GetCollectionStatus checks the status of a collection
func (p *PawapayProvider) GetCollectionStatus(ctx context.Context, referenceID string) (*CollectionStatusResponse, error) {
	log.Printf("[Pawapay] 🔍 Checking status for: %s", referenceID)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/deposits/"+referenceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("status check failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status check failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var statusResp []PawapayStatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		// Try single object response
		var singleResp PawapayStatusResponse
		if err2 := json.Unmarshal(body, &singleResp); err2 != nil {
			return nil, fmt.Errorf("failed to parse status response: %w", err)
		}
		statusResp = []PawapayStatusResponse{singleResp}
	}

	if len(statusResp) == 0 {
		return nil, fmt.Errorf("no status found for deposit %s", referenceID)
	}

	deposit := statusResp[0]

	// Map status
	status := CollectionStatusPending
	var failureReason string

	switch deposit.Status {
	case "ACCEPTED", "SUBMITTED":
		status = CollectionStatusPending
	case "COMPLETED":
		status = CollectionStatusSuccess
	case "FAILED":
		status = CollectionStatusFailed
		if deposit.FailureReason != nil {
			failureReason = deposit.FailureReason.RejectionMessage
		}
	}

	log.Printf("[Pawapay] ✅ Status for %s: %s", referenceID, deposit.Status)

	return &CollectionStatusResponse{
		ReferenceID:       referenceID,
		ProviderReference: deposit.DepositID,
		Status:            status,
		FailureReason:     failureReason,
	}, nil
}

// VerifyWebhook verifies the webhook signature from Pawapay
func (p *PawapayProvider) VerifyWebhook(signature string, payload []byte) bool {
	// Pawapay uses a different signature mechanism
	// For now, just validate that we have a webhook secret configured
	if p.config.WebhookSecret == "" {
		log.Printf("[Pawapay] ⚠️ Webhook secret not configured, skipping verification")
		return true
	}

	// TODO: Implement proper webhook signature verification
	// Pawapay typically uses HMAC-SHA256 or similar
	return true
}

// PawapayWebhookPayload represents the webhook callback payload
type PawapayWebhookPayload struct {
	DepositID        string                  `json:"depositId"`
	Status           string                  `json:"status"`
	Amount           string                  `json:"amount"`
	Currency         string                  `json:"currency"`
	Country          string                  `json:"country"`
	Correspondent    string                  `json:"correspondent"`
	Payer            *PawapayPayer           `json:"payer"`
	Created          time.Time               `json:"created"`
	ReceivedByRecipient time.Time            `json:"receivedByRecipient,omitempty"`
	Metadata         []PawapayMetadataItem   `json:"metadata"`
	FailureReason    *PawapayRejectionReason `json:"failureReason,omitempty"`
}

// ParseWebhook parses a Pawapay webhook payload
func (p *PawapayProvider) ParseWebhook(payload []byte) (*PawapayWebhookPayload, error) {
	var webhook PawapayWebhookPayload
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}
	return &webhook, nil
}
