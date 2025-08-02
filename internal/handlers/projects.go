package handlers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/services"
	"portfolio-backend/pkg/response"
)

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetAllProjects handles GET /v1/projects
func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	ctx := c.Request.Context()

	// Check if featured filter is requested
	featuredParam := c.Query("featured")
	if featuredParam == "true" {
		h.GetFeaturedProjects(c)
		return
	}

	projects, err := h.projectService.GetAllProjects(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get projects")
		response.InternalServerError(c, err, "Failed to get projects")
		return
	}

	response.Success(c, projects)
}

// GetProjectByID handles GET /v1/projects/{id}
func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn().Str("id", idParam).Msg("Invalid project ID")
		response.BadRequest(c, err, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to get project")
		
		// Check if it's a not found error
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err, "Project not found")
			return
		}
		
		response.InternalServerError(c, err, "Failed to get project")
		return
	}

	response.Success(c, project)
}

// GetFeaturedProjects handles GET /v1/projects?featured=true
func (h *ProjectHandler) GetFeaturedProjects(c *gin.Context) {
	ctx := c.Request.Context()

	projects, err := h.projectService.GetFeaturedProjects(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get featured projects")
		response.InternalServerError(c, err, "Failed to get featured projects")
		return
	}

	response.Success(c, projects)
}