package logging

import (
	"log/slog"
	"os"
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var levelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func NewLogger(ctxName string, defaultAttr ...slog.Attr) *slog.Logger {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: LevelTrace,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.LevelKey {
				level := attr.Value.Any().(slog.Level)
				levelLabel, exists := levelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				attr.Value = slog.StringValue(levelLabel)
			}

			return attr
		},
	}))

	// prepend context name
	defaultAttr = append([]slog.Attr{slog.String("name", ctxName)}, defaultAttr...)

	return logger.With(
		slog.GroupAttrs("context", defaultAttr...))
}
