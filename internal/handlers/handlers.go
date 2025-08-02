package handlers

import (
	"database/sql"

	"portfolio-backend/internal/database"
	"portfolio-backend/internal/database/repositories"
	"portfolio-backend/internal/services"
)

// Handlers holds all the HTTP handlers and their dependencies
type Handlers struct {
	Profile       *ProfileHandler
	Experience    *ExperienceHandler
	Skills        *SkillsHandler
	Education     *EducationHandler
	Certifications *CertificationsHandler
	Projects      *ProjectHandler
	Health        *HealthHandler
}

// NewHandlers creates and initializes all handlers
func NewHandlers(db *sql.DB) *Handlers {
	// Initialize repositories
	profileRepo := repositories.NewProfileRepository(db)
	experienceRepo := repositories.NewExperienceRepository(db)
	skillRepo := repositories.NewSkillRepository(db)
	educationRepo := repositories.NewEducationRepository(db)
	certificationRepo := repositories.NewCertificationRepository(db)
	projectRepo := repositories.NewProjectRepository(db)

	// Initialize services
	profileService := services.NewProfileService(profileRepo)
	experienceService := services.NewExperienceService(experienceRepo)
	projectService := services.NewProjectService(projectRepo)
	
	// Create a DB wrapper for health service
	dbWrapper := &database.DB{DB: db}
	healthService := services.NewHealthService(dbWrapper)

	return &Handlers{
		Profile:       NewProfileHandler(profileService),
		Experience:    NewExperienceHandler(experienceService),
		Skills:        NewSkillsHandler(skillRepo),
		Education:     NewEducationHandler(educationRepo),
		Certifications: NewCertificationsHandler(certificationRepo),
		Projects:      NewProjectHandler(projectService),
		Health:        NewHealthHandler(healthService),
	}
}