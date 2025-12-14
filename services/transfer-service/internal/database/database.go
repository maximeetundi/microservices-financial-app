package database

import (
	"database/sql"
	"fmt"

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
	queues := []string{"transfers", "notifications", "compliance", "transfer.wallet_updates"}
	for _, q := range queues {
		_, err = ch.QueueDeclare(q, true, false, false, false, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to declare queue %s: %w", q, err)
		}
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

	return &RabbitMQClient{conn: conn, channel: ch}, nil
}

func (r *RabbitMQClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQClient) Publish(queue string, message []byte) error {
	return r.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
