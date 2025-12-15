package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/services"
)

type AuthHandler struct {
	authService  *services.AuthService
	emailService *services.EmailService
	smsService   *services.SMSService
	totpService  *services.TOTPService
	auditService *services.AuditService
}

func NewAuthHandler(authService *services.AuthService, emailService *services.EmailService, smsService *services.SMSService, totpService *services.TOTPService, auditService *services.AuditService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		emailService: emailService,
		smsService:   smsService,
		totpService:  totpService,
		auditService: auditService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate date of birth (must be 18+)
	if time.Since(req.DateOfBirth) < time.Hour*24*365*18 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must be 18 years or older"})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	// Send verification email
	go func() {
		// Create email verification token
		// token, err := h.authService.CreateEmailVerificationToken(user.ID)
		// if err == nil {
		// 	h.emailService.SendVerificationEmail(user.Email, token)
		// }
	}()

	// Remove sensitive data
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get client info
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.authService.Login(&req, ipAddress, userAgent)
	if err != nil {
		if strings.Contains(err.Error(), "2FA code required") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
				"requires_2fa": true,
			})
			return
		}
		
		if strings.Contains(err.Error(), "locked") {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Send login notification
	go func() {
		if h.smsService.IsConfigured() {
			location := h.getLocationFromIP(ipAddress)
			h.smsService.SendLoginAlert(response.User.Phone, location)
		}
	}()

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract session token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization header"})
		return
	}

	sessionToken := tokenParts[1]
	if err := h.authService.Logout(sessionToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{"message": "If an account with this email exists, a password reset link has been sent."})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.VerifyEmail(req.Token)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) VerifyPhone(c *gin.Context) {
	var req models.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement phone verification logic
	// This would involve checking the SMS verification code stored in Redis

	c.JSON(http.StatusOK, gin.H{"message": "Phone verified successfully"})
}

func (h *AuthHandler) Enable2FA(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEmail, _ := c.Get("user_email")
	
	// Generate TOTP secret and QR code
	totpSetup, err := h.totpService.GenerateSecret(userEmail.(string), "CryptoBank")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA setup"})
		return
	}

	// Store secret temporarily (would be confirmed by VerifyTOTP)
	// TODO: Store in Redis with expiration

	c.JSON(http.StatusOK, gin.H{
		"secret":  totpSetup.Secret,
		"qr_code": totpSetup.QRCode,
		"message": "Scan the QR code with your authenticator app and verify with a code to complete setup",
	})
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req models.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Get temporary secret from Redis and validate code
	// If valid, store secret permanently and enable 2FA

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

func (h *AuthHandler) Disable2FA(c *gin.Context) {
	var req models.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Verify current 2FA code before disabling
	// TODO: Disable 2FA in database

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessions, err := h.authService.GetUserSessions(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *AuthHandler) RevokeSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.authService.RevokeSession(userID.(string), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked successfully"})
}

// Helper function to get location from IP (simplified)
func (h *AuthHandler) getLocationFromIP(ip string) string {
	// In a real implementation, you would use a GeoIP service
	// For now, return a placeholder
	return "Unknown location"
}