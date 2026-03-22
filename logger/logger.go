package logger

import (
	"log/slog"
	"os"
	"strings"
)

func InitLogger() *slog.Logger {
	logPath := os.Getenv("LOGFILE")
	logFormat := strings.ToLower(strings.TrimSpace(os.Getenv("LOGFORMAT")))
	logFile := os.Stdout
	if logPath != "" {
		var err error
		logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			logFile = os.Stdout
		}
	}

	var handler slog.Handler
	if logFormat == "text" || logFormat == "plain" {
		handler = slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slogLevelFromEnv()})
	} else {
		handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slogLevelFromEnv()})
	}

	return slog.New(handler)
}

func slogLevelFromEnv() slog.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOGLEVEL"))) {
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
