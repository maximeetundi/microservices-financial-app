package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CoinGeckoProvider fetches cryptocurrency prices from CoinGecko API
type CoinGeckoProvider struct {
	client  *http.Client
	baseURL string
}

// CoinGeckoPrice represents a price response from CoinGecko
type CoinGeckoPrice struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
}

// CoinGeckoResponse is the API response structure
type CoinGeckoResponse map[string]CoinGeckoPrice

// CryptoSymbolMapping maps our internal symbols to CoinGecko IDs
var CryptoSymbolMapping = map[string]string{
	"BTC":  "bitcoin",
	"ETH":  "ethereum",
	"BNB":  "binancecoin",
	"SOL":  "solana",
	"XRP":  "ripple",
	"DOGE": "dogecoin",
	"ADA":  "cardano",
	"DOT":  "polkadot",
	"LTC":  "litecoin",
	"USDT": "tether",
	"USDC": "usd-coin",
	"AVAX": "avalanche-2",
	"MATIC": "matic-network",
	"LINK": "chainlink",
	"UNI":  "uniswap",
	"ATOM": "cosmos",
	"XLM":  "stellar",
	"ALGO": "algorand",
	"VET":  "vechain",
	"FIL":  "filecoin",
}

// NewCoinGeckoProvider creates a new CoinGecko price provider
func NewCoinGeckoProvider() *CoinGeckoProvider {
	return &CoinGeckoProvider{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: "https://api.coingecko.com/api/v3",
	}
}

// GetCryptoPrices fetches prices for multiple cryptocurrencies
// Returns a map of symbol -> CoinGeckoPrice
func (p *CoinGeckoProvider) GetCryptoPrices(symbols []string) (map[string]CoinGeckoPrice, error) {
	// Convert symbols to CoinGecko IDs
	var coinIDs []string
	symbolToID := make(map[string]string)
	
	for _, symbol := range symbols {
		if id, ok := CryptoSymbolMapping[strings.ToUpper(symbol)]; ok {
			coinIDs = append(coinIDs, id)
			symbolToID[id] = strings.ToUpper(symbol)
		}
	}
	
	if len(coinIDs) == 0 {
		return nil, fmt.Errorf("no valid crypto symbols provided")
	}
	
	// Build API URL
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd,eur",
		p.baseURL,
		strings.Join(coinIDs, ","))
	
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices from CoinGecko: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko API returned status: %d", resp.StatusCode)
	}
	
	var result CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode CoinGecko response: %w", err)
	}
	
	// Convert back to our symbol format
	prices := make(map[string]CoinGeckoPrice)
	for coinID, price := range result {
		if symbol, ok := symbolToID[coinID]; ok {
			prices[symbol] = price
		}
	}
	
	return prices, nil
}

// GetAllCryptoPrices fetches prices for all supported cryptocurrencies
func (p *CoinGeckoProvider) GetAllCryptoPrices() (map[string]CoinGeckoPrice, error) {
	symbols := make([]string, 0, len(CryptoSymbolMapping))
	for symbol := range CryptoSymbolMapping {
		symbols = append(symbols, symbol)
	}
	return p.GetCryptoPrices(symbols)
}

// GetCryptoPrice fetches price for a single cryptocurrency
func (p *CoinGeckoProvider) GetCryptoPrice(symbol string) (*CoinGeckoPrice, error) {
	prices, err := p.GetCryptoPrices([]string{symbol})
	if err != nil {
		return nil, err
	}
	
	if price, ok := prices[strings.ToUpper(symbol)]; ok {
		return &price, nil
	}
	
	return nil, fmt.Errorf("price not found for %s", symbol)
}
