package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/config"
	"github.com/streadway/amqp"
)

// NotificationService handles all notification types
type NotificationService struct {
	channel *amqp.Channel
	config  *config.Config
}

// NewNotificationService creates a new notification service
func NewNotificationService(channel *amqp.Channel, cfg *config.Config) *NotificationService {
	return &NotificationService{
		channel: channel,
		config:  cfg,
	}
}

// Start begins consuming from all notification queues
func (s *NotificationService) Start() error {
	// Subscribe to all event exchanges
	queues := map[string][]string{
		"notification.wallet":      {"wallet.events", "#"},
		"notification.transaction": {"transaction.events", "#"},
		"notification.transfer":    {"transfer.events", "#"},
		"notification.exchange":    {"exchange.events", "#"},
		"notification.card":        {"card.events", "#"},
	}

	for queue, binding := range queues {
		// Declare queue
		_, err := s.channel.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queue, err)
		}

		// Bind to exchange
		err = s.channel.QueueBind(queue, binding[1], binding[0], false, nil)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queue, err)
		}

		// Start consumer
		go s.consume(queue)
	}

	log.Println("Notification service consumers started")
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

	log.Printf("Processing notification event: %s for user: %s", eventType, userID)

	// Determine notification based on event type
	notification := s.createNotification(eventType, event)
	if notification == nil {
		msg.Ack(false)
		return
	}

	// Send notifications via different channels
	if userID != "" {
		s.sendNotification(userID, notification)
	}

	msg.Ack(false)
}

type Notification struct {
	Title    string
	Body     string
	Type     string
	Data     map[string]interface{}
	Channels []string // email, sms, push
}

func (s *NotificationService) createNotification(eventType string, event map[string]interface{}) *Notification {
	switch eventType {
	case "wallet.created":
		return &Notification{
			Title:    "Wallet Created",
			Body:     fmt.Sprintf("Your %s wallet has been created successfully.", event["currency"]),
			Type:     "wallet",
			Data:     event,
			Channels: []string{"push", "email"},
		}
	case "transaction.completed":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		return &Notification{
			Title:    "Transaction Completed",
			Body:     fmt.Sprintf("Your transaction of %.2f %s has been completed.", amount, currency),
			Type:     "transaction",
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}
	case "transfer.completed":
		amount, _ := event["amount"].(float64)
		currency, _ := event["currency"].(string)
		return &Notification{
			Title:    "Transfer Completed",
			Body:     fmt.Sprintf("Transfer of %.2f %s completed successfully.", amount, currency),
			Type:     "transfer",
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}
	case "exchange.completed":
		fromCurrency, _ := event["from_currency"].(string)
		toCurrency, _ := event["to_currency"].(string)
		return &Notification{
			Title:    "Exchange Completed",
			Body:     fmt.Sprintf("Your exchange from %s to %s has been completed.", fromCurrency, toCurrency),
			Type:     "exchange",
			Data:     event,
			Channels: []string{"push", "email"},
		}
	case "card.created":
		return &Notification{
			Title:    "Card Created",
			Body:     "Your new card has been created successfully.",
			Type:     "card",
			Data:     event,
			Channels: []string{"push", "email"},
		}
	case "card.shipped":
		return &Notification{
			Title:    "Card Shipped",
			Body:     "Your physical card has been shipped!",
			Type:     "card",
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}
	case "wallet.frozen":
		return &Notification{
			Title:    "Wallet Frozen",
			Body:     "Your wallet has been frozen for security reasons.",
			Type:     "security",
			Data:     event,
			Channels: []string{"push", "email", "sms"},
		}
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func (s *NotificationService) sendNotification(userID string, notification *Notification) {
	for _, channel := range notification.Channels {
		switch channel {
		case "email":
			go s.sendEmail(userID, notification)
		case "sms":
			go s.sendSMS(userID, notification)
		case "push":
			go s.sendPush(userID, notification)
		}
	}

	// Log the notification
	s.logNotification(userID, notification)
}

func (s *NotificationService) sendEmail(userID string, notification *Notification) {
	// TODO: Implement email sending via SMTP
	log.Printf("[EMAIL] Sending to user %s: %s - %s", userID, notification.Title, notification.Body)

	if s.config.SMTPHost == "" || s.config.SMTPUser == "" {
		log.Println("[EMAIL] SMTP not configured, skipping")
		return
	}

	// Simulate sending
	time.Sleep(100 * time.Millisecond)
	log.Printf("[EMAIL] Sent successfully to user %s", userID)
}

func (s *NotificationService) sendSMS(userID string, notification *Notification) {
	// TODO: Implement SMS sending via Twilio
	log.Printf("[SMS] Sending to user %s: %s", userID, notification.Body)

	if s.config.TwilioAccountSID == "" {
		log.Println("[SMS] Twilio not configured, skipping")
		return
	}

	// Simulate sending
	time.Sleep(100 * time.Millisecond)
	log.Printf("[SMS] Sent successfully to user %s", userID)
}

func (s *NotificationService) sendPush(userID string, notification *Notification) {
	// TODO: Implement push notification via FCM/APNS
	log.Printf("[PUSH] Sending to user %s: %s - %s", userID, notification.Title, notification.Body)

	if s.config.FCMServerKey == "" {
		log.Println("[PUSH] FCM not configured, skipping")
		return
	}

	// Simulate sending
	time.Sleep(50 * time.Millisecond)
	log.Printf("[PUSH] Sent successfully to user %s", userID)
}

func (s *NotificationService) logNotification(userID string, notification *Notification) {
	// TODO: Save notification to database for user notification center
	log.Printf("[LOG] Notification logged for user %s: type=%s, title=%s", userID, notification.Type, notification.Title)
}
