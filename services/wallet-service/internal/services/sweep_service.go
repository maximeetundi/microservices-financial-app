package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
)

// SweepService handles automatic sweeping of user wallet funds to platform hot wallets
// This is used for currencies that don't support memo-based deposits (BTC, ETH, etc.)
type SweepService struct {
	walletRepo      *repository.WalletRepository
	cryptoService   *CryptoService
	platformService *PlatformAccountService
	isRunning       bool
	stopChan        chan struct{}
	mu              sync.Mutex
}

// SweepConfig contains configuration for the sweep job
type SweepConfig struct {
	IntervalMinutes int                // How often to run the sweep job
	MinBalances     map[string]float64 // Minimum balance to sweep (per currency, to cover fees)
}

// DefaultSweepConfig returns default sweep configuration
func DefaultSweepConfig() *SweepConfig {
	return &SweepConfig{
		IntervalMinutes: 5,
		MinBalances: map[string]float64{
			"BTC":  0.0001, // ~$5 at $50k BTC
			"ETH":  0.001,  // ~$3 at $3k ETH
			"USDT": 10.0,   // $10 minimum
			"USDC": 10.0,   // $10 minimum
			"SOL":  0.1,    // ~$15 at $150 SOL
			"TRX":  50.0,   // ~$5 at $0.10 TRX
			"BNB":  0.01,   // ~$5 at $500 BNB
			"DOGE": 100.0,  // ~$10 at $0.10 DOGE
			"LTC":  0.01,   // ~$1 at $100 LTC
		},
	}
}

func NewSweepService(
	walletRepo *repository.WalletRepository,
	cryptoService *CryptoService,
	platformService *PlatformAccountService,
) *SweepService {
	return &SweepService{
		walletRepo:      walletRepo,
		cryptoService:   cryptoService,
		platformService: platformService,
		stopChan:        make(chan struct{}),
	}
}

// Start begins the background sweep job
func (s *SweepService) Start(ctx context.Context, config *SweepConfig) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("sweep service already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	if config == nil {
		config = DefaultSweepConfig()
	}

	log.Printf("[SweepService] Starting sweep job with %d minute interval", config.IntervalMinutes)

	go s.runSweepLoop(ctx, config)

	return nil
}

// Stop stops the background sweep job
func (s *SweepService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isRunning {
		close(s.stopChan)
		s.isRunning = false
		log.Printf("[SweepService] Stopped")
	}
}

func (s *SweepService) runSweepLoop(ctx context.Context, config *SweepConfig) {
	ticker := time.NewTicker(time.Duration(config.IntervalMinutes) * time.Minute)
	defer ticker.Stop()

	// Run once immediately at startup
	s.runSweepCycle(config)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[SweepService] Context cancelled, stopping")
			return
		case <-s.stopChan:
			log.Printf("[SweepService] Stop signal received")
			return
		case <-ticker.C:
			s.runSweepCycle(config)
		}
	}
}

func (s *SweepService) runSweepCycle(config *SweepConfig) {
	log.Printf("[SweepService] Starting sweep cycle")

	// Get all wallets with pending sweep amounts
	wallets, err := s.GetPendingSweeps()
	if err != nil {
		log.Printf("[SweepService] Error getting pending sweeps: %v", err)
		return
	}

	swept := 0
	for _, wallet := range wallets {
		// Skip memo-based currencies (they don't need sweep)
		if s.cryptoService.RequiresMemo(wallet.Currency) {
			continue
		}

		// Check minimum balance
		minBalance, ok := config.MinBalances[wallet.Currency]
		if !ok {
			minBalance = 0.001 // Default minimum
		}

		if wallet.PendingSweepAmount < minBalance {
			continue
		}

		// Perform the sweep
		err := s.SweepWallet(&wallet)
		if err != nil {
			log.Printf("[SweepService] Error sweeping wallet %s: %v", wallet.ID, err)
			continue
		}
		swept++
	}

	if swept > 0 {
		log.Printf("[SweepService] Sweep cycle completed: %d wallets swept", swept)
	}
}

// GetPendingSweeps returns all wallets that need sweeping
func (s *SweepService) GetPendingSweeps() ([]models.Wallet, error) {
	return s.walletRepo.GetWalletsNeedingSweep()
}

// SweepWallet moves funds from a user wallet to the platform hot wallet
func (s *SweepService) SweepWallet(wallet *models.Wallet) error {
	if wallet.WalletAddress == nil || *wallet.WalletAddress == "" {
		return fmt.Errorf("wallet has no address")
	}

	if wallet.PendingSweepAmount <= 0 {
		return fmt.Errorf("no funds to sweep")
	}

	currency := strings.ToUpper(wallet.Currency)
	amount := wallet.PendingSweepAmount

	log.Printf("[SweepService] Sweeping %.8f %s from wallet %s", amount, currency, wallet.ID)

	// 1. Get the target hot wallet
	hotWallet, err := s.platformService.SelectBestCryptoWalletForCredit(currency, "", amount)
	if err != nil {
		return fmt.Errorf("no hot wallet available for %s: %w", currency, err)
	}

	// 2. Get the wallet's private key from vault
	if wallet.PrivateKeyEncrypted == nil || *wallet.PrivateKeyEncrypted == "" {
		return fmt.Errorf("wallet has no private key stored")
	}

	privateKey, err := s.cryptoService.DecryptStoredPrivateKey(*wallet.PrivateKeyEncrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// 3. Estimate fees to leave enough for gas
	feeEstimate, err := s.cryptoService.EstimateTransactionFee(currency, amount, "normal")
	if err != nil {
		log.Printf("[SweepService] Warning: Could not estimate fees for %s, using default", currency)
		feeEstimate = &models.CryptoTransactionEstimate{EstimatedFee: 0.0001}
	}

	sweepAmount := amount - feeEstimate.EstimatedFee
	if sweepAmount <= 0 {
		return fmt.Errorf("sweep amount after fees is negative")
	}

	// 4. Create and sign the sweep transaction
	var signedTxHex string
	gasPrice := int64(0) // Use default gas price

	switch {
	case currency == "BTC" || currency == "LTC" || currency == "DOGE":
		signedTxHex, err = s.cryptoService.signBitcoinLikeTx(privateKey, *wallet.WalletAddress, hotWallet.Address, sweepAmount, currency)
	case currency == "ETH" || currency == "USDT" || currency == "USDC" || currency == "BNB":
		signedTxHex, err = s.cryptoService.signEthereumTx(privateKey, hotWallet.Address, sweepAmount, currency, &gasPrice)
	case currency == "TRX":
		signedTxHex, err = s.cryptoService.signTronTx(privateKey, hotWallet.Address, sweepAmount)
	case currency == "SOL":
		signedTxHex, err = s.cryptoService.signSolanaTx(privateKey, hotWallet.Address, sweepAmount)
	default:
		return fmt.Errorf("sweep not implemented for currency: %s", currency)
	}

	if err != nil {
		return fmt.Errorf("failed to sign sweep transaction: %w", err)
	}

	// 5. Broadcast the transaction
	txHash, err := s.cryptoService.blockchain.BroadcastTransaction(currency, signedTxHex)
	if err != nil {
		return fmt.Errorf("failed to broadcast sweep transaction: %w", err)
	}

	log.Printf("[SweepService] ✅ Sweep broadcast successful. TxHash: %s", txHash)

	// 6. Update wallet sweep status
	err = s.walletRepo.UpdateSweepStatus(wallet.ID, "completed", sweepAmount)
	if err != nil {
		log.Printf("[SweepService] Warning: Failed to update sweep status: %v", err)
	}

	// 7. Credit the platform hot wallet (DB update)
	err = s.platformService.CreditCryptoWallet(hotWallet.ID, sweepAmount, fmt.Sprintf("Sweep from user wallet %s", wallet.ID))
	if err != nil {
		log.Printf("[SweepService] Warning: Failed to credit platform wallet: %v", err)
	}

	return nil
}

// MarkForSweep marks a wallet as needing sweep after a deposit
func (s *SweepService) MarkForSweep(walletID string, amount float64) error {
	return s.walletRepo.AddPendingSweepAmount(walletID, amount)
}
