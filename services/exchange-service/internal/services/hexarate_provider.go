package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HexarateProvider fetches fiat currency exchange rates from Hexarate API
type HexarateProvider struct {
	client  *http.Client
	baseURL string
	cache   map[string]float64
	cacheMu sync.RWMutex
	cacheTime time.Time
}

// HexarateResponse represents the API response structure
type HexarateResponse struct {
	StatusCode int `json:"status_code"`
	Data       struct {
		Base      string  `json:"base"`
		Target    string  `json:"target"`
		Mid       float64 `json:"mid"`
		Unit      int     `json:"unit"`
		Timestamp string  `json:"timestamp"`
	} `json:"data"`
}

// SupportedFiatCurrencies contains all fiat currencies we support
var SupportedFiatCurrencies = []string{
	// Americas
	"USD", "CAD", "MXN", "BRL", "ARS", "CLP", "COP", "PEN", "UYU", "VES",
	// Europe
	"EUR", "GBP", "CHF", "NOK", "SEK", "DKK", "PLN", "CZK", "HUF", "RON", "RUB", "UAH", "TRY", "BGN", "HRK", "ISK",
	// Asia
	"JPY", "CNY", "HKD", "SGD", "KRW", "INR", "IDR", "MYR", "THB", "PHP", "VND", "PKR", "BDT", "TWD", "LKR", "NPR",
	// Middle East
	"AED", "SAR", "QAR", "KWD", "BHD", "OMR", "ILS", "EGP", "JOD", "LBP", "IQD",
	// Africa
	"XAF", "XOF", "NGN", "ZAR", "KES", "GHS", "MAD", "TND", "DZD", "AOA", "ETB", "UGX", "TZS", "RWF", "MZN", "BWP", "MUR", "SCR",
	// Oceania
	"AUD", "NZD", "FJD", "PGK",
}

// NewHexarateProvider creates a new Hexarate rate provider
func NewHexarateProvider() *HexarateProvider {
	return &HexarateProvider{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: "https://hexarate.paikama.co/api/rates",
		cache:   make(map[string]float64),
	}
}

// GetFiatRate fetches the exchange rate between two fiat currencies
func (p *HexarateProvider) GetFiatRate(from, to string) (float64, error) {
	// Check cache first
	cacheKey := from + "_" + to
	p.cacheMu.RLock()
	if rate, ok := p.cache[cacheKey]; ok && time.Since(p.cacheTime) < time.Hour {
		p.cacheMu.RUnlock()
		return rate, nil
	}
	p.cacheMu.RUnlock()

	// Fetch from API
	url := fmt.Sprintf("%s/%s/%s/latest", p.baseURL, from, to)
	
	resp, err := p.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch rate from Hexarate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Hexarate API returned status: %d", resp.StatusCode)
	}

	var result HexarateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode Hexarate response: %w", err)
	}

	if result.StatusCode != 200 {
		return 0, fmt.Errorf("Hexarate API error: status %d", result.StatusCode)
	}

	// Cache the result
	p.cacheMu.Lock()
	p.cache[cacheKey] = result.Data.Mid
	p.cacheTime = time.Now()
	p.cacheMu.Unlock()

	return result.Data.Mid, nil
}

// GetAllUSDRates fetches USD rates for all supported currencies
// Returns a map of currency code -> rate (1 USD = X currency)
func (p *HexarateProvider) GetAllUSDRates() (map[string]float64, error) {
	rates := make(map[string]float64)
	rates["USD"] = 1.0 // USD to USD is always 1
	
	// Fetch rates in batches to avoid overloading the API
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(SupportedFiatCurrencies))
	
	// Limit concurrent requests
	semaphore := make(chan struct{}, 5)
	
	for _, currency := range SupportedFiatCurrencies {
		if currency == "USD" {
			continue
		}
		
		wg.Add(1)
		go func(curr string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			rate, err := p.GetFiatRate("USD", curr)
			if err != nil {
				// Log error but continue with other currencies
				fmt.Printf("Warning: Failed to fetch USD/%s rate: %v\n", curr, err)
				return
			}
			
			mu.Lock()
			rates[curr] = rate
			mu.Unlock()
		}(currency)
	}
	
	wg.Wait()
	close(errChan)
	
	// We don't fail if some currencies fail, just return what we got
	if len(rates) < 5 {
		return nil, fmt.Errorf("failed to fetch enough currency rates")
	}
	
	return rates, nil
}

// GetAllEURRates fetches EUR rates for all supported currencies
func (p *HexarateProvider) GetAllEURRates() (map[string]float64, error) {
	rates := make(map[string]float64)
	rates["EUR"] = 1.0
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, 5)
	
	for _, currency := range SupportedFiatCurrencies {
		if currency == "EUR" {
			continue
		}
		
		wg.Add(1)
		go func(curr string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			rate, err := p.GetFiatRate("EUR", curr)
			if err != nil {
				return
			}
			
			mu.Lock()
			rates[curr] = rate
			mu.Unlock()
		}(currency)
	}
	
	wg.Wait()
	return rates, nil
}
