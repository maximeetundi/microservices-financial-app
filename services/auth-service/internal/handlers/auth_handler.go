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
	totpSetup, err := h.totpService.GenerateSecret(userEmail.(string), "Zekora")
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

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Secret must be provided when enabling 2FA
	if req.Secret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Secret is required to enable 2FA"})
		return
	}

	// Validate the TOTP code against the provided secret
	if !h.totpService.ValidateCode(req.Secret, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Save the secret and enable 2FA
	if err := h.authService.Enable2FAForUser(userID.(string), req.Secret); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable 2FA"})
		return
	}

	// Generate backup/recovery codes
	recoveryCodes, err := h.totpService.GenerateBackupCodes()
	if err != nil {
		// Log error but don't fail - 2FA is still enabled
		recoveryCodes = []string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "2FA enabled successfully",
		"recovery_codes": recoveryCodes,
	})
}

func (h *AuthHandler) Disable2FA(c *gin.Context) {
	var req models.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify current 2FA code before disabling
	if err := h.authService.Disable2FAForUser(userID.(string), req.Code); err != nil {
		if strings.Contains(err.Error(), "invalid code") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable 2FA"})
		return
	}

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

	// Get current session token for marking
	currentToken := ""
	if authHeader := c.GetHeader("Authorization"); len(authHeader) > 7 {
		currentToken = authHeader[7:] // Remove "Bearer "
	}

	// Enrich sessions with parsed device info
	enrichedSessions := make([]gin.H, len(sessions))
	for i, session := range sessions {
		browser, os, deviceType := parseUserAgent(session.UserAgent)
		
		enrichedSessions[i] = gin.H{
			"id":          session.ID,
			"ip_address":  session.IPAddress,
			"user_agent":  session.UserAgent,
			"browser":     browser,
			"os":          os,
			"device_type": deviceType,
			"device_name": browser + " - " + os,
			"location":    h.getLocationFromIP(session.IPAddress),
			"created_at":  session.CreatedAt,
			"last_active": session.CreatedAt, // Use created_at as last_active for now
			"is_current":  session.SessionToken == currentToken,
			"is_active":   session.IsActive,
		}
	}

	c.JSON(http.StatusOK, gin.H{"sessions": enrichedSessions})
}

// parseUserAgent extracts browser, OS and device type from user agent string
func parseUserAgent(ua string) (browser, os, deviceType string) {
	// Default values
	browser = "Unknown browser"
	os = "Unknown OS"
	deviceType = "desktop"

	ua = strings.ToLower(ua)

	// Detect OS
	switch {
	case strings.Contains(ua, "windows nt 10"):
		os = "Windows 10/11"
	case strings.Contains(ua, "windows nt 6.3"):
		os = "Windows 8.1"
	case strings.Contains(ua, "windows nt 6.2"):
		os = "Windows 8"
	case strings.Contains(ua, "windows nt 6.1"):
		os = "Windows 7"
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac os x"):
		os = "macOS"
	case strings.Contains(ua, "android"):
		os = "Android"
		deviceType = "mobile"
	case strings.Contains(ua, "iphone"):
		os = "iOS (iPhone)"
		deviceType = "mobile"
	case strings.Contains(ua, "ipad"):
		os = "iOS (iPad)"
		deviceType = "tablet"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	case strings.Contains(ua, "ubuntu"):
		os = "Ubuntu"
	case strings.Contains(ua, "chromeos"):
		os = "Chrome OS"
	}

	// Detect Browser
	switch {
	case strings.Contains(ua, "edg/"):
		browser = "Microsoft Edge"
	case strings.Contains(ua, "opr/") || strings.Contains(ua, "opera"):
		browser = "Opera"
	case strings.Contains(ua, "brave"):
		browser = "Brave"
	case strings.Contains(ua, "vivaldi"):
		browser = "Vivaldi"
	case strings.Contains(ua, "chrome") && !strings.Contains(ua, "chromium"):
		browser = "Google Chrome"
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		browser = "Safari"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		browser = "Internet Explorer"
	}

	return browser, os, deviceType
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

func (h *AuthHandler) LookupUser(c *gin.Context) {
	identifier := c.Query("email")
	if identifier == "" {
		identifier = c.Query("phone")
	}

	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone required"})
		return
	}

	user, err := h.authService.LookupUser(identifier)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"country":    user.Country,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.authService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              user.ID,
		"email":           user.Email,
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"phone":           user.Phone,
		"country":         user.Country,
		"date_of_birth":   user.DateOfBirth,
		"kyc_level":       user.KYCLevel,
		"kyc_status":      user.KYCStatus,
		"two_fa_enabled":  user.TwoFAEnabled,
		"email_verified":  user.EmailVerified,
		"phone_verified":  user.PhoneVerified,
		"is_active":       user.IsActive,
		"has_pin":         user.HasPin,
		"created_at":      user.CreatedAt,
		"last_login_at":   user.LastLoginAt,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req map[string]interface{} // Use map for partial updates or a specific struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Prevent updating sensitive fields via this endpoint
	delete(req, "password")
	delete(req, "role")
	delete(req, "kyc_level")
	delete(req, "id")
	delete(req, "email") // Usually require separate flow for email change

	updatedUser, err := h.authService.UpdateUser(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ChangePassword(userID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		if strings.Contains(err.Error(), "incorrect password") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// SetPin - Set up the 5-digit PIN (required after registration)
func (h *AuthHandler) SetPin(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SetPinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify PIN is 5 digits only
	if len(req.Pin) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be exactly 5 digits"})
		return
	}
	for _, ch := range req.Pin {
		if ch < '0' || ch > '9' {
			c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must contain only digits"})
			return
		}
	}

	// Confirm PIN match
	if req.Pin != req.ConfirmPin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PINs do not match"})
		return
	}

	err := h.authService.SetPin(userID.(string), req.Pin)
	if err != nil {
		if strings.Contains(err.Error(), "already set") {
			c.JSON(http.StatusConflict, gin.H{"error": "PIN is already set. Use change PIN endpoint."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN set successfully"})
}

// VerifyPin - Verify the PIN for sensitive actions
func (h *AuthHandler) VerifyPin(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.VerifyPinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.VerifyPin(userID.(string), req.Pin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify PIN"})
		return
	}

	if !response.Valid {
		status := http.StatusUnauthorized
		if response.LockedUntil != nil {
			status = http.StatusTooManyRequests
		}
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ChangePin - Change the PIN (requires current PIN)
func (h *AuthHandler) ChangePin(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.ChangePinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify new PINs match
	if req.NewPin != req.ConfirmPin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New PINs do not match"})
		return
	}

	// Validate new PIN format
	if len(req.NewPin) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be exactly 5 digits"})
		return
	}

	err := h.authService.ChangePin(userID.(string), req.CurrentPin, req.NewPin)
	if err != nil {
		if strings.Contains(err.Error(), "incorrect") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current PIN"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN changed successfully"})
}

// CheckPinStatus - Check if user has set their PIN
func (h *AuthHandler) CheckPinStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	hasPin, err := h.authService.HasPin(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check PIN status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_pin": hasPin,
		"required": true,
	})
}