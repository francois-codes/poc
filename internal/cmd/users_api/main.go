package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cognyx/psychic-robot/api"
	"cognyx/psychic-robot/app"
	"cognyx/psychic-robot/persistence/db"
	"cognyx/psychic-robot/persistence/repository"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Setup logger
	watermillLogger := watermill.NewStdLogger(false, false)
	
	// Initialize database connection
	dbURL := "postgres://cognyx:cognyx@localhost:5432/cognyx?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(context.Background())
	
	// Initialize repositories
	queries := db.New(conn)
	userRepo := repository.NewUserRepository(queries)
	versionRepo := repository.NewVersionRepository(queries)
	
	// Initialize messaging service
	messagingService, err := app.NewUsersMessagingService(
		"nats://localhost:4222",
		userRepo,
		versionRepo,
		watermillLogger,
	)
	if err != nil {
		log.Fatal("Failed to create messaging service:", err)
	}
	
	// Start messaging service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	if err := messagingService.Start(ctx); err != nil {
		log.Fatal("Failed to start messaging service:", err)
	}
	defer messagingService.Stop()
	
	// Setup Fiber app
	fiberApp := fiber.New(fiber.Config{
		AppName: "Users Versioning API",
	})
	
	// Middleware
	fiberApp.Use(cors.New())
	fiberApp.Use(logger.New())
	
	// Setup API routes (optional, for debugging)
	api.SetupRoutes(fiberApp, messagingService.GetUserHandler())
	
	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		
		log.Println("Received interrupt signal, shutting down...")
		cancel()
		fiberApp.Shutdown()
	}()
	
	// Start server
	log.Println("Server starting on :8080")
	if err := fiberApp.Listen(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}