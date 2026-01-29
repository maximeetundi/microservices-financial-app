package providers

import (
	"context"
	"fmt"
	"sync"
)

// Zone represents a geographic zone for routing
type Zone string

const (
	ZoneAfrica       Zone = "africa"
	ZoneEurope       Zone = "europe"
	ZoneNorthAmerica Zone = "north_america"
	ZoneAsia         Zone = "asia"
	ZoneLatinAmerica Zone = "latin_america"
	ZoneMiddleEast   Zone = "middle_east"
	ZoneOceania      Zone = "oceania"
)

// CountryZone maps countries to zones
var CountryZone = map[string]Zone{
	// Africa
	"NG": ZoneAfrica, "GH": ZoneAfrica, "KE": ZoneAfrica, "UG": ZoneAfrica,
	"TZ": ZoneAfrica, "RW": ZoneAfrica, "ZA": ZoneAfrica, "ZM": ZoneAfrica,
	"CI": ZoneAfrica, "SN": ZoneAfrica, "CM": ZoneAfrica, "BJ": ZoneAfrica,
	"TG": ZoneAfrica, "BF": ZoneAfrica, "ML": ZoneAfrica, "NE": ZoneAfrica,
	"CD": ZoneAfrica, "GA": ZoneAfrica, "CG": ZoneAfrica, "MW": ZoneAfrica,
	"EG": ZoneAfrica, "MA": ZoneAfrica, "TN": ZoneAfrica,

	// Europe (SEPA)
	"AT": ZoneEurope, "BE": ZoneEurope, "BG": ZoneEurope, "HR": ZoneEurope,
	"CY": ZoneEurope, "CZ": ZoneEurope, "DK": ZoneEurope, "EE": ZoneEurope,
	"FI": ZoneEurope, "FR": ZoneEurope, "DE": ZoneEurope, "GR": ZoneEurope,
	"HU": ZoneEurope, "IE": ZoneEurope, "IT": ZoneEurope, "LV": ZoneEurope,
	"LT": ZoneEurope, "LU": ZoneEurope, "MT": ZoneEurope, "NL": ZoneEurope,
	"PL": ZoneEurope, "PT": ZoneEurope, "RO": ZoneEurope, "SK": ZoneEurope,
	"SI": ZoneEurope, "ES": ZoneEurope, "SE": ZoneEurope, "GB": ZoneEurope,
	"CH": ZoneEurope, "NO": ZoneEurope,

	// North America
	"US": ZoneNorthAmerica, "CA": ZoneNorthAmerica,

	// Asia
	"PH": ZoneAsia, "ID": ZoneAsia, "VN": ZoneAsia, "TH": ZoneAsia,
	"MY": ZoneAsia, "SG": ZoneAsia, "IN": ZoneAsia, "BD": ZoneAsia,
	"PK": ZoneAsia, "NP": ZoneAsia, "LK": ZoneAsia, "JP": ZoneAsia,
	"HK": ZoneAsia, "CN": ZoneAsia, "KR": ZoneAsia,

	// Latin America
	"MX": ZoneLatinAmerica, "BR": ZoneLatinAmerica, "CO": ZoneLatinAmerica,
	"PE": ZoneLatinAmerica, "CL": ZoneLatinAmerica, "AR": ZoneLatinAmerica,

	// Middle East
	"AE": ZoneMiddleEast, "SA": ZoneMiddleEast, "KW": ZoneMiddleEast,
	"QA": ZoneMiddleEast, "BH": ZoneMiddleEast, "OM": ZoneMiddleEast,

	// Oceania
	"AU": ZoneOceania, "NZ": ZoneOceania,
}

// ZoneRouter routes payouts to the appropriate provider
type ZoneRouter struct {
	providers    map[string]PayoutProvider
	zonePriority map[Zone][]string // Provider priority per zone
	mu           sync.RWMutex
}

// NewZoneRouter creates a new zone router with providers
func NewZoneRouter() *ZoneRouter {
	return &ZoneRouter{
		providers: make(map[string]PayoutProvider),
		zonePriority: map[Zone][]string{
			ZoneAfrica:       {"flutterwave", "mtn_momo", "orange_money", "chipper", "pesapal", "thunes"},
			ZoneEurope:       {"stripe", "paypal", "thunes"},
			ZoneNorthAmerica: {"stripe", "paypal", "thunes"},
			ZoneAsia:         {"thunes", "stripe", "paypal"},
			ZoneLatinAmerica: {"paypal", "thunes"},
			ZoneMiddleEast:   {"thunes", "paypal"},
			ZoneOceania:      {"stripe", "paypal", "thunes"},
		},
	}
}

// RegisterProvider registers a provider
func (r *ZoneRouter) RegisterProvider(provider PayoutProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[provider.GetName()] = provider
}

// GetZone returns the zone for a country
func (r *ZoneRouter) GetZone(country string) Zone {
	if zone, ok := CountryZone[country]; ok {
		return zone
	}
	return ZoneAsia // Default fallback
}

// GetProvider returns the best provider for a country and method
func (r *ZoneRouter) GetProvider(country string, method PayoutMethod) (PayoutProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	zone := r.GetZone(country)
	providerNames, ok := r.zonePriority[zone]
	if !ok {
		return nil, fmt.Errorf("no providers configured for zone %s", zone)
	}

	// Try each provider in priority order
	for _, name := range providerNames {
		provider, exists := r.providers[name]
		if !exists {
			continue
		}

		// Check if provider supports the country
		supported := false
		for _, c := range provider.GetSupportedCountries() {
			if c == country {
				supported = true
				break
			}
		}

		if supported {
			// Check if provider supports the method
			methods, err := provider.GetAvailableMethods(context.Background(), country)
			if err != nil {
				continue
			}

			for _, m := range methods {
				if m.Method == method {
					return provider, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no provider available for country %s with method %s", country, method)
}

// GetAvailableMethodsForCountry returns all available methods for a country
func (r *ZoneRouter) GetAvailableMethodsForCountry(ctx context.Context, country string) ([]AvailableMethod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	methodMap := make(map[PayoutMethod]AvailableMethod)

	zone := r.GetZone(country)
	providerNames := r.zonePriority[zone]

	for _, name := range providerNames {
		provider, exists := r.providers[name]
		if !exists {
			continue
		}

		methods, err := provider.GetAvailableMethods(ctx, country)
		if err != nil {
			continue
		}

		// Add methods (first provider wins for each method type)
		for _, m := range methods {
			if _, exists := methodMap[m.Method]; !exists {
				methodMap[m.Method] = m
			}
		}
	}

	result := make([]AvailableMethod, 0, len(methodMap))
	for _, m := range methodMap {
		result = append(result, m)
	}

	return result, nil
}

// GetBanksForCountry returns banks from the appropriate provider
func (r *ZoneRouter) GetBanksForCountry(ctx context.Context, country string) ([]Bank, error) {
	provider, err := r.GetProvider(country, PayoutMethodBankTransfer)
	if err != nil {
		return nil, err
	}
	return provider.GetBanks(ctx, country)
}

// GetMobileOperatorsForCountry returns mobile operators from the appropriate provider
func (r *ZoneRouter) GetMobileOperatorsForCountry(ctx context.Context, country string) ([]MobileOperator, error) {
	provider, err := r.GetProvider(country, PayoutMethodMobileMoney)
	if err != nil {
		return nil, err
	}
	return provider.GetMobileOperators(ctx, country)
}

// CreatePayout creates a payout using the appropriate provider
func (r *ZoneRouter) CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error) {
	provider, err := r.GetProvider(req.RecipientCountry, req.PayoutMethod)
	if err != nil {
		return nil, err
	}

	// Validate recipient
	if err := provider.ValidateRecipient(ctx, req); err != nil {
		return nil, err
	}

	// Create payout
	return provider.CreatePayout(ctx, req)
}

// GetPayoutStatus gets status from any provider
func (r *ZoneRouter) GetPayoutStatus(ctx context.Context, referenceID, providerName string) (*PayoutStatusResponse, error) {
	r.mu.RLock()
	provider, exists := r.providers[providerName]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	return provider.GetPayoutStatus(ctx, referenceID)
}
