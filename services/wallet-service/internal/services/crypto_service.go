package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
)

type CryptoService struct {
	config *config.Config
	tatum  *TatumProvider
}

func NewCryptoService(cfg *config.Config) *CryptoService {
	return &CryptoService{
		config: cfg,
		tatum:  NewTatumProvider(cfg),
	}
}

type CryptoWallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Currency   string `json:"currency"`
	Network    string `json:"network"`
	// Tatum specific
	AccountID  string `json:"account_id,omitempty"`
}

func (s *CryptoService) GenerateWallet(currency string) (*CryptoWallet, error) {
	// If Tatum API Key is set, use Tatum
	if s.config.TatumAPIKey != "" {
		return s.generateTatumWallet(currency)
	}

	// Fallback to Mock for Dev without API Key
	fmt.Println("[CryptoService] WARNING: No TATUM_API_KEY found. Using Mock Wallet Generation.")
	switch strings.ToUpper(currency) {
	case "BTC":
		return s.generateBitcoinWallet()
	case "ETH":
		return s.generateEthereumWallet()
	case "BSC":
		return s.generateBSCWallet()
	default:
		return nil, fmt.Errorf("unsupported currency: %s", currency)
	}
}

// generateTatumWallet creates a Virtual Account and a Deposit Address
func (s *CryptoService) generateTatumWallet(currency string) (*CryptoWallet, error) {
	// 1. Create Virtual Account
	// Note: We don't provide xpub here in this simplified version, letting Tatum manage it (Custodial)
	acc, err := s.tatum.CreateVirtualAccount(currency, "") 
	if err != nil {
		return nil, fmt.Errorf("failed to create tatum account: %w", err)
	}

	// 2. Generate Deposit Address using the account ID
	// Tatum API returns the account ID in the `id` field
	addr, err := s.tatum.GenerateDepositAddress(acc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	// 3. Subscribe to incoming transactions webhook
	// This ensures we get notified when deposits are made
	webhookURL := s.config.BaseURL + "/api/v1/webhooks/tatum"
	_, err = s.tatum.SubscribeToAccountTransactions(acc.ID, webhookURL)
	if err != nil {
		// Log but don't fail - wallet is still created, webhook can be added later
		fmt.Printf("[CryptoService] Warning: Failed to subscribe to webhook: %v\n", err)
	}

	return &CryptoWallet{
		Address:    addr.Address,
		PrivateKey: "", // Custodial - No Private Key exposed
		PublicKey:  "",
		Currency:   currency,
		Network:    "tatum-virtual",
		AccountID:  acc.ID, // Store the Tatum account ID for future reference
	}, nil
}

func (s *CryptoService) ValidateAddress(address, currency string) bool {
	// Basic regex validation still useful
	switch strings.ToUpper(currency) {
	case "BTC":
		return s.validateBitcoinAddress(address)
	case "ETH", "BSC":
		return s.validateEthereumAddress(address)
	default:
		return true // Allow others if unknown
	}
}

func (s *CryptoService) EstimateTransactionFee(currency string, amount float64, priority string) (*models.CryptoTransactionEstimate, error) {
	baseFee := s.config.NetworkFees[strings.ToUpper(currency)]
	if baseFee == 0 {
		baseFee = 0.001 // Default fee
	}

	// Adjust fee based on priority
	multiplier := 1.0
	switch priority {
	case "low":
		multiplier = 0.5
	case "normal":
		multiplier = 1.0
	case "high":
		multiplier = 2.0
	case "urgent":
		multiplier = 3.0
	default:
		multiplier = 1.0
	}

	estimatedFee := baseFee * multiplier
	estimatedTotal := amount + estimatedFee

	estimate := &models.CryptoTransactionEstimate{
		EstimatedFee:   estimatedFee,
		EstimatedTotal: estimatedTotal,
		Currency:       strings.ToUpper(currency),
	}

	// Add gas-specific fields for Ethereum-based currencies
	if strings.ToUpper(currency) == "ETH" || strings.ToUpper(currency) == "BSC" {
		gasPrice := int64(20000000000) // 20 gwei default
		gasLimit := int64(21000)       // Standard transfer

		if priority == "high" {
			gasPrice = int64(40000000000) // 40 gwei
		} else if priority == "urgent" {
			gasPrice = int64(60000000000) // 60 gwei
		}

		estimate.GasPrice = &gasPrice
		estimate.GasLimit = &gasLimit
	}

	return estimate, nil
}

func (s *CryptoService) CreateTransaction(fromWallet *models.Wallet, toAddress string, amount float64, gasPrice *int64) (string, error) {
	// If Tatum API Key is set, use Tatum
	if s.config.TatumAPIKey != "" {
		// Use Tatum to send funds (Off-chain transfer to blockchain address)
		// fromWallet.ExternalID MUST hold the Tatum Account ID
		if fromWallet.ExternalID == "" {
			return "", fmt.Errorf("wallet %s has no linked tatum account", fromWallet.ID)
		}

		resp, err := s.tatum.SendToBlockchain(fromWallet.Currency, fromWallet.ExternalID, toAddress, amount)
		if err != nil {
			return "", fmt.Errorf("tatum transfer failed: %w", err)
		}
		
		return resp.TxID, nil
	}

	// Mock Logic
	currency := strings.ToUpper(fromWallet.Currency)
	
	// Validate destination address
	if !s.ValidateAddress(toAddress, currency) {
		return "", fmt.Errorf("invalid destination address")
	}

	// Check balance
	if fromWallet.Balance < amount {
		return "", fmt.Errorf("insufficient balance")
	}

	// Generate transaction hash (simplified)
	txHash := s.generateTransactionHash(fromWallet.WalletAddress, &toAddress, amount, currency)

	return txHash, nil
}

func (s *CryptoService) GetTransactionStatus(txHash, currency string) (string, int, error) {
	// Mock implementation
	// In production, you would query the blockchain
	
	// Simulate different transaction states
	switch len(txHash) % 4 {
	case 0:
		return "pending", 0, nil
	case 1:
		return "confirmed", 1, nil
	case 2:
		return "confirmed", 6, nil
	default:
		return "confirmed", 15, nil
	}
}

func (s *CryptoService) GetBalance(address, currency string) (float64, error) {
	// Mock implementation
	// In production, query actual blockchain
	
	// Return a mock balance based on address hash
	hash := sha256.Sum256([]byte(address + currency))
	balance := float64(hash[0]) / 100.0
	
	return balance, nil
}

// Private methods for wallet generation

func (s *CryptoService) generateBitcoinWallet() (*CryptoWallet, error) {
	// Simplified Bitcoin wallet generation
	// In production, use proper Bitcoin libraries like btcd/btcutil
	
	privateKey, err := s.generatePrivateKey()
	if err != nil {
		return nil, err
	}

	publicKey := s.generatePublicKeyFromPrivate(privateKey)
	address := s.generateBitcoinAddress(publicKey)

	return &CryptoWallet{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Currency:   "BTC",
		Network:    "mainnet",
	}, nil
}

func (s *CryptoService) generateEthereumWallet() (*CryptoWallet, error) {
	// Simplified Ethereum wallet generation
	// In production, use proper Ethereum libraries like go-ethereum
	
	privateKey, err := s.generatePrivateKey()
	if err != nil {
		return nil, err
	}

	publicKey := s.generatePublicKeyFromPrivate(privateKey)
	address := s.generateEthereumAddress(publicKey)

	return &CryptoWallet{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Currency:   "ETH",
		Network:    "mainnet",
	}, nil
}

func (s *CryptoService) generateBSCWallet() (*CryptoWallet, error) {
	// BSC uses same address format as Ethereum
	wallet, err := s.generateEthereumWallet()
	if err != nil {
		return nil, err
	}
	
	wallet.Currency = "BSC"
	wallet.Network = "bsc-mainnet"
	
	return wallet, nil
}

func (s *CryptoService) generatePrivateKey() (string, error) {
	// Generate 32-byte private key
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *CryptoService) generatePublicKeyFromPrivate(privateKey string) string {
	// Simplified public key generation
	// In production, use proper elliptic curve cryptography
	hash := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(hash[:])
}

func (s *CryptoService) generateBitcoinAddress(publicKey string) string {
	// Simplified Bitcoin address generation
	// In production, use proper Base58Check encoding
	hash := sha256.Sum256([]byte(publicKey))
	return "1" + hex.EncodeToString(hash[:20])
}

func (s *CryptoService) generateEthereumAddress(publicKey string) string {
	// Simplified Ethereum address generation
	// In production, use proper Keccak256 hashing
	hash := sha256.Sum256([]byte(publicKey))
	return "0x" + hex.EncodeToString(hash[:20])
}

func (s *CryptoService) validateBitcoinAddress(address string) bool {
	// Simplified Bitcoin address validation
	if len(address) < 26 || len(address) > 35 {
		return false
	}
	return strings.HasPrefix(address, "1") || strings.HasPrefix(address, "3") || strings.HasPrefix(address, "bc1")
}

func (s *CryptoService) validateEthereumAddress(address string) bool {
	// Simplified Ethereum address validation
	if len(address) != 42 {
		return false
	}
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	// Check if hex
	_, err := hex.DecodeString(address[2:])
	return err == nil
}

func (s *CryptoService) generateTransactionHash(fromAddress, toAddress *string, amount float64, currency string) string {
	// Generate a mock transaction hash
	data := fmt.Sprintf("%s-%s-%f-%s-%d", 
		*fromAddress, *toAddress, amount, currency, 
		time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *CryptoService) EncryptPrivateKey(privateKey, password string) (string, error) {
	// If empty key (Custodial), return empty
	if privateKey == "" {
		return "", nil
	}

	// Simplified encryption
	// In production, use proper encryption like AES-256-GCM
	key := sha256.Sum256([]byte(password))
	
	// XOR encryption (for demo only)
	encrypted := make([]byte, len(privateKey))
	for i, b := range []byte(privateKey) {
		encrypted[i] = b ^ key[i%32]
	}
	
	return hex.EncodeToString(encrypted), nil
}

func (s *CryptoService) DecryptPrivateKey(encryptedKey, password string) (string, error) {
	if encryptedKey == "" {
		return "", nil
	}

	// Simplified decryption
	key := sha256.Sum256([]byte(password))
	
	encrypted, err := hex.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}
	
	// XOR decryption
	decrypted := make([]byte, len(encrypted))
	for i, b := range encrypted {
		decrypted[i] = b ^ key[i%32]
	}
	
	return string(decrypted), nil
}

func (s *CryptoService) GetMinimumConfirmations(currency string) int {
	if conf, exists := s.config.MinConfirmations[strings.ToUpper(currency)]; exists {
		return conf
	}
	return 6 // Default minimum confirmations
}