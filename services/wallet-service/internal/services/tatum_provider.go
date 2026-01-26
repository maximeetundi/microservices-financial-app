package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
)

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

// Virtual Account Models
type TatumAccount struct {
	ID          string `json:"id"`
	Currency    string `json:"currency"`
	Active      bool   `json:"active"`
	Balance     TatumBalance `json:"balance"`
	AccountCode string `json:"accountCode,omitempty"`
	AccountingCurrency string `json:"accountingCurrency,omitempty"`
	Customer    interface{} `json:"customer,omitempty"`
	Frozen      bool   `json:"frozen"`
	Xpub        string `json:"xpub,omitempty"`
}

type TatumBalance struct {
	AccountBalance string `json:"accountBalance"`
	AvailableBalance string `json:"availableBalance"`
}

type CreateAccountRequest struct {
	Currency    string `json:"currency"`
	Xpub        string `json:"xpub,omitempty"`
	Customer    interface{} `json:"customer,omitempty"`
	Compliant   bool   `json:"compliant,omitempty"`
	AccountCode string `json:"accountCode,omitempty"`
	AccountingCurrency string `json:"accountingCurrency,omitempty"`
}

type DepositAddress struct {
	Address string `json:"address"`
	Currency string `json:"currency,omitempty"`
	DerivationKey int `json:"derivationKey,omitempty"`
	Xpub string `json:"xpub,omitempty"`
}

func (p *TatumProvider) CreateVirtualAccount(currency string, xpub string) (*TatumAccount, error) {
	url := fmt.Sprintf("%s/ledger/account", p.baseURL)
	
	reqBody := CreateAccountRequest{
		Currency: currency,
		Xpub:     xpub,
		Compliant: false, // For testing
	}
	
	jsonData, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("x-api-key", p.apiKey)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("tatum error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var account TatumAccount
	// Response from create might just be the ID or the full object depending on endpoint version
	// Tatum V3 usually returns the full object
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}
	
	return &account, nil
}

func (p *TatumProvider) GenerateDepositAddress(accountID string) (*DepositAddress, error) {
	// POST /ledger/account/{id}/address
	url := fmt.Sprintf("%s/ledger/account/%s/address", p.baseURL, accountID)
	
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("x-api-key", p.apiKey)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("tatum error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var address DepositAddress
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, err
	}
	
	return &address, nil
}

// Mock XPub generator for dev/test if no real XPUB provided
// In production, you would fetch this from a secure Vault/ENV
func (p *TatumProvider) GetMasterXpub(currency string) string {
	// For Bitcoin Testnet
	if currency == "BTC" {
		return "tpubD6NzVbkrYhZ4X8..." // Placeholder
	}
	// For Ethereum (Sepolia/Mainnet)
	if currency == "ETH" {
		return "xpub6F..." // Placeholder
	}
	return ""
}
