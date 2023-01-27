package config

import (
	"flag"
	"fmt"
)

const (
	rabbitMQHostArg      = "rmqh"
	rabbitMQQueueNameArg = "rmqn"

	isEmpty = "%s is empty"
)

// RabbitMQ is RabbitMQ config.
type RabbitMQ struct {
	QueueName string
	Host      string
}

// New create new rabbit mq config.
func NewRabbitMQ() (rmq RabbitMQ, err error) {

	flag.StringVar(&rmq.Host, "rmqh", "", "RabbitMQ host")
	flag.StringVar(&rmq.QueueName, "rmqn", "", "RabbitMQ queue name")

	flag.Parse()

	if rmq.Host == "" {
		err = fmt.Errorf(isEmpty, rabbitMQHostArg)
		return
	}

	if rmq.QueueName == "" {
		err = fmt.Errorf(isEmpty, rabbitMQQueueNameArg)
		return
	}

	return
}
