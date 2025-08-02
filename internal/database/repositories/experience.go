package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"portfolio-backend/internal/models"
)

type ExperienceRepository interface {
	GetAllExperiences(ctx context.Context) ([]models.Experience, error)
	GetExperienceByID(ctx context.Context, id int) (*models.Experience, error)
}

type MySQLExperienceRepository struct {
	db *sql.DB
}

func NewExperienceRepository(db *sql.DB) ExperienceRepository {
	return &MySQLExperienceRepository{db: db}
}

func (r *MySQLExperienceRepository) GetAllExperiences(ctx context.Context) ([]models.Experience, error) {
	query := `
		SELECT id, company, position, start_date, end_date, description, location, is_current, created_at, updated_at
		FROM experiences 
		ORDER BY start_date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query experiences: %w", err)
	}
	defer rows.Close()

	var experiences []models.Experience

	for rows.Next() {
		var exp models.Experience
		var endDate sql.NullTime

		err := rows.Scan(
			&exp.ID,
			&exp.Company,
			&exp.Position,
			&exp.StartDate,
			&endDate,
			&exp.Description,
			&exp.Location,
			&exp.IsCurrent,
			&exp.CreatedAt,
			&exp.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan experience: %w", err)
		}

		// Handle nullable end_date
		if endDate.Valid {
			exp.EndDate = &endDate.Time
		}

		experiences = append(experiences, exp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over experiences: %w", err)
	}

	return experiences, nil
}

func (r *MySQLExperienceRepository) GetExperienceByID(ctx context.Context, id int) (*models.Experience, error) {
	query := `
		SELECT id, company, position, start_date, end_date, description, location, is_current, created_at, updated_at
		FROM experiences 
		WHERE id = ?`

	var exp models.Experience
	var endDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exp.ID,
		&exp.Company,
		&exp.Position,
		&exp.StartDate,
		&endDate,
		&exp.Description,
		&exp.Location,
		&exp.IsCurrent,
		&exp.CreatedAt,
		&exp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("experience with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get experience: %w", err)
	}

	// Handle nullable end_date
	if endDate.Valid {
		exp.EndDate = &endDate.Time
	}

	return &exp, nil
}