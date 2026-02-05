package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

// DepositExpiryService handles background expiration of pending deposits
type DepositExpiryService struct {
	depositRepo *repository.DepositRepository
	interval    time.Duration
	stopChan    chan struct{}
	wg          sync.WaitGroup
	running     bool
	mu          sync.Mutex
}

// NewDepositExpiryService creates a new expiry service
func NewDepositExpiryService(depositRepo *repository.DepositRepository, interval time.Duration) *DepositExpiryService {
	if interval == 0 {
		interval = 5 * time.Minute // Default: check every 5 minutes
	}

	return &DepositExpiryService{
		depositRepo: depositRepo,
		interval:    interval,
		stopChan:    make(chan struct{}),
	}
}

// Start begins the background expiry job
func (s *DepositExpiryService) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	s.wg.Add(1)
	go s.run()

	log.Printf("[DepositExpiryService] ✅ Started with interval: %v", s.interval)
}

// Stop gracefully stops the background job
func (s *DepositExpiryService) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopChan)
	s.wg.Wait()

	log.Println("[DepositExpiryService] ⛔ Stopped")
}

// run is the main loop for the background job
func (s *DepositExpiryService) run() {
	defer s.wg.Done()

	// Run immediately on start
	s.expireTransactions()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.expireTransactions()
		case <-s.stopChan:
			return
		}
	}
}

// expireTransactions marks expired pending transactions
func (s *DepositExpiryService) expireTransactions() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := s.depositRepo.ExpirePendingTransactions(ctx)
	if err != nil {
		log.Printf("[DepositExpiryService] ❌ Error expiring transactions: %v", err)
		return
	}

	if count > 0 {
		log.Printf("[DepositExpiryService] ⏰ Expired %d pending transaction(s)", count)
	}
}

// ExpireNow triggers an immediate expiration check
func (s *DepositExpiryService) ExpireNow() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.depositRepo.ExpirePendingTransactions(ctx)
}

// GetPendingExpiredCount returns the count of transactions that are pending but expired
func (s *DepositExpiryService) GetPendingExpiredCount(ctx context.Context) (int, error) {
	expired, err := s.depositRepo.GetPendingExpired(ctx)
	if err != nil {
		return 0, err
	}
	return len(expired), nil
}

// VerifyAndCompleteDeposit verifies a pending deposit with the provider and completes it if paid
// This can be called to manually verify a transaction status
func (s *DepositExpiryService) VerifyAndCompleteDeposit(ctx context.Context, transactionID string) error {
	deposit, err := s.depositRepo.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}

	if deposit == nil {
		return nil
	}

	// Only verify pending transactions
	if deposit.Status != repository.DepositStatusPending {
		return nil
	}

	// TODO: Call provider's verify endpoint to check actual payment status
	// This would require access to the provider loader and calling VerifyCollection

	return nil
}

// RetryFailedDeposits retries deposits that failed due to temporary errors
func (s *DepositExpiryService) RetryFailedDeposits(ctx context.Context) (int, error) {
	// Get failed transactions that haven't exceeded max retries
	// This would require a new repository method

	// For now, just log
	log.Println("[DepositExpiryService] RetryFailedDeposits not yet implemented")
	return 0, nil
}

// CleanupOldTransactions removes very old completed/failed transactions (data retention)
func (s *DepositExpiryService) CleanupOldTransactions(ctx context.Context, olderThan time.Duration) (int, error) {
	// This would delete transactions older than the specified duration
	// Useful for GDPR compliance or data retention policies

	log.Printf("[DepositExpiryService] CleanupOldTransactions called for records older than %v", olderThan)
	return 0, nil
}

// HealthCheck returns the service health status
func (s *DepositExpiryService) HealthCheck() map[string]interface{} {
	s.mu.Lock()
	running := s.running
	s.mu.Unlock()

	return map[string]interface{}{
		"service":  "deposit_expiry",
		"running":  running,
		"interval": s.interval.String(),
	}
}
