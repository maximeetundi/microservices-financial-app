package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func Initialize(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
	mu      sync.Mutex
}

func InitializeRabbitMQ(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare all exchanges
	exchanges := []string{
		"wallet.events",
		"transaction.events",
		"transfer.events",
		"exchange.events",
		"card.events",
		"notification.events",
		"payment.events",
	}

	for _, exchange := range exchanges {
		err = ch.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare exchange %s: %w", exchange, err)
		}
	}

	// Declare queues
	queues := []string{
		"transfers",
		"notifications",
		"compliance",
		"transfer.wallet_updates",
		"payments", // New queue for payment requests
	}
	for _, q := range queues {
		_, err = ch.QueueDeclare(q, true, false, false, false, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to declare queue %s: %w", q, err)
		}
	}

	// Bind payments queue to payment.events (for payment.request)
	err = ch.QueueBind(
		"payments",
		"payment.request",
		"payment.events",
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind payments queue: %w", err)
	}

	// Bind transfer.wallet_updates to wallet.events for balance confirmations
	err = ch.QueueBind(
		"transfer.wallet_updates",
		"wallet.*",
		"wallet.events",
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	log.Println("[RabbitMQ] Transfer-service connected successfully")
	return &RabbitMQClient{conn: conn, channel: ch, url: url}, nil
}

func (r *RabbitMQClient) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// forceReconnect forces a reconnection to RabbitMQ
func (r *RabbitMQClient) forceReconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	log.Println("[RabbitMQ] Forcing reconnection...")
	
	// Close existing connection if any
	if r.channel != nil {
		r.channel.Close()
		r.channel = nil
	}
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
	
	// Reconnect
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("failed to reconnect to RabbitMQ: %w", err)
	}
	
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to reopen RabbitMQ channel: %w", err)
	}
	
	// Redeclare exchanges
	exchanges := []string{"wallet.events", "transaction.events", "transfer.events", 
		"exchange.events", "card.events", "notification.events", "payment.events"}
	for _, exchange := range exchanges {
		err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
		if err != nil {
			channel.Close()
			conn.Close()
			return fmt.Errorf("failed to redeclare exchange %s: %w", exchange, err)
		}
	}
	
	// Redeclare queues
	queues := []string{"transfers", "notifications", "compliance", "transfer.wallet_updates", "payments"}
	for _, q := range queues {
		_, err = channel.QueueDeclare(q, true, false, false, false, nil)
		if err != nil {
			channel.Close()
			conn.Close()
			return fmt.Errorf("failed to redeclare queue %s: %w", q, err)
		}
	}
	
	// Rebind queues
	channel.QueueBind("payments", "payment.request", "payment.events", false, nil)
	channel.QueueBind("transfer.wallet_updates", "wallet.*", "wallet.events", false, nil)
	
	r.conn = conn
	r.channel = channel
	log.Println("[RabbitMQ] Reconnected successfully!")
	
	return nil
}

func (r *RabbitMQClient) Publish(queue string, message []byte) error {
	r.mu.Lock()
	err := r.channel.Publish("", queue, false, false, amqp.Publishing{ContentType: "application/json", Body: message})
	r.mu.Unlock()
	
	if err != nil {
		log.Printf("[RabbitMQ] Publish failed: %v, attempting reconnect...", err)
		if reconnErr := r.forceReconnect(); reconnErr != nil {
			return fmt.Errorf("publish failed and reconnect failed: %w", reconnErr)
		}
		r.mu.Lock()
		err = r.channel.Publish("", queue, false, false, amqp.Publishing{ContentType: "application/json", Body: message})
		r.mu.Unlock()
		if err != nil {
			return fmt.Errorf("publish failed after reconnect: %w", err)
		}
		log.Println("[RabbitMQ] Publish succeeded after reconnect")
	}
	return nil
}

// PublishToExchange publishes a message to a specific exchange with routing key
func (r *RabbitMQClient) PublishToExchange(exchange, routingKey string, message []byte) error {
	r.mu.Lock()
	err := r.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	r.mu.Unlock()
	
	if err != nil {
		log.Printf("[RabbitMQ] Publish to exchange failed: %v, attempting reconnect...", err)
		if reconnErr := r.forceReconnect(); reconnErr != nil {
			return fmt.Errorf("publish failed and reconnect failed: %w", reconnErr)
		}
		r.mu.Lock()
		err = r.channel.Publish(
			exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        message,
			},
		)
		r.mu.Unlock()
		if err != nil {
			return fmt.Errorf("publish failed after reconnect: %w", err)
		}
		log.Println("[RabbitMQ] Publish to exchange succeeded after reconnect")
	}
	return nil
}

// Consume consumes messages from a queue
func (r *RabbitMQClient) Consume(queue string) (<-chan amqp.Delivery, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	return r.channel.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}
