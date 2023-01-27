package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/google/uuid"
)

const (
	addItem     requestType = "AddItem"
	removeItem  requestType = "RemoveItem"
	getItem     requestType = "GetItem"
	getAllItems requestType = "GetAllItems"

	rabbitMQHostArg      = "rmqh"
	rabbitMQQueueNameArg = "rmqn"

	isEmpty = "%s is empty"
)

var (
	errRequestTypeDoesntNotSupport = errors.New("request type doesnt not support")

	validRequestTypesMap = map[requestType]struct{}{
		addItem:     {},
		removeItem:  {},
		getItem:     {},
		getAllItems: {},
	}
)

type requestType string

// RabbitMQ is RabbitMQ config.
type RabbitMQ struct {
	QueueName string
	Host      string
}

// App is app config struct.
type App struct {
	ClientID    uuid.UUID
	RequestType string
	ItemValue   string
}

// Config is config struct.
type Config struct {
	RabbitMQ RabbitMQ
	App      App
}

// New create new config.
func New() (cfg Config, err error) {

	var clientID string

	flag.StringVar(&cfg.RabbitMQ.Host, rabbitMQHostArg, "", "RabbitMQ host")
	flag.StringVar(&cfg.RabbitMQ.QueueName, rabbitMQQueueNameArg, "", "RabbitMQ queue name")
	flag.StringVar(&clientID, "c", "", "clientID as UUID")
	flag.StringVar(&cfg.App.RequestType, "rt", "", fmt.Sprintf(
		"request type, possible values: %s, %s, %s, %s",
		addItem,
		removeItem,
		getItem,
		getAllItems,
	))

	flag.StringVar(&cfg.App.ItemValue, "v", "", "item value")

	flag.Parse()

	if cfg.RabbitMQ.Host == "" {
		err = fmt.Errorf(isEmpty, rabbitMQHostArg)
		return
	}

	if cfg.RabbitMQ.QueueName == "" {
		err = fmt.Errorf(isEmpty, rabbitMQQueueNameArg)
		return
	}

	cfg.App.ClientID, err = uuid.Parse(clientID)
	if err != nil {
		err = fmt.Errorf("parsing clientID failed: %w", err)
		return
	}

	if cfg.App.ItemValue == "" && requestType(cfg.App.RequestType) != getAllItems {
		err = fmt.Errorf(isEmpty, "item value")
		return
	}

	if err = requestType(cfg.App.RequestType).validate(); err != nil {
		err = fmt.Errorf("parsing request type failed: %w", err)
		return
	}

	return
}

func (rt requestType) validate() error {
	if _, ok := validRequestTypesMap[rt]; !ok {
		return errRequestTypeDoesntNotSupport
	}
	return nil
}
