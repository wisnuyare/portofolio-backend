package services

import (
	"context"
	"time"

	"portfolio-backend/internal/database"
	"portfolio-backend/internal/models"
)

type HealthService interface {
	CheckHealth(ctx context.Context) (*models.HealthResponse, error)
}

type healthService struct {
	db *database.DB
}

func NewHealthService(db *database.DB) HealthService {
	return &healthService{
		db: db,
	}
}

func (s *healthService) CheckHealth(ctx context.Context) (*models.HealthResponse, error) {
	components := make(map[string]string)

	// Check database health
	if err := s.db.Health(ctx); err != nil {
		components["database"] = "unhealthy"
		return &models.HealthResponse{
			Status:     "unhealthy",
			Timestamp:  time.Now(),
			Version:    "1.0.0", // This could be injected from build flags
			Components: components,
		}, err
	}

	components["database"] = "healthy"

	return &models.HealthResponse{
		Status:     "healthy",
		Timestamp:  time.Now(),
		Version:    "1.0.0",
		Components: components,
	}, nil
}