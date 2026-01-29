package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
)

// TatumProvider implements BlockchainProvider interface for Tatum API
// Used for: GetBalance (blockchain address balance), BroadcastTransaction (raw tx)
// NOTE: Virtual Account functions have been removed - using self-custody model
type TatumProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewTatumProvider(cfg *config.Config) *TatumProvider {
	return &TatumProvider{
		apiKey:  cfg.TatumAPIKey,
		baseURL: cfg.TatumBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetBalance retrieves the blockchain balance for a specific address
// Implements BlockchainProvider interface
func (p *TatumProvider) GetBalance(currency, address string) (float64, error) {
	var endpoint string
	switch currency {
	case "BTC":
		endpoint = fmt.Sprintf("bitcoin/address/balance/%s", address)
	case "ETH":
		endpoint = fmt.Sprintf("ethereum/account/balance/%s", address)
	case "BSC", "BNB":
		endpoint = fmt.Sprintf("bsc/account/balance/%s", address)
	case "SOL":
		endpoint = fmt.Sprintf("solana/account/balance/%s", address)
	case "TRX":
		endpoint = fmt.Sprintf("tron/account/%s", address)
	case "MATIC":
		endpoint = fmt.Sprintf("polygon/account/balance/%s", address)
	case "LTC":
		endpoint = fmt.Sprintf("litecoin/address/balance/%s", address)
	case "XRP":
		endpoint = fmt.Sprintf("xrp/account/%s/balance", address)
	default:
		// Default to ETH-like for ERC20/tokens
		endpoint = fmt.Sprintf("%s/account/balance/%s", strings.ToLower(currency), address)
	}

	url := fmt.Sprintf("%s/%s", p.baseURL, endpoint)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("tatum request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("tatum balance error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response based on currency
	type GenericBalance struct {
		Balance string `json:"balance"`
	}

	if currency == "BTC" || currency == "LTC" {
		type UTXOBalance struct {
			Incoming string `json:"incoming"`
			Outgoing string `json:"outgoing"`
		}
		var btcBal UTXOBalance
		if err := json.NewDecoder(resp.Body).Decode(&btcBal); err == nil {
			var incoming, outgoing float64
			fmt.Sscanf(btcBal.Incoming, "%f", &incoming)
			fmt.Sscanf(btcBal.Outgoing, "%f", &outgoing)
			return incoming - outgoing, nil
		}
	} else {
		var bal GenericBalance
		if err := json.NewDecoder(resp.Body).Decode(&bal); err == nil {
			var val float64
			fmt.Sscanf(bal.Balance, "%f", &val)
			return val, nil
		}
	}

	return 0, nil
}

// Name returns the provider name
// Implements BlockchainProvider interface
func (p *TatumProvider) Name() string {
	return "Tatum"
}

// BroadcastTransaction broadcasts a signed raw transaction to the blockchain
// Implements BlockchainProvider interface
func (p *TatumProvider) BroadcastTransaction(currency, txHex string) (string, error) {
	log.Printf("[Tatum] Broadcasting raw tx for %s (length: %d chars)", currency, len(txHex))

	var endpoint string
	switch currency {
	case "BTC":
		endpoint = "bitcoin/broadcast"
	case "ETH", "USDT", "USDC", "DAI":
		endpoint = "ethereum/broadcast"
	case "BSC", "BNB":
		endpoint = "bsc/broadcast"
	case "SOL":
		endpoint = "solana/broadcast"
	case "TRX":
		endpoint = "tron/broadcast"
	case "MATIC":
		endpoint = "polygon/broadcast"
	case "LTC":
		endpoint = "litecoin/broadcast"
	case "DOGE":
		endpoint = "dogecoin/broadcast"
	case "BCH":
		endpoint = "bcash/broadcast"
	case "XRP":
		endpoint = "xrp/broadcast"
	default:
		// Try generic broadcast pattern
		endpoint = fmt.Sprintf("%s/broadcast", strings.ToLower(currency))
	}

	url := fmt.Sprintf("%s/%s", p.baseURL, endpoint)

	reqBody := map[string]string{
		"txData": txHex,
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("x-api-key", p.apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("tatum broadcast request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[Tatum] Broadcast error %d: %s", resp.StatusCode, string(bodyBytes))
		return "", fmt.Errorf("tatum broadcast error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	type BroadcastResponse struct {
		TxId string `json:"txId"`
	}

	var res BroadcastResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to parse broadcast response: %w", err)
	}

	log.Printf("[Tatum] ✅ Broadcast success! TxHash: %s", res.TxId)
	return res.TxId, nil
}
