package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
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

// InitializeKafka creates a new Kafka client for messaging
func InitializeKafka(brokers string, groupID string) (*messaging.KafkaClient, error) {
	brokerList := strings.Split(brokers, ",")
	
	client := messaging.NewKafkaClient(brokerList, groupID)
	
	log.Printf("[Kafka] Exchange-service connected to brokers: %s with group: %s", brokers, groupID)
	return client, nil
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
