package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"vinimad.com/media2share/logger"
)

var sugar = logger.GetLogger()
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sugar.Errorw("Error upgrading websocket connection", "error", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			sugar.Errorw("Error closing websocket connection", "error", err)
		}
		sugar.Infow("Websocket connection closed!")
	}()

	sugar.Infow("Websocket connection established!")

	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			sugar.Errorw("Error reading message", "error", err)
		}
		sugar.Infow("Message received", "type", messageType, "message", string(message))

		time.Sleep(time.Second * 3)

		if err := conn.WriteMessage(messageType, []byte("Hello world from server")); err != nil {
			sugar.Errorw("Error writing message", "error", err)
			break
		}
	}
}
