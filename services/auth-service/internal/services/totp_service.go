package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"image/png"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPService struct{}

func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

type TOTPSetup struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"` // Base64 encoded PNG
	URL    string `json:"url"`
}

func (s *TOTPService) GenerateSecret(email, issuer string) (*TOTPSetup, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
		SecretSize:  32,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Generate QR code image
	img, err := key.Image(256, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code image: %w", err)
	}

	// Convert image to base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode QR code image: %w", err)
	}

	qrCodeBase64 := fmt.Sprintf("data:image/png;base64,%s", base32.StdEncoding.EncodeToString(buf.Bytes()))

	return &TOTPSetup{
		Secret: key.Secret(),
		QRCode: qrCodeBase64,
		URL:    key.URL(),
	}, nil
}

func (s *TOTPService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

func (s *TOTPService) ValidateCodeWithWindow(secret, code string, window int) bool {
	// Allow for time drift by checking previous/next time windows
	now := time.Now()
	
	for i := -window; i <= window; i++ {
		testTime := now.Add(time.Duration(i) * 30 * time.Second)
		expectedCode, err := totp.GenerateCode(secret, testTime)
		if err != nil {
			continue
		}
		
		if code == expectedCode {
			return true
		}
	}
	
	return false
}

func (s *TOTPService) GenerateBackupCodes() ([]string, error) {
	codes := make([]string, 10)
	
	for i := 0; i < 10; i++ {
		code, err := s.generateBackupCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate backup code: %w", err)
		}
		codes[i] = code
	}
	
	return codes, nil
}

func (s *TOTPService) generateBackupCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	
	// Format as XXXX-XXXX
	return fmt.Sprintf("%s-%s", string(bytes[:4]), string(bytes[4:])), nil
}