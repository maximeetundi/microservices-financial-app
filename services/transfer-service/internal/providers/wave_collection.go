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

// WaveConfig holds Wave API configuration
type WaveConfig struct {
	APIKey  string
	BaseURL string
}

// WaveCollectionProvider implements Wave payment collection
type WaveCollectionProvider struct {
	config     WaveConfig
	httpClient *http.Client
}

// NewWaveCollectionProvider creates a new Wave provider
func NewWaveCollectionProvider(config WaveConfig) *WaveCollectionProvider {
	return &WaveCollectionProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (w *WaveCollectionProvider) GetName() string {
	return "wave"
}

func (w *WaveCollectionProvider) GetSupportedCountries() []string {
	return []string{"SN", "CI", "BF", "ML", "GM", "UG"}
}

func (w *WaveCollectionProvider) GetAvailableMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	return []CollectionMethod{
		CollectionMethodMobileMoney,
	}, nil
}

type waveCheckoutRequest struct {
	Amount     int    `json:"amount"`
	Currency   string `json:"currency"`
	ErrorURL   string `json:"error_url"`
	SuccessURL string `json:"success_url"`
}

type waveCheckoutResponse struct {
	ID            string `json:"id"`
	WaveLaunchURL string `json:"wave_launch_url"`
	Status        string `json:"status"`
}

func (w *WaveCollectionProvider) InitiateCollection(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	waveReq := waveCheckoutRequest{
		Amount:     int(req.Amount),
		Currency:   req.Currency,
		ErrorURL:   req.RedirectURL + "?error=true",
		SuccessURL: req.RedirectURL + "?success=true",
	}

	jsonData, err := json.Marshal(waveReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", w.config.BaseURL+"/checkout/sessions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+w.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("wave API error (status %d): %s", resp.StatusCode, string(body))
	}

	var waveResp waveCheckoutResponse
	if err := json.Unmarshal(body, &waveResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &CollectionResponse{
		ReferenceID:       req.ReferenceID,
		ProviderReference: waveResp.ID,
		Status:            CollectionStatusPending,
		Amount:            req.Amount,
		Currency:          req.Currency,
		PaymentLink:       waveResp.WaveLaunchURL,
		Message:           "Redirecting to Wave",
	}, nil
}

func (w *WaveCollectionProvider) VerifyCollection(ctx context.Context, referenceID string) (*CollectionResponse, error) {
	url := fmt.Sprintf("%s/checkout/sessions/%s", w.config.BaseURL, referenceID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+w.config.APIKey)

	resp, err := w.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := CollectionStatusPending
	if result.Status == "complete" {
		status = CollectionStatusSuccessful
	} else if result.Status == "failed" || result.Status == "cancelled" {
		status = CollectionStatusFailed
	}

	return &CollectionResponse{
		ReferenceID:       referenceID,
		ProviderReference: result.ID,
		Status:            status,
		Amount:            float64(result.Amount),
		Currency:          result.Currency,
	}, nil
}
