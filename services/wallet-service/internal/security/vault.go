package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// VaultService handles secure encryption/decryption of sensitive data (private keys, etc.)
type VaultService struct {
	masterKey []byte
	mu        sync.RWMutex
}

var (
	vaultInstance *VaultService
	vaultOnce     sync.Once
)

// GetVaultService returns the singleton vault instance
func GetVaultService() *VaultService {
	vaultOnce.Do(func() {
		vaultInstance = &VaultService{}
		vaultInstance.initialize()
	})
	return vaultInstance
}

// initialize sets up the master encryption key from multiple sources (priority order):
// 1. WALLET_MASTER_KEY environment variable (hex-encoded, 64 chars)
// 2. Fetch from HashiCorp Vault (secret/wallet-service/wallet_master_key)
// 3. Generate random key for development (logged as warning)
func (v *VaultService) initialize() {
	var keySource string

	// Priority 1: Direct environment variable
	masterKeyEnv := os.Getenv("WALLET_MASTER_KEY")
	if masterKeyEnv != "" {
		decoded, err := hex.DecodeString(masterKeyEnv)
		if err != nil || len(decoded) != 32 {
			log.Fatalf("[SECURITY] Invalid WALLET_MASTER_KEY: must be 64 hex characters (256 bits)")
		}
		v.masterKey = decoded
		keySource = "environment variable WALLET_MASTER_KEY"
		log.Printf("[SECURITY] Master key loaded from %s", keySource)
		return
	}

	// Priority 2: Fetch from HashiCorp Vault
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")

	if vaultAddr != "" && vaultToken != "" {
		key, err := v.fetchKeyFromVault(vaultAddr, vaultToken)
		if err == nil && key != "" {
			decoded, err := hex.DecodeString(key)
			if err == nil && len(decoded) == 32 {
				v.masterKey = decoded
				keySource = "HashiCorp Vault (secret/wallet-service)"
				log.Printf("[SECURITY] Master key loaded from %s", keySource)
				return
			}
			log.Printf("[SECURITY] Warning: Key from Vault invalid, trying fallback...")
		} else {
			log.Printf("[SECURITY] Could not fetch key from Vault: %v, trying fallback...", err)
		}
	}

	// Priority 3: PBKDF2 derivation from WALLET_SECRET/WALLET_SALT (if set)
	secret := os.Getenv("WALLET_SECRET")
	salt := os.Getenv("WALLET_SALT")

	if secret != "" && salt != "" {
		v.masterKey = pbkdf2.Key([]byte(secret), []byte(salt), 100000, 32, sha256.New)
		keySource = "PBKDF2 derivation from WALLET_SECRET/WALLET_SALT"
		log.Printf("[SECURITY] Master key derived via %s", keySource)
		return
	}

	// Priority 4: Development fallback - generate random key
	log.Printf("[SECURITY] WARNING: No secure key source found. Generating random development key.")
	log.Printf("[SECURITY] WARNING: This key will be lost on restart! Do NOT use in production.")

	randomKey := make([]byte, 32)
	if _, err := rand.Read(randomKey); err != nil {
		log.Fatalf("[SECURITY] FATAL: Failed to generate random key: %v", err)
	}

	v.masterKey = randomKey
	keySource = "auto-generated random key (DEVELOPMENT ONLY)"
	log.Printf("[SECURITY] Generated development key: %s...", hex.EncodeToString(randomKey)[:16])
	log.Printf("[SECURITY] Master key source: %s", keySource)
}

// fetchKeyFromVault retrieves the wallet_master_key from HashiCorp Vault
func (v *VaultService) fetchKeyFromVault(vaultAddr, vaultToken string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", vaultAddr+"/v1/secret/wallet-service", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Vault-Token", vaultToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault returned status %d", resp.StatusCode)
	}

	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if key, ok := result.Data["wallet_master_key"].(string); ok {
		return key, nil
	}

	return "", errors.New("wallet_master_key not found in vault response")
}

// EncryptPrivateKey encrypts a private key for secure storage
// Returns base64-encoded encrypted data with IV prepended
func (v *VaultService) EncryptPrivateKey(privateKey string) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if len(v.masterKey) == 0 {
		return "", errors.New("vault not initialized")
	}

	plaintext := []byte(privateKey)

	block, err := aes.NewCipher(v.masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// GCM provides authenticated encryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create a unique nonce for this encryption
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and prepend nonce
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	// Return as base64 for safe storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey decrypts a stored private key
func (v *VaultService) DecryptPrivateKey(encryptedData string) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if len(v.masterKey) == 0 {
		return "", errors.New("vault not initialized")
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	block, err := aes.NewCipher(v.masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and decrypt
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EncryptWalletData encrypts arbitrary wallet data (like mnemonic, seed, etc.)
func (v *VaultService) EncryptWalletData(data []byte) (string, error) {
	return v.EncryptPrivateKey(string(data))
}

// DecryptWalletData decrypts wallet data
func (v *VaultService) DecryptWalletData(encryptedData string) ([]byte, error) {
	decrypted, err := v.DecryptPrivateKey(encryptedData)
	if err != nil {
		return nil, err
	}
	return []byte(decrypted), nil
}

// GenerateRandomKey generates a cryptographically secure random key
func GenerateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashData creates a SHA-256 hash of data (for verification without storing plaintext)
func HashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
