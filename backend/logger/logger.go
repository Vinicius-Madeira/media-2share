package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		logger = zap.NewExample().Sugar()
	}
	return logger
}
