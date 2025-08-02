package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RequestLogger creates a Zerolog-based logging middleware
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Build path with query params
		if raw != "" {
			path = path + "?" + raw
		}

		// Choose log level based on status code
		var logEvent *zerolog.Event
		statusCode := c.Writer.Status()
		
		switch {
		case statusCode >= 500:
			logEvent = log.Error()
		case statusCode >= 400:
			logEvent = log.Warn()
		default:
			logEvent = log.Info()
		}

		logEvent.
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("client_ip", clientIP).
			Str("user_agent", userAgent).
			Int("body_size", c.Writer.Size()).
			Msg("HTTP Request")
	}
}

// CorrelationID middleware adds a correlation ID to each request
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = generateCorrelationID()
		}

		// Add correlation ID to context
		c.Set("correlation_id", correlationID)
		
		// Add correlation ID to response header
		c.Header("X-Correlation-ID", correlationID)

		// Update logger context
		logger := log.With().Str("correlation_id", correlationID).Logger()
		c.Set("logger", &logger)

		c.Next()
	}
}

// generateCorrelationID generates a simple correlation ID
func generateCorrelationID() string {
	// Simple implementation - in production, use UUID or similar
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().Nanosecond()%1000))
}