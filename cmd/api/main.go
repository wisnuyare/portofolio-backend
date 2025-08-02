package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/database"
	"portfolio-backend/internal/handlers"
	"portfolio-backend/internal/middleware"
)

const version = "1.0.0"

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logger
	setupLogger(cfg.Logging)

	log.Info().
		Str("version", version).
		Str("host", cfg.Server.Host).
		Int("port", cfg.Server.Port).
		Msg("Starting Portfolio Backend API")

	// Connect to database
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize handlers
	h := handlers.NewHandlers(db.DB)

	// Setup Gin router
	router := setupRouter(cfg, h)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info().
			Str("address", server.Addr).
			Msg("HTTP server starting")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

func setupLogger(loggingConfig config.LoggingConfig) {
	// Set log level
	switch loggingConfig.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log format
	if loggingConfig.Format != "json" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Add caller information in debug mode
	if loggingConfig.Level == "debug" {
		log.Logger = log.With().Caller().Logger()
	}
}

func setupRouter(cfg *config.Config, h *handlers.Handlers) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Logging.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.CORS(&cfg.CORS))
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogger())
	
	// Add rate limiting
	rateLimiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
		RequestsPerSecond: cfg.RateLimit.RequestsPerSecond,
		BurstSize:         cfg.RateLimit.BurstSize,
		CleanupInterval:   cfg.RateLimit.CleanupInterval,
	})
	router.Use(rateLimiter.RateLimit())

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Health check (no caching)
		v1.GET("/health", middleware.Cache(middleware.NoCacheConfig()), h.Health.GetHealth)

		// Profile routes (short cache)
		v1.GET("/profile", middleware.Cache(middleware.DefaultCacheConfig()), h.Profile.GetProfile)
		v1.PUT("/profile", h.Profile.UpdateProfile)

		// Experience routes (long cache - relatively static)
		v1.GET("/experience", middleware.Cache(middleware.LongCacheConfig()), h.Experience.GetAllExperiences)
		v1.GET("/experience/:id", middleware.Cache(middleware.LongCacheConfig()), h.Experience.GetExperienceByID)

		// Skills routes (long cache - relatively static)
		v1.GET("/skills", middleware.Cache(middleware.LongCacheConfig()), h.Skills.GetSkills)

		// Education routes (long cache - very static)
		v1.GET("/education", middleware.Cache(middleware.LongCacheConfig()), h.Education.GetEducation)

		// Certifications routes (default cache)
		v1.GET("/certifications", middleware.Cache(middleware.DefaultCacheConfig()), h.Certifications.GetCertifications)

		// Projects routes (default cache - may be updated occasionally)
		v1.GET("/projects", middleware.Cache(middleware.DefaultCacheConfig()), h.Projects.GetAllProjects)
		v1.GET("/projects/:id", middleware.Cache(middleware.DefaultCacheConfig()), h.Projects.GetProjectByID)
	}

	return router
}