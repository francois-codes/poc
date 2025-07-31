package api

import (
	"context"
	"strconv"

	"cognyx/psychic-robot/messaging"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userHandler messaging.UserHandler
}

func NewUserController(userHandler messaging.UserHandler) *UserController {
	return &UserController{
		userHandler: userHandler,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email  string  `json:"email" validate:"required,email"`
	Status string  `json:"status" validate:"required"`
	Role   *string `json:"role,omitempty"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email  string  `json:"email" validate:"required,email"`
	Status string  `json:"status" validate:"required"`
	Role   *string `json:"role,omitempty"`
}

// CreateUser creates a new user with versioning
// POST /api/users
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from JWT or session
	createdBy := ctx.Get("X-User-ID", "system")
	
	version, err := c.userHandler.CreateUser(context.Background(), req.Email, req.Status, req.Role, createdBy)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return ctx.Status(201).JSON(version)
}

// UpdateUser updates an existing user and creates a new version
// PUT /api/users/:id
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("id")
	if userIDStr == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from JWT or session
	updatedBy := ctx.Get("X-User-ID", "system")
	
	version, err := c.userHandler.UpdateUser(context.Background(), userID, req.Email, req.Status, req.Role, updatedBy)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return ctx.JSON(version)
}

// GetUserVersions returns all versions of a user
// GET /api/users/:id/versions
func (c *UserController) GetUserVersions(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("id")
	if userIDStr == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	versions, err := c.userHandler.GetUserVersions(context.Background(), userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"user_id": userID,
		"versions": versions,
		"total": len(versions),
	})
}

// GetLatestUserVersion returns the latest version of a user
// GET /api/users/:id
func (c *UserController) GetLatestUserVersion(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("id")
	if userIDStr == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	version, err := c.userHandler.GetLatestUserVersion(context.Background(), userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return ctx.JSON(version)
}

// GetSpecificUserVersion returns a specific version of a user
// GET /api/users/:id/versions/:version
func (c *UserController) GetSpecificUserVersion(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("id")
	versionStr := ctx.Params("version")
	
	if userIDStr == "" || versionStr == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "User ID and version number are required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	versionNum, err := strconv.ParseInt(versionStr, 10, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid version number",
		})
	}

	versions, err := c.userHandler.GetUserVersions(context.Background(), userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Find the specific version
	for _, v := range versions {
		if v.Version == int32(versionNum) {
			return ctx.JSON(v)
		}
	}

	return ctx.Status(404).JSON(fiber.Map{
		"error": "Version not found",
	})
}