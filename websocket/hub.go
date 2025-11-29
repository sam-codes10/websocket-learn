package websocket

import (
	"exp/models"
	"log"

	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients (raw JSON bytes).
	broadcast chan *models.WebSocketMsg

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan *models.WebSocketMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub loop. Call this in a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
			log.Printf("client registered: %v (total: %d)\n", client, len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send) // close send channel to stop writePump
				log.Printf("client unregistered: %v (total: %d)\n", client, len(h.clients))
			}
		case message := <-h.broadcast:
			// Broadcast to all clients. Use non-blocking sends to avoid slow clients blocking hub.
			// for _, client := range h.clients {
			// select {
			// case client.send <- message:
			// default:
			// 	// Client send buffer full — assume it's dead/unresponsive, remove it.
			// 	close(client.send)
			// 	delete(h.clients, client.ID)
			// 	log.Printf("dropped slow client: %v\n", client)
			// }
			//}
			receiverId := message.LocationMsg.ReceiverID
			receiverClient, ok := h.clients[receiverId]

			if !ok {
				logrus.Error("No reciever client found!")
			} else {
				receiverClient.send <- message
			}

		}
	}
}
