package websocket

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Messenger interface {
	SendMessage(msgType string, payload interface{})
}

// Client represents a single WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan []byte
	mu   sync.Mutex
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
				logger.Error("Unexpected close error", "error", err)
			}
			break
		}

		// handles incoming message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.Error("Error unmarshaling message", "error", err)
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
				logger.Error("Error getting next writer", "error", err)
				return
			}

			if _, err := w.Write(message); err != nil {
				logger.Error("Error writing message", "error", err)
				return
			}

			if err := w.Close(); err != nil {
				logger.Error("Error closing writer", "error", err)
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
	logger.Debug("Message received", "type", msg.Type, "payload", msg.Payload)

	// Handle different message types
	switch msg.Type {
	case "ping":
		c.SendMessage("pong", "pong")
	case "chat":
		c.SendMessage("message", msg.Payload)
	default:
		logger.Warn("Unknown message type", "type", msg.Type)
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
		logger.Error("Error marshaling message", "error", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.send <- data:
	default:
		logger.Warn("Client send buffer full")
	}
}

func (c *Client) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("RemoteAddr", c.conn.RemoteAddr().String()),
	)
}
