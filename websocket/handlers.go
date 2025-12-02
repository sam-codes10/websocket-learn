package websocket

import (
	"exp/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader config - tune ReadBufferSize / WriteBufferSize as needed.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// IMPORTANT: set origin policy appropriately for your app.
	CheckOrigin: func(r *http.Request) bool {
		// TODO: tighten this in production (e.g., allow only your domain)
		return true
	},
}

// RegisterClient upgrades the HTTP request to a WebSocket connection and registers the client with the hub.
// Example route: router.GET("/location/ws/register", func(c *gin.Context) { websocket.RegisterClient(hub, c) })
func RegisterClient(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v\n", err)
		return
	}

	// Optionally, extract identifying info (user id) from query or headers:
	userID := c.Query("user_id")

	if userID == "" {
		log.Println("user_id query parameter is required")
		conn.Close()
		c.JSON(400, gin.H{"message": "user_id query parameter is required"})
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan *models.WebSocketMsg, 256), // buffered channel
		ID:   userID,
	}

	// Register client with hub
	// client.hub.register <- client
	hub.register <- client

	// Start pumps
	go client.WritePump()
	go client.ReadPump()
}

func UnregisterClient(hub *Hub, c *gin.Context) {
	userId := c.Query("user_id")
	client := hub.clients
	if client[userId] == nil {
		log.Println("No client found with given user id")
		return
	}
	hub.unregister <- client[userId]
}
