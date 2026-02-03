package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/repository"
)

// ConfigService implements the gRPC AdminConfigService
type ConfigService struct {
	repo *repository.AdminRepository
}

// NewConfigService creates a new ConfigService
func NewConfigService(repo *repository.AdminRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

// GetConfiguration returns a single configuration by key
func (s *ConfigService) GetConfiguration(ctx context.Context, req *GetConfigRequest) (*ConfigResponse, error) {
	config, err := s.repo.GetFeeConfigByKey(req.Key)
	if err != nil {
		return nil, fmt.Errorf("configuration not found: %s", req.Key)
	}

	return &ConfigResponse{
		Key:              config.Key,
		Name:             config.Name,
		Description:      config.Description,
		Type:             string(config.Type),
		FixedAmount:      config.FixedAmount,
		PercentageAmount: config.PercentageAmount,
		Currency:         config.Currency,
		IsEnabled:        config.IsEnabled,
		UpdatedAt:        config.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedBy:        config.UpdatedBy,
	}, nil
}

// GetAllConfigurations returns all configurations with optional service filter
func (s *ConfigService) GetAllConfigurations(ctx context.Context, req *GetAllConfigsRequest) (*ConfigListResponse, error) {
	configs, err := s.repo.GetFeeConfigs()
	if err != nil {
		return nil, err
	}

	responses := make([]*ConfigResponse, 0, len(configs))

	for _, config := range configs {
		// Filter by service prefix if specified
		if req.Service != "" && !strings.HasPrefix(config.Key, req.Service) {
			continue
		}

		responses = append(responses, &ConfigResponse{
			Key:              config.Key,
			Name:             config.Name,
			Description:      config.Description,
			Type:             string(config.Type),
			FixedAmount:      config.FixedAmount,
			PercentageAmount: config.PercentageAmount,
			Currency:         config.Currency,
			IsEnabled:        config.IsEnabled,
			UpdatedAt:        config.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedBy:        config.UpdatedBy,
		})
	}

	// Apply pagination
	start := int(req.Offset)
	end := start + int(req.Limit)
	if req.Limit == 0 {
		end = len(responses)
	}
	if start > len(responses) {
		start = len(responses)
	}
	if end > len(responses) {
		end = len(responses)
	}

	return &ConfigListResponse{
		Configurations: responses[start:end],
		Total:          int32(len(responses)),
	}, nil
}

// GetConfigurationsByPrefix returns configurations matching a prefix
func (s *ConfigService) GetConfigurationsByPrefix(ctx context.Context, req *GetConfigsByPrefixRequest) (*ConfigListResponse, error) {
	configs, err := s.repo.GetFeeConfigs()
	if err != nil {
		return nil, err
	}

	responses := make([]*ConfigResponse, 0)

	for _, config := range configs {
		if strings.HasPrefix(config.Key, req.Prefix) {
			responses = append(responses, &ConfigResponse{
				Key:              config.Key,
				Name:             config.Name,
				Description:      config.Description,
				Type:             string(config.Type),
				FixedAmount:      config.FixedAmount,
				PercentageAmount: config.PercentageAmount,
				Currency:         config.Currency,
				IsEnabled:        config.IsEnabled,
				UpdatedAt:        config.UpdatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedBy:        config.UpdatedBy,
			})
		}
	}

	return &ConfigListResponse{
		Configurations: responses,
		Total:          int32(len(responses)),
	}, nil
}

// IsFeatureEnabled checks if a feature flag is enabled
func (s *ConfigService) IsFeatureEnabled(ctx context.Context, req *FeatureRequest) (*FeatureResponse, error) {
	config, err := s.repo.GetFeeConfigByKey(req.FeatureKey)
	if err != nil {
		// Default to disabled if not found
		return &FeatureResponse{IsEnabled: false}, nil
	}

	return &FeatureResponse{IsEnabled: config.IsEnabled}, nil
}

// GetFee returns fee configuration for a specific operation
func (s *ConfigService) GetFee(ctx context.Context, req *GetFeeRequest) (*FeeResponse, error) {
	// Try with currency suffix first (e.g., crypto_withdrawal_btc)
	key := req.FeeKey
	if req.Currency != "" {
		key = fmt.Sprintf("%s_%s", req.FeeKey, strings.ToLower(req.Currency))
	}

	config, err := s.repo.GetFeeConfigByKey(key)
	if err != nil {
		// Fallback to generic key
		config, err = s.repo.GetFeeConfigByKey(req.FeeKey)
		if err != nil {
			return nil, fmt.Errorf("fee not found: %s", req.FeeKey)
		}
	}

	return &FeeResponse{
		FeeKey:           config.Key,
		Type:             string(config.Type),
		FixedAmount:      config.FixedAmount,
		PercentageAmount: config.PercentageAmount,
		Currency:         config.Currency,
		IsEnabled:        config.IsEnabled,
	}, nil
}

// GetLimit returns limit configuration based on KYC level
func (s *ConfigService) GetLimit(ctx context.Context, req *GetLimitRequest) (*LimitResponse, error) {
	// Map KYC level to tier name
	tierName := "guest"
	switch req.KycLevel {
	case 1:
		tierName = "standard"
	case 2:
		tierName = "verified"
	case 3:
		tierName = "enterprise"
	}

	// Build limit key (e.g., limit_daily_standard)
	key := fmt.Sprintf("%s_%s", req.LimitKey, tierName)

	config, err := s.repo.GetFeeConfigByKey(key)
	if err != nil {
		// Fallback to guest tier
		key = fmt.Sprintf("%s_guest", req.LimitKey)
		config, err = s.repo.GetFeeConfigByKey(key)
		if err != nil {
			return nil, fmt.Errorf("limit not found: %s", req.LimitKey)
		}
	}

	return &LimitResponse{
		LimitKey:    config.Key,
		LimitAmount: config.FixedAmount,
		Currency:    config.Currency,
		TierName:    tierName,
	}, nil
}
