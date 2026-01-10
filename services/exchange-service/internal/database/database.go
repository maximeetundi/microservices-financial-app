package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	_ "github.com/lib/pq"
)

// Initialize PostgreSQL database
func Initialize(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database connected and initialized successfully")
	return db, nil
}

// Initialize Redis client
func InitializeRedis(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return client, nil
}

// Initialize RabbitMQ connection
func InitializeRabbitMQ(rabbitURL string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	// Declare exchanges and queues
	if err := declareExchangesAndQueues(channel); err != nil {
		return nil, fmt.Errorf("failed to declare exchanges and queues: %w", err)
	}

	log.Println("RabbitMQ connected and configured successfully")
	return &RabbitMQClient{conn: conn, channel: channel, url: rabbitURL}, nil
}

// Helper methods for RabbitMQ

type RabbitMQClient struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	url       string
	mu        sync.Mutex
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

// GetChannel returns the underlying amqp.Channel
func (r *RabbitMQClient) GetChannel() *amqp.Channel {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.channel
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
	
	// Redeclare exchanges and queues
	if err := declareExchangesAndQueues(channel); err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("failed to redeclare exchanges: %w", err)
	}
	
	r.conn = conn
	r.channel = channel
	log.Println("[RabbitMQ] Reconnected successfully!")
	
	return nil
}

func (r *RabbitMQClient) PublishToExchange(exchange, routingKey string, message []byte) error {
	r.mu.Lock()
	
	// Try to publish
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
	
	// If publish failed, try to reconnect and retry once
	if err != nil {
		log.Printf("[RabbitMQ] Publish failed: %v, attempting reconnect...", err)
		
		if reconnErr := r.forceReconnect(); reconnErr != nil {
			return fmt.Errorf("publish failed and reconnect failed: %w", reconnErr)
		}
		
		// Retry the publish
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
		
		log.Println("[RabbitMQ] Publish succeeded after reconnect")
	}
	
	return nil
}

func (r *RabbitMQClient) Consume(queue string) (<-chan amqp.Delivery, error) {
	// Ensure connection is alive
	if err := r.ensureConnection(); err != nil {
		return nil, err
	}
	
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


// Create necessary database tables
func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS exchanges (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR NOT NULL,
			from_wallet_id VARCHAR,
			to_wallet_id VARCHAR,
			from_currency VARCHAR NOT NULL,
			to_currency VARCHAR NOT NULL,
			from_amount DECIMAL NOT NULL,
			to_amount DECIMAL NOT NULL,
			exchange_rate DECIMAL NOT NULL,
			fee DECIMAL NOT NULL,
			fee_percentage DECIMAL NOT NULL,
			status VARCHAR NOT NULL DEFAULT 'pending',
			destination_amount DECIMAL,
			destination_currency VARCHAR,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			completed_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS exchange_rates (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			from_currency VARCHAR NOT NULL,
			to_currency VARCHAR NOT NULL,
			rate DECIMAL NOT NULL,
			bid_price DECIMAL NOT NULL,
			ask_price DECIMAL NOT NULL,
			spread DECIMAL NOT NULL,
			source VARCHAR NOT NULL,
			volume_24h DECIMAL DEFAULT 0,
			change_24h DECIMAL DEFAULT 0,
			last_updated TIMESTAMP DEFAULT NOW(),
			UNIQUE(from_currency, to_currency)
		)`,
		`CREATE TABLE IF NOT EXISTS quotes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR NOT NULL,
			from_currency VARCHAR NOT NULL,
			to_currency VARCHAR NOT NULL,
			from_amount DECIMAL NOT NULL,
			to_amount DECIMAL NOT NULL,
			exchange_rate DECIMAL NOT NULL,
			fee DECIMAL NOT NULL,
			fee_percentage DECIMAL NOT NULL,
			valid_until TIMESTAMP NOT NULL,
			estimated_delivery VARCHAR,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS trading_orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR NOT NULL,
			order_type VARCHAR NOT NULL,
			pair VARCHAR NOT NULL,
			side VARCHAR NOT NULL,
			amount DECIMAL NOT NULL,
			price DECIMAL,
			stop_price DECIMAL,
			status VARCHAR NOT NULL DEFAULT 'pending',
			filled_amount DECIMAL DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_exchanges_user_id ON exchanges(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_exchanges_status ON exchanges(status)`,
		`CREATE INDEX IF NOT EXISTS idx_rates_pair ON exchange_rates(from_currency, to_currency)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON trading_orders(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_status ON trading_orders(status)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return nil
}

// Declare RabbitMQ exchanges and queues
func declareExchangesAndQueues(ch *amqp.Channel) error {
	// Declare exchanges
	exchanges := []string{
		"exchange.events",
		"fiat_exchange.events",
		"trading.events",
		"rate.updates",
	}

	for _, exchange := range exchanges {
		err := ch.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", exchange, err)
		}
	}

	// Declare payment.events exchange
	err := ch.ExchangeDeclare(
		"payment.events", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange payment.events: %w", err)
	}

	// Declare queues
	queues := []string{
		"exchange.completed",
		"fiat_exchange.completed", 
		"trading.order_filled",
		"rate.crypto_updated",
		"rate.fiat_updated",
		"exchange.payment.updates", // Queue for payment status
	}

	for _, queue := range queues {
		_, err := ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queue, err)
		}
	}

	// Bind exchange.payment.updates to payment.events
	err = ch.QueueBind(
		"exchange.payment.updates",
		"payment.status.#",
		"payment.events",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue exchange.payment.updates: %w", err)
	}

	return nil
}


