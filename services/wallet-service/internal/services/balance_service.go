package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
	"github.com/go-redis/redis/v8"
)

type BalanceService struct {
	walletRepo  *repository.WalletRepository
	redisClient *redis.Client
}

func NewBalanceService(walletRepo *repository.WalletRepository, redisClient *redis.Client) *BalanceService {
	return &BalanceService{
		walletRepo:  walletRepo,
		redisClient: redisClient,
	}
}

func (s *BalanceService) GetBalance(walletID string) (*models.Balance, error) {
	ctx := context.Background()
	
	// Try to get from cache first
	cacheKey := fmt.Sprintf("balance:%s", walletID)
	cachedBalance, err := s.redisClient.Get(ctx, cacheKey).Result()
	
	if err == nil {
		var balance models.Balance
		if json.Unmarshal([]byte(cachedBalance), &balance) == nil {
			return &balance, nil
		}
	}

	// Get from database
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	balance := &models.Balance{
		Currency:        wallet.Currency,
		Available:       wallet.Balance,
		Frozen:          wallet.FrozenBalance,
		Total:           wallet.Balance + wallet.FrozenBalance,
		PendingDeposits: 0, // TODO: Calculate from pending transactions
	}

	// Cache the balance for 30 seconds
	balanceJSON, _ := json.Marshal(balance)
	s.redisClient.Set(ctx, cacheKey, balanceJSON, 30*time.Second)

	return balance, nil
}

func (s *BalanceService) GetUserBalances(userID string) ([]*models.Balance, error) {
	wallets, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}

	var balances []*models.Balance
	for _, wallet := range wallets {
		balance, err := s.GetBalance(wallet.ID)
		if err != nil {
			continue // Skip wallets with errors
		}
		balances = append(balances, balance)
	}

	return balances, nil
}

func (s *BalanceService) UpdateBalance(walletID string, amount float64, transactionType string) error {
	// Use database transaction for consistency
	err := s.walletRepo.UpdateBalanceWithTransaction(walletID, amount, transactionType)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("balance:%s", walletID)
	s.redisClient.Del(ctx, cacheKey)

	return nil
}

func (s *BalanceService) FreezeAmount(walletID string, amount float64) error {
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet.Balance < amount {
		return fmt.Errorf("insufficient available balance")
	}

	newBalance := wallet.Balance - amount
	newFrozenBalance := wallet.FrozenBalance + amount

	err = s.walletRepo.UpdateBalance(walletID, newBalance, newFrozenBalance)
	if err != nil {
		return fmt.Errorf("failed to freeze amount: %w", err)
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("balance:%s", walletID)
	s.redisClient.Del(ctx, cacheKey)

	return nil
}

func (s *BalanceService) UnfreezeAmount(walletID string, amount float64) error {
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet.FrozenBalance < amount {
		return fmt.Errorf("insufficient frozen balance")
	}

	newBalance := wallet.Balance + amount
	newFrozenBalance := wallet.FrozenBalance - amount

	err = s.walletRepo.UpdateBalance(walletID, newBalance, newFrozenBalance)
	if err != nil {
		return fmt.Errorf("failed to unfreeze amount: %w", err)
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("balance:%s", walletID)
	s.redisClient.Del(ctx, cacheKey)

	return nil
}

func (s *BalanceService) ValidateSufficientBalance(walletID string, amount float64, includeFrozen bool) error {
	balance, err := s.GetBalance(walletID)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	availableAmount := balance.Available
	if includeFrozen {
		availableAmount += balance.Frozen
	}

	if availableAmount < amount {
		return fmt.Errorf("insufficient balance: available %.8f, required %.8f", 
			availableAmount, amount)
	}

	return nil
}

func (s *BalanceService) GetTotalBalanceInUSD(userID string) (float64, error) {
	// This would require integration with exchange rate service
	// For now, return a mock value
	balances, err := s.GetUserBalances(userID)
	if err != nil {
		return 0, err
	}

	// Mock exchange rates
	exchangeRates := map[string]float64{
		"USD": 1.0,
		"EUR": 1.18,
		"BTC": 45000.0,
		"ETH": 3000.0,
		"BSC": 300.0,
	}

	totalUSD := 0.0
	for _, balance := range balances {
		rate := exchangeRates[balance.Currency]
		if rate == 0 {
			rate = 1.0 // Default rate
		}
		totalUSD += balance.Total * rate
	}

	return totalUSD, nil
}

func (s *BalanceService) GetBalanceHistory(walletID string, days int) ([]map[string]interface{}, error) {
	// This would require storing historical balance data
	// For now, return mock data
	
	history := make([]map[string]interface{}, 0)
	baseDate := time.Now().AddDate(0, 0, -days)

	for i := 0; i < days; i++ {
		date := baseDate.AddDate(0, 0, i)
		// Mock balance calculation based on date
		mockBalance := 1000.0 + float64(i*10)
		
		history = append(history, map[string]interface{}{
			"date":    date.Format("2006-01-02"),
			"balance": mockBalance,
		})
	}

	return history, nil
}

func (s *BalanceService) InvalidateCache(walletID string) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("balance:%s", walletID)
	s.redisClient.Del(ctx, cacheKey)
}

func (s *BalanceService) InvalidateUserCache(userID string) {
	// Get all user wallets and invalidate their cache
	wallets, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return
	}

	ctx := context.Background()
	for _, wallet := range wallets {
		cacheKey := fmt.Sprintf("balance:%s", wallet.ID)
		s.redisClient.Del(ctx, cacheKey)
	}
}