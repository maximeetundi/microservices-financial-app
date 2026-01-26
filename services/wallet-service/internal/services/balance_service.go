package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
	"github.com/go-redis/redis/v8"
)

type BalanceService struct {
	walletRepo     *repository.WalletRepository
	redisClient    *redis.Client
	exchangeClient *ExchangeClient
	kafkaClient    *messaging.KafkaClient
}

func NewBalanceService(walletRepo *repository.WalletRepository, redisClient *redis.Client, exchangeClient *ExchangeClient, kafkaClient *messaging.KafkaClient) *BalanceService {
	return &BalanceService{
		walletRepo:     walletRepo,
		redisClient:    redisClient,
		exchangeClient: exchangeClient,
		kafkaClient:    kafkaClient,
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
	wallets, err := s.walletRepo.GetByUserID(userID, false)
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

	// Trigger async total balance update
	// We need UserID. How to get it? WalletRepo GetByID again?
	// The update above used walletID.
	wallet, _ := s.walletRepo.GetByID(walletID)
	if wallet != nil {
		go s.UpdateTotalBalanceUSDAndPublish(wallet.UserID)
	}

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

// GetTotalBalanceInUSD calculates the total balance in USD across all wallets
func (s *BalanceService) GetTotalBalanceInUSD(userID string) (float64, error) {
	balances, err := s.GetUserBalances(userID)
	if err != nil {
		return 0, err
	}

	totalUSD := 0.0
	for _, balance := range balances {
		// Convert each balance to USD using ExchangeClient
		// Using ConvertAmount which handles fallback rates or API calls
		usdAmount, err := s.exchangeClient.ConvertAmount(balance.Total, balance.Currency, "USD")
		if err != nil {
			// Log error but continue with other wallets (treat as 0 or use approximation)
			fmt.Printf("Error converting %s to USD: %v\n", balance.Currency, err)
			continue
		}
		totalUSD += usdAmount
	}

	return totalUSD, nil
}

// UpdateTotalBalanceUSDAndPublish calculates total USD balance and publishes an event
func (s *BalanceService) UpdateTotalBalanceUSDAndPublish(userID string) {
	totalUSD, err := s.GetTotalBalanceInUSD(userID)
	if err != nil {
		fmt.Printf("Failed to calculate total USD balance for user %s: %v\n", userID, err)
		return
	}

	// Publish event
	if s.kafkaClient != nil {
		eventData := map[string]interface{}{
			"user_id":           userID,
			"total_balance_usd": totalUSD,
			"updated_at":        time.Now(),
		}
		
		envelope := messaging.NewEventEnvelope("user.balance_updated", "wallet-service", eventData)
		s.kafkaClient.Publish(context.Background(), messaging.TopicUserEvents, envelope)
	}
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
	wallets, err := s.walletRepo.GetByUserID(userID, true)
	if err != nil {
		return
	}

	ctx := context.Background()
	for _, wallet := range wallets {
		cacheKey := fmt.Sprintf("balance:%s", wallet.ID)
		s.redisClient.Del(ctx, cacheKey)
	}
}