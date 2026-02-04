package providers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// MTNMoMoConfig holds MTN MoMo API configuration
type MTNMoMoConfig struct {
	APIUser         string
	APIKey          string
	SubscriptionKey string
	BaseURL         string
	Environment     string // "sandbox" or "prod"
}

// MTNMoMoCollectionProvider implements MTN MoMo Collections API
type MTNMoMoCollectionProvider struct {
	config      MTNMoMoConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// NewMTNMoMoCollectionProvider creates a new MTN MoMo provider
func NewMTNMoMoCollectionProvider(config MTNMoMoConfig) *MTNMoMoCollectionProvider {
	return &MTNMoMoCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (m *MTNMoMoCollectionProvider) GetName() string {
	return "mtn_momo"
}

func (m *MTNMoMoCollectionProvider) GetSupportedCountries() []string {
	return []string{"CI", "CM", "GH", "UG", "RW", "BJ", "CG", "GN", "ZA", "NG", "ZM"}
}

func (m *MTNMoMoCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{
		CollectionMethodMobileMoney,
	}, nil
}

// getAccessToken retrieves Bearer token for MTN MoMo API
func (m *MTNMoMoCollectionProvider) getAccessToken(ctx context.Context) error {
	if time.Now().Before(m.tokenExpiry) && m.accessToken != "" {
		return nil
	}

	// Create basic auth
	auth := base64.StdEncoding.EncodeToString([]byte(m.config.APIUser + ":" + m.config.APIKey))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", m.config.BaseURL+"/collection/token/", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Ocp-Apim-Subscription-Key", m.config.SubscriptionKey)

	resp, err := m.httpClient.Do(httpReq)
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

	m.accessToken = tokenResp.AccessToken
	m.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

type mtnMoMoRequestToPayRequest struct {
	Amount       string       `json:"amount"`
	Currency     string       `json:"currency"`
	ExternalID   string       `json:"externalId"`
	Payer        mtnMoMoPayer `json:"payer"`
	PayerMessage string       `json:"payerMessage"`
	PayeeNote    string       `json:"payeeNote"`
}

type mtnMoMoPayer struct {
	PartyIdType string `json:"partyIdType"`
	PartyId     string `json:"partyId"`
}

func (m *MTNMoMoCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	if err := m.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	// Generate unique reference ID for this request
	referenceID := uuid.New().String()

	momoReq := mtnMoMoRequestToPayRequest{
		Amount:     fmt.Sprintf("%.0f", req.Amount),
		Currency:   req.Currency,
		ExternalID: req.ReferenceID,
		Payer: mtnMoMoPayer{
			PartyIdType: "MSISDN",
			PartyId:     req.PhoneNumber,
		},
		PayerMessage: "Wallet deposit",
		PayeeNote:    fmt.Sprintf("Deposit %.2f %s", req.Amount, req.Currency),
	}

	jsonData, err := json.Marshal(momoReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", m.config.BaseURL+"/collection/v1_0/requesttopay", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+m.accessToken)
	httpReq.Header.Set("X-Reference-Id", referenceID)
	httpReq.Header.Set("X-Target-Environment", m.config.Environment)
	httpReq.Header.Set("Ocp-Apim-Subscription-Key", m.config.SubscriptionKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mtn momo API error (status %d): %s", resp.StatusCode, string(body))
	}

	// MTN MoMo doesn't return a payment URL, payment is done via USSD/app
	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: referenceID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		Message:           fmt.Sprintf("Payment request sent to %s. Please approve on your phone.", req.PhoneNumber),
		Metadata: map[string]string{
			"mtn_reference_id": referenceID,
		},
	}, nil
}

func (m *MTNMoMoCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	if err := m.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	url := fmt.Sprintf("%s/collection/v1_0/requesttopay/%s", m.config.BaseURL, referenceID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+m.accessToken)
	httpReq.Header.Set("X-Target-Environment", m.config.Environment)
	httpReq.Header.Set("Ocp-Apim-Subscription-Key", m.config.SubscriptionKey)

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		Amount     string `json:"amount"`
		Currency   string `json:"currency"`
		ExternalID string `json:"externalId"`
		Status     string `json:"status"`
		Reason     string `json:"reason,omitempty"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.Status == "SUCCESSFUL" {
		status = CollectionStatusSuccessful
	} else if result.Status == "FAILED" || result.Status == "REJECTED" {
		status = CollectionStatusFailed
	}

	amount := 0.0
	fmt.Sscanf(result.Amount, "%f", &amount)

	return &CollectionResponse{
		ReferenceID:       result.ExternalID,
		ProviderReference: referenceID,
		Status:            status,
		Amount:            amount,
		Currency:          result.Currency,
		Message:           result.Reason,
	}, nil
}
