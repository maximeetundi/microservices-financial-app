package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// AuditService handles security audit logging and event publishing
type AuditService struct {
	mqChannel *amqp.Channel
}

// NewAuditService creates a new audit service
func NewAuditService(mqChannel *amqp.Channel) *AuditService {
	return &AuditService{
		mqChannel: mqChannel,
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

func (s *AuditService) publishEvent(routingKey string, event AuditEvent) {
	if s.mqChannel == nil {
		log.Printf("[AUDIT] %s: user=%s success=%v ip=%s", event.Type, event.UserID, event.Success, event.IPAddress)
		return
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal audit event: %v", err)
		return
	}

	err = s.mqChannel.Publish(
		"audit.events", // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
	if err != nil {
		log.Printf("Failed to publish audit event: %v", err)
	}
}
