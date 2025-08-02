package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/services"
	"portfolio-backend/pkg/response"
)

type HealthHandler struct {
	healthService services.HealthService
}

func NewHealthHandler(healthService services.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// GetHealth handles GET /v1/health
func (h *HealthHandler) GetHealth(c *gin.Context) {
	ctx := c.Request.Context()

	health, err := h.healthService.CheckHealth(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Health check failed")
		response.InternalServerError(c, err, "Health check failed")
		return
	}

	// If health check indicates unhealthy status, return appropriate status code
	if health.Status != "healthy" {
		c.JSON(503, health) // Service Unavailable
		return
	}

	response.Success(c, health)
}