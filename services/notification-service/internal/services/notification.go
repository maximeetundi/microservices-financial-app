package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/repository"
)

// NotificationService handles all notification types
type NotificationService struct {
	config      *config.Config
	repo        *repository.NotificationRepository
	kafkaClient *messaging.KafkaClient
}

// NewNotificationService creates a new notification service
func NewNotificationService(cfg *config.Config, repo *repository.NotificationRepository, kafkaClient *messaging.KafkaClient) *NotificationService {
	return &NotificationService{
		config:      cfg,
		repo:        repo,
		kafkaClient: kafkaClient,
	}
}

// Start begins consuming from all notification queues
func (s *NotificationService) Start() error {
	log.Println("Starting notification consumers...")
	
	if s.kafkaClient == nil {
		log.Println("⚠️ Kafka client is nil, running in API-only mode")
		return nil
	}

	topics := []string{
		messaging.TopicUserEvents,
		messaging.TopicWalletEvents,
		messaging.TopicTransferEvents,
		messaging.TopicExchangeEvents,
		messaging.TopicPaymentEvents,
		messaging.TopicCardEvents,
		messaging.TopicNotificationEvents,
	}

	for _, topic := range topics {
		if err := s.kafkaClient.Subscribe(topic, s.handleEvent); err != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
		}
	}
	
	log.Println("✅ Notification service consumers started")
	return nil
}

func (s *NotificationService) handleEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	// Parse event data to map
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	
	var eventData map[string]interface{}
	if err := json.Unmarshal(dataBytes, &eventData); err != nil {
		log.Printf("Failed to unmarshal event data: %v", err)
		return err
	}
	
	// Create notification object based on event type
	notification := s.createNotification(event.Type, eventData)
	if notification == nil {
		// Event not handled or ignored
		return nil
	}
	
	// Extract metadata
	userID, _ := eventData["user_id"].(string)
	email, _ := eventData["email"].(string)
	phone, _ := eventData["phone"].(string)
	
	// Fallback for user email if in specific admin events (user_email)
	if email == "" {
		email, _ = eventData["user_email"].(string)
	}

	s.sendNotification(userID, email, phone, notification)
	return nil
}

// NotificationPriority defines urgency levels
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityNormal   NotificationPriority = "normal"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

type Notification struct {
	Title    string
	Body     string
	Type     string
	Priority NotificationPriority
	Data     map[string]interface{}
	ActionUrl string // URL or route to redirect to
	Channels []string // email, sms, push
}

func (s *NotificationService) createNotification(eventType string, event map[string]interface{}) *Notification {
	switch eventType {
	
	// ===== WALLET EVENTS =====
	case "wallet.created":
		return &Notification{
			Title:    "Portefeuille créé",
			Body:     fmt.Sprintf("Votre portefeuille %s a été créé avec succès.", event["currency"]),
			Type:     "wallet",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
	case "wallet.frozen":
		return &Notification{
			Title:    "⚠️ Portefeuille gelé",
			Body:     "Votre portefeuille a été gelé pour des raisons de sécurité. Contactez le support.",
			Type:     "security",
			Priority: PriorityCritical,
			Data:     event,
			Channels: []string{"push", "email", "sms"}, // SMS pour sécurité
		}
	case "wallet.deposit":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		return &Notification{
			Title:    "Dépôt reçu",
			Body:     fmt.Sprintf("Vous avez reçu %.2f %s sur votre portefeuille.", amount, currency),
			Type:     "transaction",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS pour dépôt
		}
		
	// ===== TRANSACTION EVENTS =====
	case "transaction.completed":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		return &Notification{
			Title:    "Transaction réussie",
			Body:     fmt.Sprintf("Votre transaction de %.2f %s a été effectuée.", amount, currency),
			Type:     "transaction",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
	case "transaction.failed":
		return &Notification{
			Title:    "❌ Transaction échouée",
			Body:     "Votre transaction a échoué. Veuillez réessayer.",
			Type:     "transaction",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
		
	// ===== TRANSFER EVENTS =====
	case "transfer.initiated":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		return &Notification{
			Title:    "Transfert en cours",
			Body:     fmt.Sprintf("Votre transfert de %.2f %s est en cours de traitement.", amount, currency),
			Type:     "transfer",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push"}, // Juste push
		}
	case "transfer.completed":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		recipient, _ := event["recipient"].(string)
		return &Notification{
			Title:    "✅ Transfert réussi",
			Body:     fmt.Sprintf("Votre transfert de %.2f %s vers %s a été effectué.", amount, currency, recipient),
			Type:     "transfer",
			Priority: PriorityHigh,
			Data:     event,
			ActionUrl: fmt.Sprintf("/transactions/%s", event["transfer_id"]),
			Channels: []string{"push", "email", "sms"}, // SMS pour transferts réussis (important)
		}
	case "transfer.sent":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		return &Notification{
			Title:    "💸 Argent envoyé",
			Body:     fmt.Sprintf("Vous avez envoyé %.2f %s avec succès.", amount, currency),
			Type:     "transfer",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"}, // Push et email pour l'émetteur
		}
	case "transfer.failed":
		return &Notification{
			Title:    "❌ Transfert échoué",
			Body:     "Votre transfert a échoué. Vos fonds ont été remboursés.",
			Type:     "transfer",
			Priority: PriorityCritical,
			Data:     event,
			Channels: []string{"push", "email", "sms"}, // SMS pour échec
		}
	case "transfer.received":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		// Use sender_name if available, fallback to sender ID
		senderName, ok := event["sender_name"].(string)
		if !ok || senderName == "" {
			senderName, _ = event["sender"].(string)
		}
		return &Notification{
			Title:    "💰 Argent reçu",
			Body:     fmt.Sprintf("Vous avez reçu %.2f %s de %s.", amount, currency, senderName),
			Type:     "transfer",
			Priority: PriorityHigh,
			Data:     event,
			ActionUrl: fmt.Sprintf("/transactions/%s", event["transfer_id"]),
			Channels: []string{"push", "email", "sms"}, // SMS pour réception
		}

	// ===== PAYMENT (MERCHANT QR) EVENTS =====
	case "payment.sent":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		title, _ := event["title"].(string)
		return &Notification{
			Title:    "💳 Paiement effectué",
			Body:     fmt.Sprintf("Vous avez payé %.2f %s pour \"%s\".", amount, currency, title),
			Type:     "transfer",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"},
		}
	case "payment.received":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		title, _ := event["title"].(string)
		return &Notification{
			Title:    "💰 Paiement reçu",
			Body:     fmt.Sprintf("Vous avez reçu %.2f %s pour \"%s\".", amount, currency, title),
			Type:     "transfer",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email", "sms"}, // SMS pour le marchand
		}
		
	// ===== EXCHANGE EVENTS =====
	case "exchange.completed":
		fromCurrency, _ := event["from_currency"].(string)
		toCurrency, _ := event["to_currency"].(string)
		return &Notification{
			Title:    "Conversion effectuée",
			Body:     fmt.Sprintf("Votre conversion %s → %s a été effectuée.", fromCurrency, toCurrency),
			Type:     "exchange",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}

	case "fiat_exchange.completed":
		fromCurrency, _ := event["from_currency"].(string)
		toCurrency, _ := event["to_currency"].(string)
		amount, _ := event["amount"].(float64)
		return &Notification{
			Title:    "Change Fiat effectué",
			Body:     fmt.Sprintf("Vous avez converti %.2f %s en %s avec succès.", amount, fromCurrency, toCurrency),
			Type:     "exchange",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"},
		}
		
	// ===== CARD EVENTS =====
	case "card.created":
		return &Notification{
			Title:    "💳 Carte créée",
			Body:     "Votre nouvelle carte a été créée avec succès.",
			Type:     "card",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
	case "card.shipped":
		return &Notification{
			Title:    "📦 Carte expédiée",
			Body:     "Votre carte physique a été expédiée! Livraison sous 5-7 jours.",
			Type:     "card",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email", "sms"}, // SMS pour expédition
		}
	case "card.activated":
		return &Notification{
			Title:    "✅ Carte activée",
			Body:     "Votre carte est maintenant active et prête à l'emploi.",
			Type:     "card",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
	case "card.blocked":
		return &Notification{
			Title:    "🚫 Carte bloquée",
			Body:     "Votre carte a été bloquée. Contactez le support si ce n'est pas vous.",
			Type:     "security",
			Priority: PriorityCritical,
			Data:     event,
			Channels: []string{"push", "email", "sms"}, // SMS pour sécurité
		}
	case "card.transaction":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		merchant, _ := event["merchant"].(string)
		return &Notification{
			Title:    "Paiement carte",
			Body:     fmt.Sprintf("Paiement de %.2f %s chez %s", amount, currency, merchant),
			Type:     "card",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push"}, // Juste push pour transactions carte
		}
		
	// ===== AUTH/SECURITY EVENTS =====
	case "auth.login":
		device, _ := event["device"].(string)
		location, _ := event["location"].(string)
		return &Notification{
			Title:    "Nouvelle connexion",
			Body:     fmt.Sprintf("Connexion depuis %s (%s)", device, location),
			Type:     "security",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
	case "auth.login_suspicious":
		return &Notification{
			Title:    "⚠️ Connexion suspecte",
			Body:     "Une connexion depuis un nouvel appareil a été détectée. Était-ce vous?",
			Type:     "security",
			Priority: PriorityCritical,
			Data:     event,
			Channels: []string{"push", "email", "sms"}, // SMS pour sécurité
		}
	case "auth.password_changed":
		return &Notification{
			Title:    "Mot de passe modifié",
			Body:     "Votre mot de passe a été modifié avec succès.",
			Type:     "security",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email", "sms"}, // SMS pour sécurité
		}
	case "auth.2fa_enabled":
		return &Notification{
			Title:    "2FA activé",
			Body:     "L'authentification à deux facteurs a été activée sur votre compte.",
			Type:     "security",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
		
	// ===== KYC EVENTS =====
	case "kyc.approved":
		return &Notification{
			Title:    "✅ Vérification approuvée",
			Body:     "Votre identité a été vérifiée. Vous avez maintenant accès à toutes les fonctionnalités.",
			Type:     "kyc",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}
	case "kyc.rejected":
		return &Notification{
			Title:    "❌ Vérification refusée",
			Body:     "Votre vérification d'identité a été refusée. Veuillez soumettre de nouveaux documents.",
			Type:     "kyc",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"}, // Pas de SMS
		}

	// ===== USER ACCOUNT EVENTS =====
	case "user.registered":
		firstName, _ := event["first_name"].(string)
		return &Notification{
			Title:    "👋 Bienvenue sur Zekora",
			Body:     fmt.Sprintf("Bonjour %s! Votre compte a été créé avec succès. Vérifiez votre identité pour débloquer toutes les fonctionnalités.", firstName),
			Type:     "account",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push", "email"},
		}
	case "user.blocked":
		reason, _ := event["reason"].(string)
		body := "Votre compte a été temporairement suspendu."
		if reason != "" {
			body = fmt.Sprintf("Votre compte a été suspendu: %s", reason)
		}
		return &Notification{
			Title:    "🚫 Compte suspendu",
			Body:     body + " Contactez le support pour plus d'informations.",
			Type:     "security",
			Priority: PriorityCritical,
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}
	case "user.unblocked":
		return &Notification{
			Title:    "✅ Compte réactivé",
			Body:     "Bonne nouvelle! Votre compte a été réactivé. Vous pouvez maintenant vous connecter.",
			Type:     "account",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}
	case "user.pin_locked":
		return &Notification{
			Title:    "🔒 Code PIN bloqué",
			Body:     "Votre code PIN a été bloqué après plusieurs tentatives incorrectes. Contactez le support.",
			Type:     "security",
			Priority: PriorityCritical,
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}
	case "user.pin_unlocked":
		return &Notification{
			Title:    "🔓 Code PIN débloqué",
			Body:     "Votre code PIN a été débloqué. Vous pouvez maintenant l'utiliser normalement.",
			Type:     "account",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"},
		}
	case "user.pin_changed":
		return &Notification{
			Title:    "🔐 Code PIN modifié",
			Body:     "Votre code PIN a été modifié avec succès. Si ce n'était pas vous, contactez le support.",
			Type:     "security",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}

	// ===== ADMIN ACTION EVENTS =====
	case "admin.user_created":
		email, _ := event["user_email"].(string)
		return &Notification{
			Title:    "👤 Nouveau compte",
			Body:     fmt.Sprintf("Nouveau compte créé: %s", email),
			Type:     "admin",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push"},
		}
	case "admin.user_blocked":
		email, _ := event["user_email"].(string)
		return &Notification{
			Title:    "🚫 Compte bloqué",
			Body:     fmt.Sprintf("Compte bloqué: %s", email),
			Type:     "admin",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push"},
		}
	case "admin.user_unblocked":
		email, _ := event["user_email"].(string)
		return &Notification{
			Title:    "✅ Compte débloqué",
			Body:     fmt.Sprintf("Compte débloqué: %s", email),
			Type:     "admin",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push"},
		}
	case "admin.pin_unlocked":
		email, _ := event["user_email"].(string)
		return &Notification{
			Title:    "🔓 PIN débloqué",
			Body:     fmt.Sprintf("PIN débloqué pour: %s", email),
			Type:     "admin",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push"},
		}
	case "admin.transaction_blocked":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		return &Notification{
			Title:    "🛑 Transaction bloquée",
			Body:     fmt.Sprintf("Transaction de %.2f %s bloquée", amount, currency),
			Type:     "admin",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push"},
		}
	case "admin.kyc_approved":
		email, _ := event["user_email"].(string)
		return &Notification{
			Title:    "✅ KYC Approuvé",
			Body:     fmt.Sprintf("KYC approuvé pour: %s", email),
			Type:     "admin",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push"},
		}
	case "admin.kyc_rejected":
		email, _ := event["user_email"].(string)
		return &Notification{
			Title:    "❌ KYC Rejeté",
			Body:     fmt.Sprintf("KYC rejeté pour: %s", email),
			Type:     "admin",
			Priority: PriorityNormal,
			Data:     event,
			Channels: []string{"push"},
		}

	// ===== PAYMENT STATUS EVENTS =====
	case "payment.request":
		paymentType, _ := event["type"].(string)
		// Skip notification for ticket purchases (user initiated)
		if paymentType == "ticket_purchase" {
			return nil
		}
		
		amount, _ := event["debit_amount"].(float64)
		if amount == 0 {
			amount, _ = event["amount"].(float64)
		}
		currency, _ := event["currency"].(string)
		title, _ := event["title"].(string) 
		return &Notification{
			Title:    "Demande de paiement",
			Body:     fmt.Sprintf("Nouvelle demande de paiement: %s (%.2f %s)", title, amount, currency),
			Type:     "payment",
			Priority: PriorityNormal,
			Data:     event,
			ActionUrl: fmt.Sprintf("/payment-requests/%s", event["request_id"]),
			Channels: []string{"push", "email"},
		}

	case "payment.status.success":
		amount, _ := event["amount"].(float64) // This comes from PaymentStatusEvent? No, status event doesn't have amount!
		// PaymentStatusEvent only has RequestID, ReferenceID, Status, Error.
		// It does NOT have Amount. 
		// So payment.status.success will ALSO show 0.00 if it relies on "amount".
		// We need to fetch the amount or include it in status event. 
		// For now, let's just say "Votre paiement a été validé." without amount if 0.
		
		body := "Votre paiement a été validé avec succès."
		if amount > 0 {
			currency, _ := event["currency"].(string)
			body = fmt.Sprintf("Votre paiement de %.2f %s a été validé avec succès.", amount, currency)
		}
		
		return &Notification{
			Title:    "Paiement validé",
			Body:     body,
			Type:     "payment",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"},
		}

	case "payment.status.failed":
		errMsg, _ := event["error"].(string)
		return &Notification{
			Title:    "❌ Paiement échoué",
			Body:     fmt.Sprintf("Votre paiement a échoué: %s", errMsg),
			Type:     "payment",
			Priority: PriorityHigh,
			Data:     event,
			Channels: []string{"push", "email"},
		}

	// ===== DIRECT NOTIFICATION EVENTS =====
	case "notification.created":
		title, _ := event["title"].(string)
		message, _ := event["message"].(string)
		notifType, _ := event["type"].(string)
		// Extract Data if present, handled by map decoding usually?
		// event["data"] might be a map
		var data map[string]interface{}
		if d, ok := event["data"].(map[string]interface{}); ok {
			data = d
		}
		
		actionUrl, _ := event["action_url"].(string)

		return &Notification{
			Title:    title,
			Body:     message,
			Type:     notifType,
			Priority: PriorityNormal,
			Data:     data,
			ActionUrl: actionUrl,
			Channels: []string{"push", "email"},
		}
	case "ticket.created":
		ticketID, _ := event["ticket_id"].(string)
		eventID, _ := event["event_id"].(string)
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		reason, _ := event["reason"].(string)

		body := fmt.Sprintf("Nouveau ticket créé pour l'événement %s. Montant: %.2f %s.", eventID, amount, currency)
		if reason != "" {
			body += fmt.Sprintf(" Raison: %s", reason)
		}

		return &Notification{
			Title:    "🎫 Nouveau Ticket",
			Body:     body,
			Type:     "ticket",
			Priority: PriorityNormal,
			Data: map[string]interface{}{
				"event_id":  eventID,
				"amount":    amount,
				"currency":  currency,
				"reason":    reason,
			},
			ActionUrl: fmt.Sprintf("/tickets/%s", ticketID),
			Channels: []string{"push", "email"},
		}

	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func (s *NotificationService) sendNotification(userID, email, phone string, notification *Notification) {
	for _, channel := range notification.Channels {
		switch channel {
		case "email":
			if email != "" {
				go s.sendEmail(email, notification)
			}
		case "sms":
			if phone != "" {
				go s.sendSMS(phone, notification)
			}
		case "push":
			if userID != "" {
				go s.sendPush(userID, notification)
			}
		}
	}

	// Log the notification
	s.logNotification(userID, notification)
}

func (s *NotificationService) logNotification(userID string, notification *Notification) {
	if s.repo == nil {
		return
	}
	
	dataStr := s.mapToString(notification.Data)
	err := s.repo.Create(&models.Notification{
		UserID:    userID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Body,
		Data:      &dataStr,
		CreatedAt: time.Now(),
	})
	
	if err != nil {
		log.Printf("Failed to log notification to DB: %v", err)
	}
}

func (s *NotificationService) mapToString(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// ===== EMAIL IMPLEMENTATION =====

func (s *NotificationService) sendEmail(toEmail string, notification *Notification) {
	if s.config.SMTPHost == "" || s.config.SMTPUser == "" {
		log.Println("[EMAIL] SMTP not configured, skipping")
		return
	}

	// Create HTML email
	htmlBody := s.createEmailHTML(notification)
	
	// Prepare email headers
	headers := make(map[string]string)
	headers["From"] = s.config.FromEmail
	headers["To"] = toEmail
	headers["Subject"] = notification.Title
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	
	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)
	
	// Connect to SMTP server
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	
	// Use TLS for Gmail and other providers
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	}
	
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		// Fallback to regular SMTP
		err = smtp.SendMail(addr, auth, s.config.FromEmail, []string{toEmail}, []byte(message.String()))
		if err != nil {
			log.Printf("[EMAIL] Failed to send: %v", err)
			return
		}
	} else {
		defer conn.Close()
		client, err := smtp.NewClient(conn, s.config.SMTPHost)
		if err != nil {
			log.Printf("[EMAIL] Failed to create client: %v", err)
			return
		}
		
		if err = client.Auth(auth); err != nil {
			log.Printf("[EMAIL] Auth failed: %v", err)
			return
		}
		
		if err = client.Mail(s.config.FromEmail); err != nil {
			log.Printf("[EMAIL] Mail failed: %v", err)
			return
		}
		
		if err = client.Rcpt(toEmail); err != nil {
			log.Printf("[EMAIL] Rcpt failed: %v", err)
			return
		}
		
		w, err := client.Data()
		if err != nil {
			log.Printf("[EMAIL] Data failed: %v", err)
			return
		}
		
		_, err = w.Write([]byte(message.String()))
		if err != nil {
			log.Printf("[EMAIL] Write failed: %v", err)
			return
		}
		
		w.Close()
		client.Quit()
	}

	log.Printf("[EMAIL] ✅ Sent to %s: %s", toEmail, notification.Title)
}

func (s *NotificationService) createEmailHTML(notification *Notification) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 20px; }
		.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 10px; overflow: hidden; }
		.header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
		.content { padding: 30px; }
		.footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
		.priority-critical { border-left: 4px solid #dc3545; }
		.priority-high { border-left: 4px solid #ffc107; }
	</style>
</head>
<body>
	<div class="container {{if eq .Priority "critical"}}priority-critical{{else if eq .Priority "high"}}priority-high{{end}}">
		<div class="header">
			<h1>{{.Title}}</h1>
		</div>
		<div class="content">
			<p>{{.Body}}</p>
		</div>
		<div class="footer">
			<p>Zekora - Votre banque digitale</p>
			<p>Cet email a été envoyé automatiquement, ne pas répondre.</p>
		</div>
	</div>
</body>
</html>`
	
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return notification.Body
	}
	
	var buf bytes.Buffer
	t.Execute(&buf, notification)
	return buf.String()
}

// ===== SMS IMPLEMENTATION (TWILIO) =====

func (s *NotificationService) sendSMS(toPhone string, notification *Notification) {
	if s.config.TwilioAccountSID == "" || s.config.TwilioAuthToken == "" {
		log.Println("[SMS] Twilio not configured, skipping")
		return
	}

	// Twilio API endpoint
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.config.TwilioAccountSID)
	
	// Prepare message (SMS has 160 char limit)
	message := fmt.Sprintf("%s: %s", notification.Title, notification.Body)
	if len(message) > 160 {
		message = message[:157] + "..."
	}
	
	// Prepare form data
	data := url.Values{}
	data.Set("To", toPhone)
	data.Set("From", s.config.TwilioFromNumber)
	data.Set("Body", message)
	
	// Create request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("[SMS] Failed to create request: %v", err)
		return
	}
	
	req.SetBasicAuth(s.config.TwilioAccountSID, s.config.TwilioAuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[SMS] Failed to send: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		log.Printf("[SMS] Twilio error: status %d", resp.StatusCode)
		return
	}

	log.Printf("[SMS] ✅ Sent to %s: %s", toPhone, notification.Title)
}

// ===== PUSH NOTIFICATION (FCM) =====

func (s *NotificationService) sendPush(userID string, notification *Notification) {
	if s.config.FCMServerKey == "" {
		log.Println("[PUSH] FCM not configured, skipping")
		return
	}

	// In production, you would look up the user's FCM token from database
	// For now, we'll just log it
	log.Printf("[PUSH] Would send to user %s: %s", userID, notification.Title)
	
	// FCM HTTP v1 API
	// Note: In production, use the Firebase Admin SDK
	fcmURL := "https://fcm.googleapis.com/fcm/send"
	
	payload := map[string]interface{}{
		"to": "/topics/user_" + userID, // Topic-based for simplicity
		"notification": map[string]string{
			"title": notification.Title,
			"body":  notification.Body,
		},
		"data": func() map[string]interface{} {
			// Copy data and add action_url
			d := make(map[string]interface{})
			for k, v := range notification.Data {
				d[k] = v
			}
			if notification.ActionUrl != "" {
				d["action_url"] = notification.ActionUrl
			}
			return d
		}(),
	}
	
	body, _ := json.Marshal(payload)
	
	req, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[PUSH] Failed to create request: %v", err)
		return
	}
	
	req.Header.Set("Authorization", "key="+s.config.FCMServerKey)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[PUSH] Failed to send: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		log.Printf("[PUSH] FCM error: status %d", resp.StatusCode)
		return
	}
	
	log.Printf("[PUSH] ✅ Sent to user %s", userID)
}
