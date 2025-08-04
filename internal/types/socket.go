package types

// CheckpointType represents a checkpoint for replication
type CheckpointType struct {
	UpdatedAt string `json:"updated_at"`
	ID        string `json:"id"`
}

// RxReplicationPullStreamItem represents an item in a pull stream from RxDB
type RxReplicationPullStreamItem struct {
	Documents  []RxDocumentData `json:"documents"`
	Checkpoint CheckpointType   `json:"checkpoint"`
}

// RxDocumentData represents a document with RxDB metadata
type RxDocumentData struct {
	// The actual document data
	User
	// RxDB metadata
	Rev      string `json:"_rev"`
	Deleted  bool   `json:"_deleted"`
	Meta     *RxDocumentMeta `json:"_meta,omitempty"`
}

// RxDocumentMeta represents RxDB document metadata
type RxDocumentMeta struct {
	Lwt int64 `json:"lwt"` // Last write time
}

// CollectionStreamEvent represents a stream event for a collection
type CollectionStreamEvent struct {
	Data RxReplicationPullStreamItem `json:"data"`
}

// UsersStreamEvent represents a stream event specifically for users
type UsersStreamEvent struct {
	Data RxReplicationPullStreamItem `json:"data"`
}

// SocketServerEvents represents the events that can be emitted by the socket server
// In Go, this would typically be implemented as interfaces or function types
type SocketServerEvents interface {
	// Sync event handler
	OnSync(event CollectionStreamEvent)
	// Error event handler  
	OnError(error error)
	// Connection open handler
	OnOpen()
	// Connection close handler
	OnClose()
}

// SocketSyncEvent represents a sync event for the socket
type SocketSyncEvent struct {
	Event CollectionStreamEvent `json:"event"`
}

// SocketErrorEvent represents an error event for the socket
type SocketErrorEvent struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Stack   string `json:"stack,omitempty"`
}
