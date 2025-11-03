package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

// Info logs informational messages.
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

// Warn logs warning messages.
func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

// Error logs error messages with error context.
func Error(msg string, err error, args ...any) {
	Logger.Error(msg, append(args, "error", err)...)
}

// HTTPLogger returns a gin middleware that logs incoming HTTP requests.
func HTTPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		Logger.Info("http_request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
