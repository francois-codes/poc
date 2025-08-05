package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestJWTAuth_MissingHeader(t *testing.T) {
	app := fiber.New()
	app.Get("/test", JWTAuth(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestJWTAuth_InvalidFormat(t *testing.T) {
	app := fiber.New()
	app.Get("/test", JWTAuth(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	app := fiber.New()
	app.Get("/test", JWTAuth(), func(c *fiber.Ctx) error {
		userID := GetUserIDFromContext(c)
		userEmail := GetUserEmailFromContext(c)
		return c.JSON(fiber.Map{
			"message":    "success",
			"user_id":    userID,
			"user_email": userEmail,
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-jwt-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// Since verifyJWT always returns true, this should succeed
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestVerifyJWT_AlwaysTrue(t *testing.T) {
	// Test that verifyJWT currently always returns true as requested
	claims, err := verifyJWT("any-token")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if claims.UserID != "dummy-user-id" {
		t.Errorf("Expected dummy-user-id, got %s", claims.UserID)
	}

	if claims.Email != "dummy@example.com" {
		t.Errorf("Expected dummy@example.com, got %s", claims.Email)
	}
}

func TestSocketIOJWTAuth_MissingToken(t *testing.T) {
	authData := map[string]interface{}{}
	
	_, err := SocketIOJWTAuth(authData)
	if err == nil {
		t.Error("Expected error for missing token, got nil")
	}
}

func TestSocketIOJWTAuth_WithToken(t *testing.T) {
	authData := map[string]interface{}{
		"token": "valid-jwt-token",
	}
	
	claims, err := SocketIOJWTAuth(authData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if claims.UserID != "dummy-user-id" {
		t.Errorf("Expected dummy-user-id, got %s", claims.UserID)
	}
}

func TestSocketIOJWTAuth_WithAuthorization(t *testing.T) {
	authData := map[string]interface{}{
		"authorization": "Bearer valid-jwt-token",
	}
	
	claims, err := SocketIOJWTAuth(authData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if claims.UserID != "dummy-user-id" {
		t.Errorf("Expected dummy-user-id, got %s", claims.UserID)
	}
}