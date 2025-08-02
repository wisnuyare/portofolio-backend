package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/database/repositories"
	"portfolio-backend/pkg/response"
)

type EducationHandler struct {
	educationRepo repositories.EducationRepository
}

func NewEducationHandler(educationRepo repositories.EducationRepository) *EducationHandler {
	return &EducationHandler{
		educationRepo: educationRepo,
	}
}

// GetEducation handles GET /v1/education
func (h *EducationHandler) GetEducation(c *gin.Context) {
	ctx := c.Request.Context()

	education, err := h.educationRepo.GetAllEducation(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get education")
		response.InternalServerError(c, err, "Failed to get education")
		return
	}

	response.Success(c, education)
}