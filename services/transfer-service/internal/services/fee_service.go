package services

import (
	"math"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

type FeeService struct {
	repo *repository.FeeRepository
}

func NewFeeService(repo *repository.FeeRepository) *FeeService {
	return &FeeService{repo: repo}
}

func (s *FeeService) CalculateFee(key string, amount float64) (float64, error) {
	config, err := s.repo.GetByKey(key)
	if err != nil {
		return 0, err
	}
	if config == nil || !config.IsActive {
		return 0, nil // Default to free if no config found or inactive
	}

	var fee float64

	switch config.Type {
	case "percentage":
		fee = amount * (config.Percentage / 100)
	case "fixed":
		fee = config.FixedAmount
	case "hybrid":
		fee = (amount * (config.Percentage / 100)) + config.FixedAmount
	default:
		fee = 0
	}

	// Apply Min/Max constraints
	if config.MinFee > 0 && fee < config.MinFee {
		fee = config.MinFee
	}
	if config.MaxFee > 0 && fee > config.MaxFee {
		fee = config.MaxFee
	}

	// Round to 8 decimal places
	return math.Round(fee*100000000) / 100000000, nil
}
