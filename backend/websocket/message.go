package websocket

// Message represents the structure of WebSocket messages
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
