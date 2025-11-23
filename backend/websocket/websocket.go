package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"vinimad.com/media2share/logging"
)

// WebSocket connection configuration
const (
	writeWait      = 10 * time.Second    // Time allowed to write a message
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period
	maxMessageSize = 512 * 1024          // Maximum message size allowed
)

// Logger instance
var logger = logging.NewLogger("WebSocket")

// Websocket connection upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Clients Current available clients connections
var Clients []*Client

// HandleWebSocket handles the WebSocket connection
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Error upgrading connection", "error", err)
		return
	}

	client := NewClient(conn)
	logger.Info("New WebSocket connection established")
	logger.Debug("Client info", "client", client)

	Clients = append(Clients, client)

	// Start the read and write pumps
	go client.writePump()
	go client.readPump()
}
