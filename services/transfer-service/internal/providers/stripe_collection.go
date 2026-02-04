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

// StripeConfig holds Stripe API configuration
type StripeConfig struct {
	PublicKey     string
	SecretKey     string
	WebhookSecret string
	BaseURL       string
}

// StripeCollectionProvider implements Stripe payment collection
type StripeCollectionProvider struct {
	config     StripeConfig
	httpClient *http.Client
}

// NewStripeCollectionProvider creates a new Stripe provider
func NewStripeCollectionProvider(config StripeConfig) *StripeCollectionProvider {
	return &StripeCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *StripeCollectionProvider) GetName() string {
	return "stripe"
}

func (s *StripeCollectionProvider) GetSupportedCountries() []string {
	// Stripe is global
	return []string{"*"}
}

func (s *StripeCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{
		CollectionMethodCard,
	}, nil
}

type stripeCheckoutSessionRequest struct {
	PaymentMethodTypes []string         `json:"payment_method_types"`
	LineItems          []stripeLineItem `json:"line_items"`
	Mode               string           `json:"mode"`
	SuccessURL         string           `json:"success_url"`
	CancelURL          string           `json:"cancel_url"`
	ClientReferenceID  string           `json:"client_reference_id,omitempty"`
	CustomerEmail      string           `json:"customer_email,omitempty"`
}

type stripeLineItem struct {
	PriceData stripePriceData `json:"price_data"`
	Quantity  int             `json:"quantity"`
}

type stripePriceData struct {
	Currency    string            `json:"currency"`
	ProductData stripeProductData `json:"product_data"`
	UnitAmount  int64             `json:"unit_amount"` // Amount in cents
}

type stripeProductData struct {
	Name string `json:"name"`
}

type stripeCheckoutSessionResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func (s *StripeCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// Convert amount to cents
	amountInCents := int64(req.Amount * 100)

	sessionReq := stripeCheckoutSessionRequest{
		PaymentMethodTypes: []string{"card"},
		LineItems: []stripeLineItem{
			{
				PriceData: stripePriceData{
					Currency: req.Currency,
					ProductData: stripeProductData{
						Name: "Wallet Deposit",
					},
					UnitAmount: amountInCents,
				},
				Quantity: 1,
			},
		},
		Mode:              "payment",
		SuccessURL:        req.RedirectURL + "?session_id={CHECKOUT_SESSION_ID}",
		CancelURL:         req.RedirectURL + "?canceled=true",
		ClientReferenceID: req.ReferenceID,
		CustomerEmail:     req.Email,
	}

	jsonData, err := json.Marshal(sessionReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL+"/checkout/sessions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.config.SecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe API error (status %d): %s", resp.StatusCode, string(body))
	}

	var sessionResp stripeCheckoutSessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: sessionResp.ID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       sessionResp.URL,
		Message:           "Redirecting to Stripe Checkout",
	}, nil
}

func (s *StripeCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	// Retrieve checkout session
	url := fmt.Sprintf("%s/checkout/sessions/%s", s.config.BaseURL, referenceID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.config.SecretKey)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		ID                string `json:"id"`
		PaymentStatus     string `json:"payment_status"`
		ClientReferenceID string `json:"client_reference_id"`
		AmountTotal       int64  `json:"amount_total"`
		Currency          string `json:"currency"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.PaymentStatus == "paid" {
		status = CollectionStatusSuccessful
	} else if result.PaymentStatus == "unpaid" || result.PaymentStatus == "canceled" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       result.ClientReferenceID,
		ProviderReference: result.ID,
		Status:            status,
		Amount:            float64(result.AmountTotal) / 100,
		Currency:          result.Currency,
	}, nil
}
