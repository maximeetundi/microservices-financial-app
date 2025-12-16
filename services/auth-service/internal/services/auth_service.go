package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/repository"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	config      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		config:      cfg,
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

	return user, nil
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
		if req.TwoFACode == "" {
			return nil, fmt.Errorf("2FA code required")
		}
		// 2FA verification will be handled separately
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

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

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFASecret = ""

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.AccessTokenExpiry.Seconds()),
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	// Get session by refresh token
	session, err := s.sessionRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new tokens
	newAccessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update session
	session.SessionToken = newAccessToken
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(s.config.RefreshTokenExpiry)

	if err := s.sessionRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFASecret = ""

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