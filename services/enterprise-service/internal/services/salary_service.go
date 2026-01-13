package services

import (
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
)

// SalaryService handles calculation logic (Point 3)
type SalaryService struct{}

func NewSalaryService() *SalaryService {
	return &SalaryService{}
}

// CalculateNetSalary computes Gross -> Net based on config
// 3. "Calcul automatique: Salaire brut, Total des déductions, Salaire net à payer"
func (s *SalaryService) CalculateNetSalary(config *models.SalaryConfig) {
	gross := config.BaseAmount
	
	// Add Bonuses (Primes fixes)
	for i, bonus := range config.Bonuses {
		if bonus.Type == "PERCENTAGE" {
			config.Bonuses[i].Amount = gross * (bonus.Value / 100)
		} else {
			config.Bonuses[i].Amount = bonus.Value
		}
		gross += config.Bonuses[i].Amount
	}

	totalDeductions := 0.0
	// Apply Deductions (Impôts, Cotisations)
	for i, deduction := range config.Deductions {
		if deduction.Type == "PERCENTAGE" {
			config.Deductions[i].Amount = gross * (deduction.Value / 100)
		} else {
			config.Deductions[i].Amount = deduction.Value
		}
		totalDeductions += config.Deductions[i].Amount
	}

	config.NetPayable = gross - totalDeductions
}
