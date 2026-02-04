package providers

import (
	"context"
	"fmt"
)

// CollectionRouter intelligently routes collection/deposit payments to the best provider
type CollectionRouter struct {
	providers map[string]CollectionProvider
	// Country -> preferred providers in priority order
	countryRouting map[string][]string
}

// NewCollectionRouter creates a new intelligent collection provider router
func NewCollectionRouter(providers map[string]CollectionProvider) *CollectionRouter {
	// Define country-specific provider priorities for DEPOSITS/COLLECTIONS
	countryRouting := map[string][]string{
		// West Africa - Francophone (XOF)
		"CI": {"wave", "orange_money", "mtn_momo", "cinetpay", "flutterwave"},
		"SN": {"wave", "orange_money", "cinetpay", "flutterwave"},
		"BF": {"orange_money", "wave", "cinetpay", "mtn_momo", "flutterwave"},
		"ML": {"orange_money", "wave", "cinetpay", "flutterwave"},
		"NE": {"orange_money", "cinetpay", "flutterwave"},
		"TG": {"cinetpay", "orange_money", "mtn_momo", "flutterwave"},
		"BJ": {"mtn_momo", "cinetpay", "orange_money", "flutterwave"},
		"GN": {"orange_money", "cinetpay", "mtn_momo", "flutterwave"},

		// West Africa - Anglophone
		"NG": {"paystack", "flutterwave", "mtn_momo"},
		"GH": {"paystack", "flutterwave", "mtn_momo", "wave"},
		"GM": {"wave", "flutterwave"},

		// East Africa
		"KE": {"paystack", "flutterwave"},
		"UG": {"mtn_momo", "wave", "flutterwave"},
		"TZ": {"flutterwave"},
		"RW": {"mtn_momo", "flutterwave"},

		// Central Africa
		"CM": {"orange_money", "cinetpay", "mtn_momo", "flutterwave"},
		"CG": {"cinetpay", "mtn_momo", "orange_money"},
		"CD": {"orange_money", "cinetpay"},
		"GA": {"cinetpay"},
		"TD": {"cinetpay"},

		// Southern Africa
		"ZA": {"paystack", "mtn_momo", "flutterwave"},
		"ZM": {"mtn_momo", "flutterwave"},

		// Global (fallback to international providers)
		"*": {"stripe", "paypal"},
	}

	return &CollectionRouter{
		providers:      providers,
		countryRouting: countryRouting,
	}
}

// SelectProvider selects the best provider for a given country and currency
func (r *CollectionRouter) SelectProvider(ctx context.Context, country, currency, providerHint string) (CollectionProvider, error) {
	// If provider is explicitly requested, use it
	if providerHint != "" && providerHint != "demo" {
		if provider, ok := r.providers[providerHint]; ok {
			return provider, nil
		}
		return nil, fmt.Errorf("requested provider '%s' not available", providerHint)
	}

	// Get preferred providers for this country
	preferredProviders, ok := r.countryRouting[country]
	if !ok {
		// Fall back to global providers
		preferredProviders = r.countryRouting["*"]
	}

	// Try each provider in priority order
	for _, providerName := range preferredProviders {
		if provider, ok := r.providers[providerName]; ok {
			// Check if provider supports this country
			supportedCountries := provider.GetSupportedCountries()
			if r.supportsCountry(supportedCountries, country) {
				return provider, nil
			}
		}
	}

	// If no specific match, try global providers
	for _, providerName := range r.countryRouting["*"] {
		if provider, ok := r.providers[providerName]; ok {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("no provider available for country %s", country)
}

func (r *CollectionRouter) supportsCountry(supportedCountries []string, country string) bool {
	for _, c := range supportedCountries {
		if c == "*" || c == country {
			return true
		}
	}
	return false
}

// GetAvailableProviders returns all providers available for a country
func (r *CollectionRouter) GetAvailableProviders(country string) []CollectionProvider {
	var available []CollectionProvider

	preferredProviders, ok := r.countryRouting[country]
	if !ok {
		preferredProviders = r.countryRouting["*"]
	}

	for _, providerName := range preferredProviders {
		if provider, ok := r.providers[providerName]; ok {
			supportedCountries := provider.GetSupportedCountries()
			if r.supportsCountry(supportedCountries, country) {
				available = append(available, provider)
			}
		}
	}

	return available
}

// GetProviderStatus returns status of all providers
func (r *CollectionRouter) GetProviderStatus() map[string]bool {
	status := make(map[string]bool)
	for name := range r.providers {
		status[name] = true
	}
	return status
}
