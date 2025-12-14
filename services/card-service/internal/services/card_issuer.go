package services

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/config"
)

// CardIssuerService handles communication with card processors (Marqeta, etc.)
type CardIssuerService struct {
	config     *config.Config
	httpClient *http.Client
}

// NewCardIssuerService creates a new card issuer service
func NewCardIssuerService(cfg *config.Config) *CardIssuerService {
	return &CardIssuerService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateCardResponse represents the response from card creation
type CreateCardResponse struct {
	ExternalCardID string `json:"external_card_id"`
	CardNumber     string `json:"card_number"`
	CVV            string `json:"cvv"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	Status         string `json:"status"`
}

// IssueCard issues a new card through the card processor
func (s *CardIssuerService) IssueCard(userID, cardholderName, cardType, currency string) (*CreateCardResponse, error) {
	// If no API key configured, use mock mode
	if s.config.CardIssuer.APIKey == "" {
		return s.mockIssueCard(userID, cardholderName, cardType, currency)
	}

	// Real implementation for Marqeta or other processors
	url := fmt.Sprintf("%s/cards", s.config.CardIssuer.BaseURL)

	reqBody := map[string]interface{}{
		"user_token":     userID,
		"card_product_token": s.getCardProductToken(cardType),
		"fulfillment": map[string]interface{}{
			"card_personalization": map[string]interface{}{
				"text": map[string]interface{}{
					"name_line_1": cardholderName,
				},
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.config.CardIssuer.APIKey, s.config.CardIssuer.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to issue card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("card issuer returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var cardResp CreateCardResponse
	if err := json.NewDecoder(resp.Body).Decode(&cardResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &cardResp, nil
}

// ActivateCard activates a card through the processor
func (s *CardIssuerService) ActivateCard(externalCardID string) error {
	if s.config.CardIssuer.APIKey == "" {
		return nil // Mock mode
	}

	url := fmt.Sprintf("%s/cards/%s/activate", s.config.CardIssuer.BaseURL, externalCardID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.config.CardIssuer.APIKey, s.config.CardIssuer.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to activate card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("card issuer returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// BlockCard blocks a card through the processor
func (s *CardIssuerService) BlockCard(externalCardID, reason string) error {
	if s.config.CardIssuer.APIKey == "" {
		return nil // Mock mode
	}

	url := fmt.Sprintf("%s/cards/%s/block", s.config.CardIssuer.BaseURL, externalCardID)

	reqBody := map[string]string{"reason": reason}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.config.CardIssuer.APIKey, s.config.CardIssuer.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to block card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("card issuer returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// SetPIN sets the PIN for a card
func (s *CardIssuerService) SetPIN(externalCardID, encryptedPIN string) error {
	if s.config.CardIssuer.APIKey == "" {
		return nil // Mock mode
	}

	url := fmt.Sprintf("%s/cards/%s/pin", s.config.CardIssuer.BaseURL, externalCardID)

	reqBody := map[string]string{"pin": encryptedPIN}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.config.CardIssuer.APIKey, s.config.CardIssuer.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to set PIN: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("card issuer returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// OrderPhysicalCard orders a physical card for shipping
func (s *CardIssuerService) OrderPhysicalCard(externalCardID string, shippingAddress map[string]string) (string, error) {
	if s.config.CardIssuer.APIKey == "" {
		// Mock mode - return fake tracking number
		return fmt.Sprintf("TRK%d", time.Now().Unix()), nil
	}

	url := fmt.Sprintf("%s/cards/%s/ship", s.config.CardIssuer.BaseURL, externalCardID)

	body, _ := json.Marshal(shippingAddress)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.config.CardIssuer.APIKey, s.config.CardIssuer.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to order physical card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("card issuer returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		TrackingNumber string `json:"tracking_number"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.TrackingNumber, nil
}

// Private methods

func (s *CardIssuerService) getCardProductToken(cardType string) string {
	// Map card types to issuer product tokens
	products := map[string]string{
		"prepaid": "cryptobank_prepaid_card",
		"virtual": "cryptobank_virtual_card",
		"gift":    "cryptobank_gift_card",
	}
	
	if token, ok := products[cardType]; ok {
		return token
	}
	return "cryptobank_prepaid_card"
}

// mockIssueCard generates mock card data for development
func (s *CardIssuerService) mockIssueCard(userID, cardholderName, cardType, currency string) (*CreateCardResponse, error) {
	// Generate mock card number (16 digits starting with 4 for Visa)
	cardNumber := "4" + s.generateRandomDigits(15)
	
	// Generate CVV
	cvv := s.generateRandomDigits(3)
	
	// Expiry: 3 years from now
	now := time.Now()
	expiryMonth := int(now.Month())
	expiryYear := now.Year() + 3

	return &CreateCardResponse{
		ExternalCardID: fmt.Sprintf("mock_%s_%d", userID[:8], time.Now().Unix()),
		CardNumber:     cardNumber,
		CVV:            cvv,
		ExpiryMonth:    expiryMonth,
		ExpiryYear:     expiryYear,
		Status:         "active",
	}, nil
}

func (s *CardIssuerService) generateRandomDigits(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		result[i] = digits[num.Int64()]
	}
	return string(result)
}
