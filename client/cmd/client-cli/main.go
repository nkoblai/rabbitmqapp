package main

import (
	"context"
	"log"

	"github.com/nkoblai/rabbitmqapp/client/config"
	"github.com/nkoblai/rabbitmqapp/client/internal/queue"
	"github.com/nkoblai/rabbitmqapp/client/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't create zap logger instance: %v\n", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("can't parse config", zap.Error(err))
	}

	rabbitMQ, err := queue.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		logger.Fatal("can't create queue", zap.Error(err))
	}

	if err = service.New(logger, rabbitMQ).Request(context.Background(), service.Item{
		ClientID:    cfg.App.ClientID,
		RequestType: cfg.App.RequestType,
		Value:       cfg.App.ItemValue,
	}); err != nil {
		logger.Fatal("can't make request", zap.Error(err))
	}
}
