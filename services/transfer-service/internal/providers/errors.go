package providers

import (
	"fmt"
)

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	// Provider Errors
	ErrCodeProviderUnavailable ErrorCode = "PROVIDER_UNAVAILABLE"
	ErrCodeProviderTimeout     ErrorCode = "PROVIDER_TIMEOUT"
	ErrCodeProviderRejected    ErrorCode = "PROVIDER_REJECTED"
	ErrCodeProviderMaintenance ErrorCode = "PROVIDER_MAINTENANCE"
	ErrCodeAllProvidersFailed  ErrorCode = "ALL_PROVIDERS_FAILED"

	// Country/Region Errors
	ErrCodeCountryNotSupported  ErrorCode = "COUNTRY_NOT_SUPPORTED"
	ErrCodeCountryDisabled      ErrorCode = "COUNTRY_DISABLED"
	ErrCodeCurrencyNotSupported ErrorCode = "CURRENCY_NOT_SUPPORTED"

	// Transaction Errors
	ErrCodeInsufficientBalance  ErrorCode = "INSUFFICIENT_BALANCE"
	ErrCodeAmountTooLow         ErrorCode = "AMOUNT_TOO_LOW"
	ErrCodeAmountTooHigh        ErrorCode = "AMOUNT_TOO_HIGH"
	ErrCodeInvalidRecipient     ErrorCode = "INVALID_RECIPIENT"
	ErrCodeDuplicateTransaction ErrorCode = "DUPLICATE_TRANSACTION"

	// Validation Errors
	ErrCodeInvalidPhone       ErrorCode = "INVALID_PHONE"
	ErrCodeInvalidBankAccount ErrorCode = "INVALID_BANK_ACCOUNT"
	ErrCodeInvalidBankCode    ErrorCode = "INVALID_BANK_CODE"
)

// ProviderError represents a structured error with user-friendly messages
type ProviderError struct {
	Code          ErrorCode         `json:"code"`
	Message       string            `json:"message"`         // Technical message (logs)
	UserMessage   string            `json:"user_message"`    // User-friendly message
	UserMessageFR string            `json:"user_message_fr"` // French translation
	ProviderName  string            `json:"provider_name,omitempty"`
	Retryable     bool              `json:"retryable"`
	NextProvider  string            `json:"next_provider,omitempty"`
	Details       map[string]string `json:"details,omitempty"`
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.ProviderName, e.Message)
}

// NewProviderError creates a new provider error
func NewProviderError(code ErrorCode, provider, message string) *ProviderError {
	err := &ProviderError{
		Code:         code,
		ProviderName: provider,
		Message:      message,
		Retryable:    isRetryable(code),
	}
	err.UserMessage, err.UserMessageFR = getLocalizedMessages(code)
	return err
}

// WithDetails adds details to the error
func (e *ProviderError) WithDetails(details map[string]string) *ProviderError {
	e.Details = details
	return e
}

// WithNextProvider sets the failover provider
func (e *ProviderError) WithNextProvider(name string) *ProviderError {
	e.NextProvider = name
	return e
}

// isRetryable checks if the error can be retried with another provider
func isRetryable(code ErrorCode) bool {
	switch code {
	case ErrCodeProviderUnavailable, ErrCodeProviderTimeout, ErrCodeProviderMaintenance:
		return true
	default:
		return false
	}
}

// getLocalizedMessages returns user-friendly messages in EN and FR
func getLocalizedMessages(code ErrorCode) (en, fr string) {
	switch code {
	case ErrCodeProviderUnavailable:
		return "Payment provider temporarily unavailable. Trying alternative...",
			"Le prestataire de paiement est temporairement indisponible. Essai d'une alternative..."

	case ErrCodeProviderTimeout:
		return "Payment provider is taking too long to respond. Please try again.",
			"Le prestataire de paiement met trop de temps à répondre. Veuillez réessayer."

	case ErrCodeProviderRejected:
		return "Payment was rejected by the provider. Please verify your details.",
			"Le paiement a été rejeté par le prestataire. Veuillez vérifier vos informations."

	case ErrCodeProviderMaintenance:
		return "Payment provider is under maintenance. Please try later.",
			"Le prestataire de paiement est en maintenance. Veuillez réessayer plus tard."

	case ErrCodeAllProvidersFailed:
		return "All payment providers are currently unavailable. Please try again later.",
			"Tous les prestataires de paiement sont actuellement indisponibles. Veuillez réessayer plus tard."

	case ErrCodeCountryNotSupported:
		return "This payment method is not available in your country.",
			"Ce mode de paiement n'est pas disponible dans votre pays."

	case ErrCodeCountryDisabled:
		return "This payment method is temporarily disabled for your country.",
			"Ce mode de paiement est temporairement désactivé pour votre pays."

	case ErrCodeCurrencyNotSupported:
		return "This currency is not supported for this payment method.",
			"Cette devise n'est pas supportée pour ce mode de paiement."

	case ErrCodeInsufficientBalance:
		return "Insufficient balance to complete this transaction.",
			"Solde insuffisant pour effectuer cette transaction."

	case ErrCodeAmountTooLow:
		return "The amount is below the minimum required.",
			"Le montant est inférieur au minimum requis."

	case ErrCodeAmountTooHigh:
		return "The amount exceeds the maximum allowed.",
			"Le montant dépasse le maximum autorisé."

	case ErrCodeInvalidRecipient:
		return "Invalid recipient information. Please check and try again.",
			"Informations du destinataire invalides. Veuillez vérifier et réessayer."

	case ErrCodeDuplicateTransaction:
		return "A similar transaction is already being processed.",
			"Une transaction similaire est déjà en cours de traitement."

	case ErrCodeInvalidPhone:
		return "Invalid phone number format.",
			"Format de numéro de téléphone invalide."

	case ErrCodeInvalidBankAccount:
		return "Invalid bank account number.",
			"Numéro de compte bancaire invalide."

	case ErrCodeInvalidBankCode:
		return "Invalid bank code. Please select a valid bank.",
			"Code banque invalide. Veuillez sélectionner une banque valide."

	default:
		return "An unexpected error occurred. Please try again.",
			"Une erreur inattendue s'est produite. Veuillez réessayer."
	}
}

// FailoverResult contains the result of a failover attempt
type FailoverResult struct {
	Success        bool             `json:"success"`
	ProviderUsed   string           `json:"provider_used"`
	AttemptedCount int              `json:"attempted_count"`
	Errors         []*ProviderError `json:"errors,omitempty"`
	Response       interface{}      `json:"response,omitempty"`
}

// NewCountryNotSupportedError creates a country not supported error
func NewCountryNotSupportedError(country string) *ProviderError {
	return NewProviderError(ErrCodeCountryNotSupported, "",
		fmt.Sprintf("country %s is not supported", country)).WithDetails(map[string]string{
		"country": country,
	})
}

// NewCountryDisabledError creates a country disabled error
func NewCountryDisabledError(country, provider string) *ProviderError {
	return NewProviderError(ErrCodeCountryDisabled, provider,
		fmt.Sprintf("country %s is disabled for provider %s", country, provider)).WithDetails(map[string]string{
		"country":  country,
		"provider": provider,
	})
}

// NewAmountError creates an amount validation error
func NewAmountError(amount, min, max float64) *ProviderError {
	if amount < min {
		return NewProviderError(ErrCodeAmountTooLow, "",
			fmt.Sprintf("amount %.2f is below minimum %.2f", amount, min)).WithDetails(map[string]string{
			"amount":  fmt.Sprintf("%.2f", amount),
			"minimum": fmt.Sprintf("%.2f", min),
		})
	}
	return NewProviderError(ErrCodeAmountTooHigh, "",
		fmt.Sprintf("amount %.2f exceeds maximum %.2f", amount, max)).WithDetails(map[string]string{
		"amount":  fmt.Sprintf("%.2f", amount),
		"maximum": fmt.Sprintf("%.2f", max),
	})
}
