package logger

import "go.uber.org/zap"

func GetLogger() *zap.SugaredLogger {
	logger := zap.NewExample().Sugar()

	return logger
}
