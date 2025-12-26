package services

import (
	"bytes"
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

	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/repository"
	"github.com/streadway/amqp"
)

// NotificationService handles all notification types
type NotificationService struct {
	channel *amqp.Channel
	config  *config.Config
	repo    *repository.NotificationRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(channel *amqp.Channel, cfg *config.Config, repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		channel: channel,
		config:  cfg,
		repo:    repo,
	}
}

// Start begins consuming from all notification queues
func (s *NotificationService) Start() error {
	// Declare exchanges if they don't exist
	exchanges := []string{
		"wallet.events",
		"transaction.events",
		"transfer.events",
		"exchange.events",
		"card.events",
		"auth.events",
	}
	
	for _, exchange := range exchanges {
		err := s.channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
		if err != nil {
			log.Printf("Warning: failed to declare exchange %s: %v", exchange, err)
		}
	}
	
	// Subscribe to all event exchanges
	queues := map[string][]string{
		"notification.wallet":      {"wallet.events", "#"},
		"notification.transaction": {"transaction.events", "#"},
		"notification.transfer":    {"transfer.events", "#"},
		"notification.exchange":    {"exchange.events", "#"},
		"notification.card":        {"card.events", "#"},
		"notification.auth":        {"auth.events", "#"},
	}

	for queue, binding := range queues {
		// Declare queue
		_, err := s.channel.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			log.Printf("Warning: failed to declare queue %s: %v", queue, err)
			continue
		}

		// Bind to exchange
		err = s.channel.QueueBind(queue, binding[1], binding[0], false, nil)
		if err != nil {
			log.Printf("Warning: failed to bind queue %s: %v", queue, err)
			continue
		}

		// Start consumer
		go s.consume(queue)
	}

	log.Println("✅ Notification service consumers started")
	return nil
}

func (s *NotificationService) consume(queue string) {
	msgs, err := s.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to consume from %s: %v", queue, err)
		return
	}

	for msg := range msgs {
		s.handleMessage(queue, msg)
	}
}

func (s *NotificationService) handleMessage(queue string, msg amqp.Delivery) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		msg.Nack(false, true)
		return
	}

	eventType, _ := event["type"].(string)
	userID, _ := event["user_id"].(string)
	email, _ := event["email"].(string)
	phone, _ := event["phone"].(string)

	log.Printf("📬 Processing event: %s for user: %s", eventType, userID)

	// Determine notification based on event type
	notification := s.createNotification(eventType, event)
	if notification == nil {
		msg.Ack(false)
		return
	}

	// Send notifications via different channels
	if userID != "" || email != "" || phone != "" {
		s.sendNotification(userID, email, phone, notification)
	}

	msg.Ack(false)
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
		"data": notification.Data,
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

	log.Printf("[PUSH] ✅ Sent to user %s: %s", userID, notification.Title)
}

func (s *NotificationService) logNotification(userID string, notification *Notification) {
	// Save notification to database for user notification center
	if s.repo != nil && userID != "" {
		// Convert data to JSON string
		var dataStr *string
		if notification.Data != nil {
			dataBytes, err := json.Marshal(notification.Data)
			if err == nil {
				str := string(dataBytes)
				dataStr = &str
			}
		}

		dbNotification := &models.Notification{
			UserID:  userID,
			Type:    notification.Type,
			Title:   notification.Title,
			Message: notification.Body,
			Data:    dataStr,
		}

		if err := s.repo.Create(dbNotification); err != nil {
			log.Printf("[LOG] ❌ Failed to save notification: %v", err)
		} else {
			log.Printf("[LOG] ✅ Notification saved: user=%s, type=%s, title=%s", 
				userID, notification.Type, notification.Title)
		}
	} else {
		log.Printf("[LOG] 📝 Notification logged (no DB): user=%s, type=%s, priority=%s, title=%s", 
			userID, notification.Type, notification.Priority, notification.Title)
	}
}
