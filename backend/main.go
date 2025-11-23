package main

import (
	"context"
	"log/slog"
	"net/http"

	"vinimad.com/media2share/logging"
	"vinimad.com/media2share/websocket"
)

var (
	logger = logging.NewLogger("App")
	PORT   = ":9090"
)

func main() {
	http.HandleFunc("/ws", websocket.HandleWebSocket)

	logger.Info("starting server", "port", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		logger.Log(context.Background(), logging.LevelFatal, "error starting server", slog.Any("error:", err))
	}
}
