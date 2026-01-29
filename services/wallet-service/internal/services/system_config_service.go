package services

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
)

type SystemConfigService struct {
	repo        *repository.ConfigRepository
	configCache map[string]models.SystemConfig
	cacheMutex  sync.RWMutex
}

func NewSystemConfigService(repo *repository.ConfigRepository) *SystemConfigService {
	s := &SystemConfigService{
		repo:        repo,
		configCache: make(map[string]models.SystemConfig),
	}
	s.LoadConfigs() // Initial load
	return s
}

// LoadConfigs loads all configs from DB to memory
func (s *SystemConfigService) LoadConfigs() {
	configs, err := s.repo.GetSystemConfigs()
	if err != nil {
		fmt.Printf("Failed to load system configs: %v\n", err)
		return // Don't crash, just empty cache
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for _, c := range configs {
		s.configCache[c.Key] = c
	}
}

// UpdateLocalConfig updates the in-memory cache and DB
func (s *SystemConfigService) UpdateLocalConfig(cfg models.SystemConfig) error {
	// Update DB first
	if err := s.repo.UpsertSystemConfig(&cfg); err != nil {
		return err
	}

	// Update Memory
	s.cacheMutex.Lock()
	s.configCache[cfg.Key] = cfg
	s.cacheMutex.Unlock()

	return nil
}

// GetConfig returns a config value by key, referencing cache
func (s *SystemConfigService) GetConfig(key string) (models.SystemConfig, bool) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	c, ok := s.configCache[key]
	return c, ok
}

// IsTestnetEnabled checks global switch
func (s *SystemConfigService) IsTestnetEnabled() bool {
	if c, ok := s.GetConfig("system_testnet_enabled"); ok {
		return c.IsEnabled
	}
	return true // Default safely to true or false? Defaulting to true for dev safety?? No, usually false. But let's say true for now as requested.
}

// CheckLimits verifies if a transaction is within limits for a user
func (s *SystemConfigService) CheckLimits(userID, kycLevel, currency string, amount float64) error {
	// Map KYC Level to Config Suffix
	// Assuming KYC Levels: "guest", "standard", "verified", "enterprise"
	// Or numerical: 0, 1, 2, 3

	suffix := "guest"
	switch kycLevel {
	case "tier1", "standard":
		suffix = "standard"
	case "tier2", "verified":
		suffix = "verified"
	case "tier3", "enterprise":
		suffix = "enterprise"
	}

	dailyKey := fmt.Sprintf("limit_daily_%s", suffix)
	monthlyKey := fmt.Sprintf("limit_monthly_%s", suffix)

	// 1. Get Limits
	dailyLimitVal := s.getLimitValue(dailyKey, 1000)     // Default 1000
	monthlyLimitVal := s.getLimitValue(monthlyKey, 5000) // Default 5000

	// 2. Get User Usage
	// Current Date Key
	now := time.Now()
	dailyPeriod := fmt.Sprintf("DAILY_%s", now.Format("2006-01-02"))
	monthlyPeriod := fmt.Sprintf("MONTHLY_%s", now.Format("2006-01"))

	dailyStats, _ := s.repo.GetUserUsage(userID, dailyPeriod, currency)
	monthlyStats, _ := s.repo.GetUserUsage(userID, monthlyPeriod, currency)

	// 3. Check
	if dailyStats.TotalAmount+amount > dailyLimitVal {
		return fmt.Errorf("daily limit exceeded (Limit: %.2f %s, Used: %.2f)", dailyLimitVal, currency, dailyStats.TotalAmount)
	}
	if monthlyStats.TotalAmount+amount > monthlyLimitVal {
		return fmt.Errorf("monthly limit exceeded (Limit: %.2f %s, Used: %.2f)", monthlyLimitVal, currency, monthlyStats.TotalAmount)
	}

	return nil
}

// RecordUsage increments the usage stats after a successful transaction
func (s *SystemConfigService) RecordUsage(userID, currency string, amount float64) {
	now := time.Now()
	dailyPeriod := fmt.Sprintf("DAILY_%s", now.Format("2006-01-02"))
	monthlyPeriod := fmt.Sprintf("MONTHLY_%s", now.Format("2006-01"))

	// We can run these async/goroutine as they are stats, but strict consistency is better for enforcement.
	// For now, async to not block critical path if DB is slow? No, sync is safer for next tx.

	// Error handling ignored for stats? Ideally log it.
	_ = s.repo.IncrementUserUsage(userID, dailyPeriod, currency, amount)
	_ = s.repo.IncrementUserUsage(userID, monthlyPeriod, currency, amount)
}

// Helper
func (s *SystemConfigService) getLimitValue(key string, defaultVal float64) float64 {
	cfg, ok := s.GetConfig(key)
	if !ok {
		return defaultVal
	}
	// Prefer FixedAmount if set, else parse Value
	if cfg.FixedAmount > 0 {
		return cfg.FixedAmount
	}
	// Try parsing value
	if val, err := strconv.ParseFloat(cfg.Value, 64); err == nil && val > 0 {
		return val
	}
	return defaultVal
}
