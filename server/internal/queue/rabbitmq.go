package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nkoblai/rabbitmqapp/server/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Item is item struct.
type Item struct {
	ClientID    uuid.UUID `json:"clientID"`
	Value       string    `json:"value,omitempty"`
	RequestType string    `json:"requestType"`
}

// RabbitMQ is rabbitMQ struct.
type RabbitMQ struct {
	ch     *amqp.Channel
	conn   *amqp.Connection
	cfg    config.RabbitMQ
	logger *zap.Logger
}

// NewRabbitMQ returns new RabbitMQ instance.
func NewRabbitMQ(logger *zap.Logger, cfg config.RabbitMQ) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("connection creation failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel connection failed: %w", err)
	}

	if _, err = ch.QueueDeclare(
		cfg.QueueName,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("can't create queue: %s: error %w", cfg.QueueName, err)
	}

	return &RabbitMQ{
		ch:     ch,
		cfg:    cfg,
		logger: logger,
	}, nil
}

// Requests starts consuming requests.
func (rmq *RabbitMQ) Requests(ctx context.Context) (<-chan Item, error) {
	msgs, err := rmq.ch.Consume(
		rmq.cfg.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("can't consume messages: %w", err)
	}

	ch := make(chan Item)

	go func() {
		for {
			select {
			case <-ctx.Done():
				rmq.ch.Close()
				return
			case d := <-msgs:
				var item Item
				if err = json.Unmarshal(d.Body, &item); err != nil {
					rmq.logger.Warn(fmt.Sprintf("unmarshaling message: %s failed: %v", d.Body, err))
					continue
				}
				ch <- item
			}
		}
	}()

	return ch, nil
}
