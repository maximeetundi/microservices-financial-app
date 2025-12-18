package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FiatRateProvider struct {
	client *http.Client
}

type ExchangeRateAPIResponse struct {
	Result             string             `json:"result"`
	Provider           string             `json:"provider"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	BaseCode           string             `json:"base_code"`
	Rates              map[string]float64 `json:"rates"`
}

func NewFiatRateProvider() *FiatRateProvider {
	return &FiatRateProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *FiatRateProvider) GetRates(baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("https://open.er-api.com/v6/latest/%s", baseCurrency)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var result ExchangeRateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Result != "success" {
		return nil, fmt.Errorf("api returned error result: %s", result.Result)
	}

	return result.Rates, nil
}
