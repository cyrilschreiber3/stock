package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger returns a middleware that logs every HTTP request.
func GinLogger() gin.HandlerFunc {
	return GinLoggerWithConfig()
}

// GinLoggerWithConfig returns a middleware that logs HTTP requests.
// Routes listed in silentPaths are skipped on success (2xx/3xx) and logged
// at Debug level on failure, instead of the normal level.
func GinLoggerWithConfig(silentPaths ...string) gin.HandlerFunc {
	silent := make(map[string]struct{}, len(silentPaths))
	for _, p := range silentPaths {
		silent[p] = struct{}{}
	}

	return func(ctx *gin.Context) {
		startedAt := time.Now()
		rawPath := ctx.Request.URL.Path
		rawQuery := ctx.Request.URL.RawQuery

		ctx.Next()

		statusCode := ctx.Writer.Status()
		path := ctx.FullPath()
		if path == "" {
			path = rawPath
		}

		level := ginLogLevel(statusCode)

		// For silent routes, skip on success; downgrade to Debug on failure.
		if _, isSilent := silent[path]; isSilent {
			if statusCode < 400 {
				return
			}
			level = slog.LevelDebug
		}

		attrs := []slog.Attr{
			slog.Int("status", statusCode),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.String("client_ip", ctx.ClientIP()),
			slog.Duration("latency", time.Since(startedAt)),
			slog.Int("body_size", ctx.Writer.Size()),
			slog.String("user_agent", ctx.Request.UserAgent()),
		}

		if rawQuery != "" {
			attrs = append(attrs, slog.String("query", rawQuery))
		}

		if referrer := ctx.Request.Referer(); referrer != "" {
			attrs = append(attrs, slog.String("referer", referrer))
		}

		if errorMessage := ctx.Errors.ByType(gin.ErrorTypePrivate).String(); errorMessage != "" {
			attrs = append(attrs, slog.String("error", errorMessage))
		}

		slog.Default().LogAttrs(ctx.Request.Context(), level, "http_request", attrs...)
	}
}

func ginLogLevel(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
