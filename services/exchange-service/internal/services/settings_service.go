package services

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
)

type SettingsService struct {
	repo  *repository.SettingsRepository
	cache map[string]string
	mutex sync.RWMutex
}

func NewSettingsService(repo *repository.SettingsRepository) *SettingsService {
	s := &SettingsService{
		repo:  repo,
		cache: make(map[string]string),
	}
	// Load initial cache
	s.RefreshCache()
	return s
}

// RefreshCache reloads all settings from DB into memory cache
func (s *SettingsService) RefreshCache() error {
	settings, err := s.repo.GetAll()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, setting := range settings {
		s.cache[setting.Key] = setting.Value
	}
	return nil
}

// GetAll returns all system settings
func (s *SettingsService) GetAll() ([]models.SystemSetting, error) {
	return s.repo.GetAll()
}

// GetByCategory returns settings for a specific category
func (s *SettingsService) GetByCategory(category string) ([]models.SystemSetting, error) {
	return s.repo.GetByCategory(category)
}

// GetSetting returns a setting value from cache (fast)
func (s *SettingsService) GetSetting(key string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.cache[key]
}

// GetSettingWithDefault returns a setting value or default if not found
func (s *SettingsService) GetSettingWithDefault(key, defaultValue string) string {
	s.mutex.RLock()
	value, exists := s.cache[key]
	s.mutex.RUnlock()

	if !exists || value == "" {
		return defaultValue
	}
	return value
}

// UpdateSetting updates a setting value
func (s *SettingsService) UpdateSetting(setting *models.SystemSetting) error {
	err := s.repo.Update(setting)
	if err != nil {
		return err
	}

	// Update cache
	s.mutex.Lock()
	s.cache[setting.Key] = setting.Value
	s.mutex.Unlock()

	fmt.Printf("[SettingsService] Setting '%s' updated to '%s'\n", setting.Key, setting.Value)
	return nil
}

// ===== Typed Getters =====

// GetRateUpdateInterval returns the rate update interval in seconds
func (s *SettingsService) GetRateUpdateInterval() int {
	value := s.GetSettingWithDefault("rate_update_interval", "60")
	interval, err := strconv.Atoi(value)
	if err != nil || interval < 10 {
		return 60 // Minimum 10 seconds, default 60
	}
	return interval
}

// GetCryptoNetwork returns the crypto network (mainnet/testnet)
func (s *SettingsService) GetCryptoNetwork() string {
	return s.GetSettingWithDefault("crypto_network", "mainnet")
}

// IsBinanceTestMode returns whether Binance test mode is enabled
func (s *SettingsService) IsBinanceTestMode() bool {
	value := s.GetSettingWithDefault("binance_test_mode", "true")
	return value == "true" || value == "1"
}

// IsRateUpdateEnabled returns whether automatic rate updates are enabled
func (s *SettingsService) IsRateUpdateEnabled() bool {
	value := s.GetSettingWithDefault("rate_update_enabled", "true")
	return value == "true" || value == "1"
}
