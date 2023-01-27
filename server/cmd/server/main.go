package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nkoblai/rabbitmqapp/server/config"
	"github.com/nkoblai/rabbitmqapp/server/internal/queue"
	"github.com/nkoblai/rabbitmqapp/server/internal/service"
	"go.uber.org/zap"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't create logger: %e\n", err)
	}

	cfg, err := config.NewRabbitMQ()
	if err != nil {
		logger.Fatal("wrong request type", zap.Error(fmt.Errorf("can't create config: %w", err)))
	}

	rabbitMQ, err := queue.NewRabbitMQ(logger, cfg)
	if err != nil {
		logger.Fatal("can't create queue", zap.Error(err))
	}

	serverLogFile, err := os.Create("server_logs.txt")
	if err != nil {
		logger.Fatal("can't create server log file", zap.Error(err))
	}

	defer serverLogFile.Close()

	logger.Info("Server started...")
	defer logger.Info("Server stopped!")

	if err = service.New(logger, serverLogFile, rabbitMQ).Process(ctx); err != nil {
		logger.Fatal("can't start processing requests", zap.Error(err))
	}

}
