package services

import (
	"fmt"
	"net/smtp"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/config"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

func (s *EmailService) SendVerificationEmail(email, token string) error {
	subject := "Verify your email address - Zekora"
	verificationURL := fmt.Sprintf("https://app.zekora.com/verify-email?token=%s", token)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify your email</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; padding: 20px; border-radius: 10px;">
        <h1 style="color: #2c3e50; text-align: center;">Welcome to Zekora!</h1>
        <p style="color: #34495e; font-size: 16px;">
            Thank you for registering with Zekora. To complete your registration, please verify your email address by clicking the button below:
        </p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #3498db; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
                Verify Email Address
            </a>
        </div>
        <p style="color: #7f8c8d; font-size: 14px;">
            If the button doesn't work, you can copy and paste this link into your browser:
            <br><br>
            <a href="%s" style="color: #3498db;">%s</a>
        </p>
        <p style="color: #7f8c8d; font-size: 12px; margin-top: 30px;">
            This verification link will expire in 24 hours. If you didn't create this account, please ignore this email.
        </p>
        <hr style="border: none; border-top: 1px solid #ecf0f1; margin: 20px 0;">
        <p style="color: #95a5a6; font-size: 12px; text-align: center;">
            © 2024 Zekora. All rights reserved.
        </p>
    </div>
</body>
</html>`, verificationURL, verificationURL, verificationURL)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(email, token string) error {
	subject := "Reset your password - Zekora"
	resetURL := fmt.Sprintf("https://app.zekora.com/reset-password?token=%s", token)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset your password</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; padding: 20px; border-radius: 10px;">
        <h1 style="color: #e74c3c; text-align: center;">Password Reset Request</h1>
        <p style="color: #34495e; font-size: 16px;">
            We received a request to reset your password for your Zekora account. Click the button below to reset your password:
        </p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #e74c3c; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
                Reset Password
            </a>
        </div>
        <p style="color: #7f8c8d; font-size: 14px;">
            If the button doesn't work, you can copy and paste this link into your browser:
            <br><br>
            <a href="%s" style="color: #e74c3c;">%s</a>
        </p>
        <p style="color: #e74c3c; font-size: 14px; font-weight: bold;">
            ⚠️ Security Notice: If you didn't request this password reset, please ignore this email and consider changing your password.
        </p>
        <p style="color: #7f8c8d; font-size: 12px; margin-top: 30px;">
            This reset link will expire in 24 hours for your security.
        </p>
        <hr style="border: none; border-top: 1px solid #ecf0f1; margin: 20px 0;">
        <p style="color: #95a5a6; font-size: 12px; text-align: center;">
            © 2024 Zekora. All rights reserved.
        </p>
    </div>
</body>
</html>`, resetURL, resetURL, resetURL)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) SendWelcomeEmail(email, firstName string) error {
	subject := "Welcome to Zekora!"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Zekora</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; padding: 20px; border-radius: 10px;">
        <h1 style="color: #27ae60; text-align: center;">Welcome to Zekora, %s! 🎉</h1>
        <p style="color: #34495e; font-size: 16px;">
            Congratulations! Your email has been verified and your Zekora account is now active.
        </p>
        <h3 style="color: #2c3e50;">What's next?</h3>
        <ul style="color: #34495e; font-size: 14px;">
            <li>Complete your KYC verification to unlock higher transaction limits</li>
            <li>Set up two-factor authentication for enhanced security</li>
            <li>Create your first crypto or fiat wallet</li>
            <li>Explore our exchange rates and trading features</li>
        </ul>
        <div style="text-align: center; margin: 30px 0;">
            <a href="https://app.zekora.com/dashboard" style="background-color: #27ae60; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
                Go to Dashboard
            </a>
        </div>
        <p style="color: #7f8c8d; font-size: 14px;">
            Need help? Our support team is available 24/7 at support@zekora.com
        </p>
        <hr style="border: none; border-top: 1px solid #ecf0f1; margin: 20px 0;">
        <p style="color: #95a5a6; font-size: 12px; text-align: center;">
            © 2024 Zekora. All rights reserved.
        </p>
    </div>
</body>
</html>`, firstName)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) SendSecurityAlert(email, alertType, details string) error {
	subject := fmt.Sprintf("Security Alert - %s", alertType)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Security Alert</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 20px; border-radius: 10px;">
        <h1 style="color: #856404; text-align: center;">🔒 Security Alert</h1>
        <p style="color: #856404; font-size: 16px; font-weight: bold;">
            %s
        </p>
        <p style="color: #6c757d; font-size: 14px;">
            %s
        </p>
        <p style="color: #856404; font-size: 14px;">
            If this wasn't you, please:
        </p>
        <ul style="color: #6c757d; font-size: 14px;">
            <li>Change your password immediately</li>
            <li>Enable two-factor authentication</li>
            <li>Contact our support team</li>
        </ul>
        <div style="text-align: center; margin: 30px 0;">
            <a href="https://app.zekora.com/security" style="background-color: #dc3545; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
                Review Security Settings
            </a>
        </div>
        <hr style="border: none; border-top: 1px solid #ecf0f1; margin: 20px 0;">
        <p style="color: #95a5a6; font-size: 12px; text-align: center;">
            © 2024 Zekora. All rights reserved.
        </p>
    </div>
</body>
</html>`, alertType, details)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	if s.config.SMTPUsername == "" || s.config.SMTPPassword == "" {
		// Email service not configured, log instead
		fmt.Printf("EMAIL: To: %s, Subject: %s\n", to, subject)
		return nil
	}

	// Set up authentication
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	// Compose message
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", to, subject, body))

	// Send email
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
	
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailService) IsConfigured() bool {
	return s.config.SMTPUsername != "" && s.config.SMTPPassword != ""
}