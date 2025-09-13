package main

import (
	"net/http"
	"vinimad.com/media2share/logger"
	"vinimad.com/media2share/websocket"
)

var sugar = logger.GetLogger()

func main() {
	sugar.Infow("Starting websocket server")

	http.HandleFunc("/ws", websocket.HandleWebSocket)

	sugar.Infow("Websocket server started")

	if err := http.ListenAndServe(":9090", nil); err != nil {
		sugar.Fatalf("Error starting server: %s", err)
	}
}
