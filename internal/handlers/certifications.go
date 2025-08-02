package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/database/repositories"
	"portfolio-backend/pkg/response"
)

type CertificationsHandler struct {
	certificationRepo repositories.CertificationRepository
}

func NewCertificationsHandler(certificationRepo repositories.CertificationRepository) *CertificationsHandler {
	return &CertificationsHandler{
		certificationRepo: certificationRepo,
	}
}

// GetCertifications handles GET /v1/certifications
func (h *CertificationsHandler) GetCertifications(c *gin.Context) {
	ctx := c.Request.Context()

	certifications, err := h.certificationRepo.GetAllCertifications(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get certifications")
		response.InternalServerError(c, err, "Failed to get certifications")
		return
	}

	response.Success(c, certifications)
}