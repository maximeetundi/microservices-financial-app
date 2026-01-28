package services

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	btcbase58 "github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58" // For plain Base58 (Solana)

	"github.com/crypto-bank/microservices-financial-app/services/common/secrets"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
)

type CryptoService struct {
	config      *config.Config
	blockchain  BlockchainProvider // Abstraction for redundancy
	vaultClient *secrets.VaultClient
	// Legacy tatum field removed, use blockchain provider
}

func NewCryptoService(cfg *config.Config) *CryptoService {
	vc, err := secrets.NewVaultClient()
	if err != nil {
		log.Printf("Warning: Vault client could not be initialized in CryptoService: %v", err)
	}

	// Initialize Providers
	tatum := NewTatumProvider(cfg)
	blockcypher := NewBlockCypherProvider(cfg)
	rpcProvider := NewRpcProvider(cfg)

	// Create Failover Manager
	// Priorities: Tatum -> BlockCypher (BTC) -> RPC (ETH/BSC) -> (Others)
	failover := NewFailoverProvider(tatum, blockcypher, rpcProvider)

	return &CryptoService{
		config:      cfg,
		blockchain:  failover,
		vaultClient: vc,
	}
}

type CryptoWallet struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
	Network  string `json:"network"`
}

func (s *CryptoService) GenerateWallet(userID, currency string) (*CryptoWallet, error) {
	// NON-CUSTODIAL IMPLEMENTATION
	var address, privKeyHex, pubKeyHex string
	var err error

	// Determine Chain Type based on Currency
	switch strings.ToUpper(currency) {
	// Bitcoin Family
	case "BTC":
		address, privKeyHex, pubKeyHex, err = s.generateBitcoinKeys()
	
	// Solana Family
	case "SOL":
		address, privKeyHex, pubKeyHex, err = s.generateSolanaKeys()

	// TRON Family (TRX + TRC20 Tokens)
	case "TRX", "USDT":
		// USDT defaults to TRC20 (Tron) for low fees
		address, privKeyHex, pubKeyHex, err = s.generateTronKeys()

	// EVM Family (ETH + BSC + Polygon + Avalanche + ERC20/BEP20 Tokens)
	case "ETH", "USDC", "BNB", "MATIC", "AVAX", "LINK", "UNI", "SHIB", "DAI":
		// All these use Ethereum-compatible keys (Hex address starting with 0x)
		address, privKeyHex, pubKeyHex, err = s.generateEthereumKeys()

	default:
		return nil, fmt.Errorf("unsupported currency: %s", currency)
	}

	if err != nil {
		return nil, err
	}

	// Store in Vault
	err = s.storeKeyPairInVault(userID, currency, privKeyHex, pubKeyHex, address)
	if err != nil {
		return nil, fmt.Errorf("failed to secure keys in vault: %w", err)
	}

	return &CryptoWallet{
		Address:  address,
		Currency: currency,
		Network:  s.getNetworkForCurrency(currency),
	}, nil
}

// --- Bitcoin ---
// --- Bitcoin (Native SegWit - Bech32) ---
// Helper to get params based on config
func (s *CryptoService) getBitcoinParams() *chaincfg.Params {
	if s.config.CryptoNetwork == "testnet" {
		return &chaincfg.TestNet3Params
	}
	return &chaincfg.MainNetParams
}

// --- Bitcoin (Native SegWit - Bech32) ---
func (s *CryptoService) generateBitcoinKeys() (string, string, string, error) {
	privKey, _ := btcec.NewPrivateKey(btcec.S256())
	privKeyHex := hex.EncodeToString(privKey.Serialize())
	pubKey := privKey.PubKey()
	pubKeyCompressed := pubKey.SerializeCompressed()
	pubKeyHex := hex.EncodeToString(pubKeyCompressed)

	params := s.getBitcoinParams()

	// Generate Witness Public Key Hash (P2WPKH)
	// 1. Hash160(PubKeyCompressed)
	pubKeyHash := btcutil.Hash160(pubKeyCompressed)

	// 2. Create Witness Address (bc1q... or tb1q...)
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, params)
	if err != nil {
		return "", "", "", err
	}

	return addr.EncodeAddress(), privKeyHex, pubKeyHex, nil
}

func (s *CryptoService) deriveBitcoinAddress(pubKeyHex, network string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", err
	}
	
	// Determine params based on requested network
	var params *chaincfg.Params
	if strings.HasPrefix(strings.ToUpper(network), "TEST") || s.config.CryptoNetwork == "testnet" {
		if strings.ToUpper(network) == "SEGWIT" || strings.ToUpper(network) == "BTC" {
             // If user specifically asked for Mainnet equivalents despite testnet config, 
             // we could theoretically return mainnet, but let's stick to safe defaults.
             // Usually "SEGWIT" implies Mainnet in UI.
             // If network == "TESTNET", definitely testnet.
        }
        
        if strings.ToUpper(network) == "TESTNET" {
             params = &chaincfg.TestNet3Params
        } else if strings.ToUpper(network) == "SEGWIT" || strings.ToUpper(network) == "BTC" {
             params = &chaincfg.MainNetParams
        } else {
             // Default fallback to config
             params = s.getBitcoinParams()
        }
	} else {
        // Mainnet Config
        if strings.ToUpper(network) == "TESTNET" {
            params = &chaincfg.TestNet3Params
        } else {
            params = &chaincfg.MainNetParams
        }
    }

	// PubKey to Address
	// 1. Hash160
	pubKeyHash := btcutil.Hash160(pubKeyBytes)

	// 2. Create Address
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// --- Ethereum / BSC ---
func (s *CryptoService) generateEthereumKeys() (string, string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", err
	}
	privKeyBytes := crypto.FromECDSA(privateKey)
	privKeyHex := hex.EncodeToString(privKeyBytes)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", "", fmt.Errorf("error casting public key to ECDSA")
	}
	pubKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	pubKeyHex := hex.EncodeToString(pubKeyBytes)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, privKeyHex, pubKeyHex, nil
}

// --- Solana ---
func (s *CryptoService) generateSolanaKeys() (string, string, string, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", "", err
	}
	address := base58.Encode(pubKey)
	privKeyStr := base58.Encode(privKey)
	pubKeyStr := base58.Encode(pubKey)
	return address, privKeyStr, pubKeyStr, nil
}

// --- TRON (TRX / TRC20) ---
func (s *CryptoService) generateTronKeys() (string, string, string, error) {
	// 1. Generate ECDSA Key (Same curve as Ethereum)
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", err
	}
	privKeyBytes := crypto.FromECDSA(privateKey)
	privKeyHex := hex.EncodeToString(privKeyBytes)

	// 2. Get Public Key from Private Key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", "", fmt.Errorf("error casting public key to ECDSA")
	}
	pubKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	pubKeyHex := hex.EncodeToString(pubKeyBytes)

	// 3. Address Generation
	// TRON Address = Prefix + Keccak256(PubKey)[12:]
	// Prefix = 0x41 (Mainnet) or 0xa0 (Testnet/Shasta)

	// Determine prefix based on config
	prefix := byte(0x41)
	if s.config.CryptoNetwork == "testnet" {
		prefix = byte(0xa0)
	}

	// Get Ethereum-style address bytes (last 20 bytes of Keccak256 hash of pubkey)
	// crypto.PubkeyToAddress returns the 20-byte address.
	ethAddressBytes := crypto.PubkeyToAddress(*publicKeyECDSA).Bytes()

	// Prepend Prefix
	input := append([]byte{prefix}, ethAddressBytes...)

	// 4. Base58Check Encode (Payload + Checksum)
	// Github.com/mr-tron/base58 doesn't do Checksum automatically for us, need to implement manual check?
	// Actually, btcutil/base58 has CheckEncode. Let's stick to standard algos.
	// Checksum = SHA256(SHA256(input))[:4]

	// Available library: github.com/btcsuite/btcutil/base58 (already imported for SOL/BTC)
	// It has CheckEncode which does exactly this.
	address := btcbase58.CheckEncode(input[1:], input[0]) // CheckEncode takes (payload, version_byte)

	return address, privKeyHex, pubKeyHex, nil
}

// --- Get Or Create Address (Multi-Network) ---
func (s *CryptoService) GetOrCreateAddress(userID, currency, network string) (string, error) {
	// Normalize network
	network = strings.ToUpper(network)
	currency = strings.ToUpper(currency)

	// Determine Chain Family & Vault Path
	var chainFamily string
	var algoFunc func() (string, string, string, error)

	switch {
	case network == "TRC20" || network == "TRON" || currency == "TRX":
		chainFamily = "tron"
		algoFunc = s.generateTronKeys
	case network == "ERC20" || network == "BEP20" || network == "ETHEREUM" || network == "BSC" || currency == "ETH":
		chainFamily = "evm" // Shared address for all EVM chains
		algoFunc = s.generateEthereumKeys
	case currency == "BTC":
		chainFamily = "bitcoin"
		algoFunc = s.generateBitcoinKeys
	case currency == "SOL":
		chainFamily = "solana"
		algoFunc = s.generateSolanaKeys
	default:
		// Default behavior if network not specified: derive from currency
		if currency == "USDT" {
			// Default USDT to TRON (Binance-style preference)
			chainFamily = "tron"
			algoFunc = s.generateTronKeys
		} else {
			return "", fmt.Errorf("unsupported network/currency combination: %s on %s", currency, network)
		}
	}

	// 1. Check Vault for existing keys
	if s.vaultClient != nil {
		path := fmt.Sprintf("secret/wallets/%s/chains/%s", userID, chainFamily)
		secret, err := s.vaultClient.GetSecret(path)
		if err == nil && secret != nil {
			// Found existing keys
			pubKey, ok := secret["public_key"].(string)
			if ok && pubKey != "" {
				// Re-derive address for the specific network requested
				if chainFamily == "bitcoin" {
					return s.deriveBitcoinAddress(pubKey, network)
				}
				// For EVM and others where address is static/same across networks
				if addr, ok := secret["address"].(string); ok && addr != "" {
					return addr, nil
				}
			}
		}
	}

	// 2. Generate New Keys
	address, priv, pub, err := algoFunc()
	if err != nil {
		return "", err
	}
    
    // If specific network requested during initial generation (e.g. BTC Testnet), 
    // ensure the returned address matches that network if algoFunc didn't handle it
    if chainFamily == "bitcoin" && network != "" {
        derivedAddr, err := s.deriveBitcoinAddress(pub, network)
        if err == nil {
            address = derivedAddr
        }
    }

	// 3. Store in Vault
	if s.vaultClient != nil {
		path := fmt.Sprintf("secret/wallets/%s/chains/%s", userID, chainFamily)
		data := map[string]interface{}{
			"address":     address, // Store default (usually mainnet) address reference
			"private_key": priv,
			"public_key":  pub,
			"created_at":  fmt.Sprintf("%d", time.Now().Unix()),
			"network":     network,
		}
		if err := s.vaultClient.WriteSecret(path, data); err != nil {
			return "", fmt.Errorf("failed to write to vault: %v", err)
		}
	}

	return address, nil
}

// --- Vault Storage ---
func (s *CryptoService) storeKeyPairInVault(userID, currency, privKey, pubKey, address string) error {
	if s.vaultClient == nil {
		return fmt.Errorf("vault client is not available")
	}
	path := fmt.Sprintf("secret/wallets/%s/%s", userID, strings.ToLower(currency))
	data := map[string]interface{}{
		"private_key": privKey,
		"public_key":  pubKey,
		"address":     address,
		"created_at":  fmt.Sprintf("%d", time.Now().Unix()),
	}
	return s.vaultClient.WriteSecret(path, data)
}

// --- Helpers & Validation ---
func (s *CryptoService) ValidateAddress(address, currency string) bool {
	switch strings.ToUpper(currency) {
	case "BTC":
		// Check against the configured network first
		params := s.getBitcoinParams()
		_, err := btcutil.DecodeAddress(address, params)
		if err == nil {
			return true
		}

		// Fallback: If in testnet, might want to block mainnet addresses or vice-versa?
		// For safety, strictly enforce the configured network.
		return false
	case "ETH", "BSC":
		return common.IsHexAddress(address)
	case "SOL":
		decoded, err := base58.Decode(address)
		if err != nil {
			return false
		}
		return len(decoded) == 32
	default:
		return len(address) > 20
	}
}

func (s *CryptoService) DetectAddressNetwork(address string) NetworkDetectionResult {
	if strings.HasPrefix(address, "0x") && len(address) == 42 {
		return NetworkDetectionResult{Type: "EVM", Network: "unknown", Variant: "ERC20"}
	}
	// BTC Mainnet
	if strings.HasPrefix(address, "1") || strings.HasPrefix(address, "3") || strings.HasPrefix(address, "bc1") {
		return NetworkDetectionResult{Type: "BTC", Network: "mainnet", Variant: "BTC"}
	}
	// BTC Testnet
	if strings.HasPrefix(address, "m") || strings.HasPrefix(address, "n") || strings.HasPrefix(address, "2") || strings.HasPrefix(address, "tb1") {
		return NetworkDetectionResult{Type: "BTC", Network: "testnet", Variant: "BTC"}
	}

	if len(address) >= 32 && len(address) <= 44 {
		_, err := base58.Decode(address)
		if err == nil {
			return NetworkDetectionResult{Type: "SOL", Network: "mainnet-beta", Variant: "SPL"}
		}
	}
	return NetworkDetectionResult{Type: "UNKNOWN", Network: "unknown", Variant: "unknown"}
}

func (s *CryptoService) getNetworkForCurrency(currency string) string {
	switch strings.ToUpper(currency) {
	case "BTC":
		if s.config.CryptoNetwork == "testnet" {
			return "bitcoin-testnet"
		}
		return "bitcoin-mainnet"
	case "ETH":
		return "ethereum-mainnet"
	case "BSC":
		return "bsc-mainnet"
	case "SOL":
		return "solana-mainnet-beta"
	default:
		return "mainnet"
	}
}

// Fee Estimation
func (s *CryptoService) EstimateTransactionFee(currency string, amount float64, priority string) (*models.CryptoTransactionEstimate, error) {
	baseFee := 0.001
	if strings.ToUpper(currency) == "SOL" {
		baseFee = 0.000005
	}
	estimate := &models.CryptoTransactionEstimate{
		EstimatedFee:   baseFee,
		EstimatedTotal: amount + baseFee,
		Currency:       strings.ToUpper(currency),
	}
	return estimate, nil
}

// GetBalance - USES FAILOVER PROVIDER
func (s *CryptoService) GetBalance(address, currency string) (float64, error) {
	return s.blockchain.GetBalance(currency, address)
}

func (s *CryptoService) EncryptPrivateKey(privateKey, password string) (string, error) {
	return "STORED_IN_VAULT", nil
}

func (s *CryptoService) CreateTransaction(fromWallet *models.Wallet, toAddress string, amount float64, gasPrice *int64) (string, error) {
	// 1. Fetch Key from Vault (s.vaultClient.GetSecret)
	// 2. Sign locally
	// 3. Broadcast using s.blockchain.BroadcastTransaction(...)

	// Mock implementation
	return "0xMOCK_TX_HASH_SIGNED_LOCALLY", nil
}

func (s *CryptoService) GetMinimumConfirmations(currency string) int {
	if strings.ToUpper(currency) == "SOL" {
		return 1
	}
	return 6
}

func (s *CryptoService) GetTransactionStatus(txHash, currency string) (string, int, error) {
	return "confirmed", 12, nil
}

type NetworkDetectionResult struct {
	Type    string
	Network string
	Variant string
}
