package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
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

// Logger instance
var sugar = logger.GetLogger()

// Websocket connection upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
