package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheConfig holds caching configuration
type CacheConfig struct {
	MaxAge     int  // Cache max age in seconds
	Public     bool // Whether cache is public (CDN cacheable)
	ETagEnable bool // Whether to generate ETags
}

// Cache returns a middleware that adds caching headers
func Cache(config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Set cache control headers
		cacheControl := "no-cache"
		if config.MaxAge > 0 {
			if config.Public {
				cacheControl = fmt.Sprintf("public, max-age=%d", config.MaxAge)
			} else {
				cacheControl = fmt.Sprintf("private, max-age=%d", config.MaxAge)
			}
		}
		c.Header("Cache-Control", cacheControl)

		// Add Last-Modified header (current time for dynamic content)
		c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		// Check if client sent If-None-Match header for ETag validation
		if config.ETagEnable {
			clientETag := c.GetHeader("If-None-Match")
			
			// Store original writer
			writer := &etagResponseWriter{
				ResponseWriter: c.Writer,
				clientETag:     clientETag,
				context:        c,
			}
			c.Writer = writer
		}

		c.Next()
	}
}

// etagResponseWriter wraps gin.ResponseWriter to generate ETags
type etagResponseWriter struct {
	gin.ResponseWriter
	clientETag string
	context    *gin.Context
	body       []byte
}

func (w *etagResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return len(data), nil
}

func (w *etagResponseWriter) WriteString(s string) (int, error) {
	w.body = append(w.body, []byte(s)...)
	return len(s), nil
}

func (w *etagResponseWriter) WriteHeader(statusCode int) {
	// Generate ETag from response body hash
	if len(w.body) > 0 && statusCode == http.StatusOK {
		etag := fmt.Sprintf(`"%x"`, md5.Sum(w.body))
		w.ResponseWriter.Header().Set("ETag", etag)
		
		// Check if client ETag matches
		if w.clientETag == etag {
			w.ResponseWriter.WriteHeader(http.StatusNotModified)
			return
		}
	}
	
	w.ResponseWriter.WriteHeader(statusCode)
	if len(w.body) > 0 {
		_, _ = w.ResponseWriter.Write(w.body)
	}
}

// DefaultCacheConfig returns a reasonable default caching configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		MaxAge:     300,   // 5 minutes
		Public:     true,  // Allow CDN caching
		ETagEnable: false, // Disable ETag for now
	}
}

// LongCacheConfig returns configuration for long-term caching (for static-like data)
func LongCacheConfig() CacheConfig {
	return CacheConfig{
		MaxAge:     3600,  // 1 hour
		Public:     true,  // Allow CDN caching
		ETagEnable: false, // Disable ETag for now
	}
}

// NoCacheConfig returns configuration that disables caching
func NoCacheConfig() CacheConfig {
	return CacheConfig{
		MaxAge:     0,     // No caching
		Public:     false, // Private
		ETagEnable: false, // No ETag
	}
}