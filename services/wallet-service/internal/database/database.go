package database

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	_ "github.com/lib/pq"
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

func InitializeRabbitMQ(rabbitMQURL string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchanges
	exchanges := []string{
		"wallet.events",
		"transaction.events",
		"transfer.events",
	}

	for _, exchange := range exchanges {
		err = channel.ExchangeDeclare(
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

	return channel, nil
}