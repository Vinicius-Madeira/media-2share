package main

import (
	"vinimad.com/media2share/logger"
)

func main() {
	var sugar = logger.GetLogger()

	sugar.Infow("Hello World", "Logger", "hey")
}
