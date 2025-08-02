package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/database/repositories"
	"portfolio-backend/internal/models"
)

type ProfileService interface {
	GetProfile(ctx context.Context) (*models.Profile, error)
	UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) (*models.Profile, error)
}

type profileService struct {
	profileRepo repositories.ProfileRepository
}

func NewProfileService(profileRepo repositories.ProfileRepository) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
	}
}

func (s *profileService) GetProfile(ctx context.Context) (*models.Profile, error) {
	log.Debug().Msg("Getting profile")

	profile, err := s.profileRepo.GetProfile(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get profile from repository")
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	log.Debug().
		Str("name", profile.Name).
		Str("title", profile.Title).
		Msg("Profile retrieved successfully")

	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) (*models.Profile, error) {
	log.Debug().
		Str("name", req.Name).
		Str("title", req.Title).
		Msg("Updating profile")

	// Business logic validation can be added here
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	profile, err := s.profileRepo.UpdateProfile(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update profile in repository")
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	log.Info().
		Str("name", profile.Name).
		Str("title", profile.Title).
		Msg("Profile updated successfully")

	return profile, nil
}