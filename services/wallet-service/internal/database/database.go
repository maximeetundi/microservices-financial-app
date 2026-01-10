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

func Initialize(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

func InitializeRedis(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)
	
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}

// InitializeKafka creates a new Kafka client for messaging
func InitializeKafka(brokers string, groupID string) (*messaging.KafkaClient, error) {
	brokerList := strings.Split(brokers, ",")
	
	client := messaging.NewKafkaClient(brokerList, groupID)
	
	log.Printf("[Kafka] Wallet-service connected to brokers: %s with group: %s", brokers, groupID)
	return client, nil
}