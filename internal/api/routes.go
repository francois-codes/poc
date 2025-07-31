package api

import (
	"cognyx/psychic-robot/messaging"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler messaging.UserHandler) {
	userController := NewUserController(userHandler)
	
	api := app.Group("/api")
	users := api.Group("/users")
	
	// User management endpoints
	users.Post("/", userController.CreateUser)
	users.Get("/:id", userController.GetLatestUserVersion)
	users.Put("/:id", userController.UpdateUser)
	users.Get("/:id/versions", userController.GetUserVersions)
	users.Get("/:id/versions/:version", userController.GetSpecificUserVersion)
}