package services

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/security"
)

type CryptoService struct {
	config       *config.Config
	blockchain   BlockchainProvider // FailoverProvider (Tatum -> BlockCypher -> RPC)
	vaultClient  *secrets.VaultClient
	systemConfig *SystemConfigService
}

func NewCryptoService(cfg *config.Config, systemConfig *SystemConfigService) *CryptoService {
	vc, err := secrets.NewVaultClient()
	if err != nil {
		log.Printf("Warning: Vault client could not be initialized in CryptoService: %v", err)
	}

	// Initialize Providers for FailoverProvider
	tatum := NewTatumProvider(cfg)
	blockcypher := NewBlockCypherProvider(cfg)
	rpcProvider := NewRpcProvider(cfg)

	// Create Failover Manager - tries providers in order until one succeeds
	failover := NewFailoverProvider(tatum, blockcypher, rpcProvider)

	return &CryptoService{
		config:       cfg,
		blockchain:   failover,
		vaultClient:  vc,
		systemConfig: systemConfig,
	}
}

type CryptoWallet struct {
	Address             string `json:"address"`
	Currency            string `json:"currency"`
	Network             string `json:"network"`
	EncryptedPrivateKey string `json:"-"` // Internal use only
}

// isTestnet checks dynamic system config first, falls back to static config
func (s *CryptoService) isTestnet() bool {
	if s.systemConfig != nil {
		if s.systemConfig.IsTestnetEnabled() {
			return true
		}
	}
	return s.config.CryptoNetwork == "testnet"
}

// GetAvailableNetworks returns the list of valid networks for a currency based on system config
func (s *CryptoService) GetAvailableNetworks(currency string) []string {
	isTestnet := s.isTestnet()
	currency = strings.ToUpper(currency)

	switch currency {
	case "USDT", "USDC":
		if isTestnet {
			return []string{"ERC20", "BEP20", "TRC20", "SEPOLIA (Testnet)", "BSC-TESTNET (Testnet)", "SHASTA (Testnet)"}
		}
		return []string{"ERC20", "BEP20", "TRC20"}
	case "ETH":
		if isTestnet {
			return []string{"ERC20", "BEP20", "SEPOLIA (Testnet)", "GOERLI (Testnet)"}
		}
		return []string{"ERC20", "BEP20"}
	case "BTC":
		if isTestnet {
			return []string{"BTC", "SEGWIT", "TESTNET (Testnet)"}
		}
		return []string{"BTC", "SEGWIT"}
	case "TRX":
		if isTestnet {
			return []string{"TRC20", "SHASTA (Testnet)"}
		}
		return []string{"TRC20"}
	case "SOL":
		if isTestnet {
			return []string{"SOLANA", "DEVNET (Testnet)"}
		}
		return []string{"SOLANA"}
	case "BNB":
		if isTestnet {
			return []string{"BEP20", "BSC-TESTNET (Testnet)"}
		}
		return []string{"BEP20"}
	}

	// Default for others (MATIC, AVAX, etc)
	if isTestnet {
		return []string{"MAINNET", "TESTNET (Testnet)"}
	}
	// For mainnet-only or unsupported, maybe empty or just MAINNET?
	// Sticking to frontend logic:
	return []string{}
}

// ==================== Hybrid Deposit System Helpers ====================

// RequiresMemo returns true for currencies that support memo/tag deposits
// These currencies use a single hot wallet address + unique memo per user
func (s *CryptoService) RequiresMemo(currency string) bool {
	currency = strings.ToUpper(currency)
	switch currency {
	case "XRP", "XLM", "TON", "EOS", "ATOM", "HBAR":
		return true
	default:
		return false
	}
}

// GenerateDepositMemo generates a unique deposit memo/tag for a user
// XRP uses numeric destination tags, others use alphanumeric
func (s *CryptoService) GenerateDepositMemo(userID, currency string) string {
	currency = strings.ToUpper(currency)

	// Create a deterministic but unique memo based on userID
	// Hash the userID to create a consistent memo
	hash := sha256.Sum256([]byte(userID + "_" + currency + "_deposit"))

	switch currency {
	case "XRP":
		// XRP destination tags are 32-bit unsigned integers
		// Use first 4 bytes of hash, ensure it's non-zero
		tag := uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])
		if tag == 0 {
			tag = 1 // Destination tag 0 is often invalid
		}
		return fmt.Sprintf("%d", tag)
	default:
		// For XLM, TON, etc. use alphanumeric (first 12 chars of hex hash)
		return hex.EncodeToString(hash[:6])
	}
}

// GetDepositInfo returns the deposit address and memo for a currency
// For memo-based currencies, returns the hot wallet address + unique memo
// For address-based currencies, returns the user's individual deposit address
func (s *CryptoService) GetDepositInfo(currency, userAddress, userID string) (address string, memo string, usesMemo bool) {
	if s.RequiresMemo(currency) {
		// Get hot wallet address for this currency from platform accounts
		// In production, this would query the platform_crypto_wallets table
		// For now, return a placeholder that should be configured
		hotWalletAddress := s.getHotWalletAddress(currency)
		depositMemo := s.GenerateDepositMemo(userID, currency)
		return hotWalletAddress, depositMemo, true
	}

	// For non-memo currencies, use the user's individual address
	return userAddress, "", false
}

// getHotWalletAddress returns the platform hot wallet address for a currency
// This should be configured in the database/config
func (s *CryptoService) getHotWalletAddress(currency string) string {
	// In production, this would query platform_crypto_wallets table
	// For now, return environment-based or configured addresses
	// TODO: Integrate with PlatformAccountService
	return fmt.Sprintf("HOT_WALLET_%s_NOT_CONFIGURED", currency)
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

	// TON (The Open Network)
	case "TON":
		address, privKeyHex, pubKeyHex, err = s.generateTonKeys()

	// XRP (Ripple)
	case "XRP":
		address, privKeyHex, pubKeyHex, err = s.generateXrpKeys()

	// Litecoin / Dogecoin / BCH (Bitcoin-like)
	case "LTC", "DOGE", "BCH":
		// For now, using Bitcoin-like generation (compatible curve)
		address, privKeyHex, pubKeyHex, err = s.generateBitcoinKeys() // Simplified for prototype

	// EVM Family (ETH + BSC + Polygon + Avalanche + ERC20/BEP20 Tokens)
	case "ETH", "USDC", "BNB", "MATIC", "AVAX", "LINK", "UNI", "SHIB", "DAI", "CRO", "FTM":
		// All these use Ethereum-compatible keys (Hex address starting with 0x)
		address, privKeyHex, pubKeyHex, err = s.generateEthereumKeys()

	default:
		return nil, fmt.Errorf("unsupported currency: %s", currency)
	}

	if err != nil {
		return nil, err
	}

	// Encrypt Private Key for DB Storage (Envelope Encryption)
	encryptedKey, err := s.EncryptPrivateKey(privKeyHex, "")
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	return &CryptoWallet{
		Address:             address,
		Currency:            currency,
		Network:             s.getNetworkForCurrency(currency),
		EncryptedPrivateKey: encryptedKey,
	}, nil
}

// --- Bitcoin ---
// --- Bitcoin (Native SegWit - Bech32) ---
// Helper to get params based on config
func (s *CryptoService) getBitcoinParams() *chaincfg.Params {
	if s.isTestnet() {
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

	isTest := s.isTestnet()
	log.Printf("[CryptoService] Generating Bitcoin Keys. IsTestnet=%v (Config: %s)", isTest, s.config.CryptoNetwork)
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
	if strings.HasPrefix(strings.ToUpper(network), "TEST") || s.isTestnet() {
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

// --- TON (The Open Network) ---
func (s *CryptoService) generateTonKeys() (string, string, string, error) {
	// TON uses Ed25519
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", "", err
	}

	// Standard Ed25519 keys
	privKeyHex := hex.EncodeToString(privKey)
	pubKeyHex := hex.EncodeToString(pubKey)

	// TON Address Generation (Friendly format)
	// Format: [flags:1][workchain:1][hash:32][crc16:2]
	// Total 36 bytes, encoded in Base64url

	// 1. Hash Public Key (SHA256) - 32 bytes
	hashBytes := sha256.Sum256(pubKey)

	// 2. Build raw address
	// Flags: 0x11 = bounceable, 0x51 = non-bounceable, 0x91 = test bounceable
	flags := byte(0x11)     // Bounceable mainnet
	workchain := byte(0x00) // Basechain (workchain 0)

	if s.isTestnet() {
		flags = byte(0x91) // Test bounceable
	}

	// Build 34-byte address buffer
	addressData := make([]byte, 34)
	addressData[0] = flags
	addressData[1] = workchain
	copy(addressData[2:], hashBytes[:])

	// 3. Calculate CRC16-CCITT checksum
	crc := s.crc16CCITT(addressData)

	// 4. Append CRC16 (big-endian)
	fullAddress := append(addressData, byte(crc>>8), byte(crc&0xff))

	// 5. Encode to Base64url (RFC 4648)
	address := base64.RawURLEncoding.EncodeToString(fullAddress)

	// Add standard EQ prefix indicator (the full address contains it in encoding)
	// TON friendly addresses look like: EQDtFpEwcFAEcRe5mLVh2N6C0x-P_ZM...

	return address, privKeyHex, pubKeyHex, nil
}

// crc16CCITT calculates CRC16-CCITT checksum (used by TON)
func (s *CryptoService) crc16CCITT(data []byte) uint16 {
	crc := uint16(0x0000)
	polynomial := uint16(0x1021)

	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

// --- XRP (Ripple) ---
// XRP uses secp256k1 with a specific Base58 alphabet (rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz)
func (s *CryptoService) generateXrpKeys() (string, string, string, error) {
	// 1. Generate secp256k1 key pair
	privKey, _ := btcec.NewPrivateKey(btcec.S256())
	privKeyHex := hex.EncodeToString(privKey.Serialize())
	pubKey := privKey.PubKey()
	pubKeyCompressed := pubKey.SerializeCompressed()
	pubKeyHex := hex.EncodeToString(pubKeyCompressed)

	// 2. XRP Address = Base58Check(AccountID)
	// AccountID = RIPEMD160(SHA256(PubKeyCompressed))
	accountID := btcutil.Hash160(pubKeyCompressed) // 20 bytes

	// 3. Build payload: [prefix:1][accountID:20][checksum:4]
	// XRP account prefix = 0x00
	payload := make([]byte, 25)
	payload[0] = 0x00 // Account prefix
	copy(payload[1:21], accountID)

	// 4. Calculate checksum: first 4 bytes of SHA256(SHA256(prefix + accountID))
	hash1 := sha256.Sum256(payload[:21])
	hash2 := sha256.Sum256(hash1[:])
	copy(payload[21:], hash2[:4])

	// 5. Encode using Ripple's Base58 alphabet
	address := s.encodeXrpBase58(payload)

	return address, privKeyHex, pubKeyHex, nil
}

// encodeXrpBase58 encodes bytes using Ripple's Base58 alphabet
// Ripple alphabet: rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz
func (s *CryptoService) encodeXrpBase58(input []byte) string {
	const alphabet = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"

	// Count leading zeros
	leadingZeros := 0
	for _, b := range input {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	// Convert to big integer and encode
	var result []byte
	for len(input) > 0 {
		var carry int
		var newInput []byte
		for _, b := range input {
			carry = carry*256 + int(b)
			if len(newInput) > 0 || carry >= 58 {
				newInput = append(newInput, byte(carry/58))
			}
			carry = carry % 58
		}
		result = append([]byte{alphabet[carry]}, result...)
		input = newInput
	}

	// Add leading 'r' characters for leading zeros (in Ripple alphabet, 'r' = 0)
	for i := 0; i < leadingZeros; i++ {
		result = append([]byte{'r'}, result...)
	}

	return string(result)
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
	if s.isTestnet() {
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
		if s.isTestnet() {
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

// GetBlockchainBalance fetches the current blockchain balance for an address
// This is an alias of GetBalance for semantic clarity in admin operations
func (s *CryptoService) GetBlockchainBalance(currency, address string) (float64, error) {
	return s.blockchain.GetBalance(currency, address)
}

// EncryptPrivateKey encrypts a private key for database storage
// This provides defense-in-depth alongside the external vault
func (s *CryptoService) EncryptPrivateKey(privateKey, _ string) (string, error) {
	vault := security.GetVaultService()
	return vault.EncryptPrivateKey(privateKey)
}

// GetEncryptedKeyWithHash returns an encrypted private key and its hash for database storage
func (s *CryptoService) GetEncryptedKeyWithHash(privateKey string) (encryptedKey, keyHash string, err error) {
	vault := security.GetVaultService()

	encryptedKey, err = vault.EncryptPrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt private key: %w", err)
	}

	keyHash = security.HashData(privateKey)
	return encryptedKey, keyHash, nil
}

// DecryptPrivateKey decrypts a stored private key
func (s *CryptoService) DecryptPrivateKey(encryptedKey string) (string, error) {
	vault := security.GetVaultService()
	return vault.DecryptPrivateKey(encryptedKey)
}

func (s *CryptoService) CreateTransaction(fromWallet *models.Wallet, toAddress string, amount float64, gasPrice *int64) (string, error) {
	// Validate inputs
	if fromWallet == nil {
		return "", fmt.Errorf("source wallet is nil")
	}
	if toAddress == "" {
		return "", fmt.Errorf("destination address is required")
	}
	if amount <= 0 {
		return "", fmt.Errorf("amount must be positive")
	}

	currency := strings.ToUpper(fromWallet.Currency)
	log.Printf("[CryptoService] CreateTransaction: %f %s from wallet %s to %s", amount, currency, fromWallet.ID, toAddress)

	// Validate address for non-memo currencies that need UTXOs (BTC etc)
	if fromWallet.WalletAddress == nil || *fromWallet.WalletAddress == "" {
		return "", fmt.Errorf("source wallet %s has no address", fromWallet.ID)
	}

	// Self-custody mode: Sign locally and broadcast via FailoverProvider (Tatum -> BlockCypher -> RPC)
	// The wallet must have an encrypted private key stored
	if fromWallet.PrivateKeyEncrypted == nil || *fromWallet.PrivateKeyEncrypted == "" {
		return "", fmt.Errorf("wallet %s has no encrypted private key - cannot sign transaction", fromWallet.ID)
	}

	log.Printf("[CryptoService] Using self-custody signing for %s", currency)

	// Step 1: Decrypt private key
	privateKey, err := s.DecryptStoredPrivateKey(*fromWallet.PrivateKeyEncrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// Step 2: Build and sign transaction based on currency type
	var signedTxHex string
	switch {
	case currency == "BTC" || currency == "LTC" || currency == "BCH" || currency == "DOGE":
		signedTxHex, err = s.signBitcoinLikeTx(privateKey, *fromWallet.WalletAddress, toAddress, amount, currency)
	case currency == "ETH" || currency == "USDT" || currency == "USDC" || currency == "BNB":
		signedTxHex, err = s.signEthereumTx(privateKey, toAddress, amount, currency, gasPrice)
	case currency == "TRX":
		signedTxHex, err = s.signTronTx(privateKey, toAddress, amount)
	case currency == "SOL":
		signedTxHex, err = s.signSolanaTx(privateKey, toAddress, amount)
	case currency == "XRP":
		signedTxHex, err = s.signXrpTx(privateKey, toAddress, amount)
	case currency == "TON":
		signedTxHex, err = s.signTonTx(privateKey, toAddress, amount)
	default:
		return "", fmt.Errorf("direct signing not implemented for currency: %s", currency)
	}

	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Step 3: Broadcast via FailoverProvider (tries Tatum -> BlockCypher -> RPC)
	log.Printf("[CryptoService] Broadcasting signed transaction via FailoverProvider")
	txHash, err := s.blockchain.BroadcastTransaction(currency, signedTxHex)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	log.Printf("[CryptoService] Transaction broadcast successful. TxHash: %s", txHash)
	return txHash, nil
}

// Signing stubs - These would need full implementation with proper blockchain libraries
// For now, they return errors indicating the feature is not yet implemented

func (s *CryptoService) signBitcoinLikeTx(privateKey string, fromAddr string, toAddr string, amount float64, currency string) (string, error) {
	// TODO: Implement using btcd/wire and btcutil for transaction building
	// Requires: UTXO lookup, fee estimation, script building
	return "", fmt.Errorf("Bitcoin-like signing not yet implemented - use Tatum virtual accounts for %s", currency)
}

func (s *CryptoService) signEthereumTx(privateKey, toAddr string, amount float64, currency string, gasPrice *int64) (string, error) {
	// TODO: Implement using go-ethereum for transaction signing
	// This is the most straightforward to implement with crypto.Sign
	return "", fmt.Errorf("Ethereum signing not yet implemented - use Tatum virtual accounts for %s", currency)
}

func (s *CryptoService) signTronTx(privateKey, toAddr string, amount float64) (string, error) {
	// TODO: Implement using TRON protobuf transaction format
	return "", fmt.Errorf("TRON signing not yet implemented - use Tatum virtual accounts")
}

func (s *CryptoService) signSolanaTx(privateKey, toAddr string, amount float64) (string, error) {
	// TODO: Implement using solana-go library
	return "", fmt.Errorf("Solana signing not yet implemented - use Tatum virtual accounts")
}

func (s *CryptoService) signXrpTx(privateKey, toAddr string, amount float64) (string, error) {
	// TODO: Implement using XRP transaction format
	return "", fmt.Errorf("XRP signing not yet implemented - use Tatum virtual accounts")
}

func (s *CryptoService) signTonTx(privateKey, toAddr string, amount float64) (string, error) {
	// TODO: Implement using TON cell serialization
	return "", fmt.Errorf("TON signing not yet implemented - use Tatum virtual accounts")
}

// DecryptStoredPrivateKey decrypts an encrypted private key from database storage
func (s *CryptoService) DecryptStoredPrivateKey(encryptedKey string) (string, error) {
	vault := security.GetVaultService()
	return vault.DecryptPrivateKey(encryptedKey)
}

// CreateTransactionFromPlatformWallet creates a blockchain transaction from a platform hot wallet
// This is used for external sends where the user's balance is debited in DB
// but the actual blockchain transaction comes from the platform's hot wallet
func (s *CryptoService) CreateTransactionFromPlatformWallet(
	hotWallet *models.PlatformCryptoWallet,
	toAddress string,
	amount float64,
	gasPrice *int64,
) (string, error) {
	currency := strings.ToUpper(hotWallet.Currency)
	log.Printf("[CryptoService] CreateTransactionFromPlatformWallet: %.8f %s to %s", amount, currency, toAddress)

	if hotWallet.EncryptedPrivateKey == "" {
		return "", fmt.Errorf("platform wallet %s has no encrypted private key", hotWallet.ID)
	}

	// Decrypt the platform wallet's private key
	privateKey, err := s.DecryptStoredPrivateKey(hotWallet.EncryptedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt platform wallet private key: %w", err)
	}

	// Build and sign the transaction
	var signedTxHex string
	switch {
	case currency == "BTC" || currency == "LTC" || currency == "BCH" || currency == "DOGE":
		signedTxHex, err = s.signBitcoinLikeTx(privateKey, hotWallet.Address, toAddress, amount, currency)
	case currency == "ETH" || currency == "USDT" || currency == "USDC" || currency == "BNB":
		signedTxHex, err = s.signEthereumTx(privateKey, toAddress, amount, currency, gasPrice)
	case currency == "TRX":
		signedTxHex, err = s.signTronTx(privateKey, toAddress, amount)
	case currency == "SOL":
		signedTxHex, err = s.signSolanaTx(privateKey, toAddress, amount)
	case currency == "XRP":
		signedTxHex, err = s.signXrpTx(privateKey, toAddress, amount)
	case currency == "TON":
		signedTxHex, err = s.signTonTx(privateKey, toAddress, amount)
	default:
		return "", fmt.Errorf("unsupported currency for platform wallet transaction: %s", currency)
	}

	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Broadcast via FailoverProvider
	log.Printf("[CryptoService] Broadcasting platform wallet transaction via FailoverProvider")
	txHash, err := s.blockchain.BroadcastTransaction(currency, signedTxHex)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	log.Printf("[CryptoService] ✅ Platform wallet transaction broadcast. TxHash: %s", txHash)
	return txHash, nil
}

func (s *CryptoService) GetMinimumConfirmations(currency string) int {
	switch strings.ToUpper(currency) {
	case "SOL":
		return 1
	case "XRP", "TRX":
		return 1
	case "ETH", "USDT", "USDC", "BNB":
		return 12
	case "BTC":
		return 6
	default:
		return 6
	}
}

func (s *CryptoService) GetTransactionStatus(txHash, currency string) (string, int, error) {
	// TODO: Implement actual blockchain status check via Tatum or RPC
	// For now, return mock values
	log.Printf("[CryptoService] GetTransactionStatus called for %s (currency: %s) - returning mock data", txHash, currency)
	return "confirmed", 12, nil
}

type NetworkDetectionResult struct {
	Type    string
	Network string
	Variant string
}
