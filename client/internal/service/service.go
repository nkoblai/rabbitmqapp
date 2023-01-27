package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Item is item struct.
type Item struct {
	ClientID    uuid.UUID `json:"clientID"`
	Value       string    `json:"value,omitempty"`
	RequestType string    `json:"requestType"`
}

// Queue is queue interface
type Queue interface {
	Add(ctx context.Context, body []byte) error
}

// Service is service struct.
type Service struct {
	logger *zap.Logger
	queue  Queue
}

// New return new service instance.
func New(logger *zap.Logger, queue Queue) *Service {
	return &Service{
		logger: logger,
		queue:  queue,
	}
}

// Request makes request.
func (s *Service) Request(ctx context.Context, item Item) error {
	body, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshaling failed: %w", err)
	}
	if err = s.queue.Add(ctx, body); err != nil {
		return fmt.Errorf("add item to queue failed: %w", err)
	}
	return nil
}
