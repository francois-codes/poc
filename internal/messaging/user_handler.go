package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cognyx/psychic-robot/persistence/db"
	"cognyx/psychic-robot/persistence/repository"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/google/uuid"
)

type UserHandler interface {
	HandleUserUpdate(ctx context.Context, event UserUpdateEvent) error
	CreateUser(ctx context.Context, email, status string, role *string, createdBy string) (*UserVersionDTO, error)
	UpdateUser(ctx context.Context, userID int64, email, status string, role *string, updatedBy string) (*UserVersionDTO, error)
	GetUserVersions(ctx context.Context, userID int64) ([]UserVersionDTO, error)
	GetLatestUserVersion(ctx context.Context, userID int64) (*UserVersionDTO, error)
}

type UserHandlerImpl struct {
	publisher   *NATSPublisher
	userRepo    repository.UserRepository
	versionRepo repository.VersionRepository
	logger      watermill.LoggerAdapter
}

func NewUserHandler(
	publisher *NATSPublisher, 
	userRepo repository.UserRepository,
	versionRepo repository.VersionRepository,
	logger watermill.LoggerAdapter,
) *UserHandlerImpl {
	return &UserHandlerImpl{
		publisher:   publisher,
		userRepo:    userRepo,
		versionRepo: versionRepo,
		logger:      logger,
	}
}

func (h *UserHandlerImpl) CreateUser(ctx context.Context, email, status string, role *string, createdBy string) (*UserVersionDTO, error) {
	// Create user using existing repository
	createParams := db.CreateUserParams{
		Email:  email,
		Status: status,
	}
	if role != nil {
		createParams.Role.String = *role
		createParams.Role.Valid = true
	}

	user, err := h.userRepo.Create(ctx, createParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Serialize user data to JSON for version storage
	userDataJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Create version using existing repository
	versionParams := db.CreateVersionParams{
		ObjectType: "user",
		ObjectID:   user.ID,
		Json:       userDataJSON,
		Version:    1,
		Action:     "create",
		Actor:      createdBy,
	}

	version, err := h.versionRepo.Create(ctx, versionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Convert to DTO
	versionDTO := ConvertVersionToDTO(version, user)

	// Publish NATS event
	event := UserUpdateEvent{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Operation: "create",
		Version:   version.Version,
		UserData:  user,
		Timestamp: time.Now(),
		CreatedBy: createdBy,
	}

	if err := h.publisher.PublishUserUpdate(ctx, event); err != nil {
		h.logger.Error("Failed to publish user create event", err, watermill.LogFields{
			"user_id": user.ID,
		})
		// Don't fail the operation if publishing fails
	}

	h.logger.Info("User created", watermill.LogFields{
		"user_id":    user.ID,
		"version":    version.Version,
		"created_by": createdBy,
	})

	return &versionDTO, nil
}

func (h *UserHandlerImpl) UpdateUser(ctx context.Context, userID int64, email, status string, role *string, updatedBy string) (*UserVersionDTO, error) {
	// Get existing user for previous data
	existingUser, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user: %w", err)
	}

	// Update user using existing repository
	updateParams := db.UpdateUserParams{
		ID:     userID,
		Email:  email,
		Status: status,
	}
	if role != nil {
		updateParams.Role.String = *role
		updateParams.Role.Valid = true
	}

	updatedUser, err := h.userRepo.Update(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Get the next version number
	existingVersions, err := h.versionRepo.ListByObject(ctx, "user", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing versions: %w", err)
	}

	nextVersion := int32(len(existingVersions) + 1)

	// Serialize updated user data to JSON
	userDataJSON, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Create new version
	versionParams := db.CreateVersionParams{
		ObjectType: "user",
		ObjectID:   userID,
		Json:       userDataJSON,
		Version:    nextVersion,
		Action:     "update",
		Actor:      updatedBy,
	}

	version, err := h.versionRepo.Create(ctx, versionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Convert to DTO
	versionDTO := ConvertVersionToDTO(version, updatedUser)

	// Publish NATS event
	event := UserUpdateEvent{
		ID:           uuid.New().String(),
		UserID:       userID,
		Operation:    "update",
		Version:      nextVersion,
		UserData:     updatedUser,
		PreviousData: &existingUser,
		Timestamp:    time.Now(),
		CreatedBy:    updatedBy,
	}

	if err := h.publisher.PublishUserUpdate(ctx, event); err != nil {
		h.logger.Error("Failed to publish user update event", err, watermill.LogFields{
			"user_id": userID,
		})
		// Don't fail the operation if publishing fails
	}

	h.logger.Info("User updated", watermill.LogFields{
		"user_id":    userID,
		"version":    nextVersion,
		"updated_by": updatedBy,
	})

	return &versionDTO, nil
}

func (h *UserHandlerImpl) GetUserVersions(ctx context.Context, userID int64) ([]UserVersionDTO, error) {
	// Get all versions for the user
	versions, err := h.versionRepo.ListByObject(ctx, "user", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	// Get current user data
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert versions to DTOs
	versionDTOs := make([]UserVersionDTO, len(versions))
	for i, version := range versions {
		// Unmarshal the JSON data to get the user data at that version
		var versionUserData db.User
		if err := json.Unmarshal(version.Json, &versionUserData); err != nil {
			// If unmarshal fails, use current user data
			versionUserData = user
		}
		versionDTOs[i] = ConvertVersionToDTO(version, versionUserData)
	}

	return versionDTOs, nil
}

func (h *UserHandlerImpl) GetLatestUserVersion(ctx context.Context, userID int64) (*UserVersionDTO, error) {
	// Get current user
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get latest version
	versions, err := h.versionRepo.ListByObject(ctx, "user", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found for user: %d", userID)
	}

	// Get the latest version (versions should be ordered by creation time)
	latestVersion := versions[len(versions)-1]
	versionDTO := ConvertVersionToDTO(latestVersion, user)

	return &versionDTO, nil
}

func (h *UserHandlerImpl) HandleUserUpdate(ctx context.Context, event UserUpdateEvent) error {
	// This handles incoming events from NATS users.update channel
	h.logger.Info("Received user update event from NATS", watermill.LogFields{
		"user_id":   event.UserID,
		"operation": event.Operation,
		"version":   event.Version,
	})

	var result *UserVersionDTO
	var err error

	switch event.Operation {
	case "create":
		// Extract user data from the event
		userData := event.UserData
		var role *string
		if userData.Role.Valid {
			role = &userData.Role.String
		}
		
		result, err = h.createUserInternal(ctx, userData.Email, userData.Status, role, event.CreatedBy)
		if err != nil {
			h.logger.Error("Failed to create user from NATS event", err, watermill.LogFields{
				"event_id": event.ID,
				"email":    userData.Email,
			})
			return err
		}

	case "update":
		if event.UserID == 0 {
			return fmt.Errorf("user_id is required for update operation")
		}
		
		// Extract user data from the event
		userData := event.UserData
		var role *string
		if userData.Role.Valid {
			role = &userData.Role.String
		}
		
		result, err = h.updateUserInternal(ctx, event.UserID, userData.Email, userData.Status, role, event.CreatedBy)
		if err != nil {
			h.logger.Error("Failed to update user from NATS event", err, watermill.LogFields{
				"event_id": event.ID,
				"user_id":  event.UserID,
				"email":    userData.Email,
			})
			return err
		}

	default:
		h.logger.Error("Unknown operation in user update event", nil, watermill.LogFields{
			"operation": event.Operation,
			"event_id":  event.ID,
		})
		return fmt.Errorf("unknown operation: %s", event.Operation)
	}

	h.logger.Info("Successfully processed user update event", watermill.LogFields{
		"event_id":  event.ID,
		"operation": event.Operation,
		"user_id":   result.UserID,
		"version":   result.Version,
	})

	return nil
}

// Internal method that creates user without publishing to NATS (to avoid loops)
func (h *UserHandlerImpl) createUserInternal(ctx context.Context, email, status string, role *string, createdBy string) (*UserVersionDTO, error) {
	// Create user using existing repository
	createParams := db.CreateUserParams{
		Email:  email,
		Status: status,
	}
	if role != nil {
		createParams.Role.String = *role
		createParams.Role.Valid = true
	}

	user, err := h.userRepo.Create(ctx, createParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Serialize user data to JSON for version storage
	userDataJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Create version using existing repository
	versionParams := db.CreateVersionParams{
		ObjectType: "user",
		ObjectID:   user.ID,
		Json:       userDataJSON,
		Version:    1,
		Action:     "create",
		Actor:      createdBy,
	}

	version, err := h.versionRepo.Create(ctx, versionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Convert to DTO
	versionDTO := ConvertVersionToDTO(version, user)

	// Publish NATS event to broadcast channel
	event := UserUpdateEvent{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Operation: "create",
		Version:   version.Version,
		UserData:  user,
		Timestamp: time.Now(),
		CreatedBy: createdBy,
	}

	if err := h.publisher.PublishUserUpdate(ctx, event); err != nil {
		h.logger.Error("Failed to publish user create broadcast", err, watermill.LogFields{
			"user_id": user.ID,
		})
		// Don't fail the operation if publishing fails
	}

	h.logger.Info("User created and broadcasted", watermill.LogFields{
		"user_id":    user.ID,
		"version":    version.Version,
		"created_by": createdBy,
	})

	return &versionDTO, nil
}

// Internal method that updates user without publishing to NATS (to avoid loops)
func (h *UserHandlerImpl) updateUserInternal(ctx context.Context, userID int64, email, status string, role *string, updatedBy string) (*UserVersionDTO, error) {
	// Get existing user for previous data
	existingUser, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user: %w", err)
	}

	// Update user using existing repository
	updateParams := db.UpdateUserParams{
		ID:     userID,
		Email:  email,
		Status: status,
	}
	if role != nil {
		updateParams.Role.String = *role
		updateParams.Role.Valid = true
	}

	updatedUser, err := h.userRepo.Update(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Get the next version number
	existingVersions, err := h.versionRepo.ListByObject(ctx, "user", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing versions: %w", err)
	}

	nextVersion := int32(len(existingVersions) + 1)

	// Serialize updated user data to JSON
	userDataJSON, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Create new version
	versionParams := db.CreateVersionParams{
		ObjectType: "user",
		ObjectID:   userID,
		Json:       userDataJSON,
		Version:    nextVersion,
		Action:     "update",
		Actor:      updatedBy,
	}

	version, err := h.versionRepo.Create(ctx, versionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Convert to DTO
	versionDTO := ConvertVersionToDTO(version, updatedUser)

	// Publish NATS event to broadcast channel
	event := UserUpdateEvent{
		ID:           uuid.New().String(),
		UserID:       userID,
		Operation:    "update",
		Version:      nextVersion,
		UserData:     updatedUser,
		PreviousData: &existingUser,
		Timestamp:    time.Now(),
		CreatedBy:    updatedBy,
	}

	if err := h.publisher.PublishUserUpdate(ctx, event); err != nil {
		h.logger.Error("Failed to publish user update broadcast", err, watermill.LogFields{
			"user_id": userID,
		})
		// Don't fail the operation if publishing fails
	}

	h.logger.Info("User updated and broadcasted", watermill.LogFields{
		"user_id":    userID,
		"version":    nextVersion,
		"updated_by": updatedBy,
	})

	return &versionDTO, nil
}