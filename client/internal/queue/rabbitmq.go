package queue

import (
	"context"
	"fmt"

	"github.com/nkoblai/rabbitmqapp/client/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ is rabbitMQ struct.
type RabbitMQ struct {
	ch  *amqp.Channel
	cfg config.RabbitMQ
}

// NewRabbitMQ returns new RabbitMQ instance.
func NewRabbitMQ(cfg config.RabbitMQ) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("connection creation failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel connection failed: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("channel connection failed: %w", err)
	}

	return &RabbitMQ{
		ch:  ch,
		cfg: cfg,
	}, nil
}

// Add adds item.
func (rmq *RabbitMQ) Add(ctx context.Context, body []byte) error {
	if err := rmq.ch.PublishWithContext(
		ctx,
		"",
		rmq.cfg.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return fmt.Errorf("publishing failed: %w", err)
	}
	return nil
}
