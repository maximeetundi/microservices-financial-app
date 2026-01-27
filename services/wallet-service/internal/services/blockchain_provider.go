package services

import (
	"fmt"
	"log"
)

// BlockchainProvider defines the interface for interacting with blockchain nodes
// It abstracts the differences between Tatum, BlockCypher, Infura, etc.
type BlockchainProvider interface {
	GetBalance(currency, address string) (float64, error)
	BroadcastTransaction(currency, txHex string) (string, error)
	Name() string
}

// FailoverProvider manages multiple providers and retries on failure
type FailoverProvider struct {
	providers []BlockchainProvider
}

func NewFailoverProvider(providers ...BlockchainProvider) *FailoverProvider {
	return &FailoverProvider{
		providers: providers,
	}
}

func (f *FailoverProvider) GetBalance(currency, address string) (float64, error) {
	var lastErr error
	for _, p := range f.providers {
		bal, err := p.GetBalance(currency, address)
		if err == nil {
			return bal, nil // Success
		}
		log.Printf("[%s] GetBalance failed: %v. Switching to next provider...", p.Name(), err)
		lastErr = err
	}
	return 0, fmt.Errorf("all providers failed. last error: %v", lastErr)
}

func (f *FailoverProvider) BroadcastTransaction(currency, txHex string) (string, error) {
	var lastErr error
	for _, p := range f.providers {
		txID, err := p.BroadcastTransaction(currency, txHex)
		if err == nil {
			return txID, nil // Success
		}
		log.Printf("[%s] BroadcastTransaction failed: %v. Switching to next provider...", p.Name(), err)
		lastErr = err
	}
	return "", fmt.Errorf("all providers failed. last error: %v", lastErr)
}

func (f *FailoverProvider) Name() string {
	return "FailoverManager"
}
