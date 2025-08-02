package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"

	"portfolio-backend/pkg/response"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int           // Number of requests per second allowed
	BurstSize         int           // Maximum burst size
	CleanupInterval   time.Duration // How often to clean up expired limiters
}

// ClientLimiter holds rate limiter for a specific client
type ClientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter manages rate limiting for multiple clients
type RateLimiter struct {
	clients map[string]*ClientLimiter
	mu      sync.RWMutex
	config  RateLimitConfig
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		config:  config,
	}

	// Start cleanup routine
	go rl.cleanupExpiredClients()

	return rl
}

// RateLimit returns a Gin middleware for rate limiting
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := rl.getClientID(c)
		
		if !rl.allow(clientID) {
			log.Warn().
				Str("client_id", clientID).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("Rate limit exceeded")

			response.TooManyRequests(c, "Too many requests, please try again later")
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientID extracts client identifier (IP address for now)
func (rl *RateLimiter) getClientID(c *gin.Context) string {
	// Try to get real IP from headers (in case of proxy/load balancer)
	if xRealIP := c.GetHeader("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}
	if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
		return xForwardedFor
	}
	return c.ClientIP()
}

// allow checks if the client is allowed to make a request
func (rl *RateLimiter) allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Get or create limiter for this client
	limiter, exists := rl.clients[clientID]
	if !exists {
		limiter = &ClientLimiter{
			limiter:  rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize),
			lastSeen: time.Now(),
		}
		rl.clients[clientID] = limiter
	} else {
		limiter.lastSeen = time.Now()
	}

	return limiter.limiter.Allow()
}

// cleanupExpiredClients removes limiters that haven't been used recently
func (rl *RateLimiter) cleanupExpiredClients() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		
		cutoff := time.Now().Add(-rl.config.CleanupInterval * 2)
		for clientID, limiter := range rl.clients {
			if limiter.lastSeen.Before(cutoff) {
				delete(rl.clients, clientID)
			}
		}
		
		rl.mu.Unlock()
	}
}

// DefaultRateLimitConfig returns a reasonable default configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond: 10,              // 10 requests per second
		BurstSize:         20,              // Allow bursts up to 20 requests
		CleanupInterval:   5 * time.Minute, // Clean up every 5 minutes
	}
}

// StrictRateLimitConfig returns a more restrictive configuration
func StrictRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond: 5,               // 5 requests per second
		BurstSize:         10,              // Allow bursts up to 10 requests
		CleanupInterval:   5 * time.Minute, // Clean up every 5 minutes
	}
}