package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	log.Printf("[Tatum] Creating virtual account for currency: %s", currency)
	
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
		log.Printf("[Tatum] Error calling API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[Tatum] API error %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("tatum error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var account TatumAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}
	
	log.Printf("[Tatum] ✅ Created virtual account ID: %s for currency: %s", account.ID, currency)
	return &account, nil
}

func (p *TatumProvider) GenerateDepositAddress(accountID string) (*DepositAddress, error) {
	log.Printf("[Tatum] Generating deposit address for account: %s", accountID)
	
	url := fmt.Sprintf("%s/ledger/account/%s/address", p.baseURL, accountID)
	
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("x-api-key", p.apiKey)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("[Tatum] Error generating address: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[Tatum] API error %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("tatum error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var address DepositAddress
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, err
	}
	
	log.Printf("[Tatum] ✅ Generated deposit address: %s", address.Address)
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

// WebhookSubscription represents a Tatum webhook subscription
type WebhookSubscription struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Attr struct {
		AccountID string `json:"id"`
		URL       string `json:"url"`
	} `json:"attr"`
}

// SubscribeToAccountTransactions creates a webhook subscription for incoming transactions
// This is called automatically when creating a crypto wallet to receive deposit notifications
func (p *TatumProvider) SubscribeToAccountTransactions(accountID string, webhookURL string) (*WebhookSubscription, error) {
	log.Printf("[Tatum] Subscribing to transactions for account: %s", accountID)
	
	url := fmt.Sprintf("%s/subscription", p.baseURL)
	
	// Request body for subscription
	reqBody := map[string]interface{}{
		"type": "ACCOUNT_INCOMING_BLOCKCHAIN_TRANSACTION",
		"attr": map[string]string{
			"id":  accountID,
			"url": webhookURL,
		},
	}
	
	jsonData, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("x-api-key", p.apiKey)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("[Tatum] Error subscribing to webhook: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[Tatum] Webhook subscription error %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("tatum error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var subscription WebhookSubscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, err
	}
	
	log.Printf("[Tatum] ✅ Webhook subscription created: %s for account %s", subscription.ID, accountID)
	return &subscription, nil
}

// GetAccountBalance retrieves the balance from Tatum virtual account
func (p *TatumProvider) GetAccountBalance(accountID string) (*TatumBalance, error) {
	url := fmt.Sprintf("%s/ledger/account/%s/balance", p.baseURL, accountID)
	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", p.apiKey)
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("tatum error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var balance TatumBalance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, err
	}
	
	return &balance, nil
}

// WithdrawalRequest represents a request to send funds to blockchain
type WithdrawalRequest struct {
	SenderAccountId string `json:"senderAccountId"`
	Address         string `json:"address"`
	Amount          string `json:"amount"`
	Fee             string `json:"fee,omitempty"`
	Compliant       bool   `json:"compliant,omitempty"`
}

// WithdrawalResponse represents the response from a withdrawal
type WithdrawalResponse struct {
	ID   string `json:"id"`   // Internal withdrawal ID
	TxID string `json:"txId"` // Blockchain transaction hash
}

// SendToBlockchain transfers crypto from virtual account to external blockchain address
func (p *TatumProvider) SendToBlockchain(currency, senderAccountID, toAddress string, amount float64) (*WithdrawalResponse, error) {
	log.Printf("[Tatum] Sending %f %s from account %s to address %s", amount, currency, senderAccountID, toAddress)
	
	// Determine endpoint based on currency
	var endpoint string
	switch currency {
	case "BTC":
		endpoint = "bitcoin/transfer"
	case "ETH":
		endpoint = "ethereum/transfer"
	case "BSC":
		endpoint = "bsc/transfer"
	case "LTC":
		endpoint = "litecoin/transfer"
	default:
		return nil, fmt.Errorf("unsupported currency for auto-withdrawal: %s", currency)
	}
	
	url := fmt.Sprintf("%s/offchain/%s", p.baseURL, endpoint)
	
	reqBody := WithdrawalRequest{
		SenderAccountId: senderAccountID,
		Address:         toAddress,
		Amount:          fmt.Sprintf("%.8f", amount), // Ensure correct precision
		Compliant:       false,
	}
	
	jsonData, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("x-api-key", p.apiKey)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("[Tatum] Error sending extraction request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[Tatum] Withdrawal error %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("tatum error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var response WithdrawalResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	
	log.Printf("[Tatum] ✅ Withdrawal success! TxHash: %s", response.TxID)
	return &response, nil
}
