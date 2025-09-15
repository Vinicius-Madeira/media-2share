package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
	"vinimad.com/media2share/logger"
)

// WebSocket connection configuration
const (
	writeWait      = 10 * time.Second    // Time allowed to write a message
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period
	maxMessageSize = 512 * 1024          // Maximum message size allowed
)

var sugar = logger.GetLogger()
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a single WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan []byte
	mu   sync.Mutex
}

// Message represents the structure of WebSocket messages
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// NewClient creates a new client instance
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256), // Buffer for outgoing messages
	}
}

// readPump handles incoming messages
func (c *Client) readPump() {
	defer func() {
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
				sugar.Errorw("Unexpected close error", "error", err)
			}
			break
		}

		// handles incoming message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			sugar.Errorw("Error unmarshaling message", "error", err)
			continue
		}

		// process message based on type
		c.handleMessage(msg)
	}
}

// writePump handles outgoing messages
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				sugar.Errorw("Error getting next writer", "error", err)
				return
			}

			if _, err := w.Write(message); err != nil {
				sugar.Errorw("Error writing message", "error", err)
				return
			}

			if err := w.Close(); err != nil {
				sugar.Errorw("Error closing writer", "error", err)
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

// handleMessage processes different types of messages
func (c *Client) handleMessage(msg Message) {
	sugar.Infow("Message received", "type", msg.Type, "payload", msg.Payload)

	// Handle different message types
	switch msg.Type {
	case "ping":
		c.SendMessage("pong", "pong")
	case "chat":
		c.SendMessage("message", msg.Payload)
	default:
		sugar.Warnw("Unknown message type", "type", msg.Type)
	}
}

// SendMessage sends a message to the client
func (c *Client) SendMessage(msgType string, payload interface{}) {
	msg := Message{
		Type:    msgType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		sugar.Errorw("Error marshaling message", "error", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.send <- data:
	default:
		sugar.Warnw("Client send buffer full")
	}
}

// HandleWebSocket handles the WebSocket connection
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sugar.Errorw("Error upgrading connection", "error", err)
		return
	}

	client := NewClient(conn)
	sugar.Infow("New WebSocket connection established")

	// Start the read and write pumps
	go client.writePump()
	go client.readPump()
}
