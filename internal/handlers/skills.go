package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/database/repositories"
	"portfolio-backend/pkg/response"
)

type SkillsHandler struct {
	skillRepo repositories.SkillRepository
}

func NewSkillsHandler(skillRepo repositories.SkillRepository) *SkillsHandler {
	return &SkillsHandler{
		skillRepo: skillRepo,
	}
}

// GetSkills handles GET /v1/skills
func (h *SkillsHandler) GetSkills(c *gin.Context) {
	ctx := c.Request.Context()

	// Check if client wants skills grouped by category
	groupBy := c.Query("group_by")
	
	if groupBy == "category" {
		categories, err := h.skillRepo.GetSkillsByCategory(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get skills by category")
			response.InternalServerError(c, err, "Failed to get skills")
			return
		}
		response.Success(c, categories)
		return
	}

	// Default: return all skills as a flat list
	skills, err := h.skillRepo.GetAllSkills(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get skills")
		response.InternalServerError(c, err, "Failed to get skills")
		return
	}

	response.Success(c, skills)
}