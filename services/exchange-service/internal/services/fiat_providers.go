package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FiatRateProviderInterface defines the interface for fetching fiat rates
type FiatRateProviderInterface interface {
	GetRates(baseCurrency string) (map[string]float64, error)
	Name() string
}

// FixerProvider implementation
type FixerProvider struct {
	client *http.Client
	apiKey string
}

func NewFixerProvider(apiKey string) *FixerProvider {
	return &FixerProvider{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}

func (p *FixerProvider) Name() string { return "Fixer.io" }

func (p *FixerProvider) GetRates(baseCurrency string) (map[string]float64, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key missing")
	}
	// Fixer free tier restricts base to EUR
	url := fmt.Sprintf("http://data.fixer.io/api/latest?access_key=%s", p.apiKey)
	
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var result struct {
		Success bool               `json:"success"`
		Rates   map[string]float64 `json:"rates"`
		Error   struct {
			Type string `json:"type"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("api error: %s", result.Error.Type)
	}

	return result.Rates, nil
}

// CurrencyLayerProvider implementation
type CurrencyLayerProvider struct {
	client *http.Client
	apiKey string
}

func NewCurrencyLayerProvider(apiKey string) *CurrencyLayerProvider {
	return &CurrencyLayerProvider{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}

func (p *CurrencyLayerProvider) Name() string { return "CurrencyLayer" }

func (p *CurrencyLayerProvider) GetRates(baseCurrency string) (map[string]float64, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key missing")
	}
	// Free tier restricts source to USD
	url := fmt.Sprintf("http://apilayer.net/api/live?access_key=%s", p.apiKey)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var result struct {
		Success bool               `json:"success"`
		Quotes  map[string]float64 `json:"quotes"`
		Error   struct {
			Info string `json:"info"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("api error: %s", result.Error.Info)
	}

	// Normalize keys (remove "USD" prefix)
	rates := make(map[string]float64)
	for k, v := range result.Quotes {
		if len(k) == 6 {
			rates[k[3:]] = v
		} else {
			rates[k] = v
		}
	}

	return rates, nil
}

// FailoverFiatProvider orchestrator
type FailoverFiatProvider struct {
	providers []FiatRateProviderInterface
}

func NewFailoverFiatProvider(fixerKey, layerKey string) *FailoverFiatProvider {
	providers := []FiatRateProviderInterface{}

	if fixerKey != "" {
		providers = append(providers, NewFixerProvider(fixerKey))
	}
	if layerKey != "" {
		providers = append(providers, NewCurrencyLayerProvider(layerKey))
	}
	// Always include the free fallback
	providers = append(providers, NewFiatRateProvider()) // Original open.er-api

	return &FailoverFiatProvider{
		providers: providers,
	}
}

func (p *FailoverFiatProvider) GetRates(baseCurrency string) (map[string]float64, string, error) {
	var lastErr error
	for _, provider := range p.providers {
		rates, err := provider.GetRates(baseCurrency)
		if err == nil && len(rates) > 0 {
			return rates, provider.Name(), nil
		}
		lastErr = err
		fmt.Printf("[FiatFailover] Provider %s failed: %v\n", provider.Name(), err)
	}
	return nil, "", fmt.Errorf("all providers failed, last error: %v", lastErr)
}

// Update existing FiatRateProvider to implement interface if needed, 
// allows it to be used in the list.
func (p *FiatRateProvider) Name() string { return "OpenER (Free)" }
