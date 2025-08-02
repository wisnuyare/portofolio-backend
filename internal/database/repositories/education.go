package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"portfolio-backend/internal/models"
)

type EducationRepository interface {
	GetAllEducation(ctx context.Context) ([]models.Education, error)
}

type MySQLEducationRepository struct {
	db *sql.DB
}

func NewEducationRepository(db *sql.DB) EducationRepository {
	return &MySQLEducationRepository{db: db}
}

func (r *MySQLEducationRepository) GetAllEducation(ctx context.Context) ([]models.Education, error) {
	query := `
		SELECT id, institution, degree, field, start_date, end_date, gpa, description, created_at, updated_at
		FROM education 
		ORDER BY start_date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query education: %w", err)
	}
	defer rows.Close()

	var educations []models.Education

	for rows.Next() {
		var edu models.Education
		var endDate sql.NullTime
		var gpa sql.NullFloat64
		var description sql.NullString

		err := rows.Scan(
			&edu.ID,
			&edu.Institution,
			&edu.Degree,
			&edu.Field,
			&edu.StartDate,
			&endDate,
			&gpa,
			&description,
			&edu.CreatedAt,
			&edu.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan education: %w", err)
		}

		// Handle nullable fields
		if endDate.Valid {
			edu.EndDate = &endDate.Time
		}
		if gpa.Valid {
			edu.GPA = &gpa.Float64
		}
		if description.Valid {
			edu.Description = &description.String
		}

		educations = append(educations, edu)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over education: %w", err)
	}

	return educations, nil
}