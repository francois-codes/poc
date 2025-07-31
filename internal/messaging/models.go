package messaging

import (
	"time"

	"cognyx/psychic-robot/persistence/db"
)

// UserUpdateEvent represents a user change event for NATS
type UserUpdateEvent struct {
	ID           string    `json:"id"`
	UserID       int64     `json:"user_id"`     // Using existing int64 ID
	Operation    string    `json:"operation"`   // "create", "update"
	Version      int32     `json:"version"`     // Version number from versions table
	UserData     db.User   `json:"user_data"`   // Complete user data snapshot
	PreviousData *db.User  `json:"previous_data,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	CreatedBy    string    `json:"created_by"`
}

// UserVersionDTO represents a user version for API responses
type UserVersionDTO struct {
	ID        int64     `json:"id"`         // Version ID
	UserID    int64     `json:"user_id"`    // Reference to user
	Version   int32     `json:"version"`    // Version number
	UserData  db.User   `json:"user_data"`  // Complete user data snapshot
	Action    string    `json:"action"`     // create, update
	CreatedAt time.Time `json:"created_at"` // When this version was created
	CreatedBy string    `json:"created_by"` // Who created this version
}

// ConvertVersionToDTO converts db.Version to UserVersionDTO
func ConvertVersionToDTO(version db.Version, userData db.User) UserVersionDTO {
	return UserVersionDTO{
		ID:        version.ID,
		UserID:    version.ObjectID,
		Version:   version.Version,
		UserData:  userData,
		Action:    version.Action,
		CreatedAt: version.CreatedAt,
		CreatedBy: version.Actor,
	}
}