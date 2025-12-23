package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/config"
)

type SMSService struct {
	config *config.Config
}

type VerificationCode struct {
	Phone     string    `json:"phone"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewSMSService(cfg *config.Config) *SMSService {
	return &SMSService{config: cfg}
}

func (s *SMSService) SendVerificationCode(phone string) (string, error) {
	// Generate 6-digit verification code
	code, err := s.generateVerificationCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate verification code: %w", err)
	}

	message := fmt.Sprintf("Your Zekora verification code is: %s. This code will expire in 5 minutes. Do not share this code with anyone.", code)

	if err := s.sendSMS(phone, message); err != nil {
		return "", fmt.Errorf("failed to send SMS: %w", err)
	}

	return code, nil
}

func (s *SMSService) SendSecurityAlert(phone, alertType string) error {
	message := fmt.Sprintf("Zekora Security Alert: %s. If this wasn't you, please contact support immediately or secure your account at app.zekora.com", alertType)
	
	return s.sendSMS(phone, message)
}

func (s *SMSService) SendTransactionAlert(phone, transactionType, amount, currency string) error {
	message := fmt.Sprintf("Zekora: %s of %s %s completed. If you didn't authorize this, contact support immediately.", transactionType, amount, currency)
	
	return s.sendSMS(phone, message)
}

func (s *SMSService) SendLoginAlert(phone, location string) error {
	message := fmt.Sprintf("Zekora: New login detected from %s. If this wasn't you, secure your account immediately.", location)
	
	return s.sendSMS(phone, message)
}

func (s *SMSService) generateVerificationCode() (string, error) {
	// Generate a random 6-digit number
	max := big.NewInt(999999)
	min := big.NewInt(100000)
	
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	
	// Ensure minimum value
	code := n.Int64() + min.Int64()
	if code > 999999 {
		code = code - 900000
	}
	
	return fmt.Sprintf("%06d", code), nil
}

func (s *SMSService) sendSMS(phone, message string) error {
	if s.config.TwilioAccountSID == "" || s.config.TwilioAuthToken == "" {
		// SMS service not configured, log instead
		fmt.Printf("SMS: To: %s, Message: %s\n", phone, message)
		return nil
	}

	// TODO: Implement actual Twilio SMS sending
	// For now, we'll just log the message
	fmt.Printf("SMS TO %s: %s\n", phone, message)
	
	/* 
	// Twilio implementation would look like this:
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: s.config.TwilioAccountSID,
		Password: s.config.TwilioAuthToken,
	})

	params := &api.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(s.config.TwilioFromNumber)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	return err
	*/

	return nil
}

func (s *SMSService) IsConfigured() bool {
	return s.config.TwilioAccountSID != "" && s.config.TwilioAuthToken != ""
}