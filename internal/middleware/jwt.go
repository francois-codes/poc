package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
}

// JWTAuth creates a JWT authentication middleware for HTTP requests
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is required",
			})
		}

		// Verify the JWT token (currently returns true as requested)
		claims, err := verifyJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)

		return c.Next()
	}
}

// verifyJWT verifies the JWT token and returns claims
// For now, this always returns true as requested
func verifyJWT(token string) (*JWTClaims, error) {
	// TODO: Implement actual JWT verification when ready
	// For now, return dummy claims to satisfy the interface
	log.Println("🚀 VERIFY JWT")
	return &JWTClaims{
		UserID: "dummy-user-id",
		Email:  "dummy@example.com",
		Exp:    999999999999, // Far future expiry
	}, nil
}

// GetUserIDFromContext extracts user ID from Fiber context
func GetUserIDFromContext(c *fiber.Ctx) string {
	if userID, ok := c.Locals("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUserEmailFromContext extracts user email from Fiber context
func GetUserEmailFromContext(c *fiber.Ctx) string {
	if email, ok := c.Locals("user_email").(string); ok {
		return email
	}
	return ""
}

// WSJWTAuth authenticates WebSocket connections via query parameter or header
func WSJWTAuth(c *fiber.Ctx) error {
	var token string

	// Try to get token from query parameter first
	token = c.Query("token")

	// If not in query, try Authorization header
	if token == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is required for WebSocket connection",
		})
	}

	// Verify the JWT token
	claims, err := verifyJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Set user information in context for WebSocket connection
	c.Locals("user_id", claims.UserID)
	c.Locals("user_email", claims.Email)

	return c.Next()
}

// SocketIOJWTAuth authenticates Socket.IO connections
func SocketIOJWTAuth(handshakeAuth map[string]interface{}) (*JWTClaims, error) {
	// Try to get token from auth data
	var token string
	if tokenInterface, exists := handshakeAuth["token"]; exists {
		if tokenStr, ok := tokenInterface.(string); ok {
			token = tokenStr
		}
	}

	if token == "" {
		// Try to get from authorization field
		if authInterface, exists := handshakeAuth["authorization"]; exists {
			if authStr, ok := authInterface.(string); ok {
				if strings.HasPrefix(authStr, "Bearer ") {
					token = strings.TrimPrefix(authStr, "Bearer ")
				}
			}
		}
	}

	if token == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Token is required for Socket.IO connection")
	}

	// Verify the JWT token
	return verifyJWT(token)
}

// WithUserContext adds user information to a context
func WithUserContext(ctx context.Context, userID, email string) context.Context {
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "user_email", email)
	return ctx
}

// GetUserIDFromGoContext extracts user ID from Go context
func GetUserIDFromGoContext(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUserEmailFromGoContext extracts user email from Go context
func GetUserEmailFromGoContext(ctx context.Context) string {
	if email, ok := ctx.Value("user_email").(string); ok {
		return email
	}
	return ""
}
