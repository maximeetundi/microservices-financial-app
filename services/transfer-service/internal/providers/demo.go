package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DemoConfig holds Demo provider configuration
type DemoConfig struct {
	SimulateDelay time.Duration
	SuccessRate   float64 // 0.0 to 1.0
	DefaultFee    float64
}

// DemoProvider implements PayoutProvider for sandbox/testing
type DemoProvider struct {
	config DemoConfig
}

// NewDemoProvider creates a new Demo provider
func NewDemoProvider(config DemoConfig) *DemoProvider {
	if config.SimulateDelay == 0 {
		config.SimulateDelay = 500 * time.Millisecond
	}
	if config.SuccessRate == 0 {
		config.SuccessRate = 0.95 // 95% success rate
	}
	if config.DefaultFee == 0 {
		config.DefaultFee = 1.5 // 1.5% default fee
	}
	return &DemoProvider{
		config: config,
	}
}

func (p *DemoProvider) GetName() string {
	return "demo"
}

func (p *DemoProvider) GetSupportedCountries() []string {
	return []string{
		"CI", "SN", "CM", "NG", "GH", "KE", "US", "FR", "GB",
		"BF", "TG", "BJ", "ML", "GN", "ZA", "UG", "TZ",
	}
}

func (p *DemoProvider) GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	time.Sleep(p.config.SimulateDelay)

	return []AvailableMethod{
		{
			Code:           "mobile_money",
			Name:           "Mobile Money (Demo)",
			Type:           "mobile",
			Countries:      []string{"CI", "SN", "CM", "GH", "KE", "UG"},
			Currencies:     []string{"XOF", "XAF", "GHS", "KES", "UGX"},
			RequiredFields: []string{"phone"},
		},
		{
			Code:           "bank_transfer",
			Name:           "Bank Transfer (Demo)",
			Type:           "bank",
			Countries:      []string{"NG", "GH", "ZA", "US", "FR", "GB"},
			Currencies:     []string{"NGN", "GHS", "ZAR", "USD", "EUR", "GBP"},
			RequiredFields: []string{"bank_code", "account_number"},
		},
		{
			Code:           "card",
			Name:           "Card Payment (Demo)",
			Type:           "card",
			Countries:      []string{"US", "FR", "GB"},
			Currencies:     []string{"USD", "EUR", "GBP"},
			RequiredFields: []string{"card_number"},
		},
	}, nil
}

func (p *DemoProvider) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	time.Sleep(p.config.SimulateDelay)

	banks := map[string][]Bank{
		"NG": {
			{Code: "DEMO_001", Name: "Demo First Bank", Country: "NG"},
			{Code: "DEMO_002", Name: "Demo Access Bank", Country: "NG"},
			{Code: "DEMO_003", Name: "Demo GTBank", Country: "NG"},
		},
		"CI": {
			{Code: "DEMO_BICICI", Name: "Demo BICICI", Country: "CI"},
			{Code: "DEMO_SGBCI", Name: "Demo Société Générale", Country: "CI"},
		},
	}

	if countryBanks, ok := banks[country]; ok {
		return countryBanks, nil
	}
	return []Bank{
		{Code: "DEMO_DEFAULT", Name: "Demo Universal Bank", Country: country},
	}, nil
}

func (p *DemoProvider) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	time.Sleep(p.config.SimulateDelay)

	return []MobileOperator{
		{Code: "DEMO_ORANGE", Name: "Demo Orange Money", Countries: []string{country}, NumberPrefix: []string{"07", "77"}},
		{Code: "DEMO_MTN", Name: "Demo MTN Money", Countries: []string{country}, NumberPrefix: []string{"05", "55"}},
		{Code: "DEMO_WAVE", Name: "Demo Wave", Countries: []string{country}, NumberPrefix: []string{"70"}},
	}, nil
}

func (p *DemoProvider) ValidateRecipient(ctx context.Context, req *PayoutRequest) error {
	time.Sleep(p.config.SimulateDelay)

	// Demo always validates successfully
	return nil
}

func (p *DemoProvider) GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	time.Sleep(p.config.SimulateDelay)

	fee := req.Amount * (p.config.DefaultFee / 100)
	return &PayoutResponse{
		ProviderName:      "demo",
		ProviderReference: "",
		AmountReceived:    req.Amount,
		ReceivedCurrency:  req.Currency,
		Fee:               fee,
		Status:            "quote", // Use string literal or cast to PayoutStatus if needed, but wait: PayoutStatus is a type.
	}, nil
}

func (p *DemoProvider) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	time.Sleep(p.config.SimulateDelay * 2)

	// Generate a demo transaction reference
	transactionID := fmt.Sprintf("DEMO_%s", uuid.New().String()[:8])

	fee := req.Amount * (p.config.DefaultFee / 100)

	return &PayoutResponse{
		ProviderName:      "demo",
		ProviderReference: transactionID,
		AmountReceived:    req.Amount,
		ReceivedCurrency:  req.Currency,
		Fee:               fee,
		TotalAmount:       req.Amount + fee,
		Status:            PayoutStatusCompleted, // Demo always succeeds
	}, nil
}

func (p *DemoProvider) GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error) {
	time.Sleep(p.config.SimulateDelay)

	// Demo transactions are always completed
	return &PayoutStatusResponse{
		Status:    "completed",
		UpdatedAt: time.Now(),
	}, nil
}

func (p *DemoProvider) CancelPayout(ctx context.Context, referenceID string) error {
	time.Sleep(p.config.SimulateDelay)

	// Demo always allows cancellation
	return nil
}
