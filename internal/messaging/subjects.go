package messaging

const (
	// User management channels
	SubjectUsersUpdate   = "users.update"    // Frontend → Backend (user operations)
	SubjectUsersBroadcast = "users.broadcast" // Backend → Frontend (real-time updates)
	
	// Health and system events  
	SubjectHealth = "system.health"
	SubjectMetrics = "system.metrics"
)