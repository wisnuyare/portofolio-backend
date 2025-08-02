package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/database/repositories"
	"portfolio-backend/internal/models"
)

type ProjectService interface {
	GetAllProjects(ctx context.Context) ([]models.Project, error)
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	GetFeaturedProjects(ctx context.Context) ([]models.Project, error)
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
	}
}

func (s *projectService) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	log.Debug().Msg("Getting all projects")

	projects, err := s.projectRepo.GetAllProjects(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get projects from repository")
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	log.Debug().
		Int("count", len(projects)).
		Msg("Projects retrieved successfully")

	return projects, nil
}

func (s *projectService) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	log.Debug().
		Int("id", id).
		Msg("Getting project by ID")

	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", id)
	}

	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		log.Error().
			Err(err).
			Int("id", id).
			Msg("Failed to get project from repository")
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	log.Debug().
		Int("id", project.ID).
		Str("title", project.Title).
		Str("status", project.Status).
		Bool("featured", project.Featured).
		Msg("Project retrieved successfully")

	return project, nil
}

func (s *projectService) GetFeaturedProjects(ctx context.Context) ([]models.Project, error) {
	log.Debug().Msg("Getting featured projects")

	projects, err := s.projectRepo.GetFeaturedProjects(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get featured projects from repository")
		return nil, fmt.Errorf("failed to get featured projects: %w", err)
	}

	log.Debug().
		Int("count", len(projects)).
		Msg("Featured projects retrieved successfully")

	return projects, nil
}