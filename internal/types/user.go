package types

import "time"

// User represents a user entity
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Role      *string   `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   bool      `json:"_deleted"`
}

// GetUsersParams represents parameters for getting users
type GetUsersParams struct {
	MinUpdatedAt *string `json:"minUpdatedAt,omitempty"`
	Limit        *int    `json:"limit,omitempty"`
}

// GetUsersResponse represents the response for getting users
type GetUsersResponse struct {
	Documents []User `json:"documents"`
}

// PostUsersBody represents the request body for creating/updating users
// This corresponds to RxReplicationWriteToMasterRow<User>[] in TypeScript
type PostUsersBody struct {
	Documents []RxReplicationWriteToMasterRow `json:"documents"`
}

// PostUsersResponse represents the response for creating/updating users
// This corresponds to ReplicationPushHandlerResult<User> in TypeScript
type PostUsersResponse struct {
	Documents ReplicationPushHandlerResult `json:"documents"`
}

// RxReplicationWriteToMasterRow represents a write operation to master
type RxReplicationWriteToMasterRow struct {
	NewDocumentState User   `json:"newDocumentState"`
	PreviousRevision string `json:"previousRevision"`
	AssumedMasterState *User `json:"assumedMasterState,omitempty"`
}

// ReplicationPushHandlerResult represents the result of a push operation
type ReplicationPushHandlerResult struct {
	// Array of successfully processed documents
	Documents []User `json:"documents,omitempty"`
	// Array of conflicts that occurred during push
	Conflicts []ReplicationConflict `json:"conflicts,omitempty"`
	// Errors that occurred during processing
	Errors []ReplicationError `json:"errors,omitempty"`
}

// ReplicationConflict represents a conflict during replication
type ReplicationConflict struct {
	DocumentID       string `json:"documentId"`
	NewDocumentState User   `json:"newDocumentState"`
	RealMasterState  User   `json:"realMasterState"`
}

// ReplicationError represents an error during replication
type ReplicationError struct {
	DocumentID string `json:"documentId"`
	Error      string `json:"error"`
	Status     int    `json:"status,omitempty"`
}
