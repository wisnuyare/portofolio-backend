package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/database/repositories"
	"portfolio-backend/internal/models"
)

type ExperienceService interface {
	GetAllExperiences(ctx context.Context) ([]models.Experience, error)
	GetExperienceByID(ctx context.Context, id int) (*models.Experience, error)
}

type experienceService struct {
	experienceRepo repositories.ExperienceRepository
}

func NewExperienceService(experienceRepo repositories.ExperienceRepository) ExperienceService {
	return &experienceService{
		experienceRepo: experienceRepo,
	}
}

func (s *experienceService) GetAllExperiences(ctx context.Context) ([]models.Experience, error) {
	log.Debug().Msg("Getting all experiences")

	experiences, err := s.experienceRepo.GetAllExperiences(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get experiences from repository")
		return nil, fmt.Errorf("failed to get experiences: %w", err)
	}

	log.Debug().
		Int("count", len(experiences)).
		Msg("Experiences retrieved successfully")

	return experiences, nil
}

func (s *experienceService) GetExperienceByID(ctx context.Context, id int) (*models.Experience, error) {
	log.Debug().
		Int("id", id).
		Msg("Getting experience by ID")

	if id <= 0 {
		return nil, fmt.Errorf("invalid experience ID: %d", id)
	}

	experience, err := s.experienceRepo.GetExperienceByID(ctx, id)
	if err != nil {
		log.Error().
			Err(err).
			Int("id", id).
			Msg("Failed to get experience from repository")
		return nil, fmt.Errorf("failed to get experience: %w", err)
	}

	log.Debug().
		Int("id", experience.ID).
		Str("company", experience.Company).
		Str("position", experience.Position).
		Msg("Experience retrieved successfully")

	return experience, nil
}