package main

import (
	"cognyx/psychic-robot/middleware"
	"cognyx/psychic-robot/persistence/db"
	"cognyx/psychic-robot/persistence/repository"
	"cognyx/psychic-robot/types"
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zishang520/socket.io/v2/socket"
	"log"
	"strings"
	"time"
)

const GENERATE = false

func main() {

	dbconn := InitDB()
	defer dbconn.Close()

	// Repository
	queries := db.New(dbconn) // 👈 conversion pool → Queries
	userRepo := repository.NewUserRepository(queries)

	// SOCKET.IO
	c := socket.DefaultServerOptions()
	c.SetServeClient(true)
	c.SetPingInterval(300 * time.Millisecond)
	c.SetPingTimeout(200 * time.Millisecond)
	c.SetMaxHttpBufferSize(1000000)
	c.SetConnectTimeout(1000 * time.Millisecond)

	socketio := socket.NewServer(nil, nil)
	socketio.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)

		// Authenticate the Socket.IO connection
		authData, ok := client.Handshake().Auth.(map[string]interface{})
		if !ok {
			log.Printf("Socket.IO authentication failed: invalid auth data format")
			client.Disconnect(true)
			return
		}
		
		claims, err := middleware.SocketIOJWTAuth(authData)
		if err != nil {
			log.Printf("Socket.IO authentication failed: %v", err)
			client.Disconnect(true)
			return
		}

		log.Printf("Socket.IO connection authenticated for user: %s (%s)", claims.UserID, claims.Email)

		client.On("message", func(args ...interface{}) {
			log.Printf("Message from user %s: %v", claims.UserID, args)
			client.Emit("message-back", args...)
		})
		
		// Emit successful authentication
		client.Emit("auth", map[string]interface{}{
			"authenticated": true,
			"user_id":       claims.UserID,
			"email":         claims.Email,
		})

		client.On("message-with-ack", func(args ...interface{}) {
			log.Printf("Message with ACK from user %s: %v", claims.UserID, args)
			ack := args[len(args)-1].(socket.Ack)
			ack(args[:len(args)-1], nil)
		})
	})

	socketio.Of("/custom", nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		
		// Authenticate the Socket.IO connection for custom namespace
		authData, ok := client.Handshake().Auth.(map[string]interface{})
		if !ok {
			log.Printf("Socket.IO /custom authentication failed: invalid auth data format")
			client.Disconnect(true)
			return
		}
		
		claims, err := middleware.SocketIOJWTAuth(authData)
		if err != nil {
			log.Printf("Socket.IO /custom authentication failed: %v", err)
			client.Disconnect(true)
			return
		}

		log.Printf("Socket.IO /custom connection authenticated for user: %s (%s)", claims.UserID, claims.Email)
		
		client.Emit("auth", map[string]interface{}{
			"authenticated": true,
			"user_id":       claims.UserID,
			"email":         claims.Email,
		})
	})

	///////////

	app := fiber.New()

	// Active CORS avec les options par défaut (autorise tout)
	app.Use(cors.New())

	// SocketIO endpoints
	app.Get("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(c)))
	app.Post("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(c)))

	// Example User endpoints with JWT authentication
	app.Get("/api/users", middleware.JWTAuth(), func(c *fiber.Ctx) error {
		userID := middleware.GetUserIDFromContext(c)
		userEmail := middleware.GetUserEmailFromContext(c)
		
		log.Printf("🚀 GET REQUEST ON /api/users from user: %s (%s)", userID, userEmail)
		
		users, err := userRepo.List(c.Context(), 25, 0)
		resp := types.GetUsersResponse{}
		resp.Documents = mapUsersToUsers(users)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}
		log.Println("🚀 GET REQUEST ON http://localhost:4000/api/users ---> SUCCESS")
		return c.JSON(resp)
	})

	// Example User endpoints with JWT authentication
	app.Post("/api/users", middleware.JWTAuth(), func(c *fiber.Ctx) error {
		var input types.PostUsersBody
		userID := middleware.GetUserIDFromContext(c)
		userEmail := middleware.GetUserEmailFromContext(c)
		
		log.Printf("🚀 POST REQUEST ON /api/users from user: %s (%s)", userID, userEmail)

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid JSON body",
			})
		}

		resp := types.ReplicationPushHandlerResult{}
		resp.Documents = make([]types.User, 0)
		resp.Errors = make([]types.ReplicationError, 0)

		for _, rxReplicationWriteToMasterRow := range input.Documents {
			user := db.User{
				ID:        rxReplicationWriteToMasterRow.NewDocumentState.ID,
				Name:      rxReplicationWriteToMasterRow.NewDocumentState.Status,
				Email:     rxReplicationWriteToMasterRow.NewDocumentState.Email,
				Roles:     []string{},
				CreatedAt: rxReplicationWriteToMasterRow.NewDocumentState.CreatedAt,
				UpdatedAt: rxReplicationWriteToMasterRow.NewDocumentState.UpdatedAt,
			}
			user, err := userRepo.Create(c.Context(), user)
			if err != nil {
				resp.Errors = append(resp.Errors, types.ReplicationError{
					DocumentID: rxReplicationWriteToMasterRow.NewDocumentState.ID,
					Error:      err.Error(),
					Status:     0,
				})
				continue
			}
			resp.Documents = append(resp.Documents, mapUserToUser(user))
		}
		log.Println("🚀 POST REQUEST ON http://localhost:4000/api/users ---> SUCCESS")

		// 🔥 Emit update to all connected clients
		rxDocumentData := mapDocumentsToRxDocumentData(resp.Documents)
		toStream := types.UsersStreamEvent{Data: types.RxReplicationPullStreamItem{
			Documents: rxDocumentData,
			Checkpoint: types.CheckpointType{
				UpdatedAt: time.Now().String(),
				ID:        "titi",
			},
		}}
		socketio.Sockets().Emit("sync", toStream)
		return c.Status(fiber.StatusCreated).JSON(resp)
	})

	// WebSocket endpoint with JWT authentication
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// Apply JWT authentication for WebSocket upgrade
			if err := middleware.WSJWTAuth(c); err != nil {
				return err
			}
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Get user information from Fiber locals (set during upgrade)
		userID, _ := c.Locals("user_id").(string)
		userEmail, _ := c.Locals("user_email").(string)
		
		if userID == "" {
			userID = "unknown"
		}
		if userEmail == "" {
			userEmail = "unknown"
		}
		
		log.Printf("WebSocket connection established for user: %s (%s)", userID, userEmail)
		
		// Send welcome message
		c.WriteJSON(map[string]interface{}{
			"type": "welcome",
			"message": "WebSocket connection authenticated",
			"user_id": userID,
			"email": userEmail,
		})

		// Handle incoming messages
		for {
			var msg map[string]interface{}
			if err := c.ReadJSON(&msg); err != nil {
				log.Printf("Error reading WebSocket message from user %s: %v", userID, err)
				break
			}

			log.Printf("Received WebSocket message from user %s: %v", userID, msg)

			// Echo message back with user info
			response := map[string]interface{}{
				"type": "echo",
				"original_message": msg,
				"from_user": userID,
				"timestamp": time.Now().Format(time.RFC3339),
			}
			
			if err := c.WriteJSON(response); err != nil {
				log.Printf("Error sending WebSocket message to user %s: %v", userID, err)
				break
			}
		}
	}, websocket.Config{
		Filter: func(c *fiber.Ctx) bool {
			return c.Locals("allowed") == true
		},
	}))

	log.Println("🚀 Server started on http://localhost:4000")
	log.Fatal(app.Listen(":4000"))
}

func mapDocumentsToRxDocumentData(users []types.User) []types.RxDocumentData {
	result := make([]types.RxDocumentData, len(users))

	for i, user := range users {
		result[i] = types.RxDocumentData{
			User: types.User{
				ID:        user.ID,
				Email:     user.Email,
				Status:    user.Status,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Deleted:   user.Deleted,
			},
			Rev:     "1-" + user.ID, // dummy rev, replace with real revision logic if needed
			Deleted: user.Deleted,
			Meta:    nil, // optional: set if you have metadata
		}
	}

	return result
}

func mapUserToUser(user db.User) types.User {

	tmp := strings.Join(user.Roles, ",")
	result := types.User{
		ID:        user.ID,
		Status:    user.Name,
		Email:     user.Email,
		Role:      &tmp,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Deleted:   false,
	}
	return result
}

func mapUsersToUsers(users []db.User) []types.User {
	result := make([]types.User, len(users))
	for i, user := range users {
		tmp := strings.Join(user.Roles, ",")
		result[i] = types.User{
			ID:        user.ID,
			Status:    user.Name,
			Email:     user.Email,
			Role:      &tmp,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Deleted:   false,
		}
	}
	return result
}

func InitDB() *pgxpool.Pool {
	dsn := "postgres://cognyx:cognyx@localhost:5432/cognyx"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Erreur de création du pool PostgreSQL : %v", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("PostgreSQL inaccessible : %v", err)
	}

	return pool
}
