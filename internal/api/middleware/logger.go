package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger creates a structured logging middleware
func Logger() gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(startTime)

		// Get status
		statusCode := c.Writer.Status()

		// Log attributes
		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.RequestURI),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
		}

		// Log based on status code
		if statusCode >= 500 {
			logger.LogAttrs(c.Request.Context(), slog.LevelError, "Server Error", attrs...)
		} else if statusCode >= 400 {
			logger.LogAttrs(c.Request.Context(), slog.LevelWarn, "Client Error", attrs...)
		} else {
			logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "Request Processed", attrs...)
		}
	}
}
