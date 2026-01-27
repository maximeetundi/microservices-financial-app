package services

import (
	"context"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

// AuditService handles security audit logging and event publishing via Kafka
type AuditService struct {
	kafkaClient *messaging.KafkaClient
}

// NewAuditService creates a new audit service
func NewAuditService(kafkaClient *messaging.KafkaClient) *AuditService {
	return &AuditService{
		kafkaClient: kafkaClient,
	}
}

// AuditEvent represents an audit event to be published
type AuditEvent struct {
	Type      string                 `json:"type"`
	UserID    string                 `json:"user_id"`
	Email     string                 `json:"email,omitempty"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
	Success   bool                   `json:"success"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// LogLogin logs a login attempt
func (s *AuditService) LogLogin(userID, email, ipAddress, userAgent string, success bool, failReason string) {
	details := map[string]interface{}{}
	if failReason != "" {
		details["reason"] = failReason
	}
	
	s.publishEvent("auth.login", AuditEvent{
		Type:      "login",
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
		Details:   details,
		Timestamp: time.Now(),
	})
}

// LogLogout logs a logout event
func (s *AuditService) LogLogout(userID, ipAddress, userAgent string) {
	s.publishEvent("auth.logout", AuditEvent{
		Type:      "logout",
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   true,
		Timestamp: time.Now(),
	})
}

// LogPasswordChange logs a password change
func (s *AuditService) LogPasswordChange(userID, ipAddress, userAgent string, success bool) {
	s.publishEvent("auth.password_change", AuditEvent{
		Type:      "password_change",
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
		Timestamp: time.Now(),
	})
}

// Log2FAEnable logs 2FA enabling
func (s *AuditService) Log2FAEnable(userID, ipAddress, userAgent string, success bool) {
	s.publishEvent("auth.2fa_enable", AuditEvent{
		Type:      "2fa_enable",
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
		Timestamp: time.Now(),
	})
}

// Log2FADisable logs 2FA disabling
func (s *AuditService) Log2FADisable(userID, ipAddress, userAgent string, success bool) {
	s.publishEvent("auth.2fa_disable", AuditEvent{
		Type:      "2fa_disable",
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
		Timestamp: time.Now(),
	})
}

// LogRegistration logs a new user registration
func (s *AuditService) LogRegistration(userID, email, ipAddress, userAgent string, success bool) {
	s.publishEvent("auth.register", AuditEvent{
		Type:      "register",
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
		Timestamp: time.Now(),
	})
}

// LogAccountLockout logs an account lockout due to failed attempts
func (s *AuditService) LogAccountLockout(userID, email, ipAddress string, failedAttempts int) {
	s.publishEvent("auth.lockout", AuditEvent{
		Type:      "account_lockout",
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
		Success:   false,
		Details: map[string]interface{}{
			"failed_attempts": failedAttempts,
		},
		Timestamp: time.Now(),
	})
}

// LogSessionRevoked logs when a session is revoked
func (s *AuditService) LogSessionRevoked(userID, sessionID, ipAddress, userAgent string) {
	s.publishEvent("auth.session_revoked", AuditEvent{
		Type:      "session_revoked",
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   true,
		Details: map[string]interface{}{
			"session_id": sessionID,
		},
		Timestamp: time.Now(),
	})
}

// LogUserBlocked logs when a user account is blocked
func (s *AuditService) LogUserBlocked(userID, email, reason, adminID string) {
	// Send notification via Kafka
	s.publishNotification("user.blocked", map[string]interface{}{
		"type":     "user.blocked",
		"user_id":  userID,
		"email":    email,
		"phone":    "",
		"reason":   reason,
		"admin_id": adminID,
	})
	
	// Also log audit event
	s.publishEvent("auth.user_blocked", AuditEvent{
		Type:    "user_blocked",
		UserID:  userID,
		Email:   email,
		Success: true,
		Details: map[string]interface{}{
			"reason":   reason,
			"admin_id": adminID,
		},
		Timestamp: time.Now(),
	})
}

// LogUserUnblocked logs when a user account is unblocked
func (s *AuditService) LogUserUnblocked(userID, email, adminID string) {
	// Send notification via Kafka
	s.publishNotification("user.unblocked", map[string]interface{}{
		"type":     "user.unblocked",
		"user_id":  userID,
		"email":    email,
		"phone":    "",
		"admin_id": adminID,
	})
	
	// Also log audit event
	s.publishEvent("auth.user_unblocked", AuditEvent{
		Type:    "user_unblocked",
		UserID:  userID,
		Email:   email,
		Success: true,
		Details: map[string]interface{}{
			"admin_id": adminID,
		},
		Timestamp: time.Now(),
	})
}

// LogPinLocked logs when a user's PIN is locked due to failed attempts
func (s *AuditService) LogPinLocked(userID, email, phone string, failedAttempts int) {
	s.publishNotification("user.pin_locked", map[string]interface{}{
		"type":            "user.pin_locked",
		"user_id":         userID,
		"email":           email,
		"phone":           phone,
		"failed_attempts": failedAttempts,
	})
	
	s.publishEvent("auth.pin_locked", AuditEvent{
		Type:    "pin_locked",
		UserID:  userID,
		Email:   email,
		Success: false,
		Details: map[string]interface{}{
			"failed_attempts": failedAttempts,
		},
		Timestamp: time.Now(),
	})
}

// LogPinUnlocked logs when a user's PIN is unlocked by admin
func (s *AuditService) LogPinUnlocked(userID, email, adminID string) {
	s.publishNotification("user.pin_unlocked", map[string]interface{}{
		"type":     "user.pin_unlocked",
		"user_id":  userID,
		"email":    email,
		"phone":    "",
		"admin_id": adminID,
	})
	
	s.publishEvent("auth.pin_unlocked", AuditEvent{
		Type:    "pin_unlocked",
		UserID:  userID,
		Email:   email,
		Success: true,
		Details: map[string]interface{}{
			"admin_id": adminID,
		},
		Timestamp: time.Now(),
	})
}

// LogUserRegistered logs when a new user registers and publishes to Kafka for wallet creation
func (s *AuditService) LogUserRegistered(user *models.User, currency string) {
	// Publish user.registered event to Kafka for wallet-service to create default wallet
	if s.kafkaClient != nil {
		userEvent := messaging.UserRegisteredEvent{
			UserID:   user.ID,
			Email:    user.Email,
			Country:  user.Country,
			Currency: currency,
		}
		envelope := messaging.NewEventEnvelope(messaging.EventUserRegistered, "auth-service", userEvent)
		if err := s.kafkaClient.Publish(context.Background(), messaging.TopicUserEvents, envelope); err != nil {
			log.Printf("Failed to publish user.registered to Kafka: %v", err)
		} else {
			log.Printf("[Kafka] ✅ Published user.registered event for user %s (country: %s, currency: %s)", user.ID, user.Country, currency)
		}
	} else {
		log.Printf("Warning: kafkaClient is nil, wallet-service won't receive user.registered event")
	}

	// Also publish notification for welcome email etc.
	s.publishNotification("user.registered", map[string]interface{}{
		"type":       "user.registered",
		"user_id":    user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"country":    user.Country,
		"currency":   currency,
		"timestamp":  time.Now().Unix(),
	})

	s.publishEvent("auth.registered", AuditEvent{
		Type:      "user_registered",
		UserID:    user.ID,
		Email:     user.Email,
		Success:   true,
		IPAddress: "", 
		Details: map[string]interface{}{
			"country":  user.Country,
			"currency": currency,
		},
		Timestamp: time.Now(),
	})
}

// LogPinChanged logs when a user changes their PIN
func (s *AuditService) LogPinChanged(userID, email, phone string) {
	s.publishNotification("user.pin_changed", map[string]interface{}{
		"type":    "user.pin_changed",
		"user_id": userID,
		"email":   email,
		"phone":   phone,
	})
}

// publishNotification publishes notification events to Kafka
func (s *AuditService) publishNotification(eventType string, event map[string]interface{}) {
	if s.kafkaClient == nil {
		log.Printf("[NOTIFICATION] %s: %v", eventType, event)
		return
	}

	envelope := messaging.NewEventEnvelope(eventType, "auth-service", event)
	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicNotificationEvents, envelope); err != nil {
		log.Printf("Failed to publish notification event: %v", err)
	} else {
		log.Printf("[NOTIFICATION] ✅ Published %s to Kafka", eventType)
	}
}

// publishEvent publishes audit events to Kafka
func (s *AuditService) publishEvent(eventType string, event AuditEvent) {
	if s.kafkaClient == nil {
		log.Printf("[AUDIT] %s: user=%s success=%v ip=%s", event.Type, event.UserID, event.Success, event.IPAddress)
		return
	}

	envelope := messaging.NewEventEnvelope(eventType, "auth-service", event)
	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicAuditEvents, envelope); err != nil {
		log.Printf("Failed to publish audit event: %v", err)
	}
}
