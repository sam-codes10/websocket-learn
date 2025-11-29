package websocket

import (
	"encoding/json"
	"exp/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer (adjust as needed).
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *models.WebSocketMsg
	// optional: an identifier for the client (like userID or socket ID)
	ID string
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		// On exit, unregister and close connection
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v\n", err)
			}
			break
		}

		var  webSocketMsg models.WebSocketMsg
		if err := json.Unmarshal(message, &webSocketMsg); err != nil {
			log.Println("invalid json: ", err)
			logrus.Error("invalid json: ", err)
			continue
		}

		locationMsg := webSocketMsg.LocationMsg

		locationMsg.Timestamp = time.Now().Unix()
		//msgBytes, _ := json.Marshal(locationMsg)
		c.hub.broadcast <- &webSocketMsg

		// Optionally: you can validate/transform the message here (e.g., ensure it's valid JSON, attach client ID, etc.)
		// For simplicity: forward raw message to hub for broadcast.
		//c.hub.broadcast <- message
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			msgBytes, err  := json.Marshal(message)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Use NextWriter to send message as a single frame
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(msgBytes); err != nil {
				w.Close()
				return
			}

			// w.Write(message)

			// If there are queued messages, send them in the same WebSocket message to reduce overhead.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write([]byte{'\n'})
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
