package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

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

// initialize sets up the master encryption key from environment or generates one
func (v *VaultService) initialize() {
	// Get master key from environment (should be set securely in production)
	masterKeyEnv := os.Getenv("WALLET_MASTER_KEY")

	if masterKeyEnv == "" {
		// For development: use a derived key from a secret
		// In production, this MUST be set via environment variable or external vault (HashiCorp Vault, AWS KMS, etc.)
		secret := os.Getenv("WALLET_SECRET")
		if secret == "" {
			secret = "default-dev-secret-change-in-production"
		}
		salt := os.Getenv("WALLET_SALT")
		if salt == "" {
			salt = "wallet-service-salt-v1"
		}
		// Derive a 256-bit key using PBKDF2
		v.masterKey = pbkdf2.Key([]byte(secret), []byte(salt), 100000, 32, sha256.New)
	} else {
		// Decode the hex-encoded master key
		decoded, err := hex.DecodeString(masterKeyEnv)
		if err != nil || len(decoded) != 32 {
			panic("Invalid WALLET_MASTER_KEY: must be 64 hex characters (256 bits)")
		}
		v.masterKey = decoded
	}
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
