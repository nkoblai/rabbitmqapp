package service

import (
	"context"
	"fmt"
	"io"
	"sync"

	orderedmap "github.com/elliotchance/orderedmap/v2"
	"github.com/google/uuid"
	"github.com/nkoblai/rabbitmqapp/server/internal/queue"
	"go.uber.org/zap"
)

const (
	addItem     string = "AddItem"
	removeItem  string = "RemoveItem"
	getItem     string = "GetItem"
	getAllItems string = "GetAllItems"

	logWithValueMsg = "clientID: '%s', request type: '%s', item value: '%s', current state: '%v'\n"
	logMsg          = "clientID: '%s', request type: '%s', current state: '%v'\n"
)

// Queue is queue interface.
type Queue interface {
	Requests(ctx context.Context) (<-chan queue.Item, error)
}

type clientsItems struct {
	m    *sync.RWMutex
	dict map[uuid.UUID]*orderedmap.OrderedMap[string, struct{}]
}

// Service is service struct.
type Service struct {
	logger       *zap.Logger
	writer       io.Writer
	queue        Queue
	clientsItems clientsItems
	handlerMap   map[string]func(item queue.Item) error
}

// New returns new service.
func New(l *zap.Logger, writer io.Writer, q Queue) *Service {
	s := &Service{
		logger: l,
		queue:  q,
		clientsItems: clientsItems{
			m:    &sync.RWMutex{},
			dict: make(map[uuid.UUID]*orderedmap.OrderedMap[string, struct{}]),
		},
		writer: writer,
	}
	s.handlerMap = map[string]func(item queue.Item) error{
		addItem:     s.addItem,
		removeItem:  s.removeItem,
		getItem:     s.getItem,
		getAllItems: s.getAllItems,
	}
	return s
}

// Process processes incoming requests.
func (s *Service) Process(ctx context.Context) error {
	requestsCh, err := s.queue.Requests(ctx)
	if err != nil {
		return fmt.Errorf("can't get requests channel: %w", err)
	}

	wg := &sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil

		case req := <-requestsCh:
			go func() {
				handle, ok := s.handlerMap[req.RequestType]
				if !ok {
					s.logger.Warn(fmt.Sprintf("recieved not supported request type: %s", req.RequestType))
					return
				}

				wg.Add(1)
				defer wg.Done()

				if err := handle(req); err != nil {
					s.logger.Warn("writing results of processing request failed", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Service) addItem(item queue.Item) error {
	s.clientsItems.m.Lock()
	defer s.clientsItems.m.Unlock()
	c, ok := s.clientsItems.dict[item.ClientID]
	if !ok {
		s.clientsItems.dict[item.ClientID] = orderedmap.NewOrderedMap[string, struct{}]()
		c = s.clientsItems.dict[item.ClientID]
	}
	c.Set(item.Value, struct{}{})

	return s.log(item)
}

func (s *Service) getItem(item queue.Item) error {
	s.clientsItems.m.RLock()
	defer s.clientsItems.m.RUnlock()
	s.clientsItems.dict[item.ClientID].Get(item.Value)

	return s.log(item)
}

func (s *Service) removeItem(item queue.Item) error {
	s.clientsItems.m.Lock()
	defer s.clientsItems.m.Unlock()
	s.clientsItems.dict[item.ClientID].Delete(item.Value)

	return s.log(item)
}

func (s *Service) getAllItems(item queue.Item) error {
	s.clientsItems.m.RLock()
	defer s.clientsItems.m.RUnlock()
	s.clientsItems.dict[item.ClientID].Keys()

	return s.log(item)
}

func (s *Service) log(item queue.Item) error {
	if item.Value == "" {
		if _, err := fmt.Fprintf(
			s.writer,
			logMsg,
			item.ClientID,
			item.RequestType,
			s.clientsItems.dict[item.ClientID].Keys(),
		); err != nil {
			return err
		}
		return nil
	}

	if _, err := fmt.Fprintf(
		s.writer,
		logWithValueMsg,
		item.ClientID,
		item.RequestType,
		item.Value,
		s.clientsItems.dict[item.ClientID].Keys(),
	); err != nil {
		return err
	}

	return nil
}
