package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/models"
	"portfolio-backend/internal/services"
	"portfolio-backend/pkg/response"
	"portfolio-backend/pkg/validator"
)

type ProfileHandler struct {
	profileService services.ProfileService
}

func NewProfileHandler(profileService services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile handles GET /v1/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	profile, err := h.profileService.GetProfile(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get profile")
		response.InternalServerError(c, err, "Failed to get profile")
		return
	}

	response.Success(c, profile)
}

// UpdateProfile handles PUT /v1/profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		response.BadRequest(c, err, "Invalid request body")
		return
	}

	// Validate the request
	if validationErrors := validator.ValidateStruct(req); validationErrors != nil {
		response.ValidationError(c, validationErrors)
		return
	}

	profile, err := h.profileService.UpdateProfile(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update profile")
		response.InternalServerError(c, err, "Failed to update profile")
		return
	}

	response.Success(c, profile, "Profile updated successfully")
}