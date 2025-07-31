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

type NATSPublisher struct {
	publisher message.Publisher
	logger    watermill.LoggerAdapter
}

func NewNATSPublisher(natsURL string, logger watermill.LoggerAdapter) (*NATSPublisher, error) {
	natsConfig := nats.PublisherConfig{
		URL: natsURL,
		NatsOptions: []natsgo.Option{
			natsgo.MaxReconnects(10),
			natsgo.ReconnectWait(2 * time.Second),
		},
		Marshaler: nats.JSONMarshaler{}, // Use JSON for better compatibility
	}

	publisher, err := nats.NewPublisher(natsConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS publisher: %w", err)
	}

	return &NATSPublisher{
		publisher: publisher,
		logger:    logger,
	}, nil
}

func (p *NATSPublisher) PublishUserUpdate(ctx context.Context, event UserUpdateEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal user update event: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.Metadata.Set("timestamp", time.Now().Format(time.RFC3339))
	msg.Metadata.Set("user_id", fmt.Sprintf("%d", event.UserID))
	msg.Metadata.Set("version", fmt.Sprintf("%d", event.Version))

	return p.publisher.Publish(SubjectUsersBroadcast, msg)
}

func (p *NATSPublisher) Close() error {
	return p.publisher.Close()
}