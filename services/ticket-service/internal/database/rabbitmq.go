package database

import (
	"fmt"
	"github.com/streadway/amqp"
)

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

	// Declare payment.events exchange
	err = ch.ExchangeDeclare(
		"payment.events", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue for payment status updates (unique for this service instance? or shared?)
	// If we use load balancing, we want a queue that all instances share to process results.
	// But `payment.events` might be broadcast.
	// We want to listen to `payment.status.*` but specifically for OUR transactions?
	// or just all and filter?
	// Let's create a queue "ticket.payment.updates" bound to "payment.status.#"
	
	_, err = ch.QueueDeclare(
		"ticket.payment.updates", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		"ticket.payment.updates",
		"payment.status.#",
		"payment.events",
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

func (r *RabbitMQClient) PublishToExchange(exchange, routingKey string, message []byte) error {
	return r.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

func (r *RabbitMQClient) Consume(queue string) (<-chan amqp.Delivery, error) {
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
