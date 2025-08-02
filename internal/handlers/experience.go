package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/services"
	"portfolio-backend/pkg/response"
)

type ExperienceHandler struct {
	experienceService services.ExperienceService
}

func NewExperienceHandler(experienceService services.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{
		experienceService: experienceService,
	}
}

// GetAllExperiences handles GET /v1/experience
func (h *ExperienceHandler) GetAllExperiences(c *gin.Context) {
	ctx := c.Request.Context()

	experiences, err := h.experienceService.GetAllExperiences(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get experiences")
		response.InternalServerError(c, err, "Failed to get experiences")
		return
	}

	response.Success(c, experiences)
}

// GetExperienceByID handles GET /v1/experience/{id}
func (h *ExperienceHandler) GetExperienceByID(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn().Str("id", idParam).Msg("Invalid experience ID")
		response.BadRequest(c, err, "Invalid experience ID")
		return
	}

	experience, err := h.experienceService.GetExperienceByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to get experience")
		
		// Check if it's a not found error
		if err.Error() == "experience with id "+idParam+" not found" {
			response.NotFound(c, err, "Experience not found")
			return
		}
		
		response.InternalServerError(c, err, "Failed to get experience")
		return
	}

	response.Success(c, experience)
}