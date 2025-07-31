package websocket

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	natsgo "github.com/nats-io/nats.go"
)

type NATSBridge struct {
	natsConn *natsgo.Conn
}

type WebSocketMessage struct {
	Subject string          `json:"subject"`
	Data    json.RawMessage `json:"data"`
}

func NewNATSBridge(natsURL string) (*NATSBridge, error) {
	conn, err := natsgo.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &NATSBridge{
		natsConn: conn,
	}, nil
}

func (bridge *NATSBridge) HandleWebSocket(c *websocket.Conn) {
	defer c.Close()

	// Subscribe to users.broadcast to forward to WebSocket client
	sub, err := bridge.natsConn.Subscribe("users.broadcast", func(msg *natsgo.Msg) {
		// Forward NATS message to WebSocket client
		wsMsg := WebSocketMessage{
			Subject: msg.Subject,
			Data:    json.RawMessage(msg.Data),
		}
		
		if err := c.WriteJSON(wsMsg); err != nil {
			log.Printf("Error sending message to WebSocket client: %v", err)
		}
	})
	if err != nil {
		log.Printf("Error subscribing to NATS: %v", err)
		return
	}
	defer sub.Unsubscribe()

	// Handle incoming WebSocket messages
	for {
		var wsMsg WebSocketMessage
		if err := c.ReadJSON(&wsMsg); err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		// Publish to NATS
		if err := bridge.natsConn.Publish(wsMsg.Subject, wsMsg.Data); err != nil {
			log.Printf("Error publishing to NATS: %v", err)
		}
	}
}

func (bridge *NATSBridge) SetupRoutes(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(bridge.HandleWebSocket))
}