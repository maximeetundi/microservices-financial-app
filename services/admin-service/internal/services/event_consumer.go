package services

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// EventConsumer listens for events from other services and creates admin notifications
type EventConsumer struct {
	channel *amqp.Channel
	db      *sql.DB
}

// Event represents an incoming event from RabbitMQ
type Event struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(channel *amqp.Channel, db *sql.DB) *EventConsumer {
	return &EventConsumer{channel: channel, db: db}
}

// StartConsuming starts consuming events from various exchanges
func (c *EventConsumer) StartConsuming() error {
	// Declare queue for admin notifications
	queue, err := c.channel.QueueDeclare(
		"admin.notifications", // queue name
		true,                  // durable
		false,                 // auto-delete
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return err
	}

	// Bind to various events we want to receive
	bindings := []struct {
		exchange   string
		routingKey string
	}{
		// KYC events
		{"user.events", "kyc.submitted"},
		{"user.events", "kyc.updated"},
		// PIN events
		{"user.events", "pin.blocked"},
		{"user.events", "pin.unlocked"},
		// Support events
		{"support.events", "ticket.created"},
		{"support.events", "ticket.message"},
		{"support.events", "ticket.escalated"},
		// Transaction events (for suspicious activity)
		{"transaction.events", "transaction.flagged"},
		{"transaction.events", "transaction.large"},
	}

	for _, b := range bindings {
		// Ensure exchange exists
		c.channel.ExchangeDeclare(b.exchange, "topic", true, false, false, false, nil)
		
		err := c.channel.QueueBind(queue.Name, b.routingKey, b.exchange, false, nil)
		if err != nil {
			log.Printf("Warning: Could not bind to %s/%s: %v", b.exchange, b.routingKey, err)
		}
	}

	// Start consuming
	msgs, err := c.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			c.processEvent(msg.RoutingKey, msg.Body)
		}
	}()

	log.Println("Admin notification consumer started, listening for events...")
	return nil
}

// processEvent processes an incoming event and creates a notification
func (c *EventConsumer) processEvent(routingKey string, body []byte) {
	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to parse event: %v", err)
		return
	}

	var notifType, title, message string
	var data map[string]interface{}

	switch routingKey {
	case "kyc.submitted":
		notifType = "kyc"
		userName := getStringFromMap(event.Data, "user_name", "Un utilisateur")
		userEmail := getStringFromMap(event.Data, "user_email", "")
		title = "📋 Nouvelle demande KYC"
		message = userName + " a soumis des documents KYC pour vérification"
		if userEmail != "" {
			message += " (" + userEmail + ")"
		}
		data = event.Data

	case "pin.blocked":
		notifType = "security"
		userName := getStringFromMap(event.Data, "user_name", "Un utilisateur")
		userEmail := getStringFromMap(event.Data, "user_email", "")
		title = "🔒 PIN bloqué"
		message = "Le PIN de " + userName + " a été bloqué après trop de tentatives"
		if userEmail != "" {
			message += " (" + userEmail + ")"
		}
		data = event.Data

	case "ticket.created":
		notifType = "support"
		userName := getStringFromMap(event.Data, "user_name", "Un utilisateur")
		subject := getStringFromMap(event.Data, "subject", "")
		title = "💬 Nouveau ticket support"
		message = userName + " a créé un nouveau ticket"
		if subject != "" {
			message += ": " + subject
		}
		data = event.Data

	case "ticket.message":
		notifType = "support"
		userName := getStringFromMap(event.Data, "user_name", "Un utilisateur")
		title = "💬 Nouveau message support"
		message = userName + " a répondu à un ticket"
		data = event.Data

	case "transaction.flagged":
		notifType = "security"
		amount := getStringFromMap(event.Data, "amount", "")
		title = "⚠️ Transaction suspecte"
		message = "Une transaction a été signalée pour vérification"
		if amount != "" {
			message += " (montant: " + amount + ")"
		}
		data = event.Data

	case "transaction.large":
		notifType = "transaction"
		amount := getStringFromMap(event.Data, "amount", "")
		title = "💰 Transaction importante"
		message = "Une transaction de montant élevé a été effectuée"
		if amount != "" {
			message += ": " + amount
		}
		data = event.Data

	case "ticket.escalated":
		notifType = "support"
		userName := getStringFromMap(event.Data, "user_name", "Un utilisateur")
		title = "🚨 Escalade support"
		message = userName + " demande à parler à un agent humain"
		data = event.Data

	default:
		log.Printf("Unknown event type: %s", routingKey)
		return
	}

	// Create notification in database
	c.createNotification(notifType, title, message, data)
}

// createNotification inserts a notification into the database
func (c *EventConsumer) createNotification(notifType, title, message string, data map[string]interface{}) {
	dataJSON, _ := json.Marshal(data)
	id := uuid.New().String()

	_, err := c.db.Exec(`
		INSERT INTO admin_notifications (id, type, title, message, data)
		VALUES ($1, $2, $3, $4, $5::jsonb)
	`, id, notifType, title, message, string(dataJSON))

	if err != nil {
		log.Printf("Failed to create notification: %v", err)
	} else {
		log.Printf("Created admin notification: [%s] %s", notifType, title)
	}
}

// Helper function to safely get string from map
func getStringFromMap(m map[string]interface{}, key, defaultValue string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}
