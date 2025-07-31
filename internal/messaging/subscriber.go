package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	natsgo "github.com/nats-io/nats.go"
)

type NATSSubscriber struct {
	subscriber message.Subscriber
	logger     watermill.LoggerAdapter
}

func NewNATSSubscriber(natsURL string, logger watermill.LoggerAdapter) (*NATSSubscriber, error) {
	natsConfig := nats.SubscriberConfig{
		URL: natsURL,
		NatsOptions: []natsgo.Option{
			natsgo.MaxReconnects(10),
			natsgo.ReconnectWait(2 * time.Second),
		},
		Unmarshaler:    nats.JSONMarshaler{},
		SubscribersCount: 1,
		CloseTimeout:   30 * time.Second,
		// JetStream configuration
		QueueGroupPrefix: "users-api",
	}

	subscriber, err := nats.NewSubscriber(natsConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS subscriber: %w", err)
	}

	return &NATSSubscriber{
		subscriber: subscriber,
		logger:     logger,
	}, nil
}

func (s *NATSSubscriber) SubscribeToUserUpdates(ctx context.Context, handler UserHandler) error {
	messages, err := s.subscriber.Subscribe(ctx, SubjectUsersUpdate)
	if err != nil {
		return fmt.Errorf("failed to subscribe to users.update: %w", err)
	}

	go func() {
		for msg := range messages {
			s.handleUserUpdate(msg, handler)
		}
	}()

	return nil
}

func (s *NATSSubscriber) handleUserUpdate(msg *message.Message, handler UserHandler) {
	var event UserUpdateEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		s.logger.Error("Failed to unmarshal user update event", err, nil)
		msg.Nack()
		return
	}

	if err := handler.HandleUserUpdate(context.Background(), event); err != nil {
		s.logger.Error("Failed to handle user update", err, watermill.LogFields{
			"user_id": event.UserID,
			"version": event.Version,
		})
		msg.Nack()
		return
	}

	msg.Ack()
}