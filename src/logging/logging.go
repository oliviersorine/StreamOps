package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/oliviersorine/StreamOps/src/config"
)

func New(cfg config.LoggingConfig) *slog.Logger {
	level := parseLevel(cfg.Level)

	options := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, options)
	default:
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
