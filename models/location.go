package models

// LocationMessage defines the structure of location data
// sent and received through WebSocket.
type LocationMessage struct {
	SenderID   string  `json:"senderId"`
	ReceiverID string  `json:"receiverId"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lng"`
	Timestamp  int64   `json:"ts"`
	// Event      string  `json:"event"` // e.g. "location_update"
}
