package app

import (
	"context"

	"cognyx/psychic-robot/messaging"
	"cognyx/psychic-robot/persistence/repository"
	"github.com/ThreeDotsLabs/watermill"
)

type UsersMessagingService struct {
	publisher   *messaging.NATSPublisher
	subscriber  *messaging.NATSSubscriber
	userHandler *messaging.UserHandlerImpl
	logger      watermill.LoggerAdapter
}

func NewUsersMessagingService(
	natsURL string,
	userRepo repository.UserRepository,
	versionRepo repository.VersionRepository,
	logger watermill.LoggerAdapter,
) (*UsersMessagingService, error) {
	
	publisher, err := messaging.NewNATSPublisher(natsURL, logger)
	if err != nil {
		return nil, err
	}

	subscriber, err := messaging.NewNATSSubscriber(natsURL, logger)
	if err != nil {
		return nil, err
	}

	userHandler := messaging.NewUserHandler(publisher, userRepo, versionRepo, logger)

	return &UsersMessagingService{
		publisher:   publisher,
		subscriber:  subscriber,
		userHandler: userHandler,
		logger:      logger,
	}, nil
}

func (ums *UsersMessagingService) Start(ctx context.Context) error {
	// Start user update subscription
	if err := ums.subscriber.SubscribeToUserUpdates(ctx, ums.userHandler); err != nil {
		return err
	}

	ums.logger.Info("Users messaging service started", nil)
	return nil
}

func (ums *UsersMessagingService) Stop() error {
	if err := ums.publisher.Close(); err != nil {
		ums.logger.Error("Failed to close publisher", err, nil)
	}
	
	ums.logger.Info("Users messaging service stopped", nil)
	return nil
}

func (ums *UsersMessagingService) GetUserHandler() messaging.UserHandler {
	return ums.userHandler
}