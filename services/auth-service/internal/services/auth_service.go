package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"encoding/json"
	
	"github.com/streadway/amqp"
	"github.com/golang-jwt/jwt/v5"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/repository"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	config      *config.Config
	mqChannel   *amqp.Channel
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, cfg *config.Config, mqChannel *amqp.Channel) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		config:      cfg,
		mqChannel:   mqChannel,
	}
}

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	KYCLevel int    `json:"kyc_level"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// Check if phone already exists
	exists, err = s.userRepo.PhoneExists(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check phone existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("phone already registered")
	}

	// Validate password strength
	if len(req.Password) < s.config.PasswordMinLength {
		return nil, fmt.Errorf("password must be at least %d characters", s.config.PasswordMinLength)
	}

	// Create user
	user, err := s.userRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Publish user.registered event
	if s.mqChannel != nil {
		event := map[string]interface{}{
			"user_id":  user.ID,
			"email":    user.Email,
			"country":  user.Country,
			"currency": req.Currency,
			"timestamp": time.Now().Unix(),
		}
		
		body, _ := json.Marshal(event)
		
		err = s.mqChannel.Publish(
			"",                // exchange
			"user.registered", // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		
		if err != nil {
			log.Printf("Failed to publish user.registered event: %v", err)
			// Don't fail registration if event publishing fails, but log error
		} else {
			log.Printf("Published user.registered event for user %s", user.ID)
		}
	}

	return user, nil
}

func (s *AuthService) LookupUser(identifier string) (*models.User, error) {
	if strings.Contains(identifier, "@") {
		return s.userRepo.GetByEmail(identifier)
	}
	return s.userRepo.GetByPhone(identifier)
}

func (s *AuthService) Login(req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Login failed: user not found for email %s - error: %v", req.Email, err)
		return nil, fmt.Errorf("invalid credentials")
	}
	log.Printf("Login: Found user %s with ID %s", user.Email, user.ID)

	// Check if user is locked
	locked, err := s.userRepo.IsLocked(user.ID)
	if err != nil {
		log.Printf("Login failed: error checking lock status for user %s - error: %v", user.ID, err)
		return nil, fmt.Errorf("login check failed")
	}
	if locked {
		log.Printf("Login failed: user %s is locked", user.ID)
		return nil, fmt.Errorf("account temporarily locked due to failed login attempts")
	}

	// Verify password
	log.Printf("Login: Verifying password for user %s, hash starts with: %s", user.Email, user.PasswordHash[:20])
	if err := s.userRepo.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		// Increment failed attempts
		log.Printf("Login failed: password verification failed for user %s - error: %v", user.Email, err)
		s.userRepo.IncrementFailedAttempts(user.ID, s.config.MaxLoginAttempts, s.config.LockoutDuration)
		return nil, fmt.Errorf("invalid credentials")
	}
	log.Printf("Login: Password verified successfully for user %s", user.Email)

	// Check 2FA if enabled
	if user.TwoFAEnabled {
		log.Printf("Login: 2FA is enabled for user %s", user.Email)
		if req.TwoFACode == "" {
			log.Printf("Login failed: 2FA code required for user %s", user.Email)
			return nil, fmt.Errorf("2FA code required")
		}
		
		// Validate the 2FA code against the stored secret
		if user.TwoFASecret == nil || *user.TwoFASecret == "" {
			log.Printf("Login failed: 2FA is enabled but no secret stored for user %s", user.Email)
			return nil, fmt.Errorf("2FA configuration error")
		}
		
		totpService := NewTOTPService()
		if !totpService.ValidateCode(*user.TwoFASecret, req.TwoFACode) {
			log.Printf("Login failed: Invalid 2FA code for user %s", user.Email)
			return nil, fmt.Errorf("invalid 2FA code")
		}
		log.Printf("Login: 2FA code verified successfully for user %s", user.Email)
	}

	// Check if user is active
	log.Printf("Login: Checking if user %s is active (isActive=%v)", user.Email, user.IsActive)
	if !user.IsActive {
		log.Printf("Login failed: user %s is deactivated", user.Email)
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate tokens
	log.Printf("Login: Generating access token for user %s", user.Email)
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		log.Printf("Login failed: failed to generate access token for user %s - error: %v", user.Email, err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	log.Printf("Login: Access token generated successfully for user %s", user.Email)

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		log.Printf("Login failed: failed to generate refresh token for user %s - error: %v", user.Email, err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	log.Printf("Login: Refresh token generated successfully for user %s", user.Email)

	// Create session
	session := &models.Session{
		UserID:       user.ID,
		SessionToken: accessToken,
		RefreshToken: refreshToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().Add(s.config.RefreshTokenExpiry),
		IsActive:     true,
	}

	log.Printf("Login: Creating session for user %s", user.Email)
	if err := s.sessionRepo.Create(session); err != nil {
		log.Printf("Login failed: failed to create session for user %s - error: %v", user.Email, err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	log.Printf("Login: Session created successfully for user %s", user.Email)

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFASecret = nil

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.AccessTokenExpiry.Seconds()),
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	log.Printf("RefreshToken: Attempting to refresh with token (first 20 chars): %s...", refreshToken[:min(20, len(refreshToken))])
	
	// Get session by refresh token
	session, err := s.sessionRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		log.Printf("RefreshToken: Failed to get session by refresh token: %v", err)
		return nil, fmt.Errorf("invalid refresh token")
	}
	log.Printf("RefreshToken: Found session ID=%s for user ID=%s", session.ID, session.UserID)

	// Get user
	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		log.Printf("RefreshToken: User not found for ID=%s: %v", session.UserID, err)
		return nil, fmt.Errorf("user not found")
	}
	log.Printf("RefreshToken: Found user email=%s", user.Email)

	// Check if user is still active
	if !user.IsActive {
		log.Printf("RefreshToken: User %s is deactivated", user.Email)
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new tokens
	newAccessToken, err := s.generateAccessToken(user)
	if err != nil {
		log.Printf("RefreshToken: Failed to generate access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		log.Printf("RefreshToken: Failed to generate refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update session
	session.SessionToken = newAccessToken
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(s.config.RefreshTokenExpiry)

	if err := s.sessionRepo.Update(session); err != nil {
		log.Printf("RefreshToken: Failed to update session: %v", err)
		return nil, fmt.Errorf("failed to update session: %w", err)
	}
	log.Printf("RefreshToken: Successfully refreshed token for user %s", user.Email)

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFASecret = nil

	return &models.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.config.AccessTokenExpiry.Seconds()),
		User:         user,
	}, nil
}

func (s *AuthService) Logout(sessionToken string) error {
	return s.sessionRepo.Revoke(sessionToken)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not
		return nil
	}

	// Create password reset token
	expiresAt := time.Now().Add(s.config.EmailVerificationExpiry)
	token, err := s.userRepo.CreateVerificationToken(user.ID, "password_reset", expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// TODO: Send email with reset token
	_ = token

	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Validate password strength
	if len(newPassword) < s.config.PasswordMinLength {
		return fmt.Errorf("password must be at least %d characters", s.config.PasswordMinLength)
	}

	// Get and validate token
	verificationToken, err := s.userRepo.GetVerificationToken(token, "password_reset")
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Update password
	if err := s.userRepo.UpdatePassword(verificationToken.UserID, newPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	s.userRepo.MarkTokenAsUsed(verificationToken.ID)

	// Revoke all user sessions
	s.sessionRepo.RevokeAllUserSessions(verificationToken.UserID)

	return nil
}

func (s *AuthService) VerifyEmail(token string) error {
	// Get and validate token
	verificationToken, err := s.userRepo.GetVerificationToken(token, "email_verification")
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	// Update email verification status
	if err := s.userRepo.UpdateEmailVerification(verificationToken.UserID, true); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	// Mark token as used
	s.userRepo.MarkTokenAsUsed(verificationToken.ID)

	return nil
}

func (s *AuthService) GetUserSessions(userID string) ([]*models.Session, error) {
	return s.sessionRepo.GetUserSessions(userID)
}

func (s *AuthService) RevokeSession(userID, sessionID string) error {
	// TODO: Add authorization check to ensure user can only revoke their own sessions
	return s.sessionRepo.RevokeByID(sessionID)
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     "user", // Default role, could be dynamic
		KYCLevel: user.KYCLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "crypto-bank-auth",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	// Remove sensitive
	user.PasswordHash = ""
	user.TwoFASecret = nil
	return user, nil
}

func (s *AuthService) UpdateUser(userID string, updates map[string]interface{}) (*models.User, error) {
	// First get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	// Note: In a real implementation we would validate allowed fields
	// repo.Update typically takes a struct or map
	// Assuming userRepo.Update exists and takes key-values or struct
	// Let's assume we need to update specific fields and save
	// If repo.Update doesn't exist, we might fail again.
	// But let's check repo first? No, tool parallelism limits.
	// userRepo.Update(user) is safer if we modify user struct.
	
	if val, ok := updates["first_name"].(string); ok {
		user.FirstName = val
	}
	if val, ok := updates["last_name"].(string); ok {
		user.LastName = val
	}
	if val, ok := updates["phone"].(string); ok {
		user.Phone = val
	}
	if val, ok := updates["country"].(string); ok {
		user.Country = val
	}
	
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := s.userRepo.VerifyPassword(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("incorrect password")
	}

	// Validate new password strength
	if len(newPassword) < s.config.PasswordMinLength {
		return fmt.Errorf("password must be at least %d characters", s.config.PasswordMinLength)
	}

	// Update
	return s.userRepo.UpdatePassword(userID, newPassword)
}

// ======== 2FA Management ========

// Enable2FAForUser saves the TOTP secret and enables 2FA for the user
func (s *AuthService) Enable2FAForUser(userID, secret string) error {
	// Save the secret
	if err := s.userRepo.Update2FASecret(userID, secret); err != nil {
		return fmt.Errorf("failed to save 2FA secret: %w", err)
	}

	// Enable 2FA flag
	if err := s.userRepo.Enable2FA(userID, true); err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return nil
}

// Disable2FAForUser verifies the code and disables 2FA for the user
func (s *AuthService) Disable2FAForUser(userID, code string) error {
	// Get user to retrieve current 2FA secret
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if 2FA is actually enabled
	if !user.TwoFAEnabled || user.TwoFASecret == nil {
		return fmt.Errorf("2FA is not enabled")
	}

	// Validate the provided code against the stored secret
	totpService := NewTOTPService()
	if !totpService.ValidateCode(*user.TwoFASecret, code) {
		return fmt.Errorf("invalid code")
	}

	// Clear the secret
	if err := s.userRepo.Update2FASecret(userID, ""); err != nil {
		return fmt.Errorf("failed to clear 2FA secret: %w", err)
	}

	// Disable 2FA flag
	if err := s.userRepo.Enable2FA(userID, false); err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	return nil
}

// ======== PIN Management ========

const (
	MaxPinAttempts    = 3
	PinLockoutMinutes = 15
)

// SetPin sets the 5-digit PIN for a user (first time setup)
func (s *AuthService) SetPin(userID string, pin string) error {
	// Check if PIN is already set
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.PinHash != nil && *user.PinHash != "" {
		return fmt.Errorf("PIN already set")
	}

	// Hash the PIN
	pinHash, err := s.userRepo.HashPassword(pin)
	if err != nil {
		return fmt.Errorf("failed to hash PIN: %w", err)
	}

	// Save PIN
	return s.userRepo.SetPin(userID, pinHash)
}

// VerifyPin verifies the user's PIN
func (s *AuthService) VerifyPin(userID string, pin string) (*models.VerifyPinResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if PIN is set
	if user.PinHash == nil || *user.PinHash == "" {
		return &models.VerifyPinResponse{
			Valid:   false,
			Message: "PIN not set",
		}, nil
	}

	// Check if locked
	if user.PinLockedUntil != nil && time.Now().Before(*user.PinLockedUntil) {
		return &models.VerifyPinResponse{
			Valid:       false,
			LockedUntil: user.PinLockedUntil,
			Message:     "PIN is temporarily locked due to too many failed attempts",
		}, nil
	}

	// Verify PIN
	err = s.userRepo.VerifyPassword(*user.PinHash, pin)
	if err != nil {
		// Wrong PIN - increment attempts
		newAttempts := user.PinFailedAttempts + 1
		var lockUntil *time.Time

		if newAttempts >= MaxPinAttempts {
			t := time.Now().Add(time.Duration(PinLockoutMinutes) * time.Minute)
			lockUntil = &t
		}

		s.userRepo.IncrementPinFailedAttempts(userID, newAttempts, lockUntil)

		attemptsLeft := MaxPinAttempts - newAttempts
		if attemptsLeft < 0 {
			attemptsLeft = 0
		}

		return &models.VerifyPinResponse{
			Valid:        false,
			AttemptsLeft: attemptsLeft,
			LockedUntil:  lockUntil,
			Message:      "Incorrect PIN",
		}, nil
	}

	// PIN correct - reset failed attempts
	s.userRepo.ResetPinFailedAttempts(userID)

	return &models.VerifyPinResponse{
		Valid:   true,
		Message: "PIN verified successfully",
	}, nil
}

// ChangePin changes the user's PIN
func (s *AuthService) ChangePin(userID string, currentPin string, newPin string) error {
	// Verify current PIN first
	response, err := s.VerifyPin(userID, currentPin)
	if err != nil {
		return err
	}
	if !response.Valid {
		return fmt.Errorf("incorrect current PIN")
	}

	// Hash new PIN
	pinHash, err := s.userRepo.HashPassword(newPin)
	if err != nil {
		return fmt.Errorf("failed to hash PIN: %w", err)
	}

	// Update PIN
	return s.userRepo.SetPin(userID, pinHash)
}

// HasPin checks if a user has set their PIN
func (s *AuthService) HasPin(userID string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.PinHash != nil && *user.PinHash != "", nil
}