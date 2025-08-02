package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"portfolio-backend/internal/models"
)

type ProjectRepository interface {
	GetAllProjects(ctx context.Context) ([]models.Project, error)
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	GetFeaturedProjects(ctx context.Context) ([]models.Project, error)
}

type MySQLProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &MySQLProjectRepository{db: db}
}

func (r *MySQLProjectRepository) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	query := `
		SELECT id, title, description, short_description, technologies, github_url, live_url, image_url, 
		       start_date, end_date, status, featured, sort_order, created_at, updated_at
		FROM projects 
		ORDER BY sort_order ASC, start_date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		project, err := r.scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over projects: %w", err)
	}

	return projects, nil
}

func (r *MySQLProjectRepository) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	query := `
		SELECT id, title, description, short_description, technologies, github_url, live_url, image_url, 
		       start_date, end_date, status, featured, sort_order, created_at, updated_at
		FROM projects 
		WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	project, err := r.scanProjectRow(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (r *MySQLProjectRepository) GetFeaturedProjects(ctx context.Context) ([]models.Project, error) {
	query := `
		SELECT id, title, description, short_description, technologies, github_url, live_url, image_url, 
		       start_date, end_date, status, featured, sort_order, created_at, updated_at
		FROM projects 
		WHERE featured = true
		ORDER BY sort_order ASC, start_date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query featured projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		project, err := r.scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over featured projects: %w", err)
	}

	return projects, nil
}

// scanProject scans a project from sql.Rows
func (r *MySQLProjectRepository) scanProject(rows *sql.Rows) (*models.Project, error) {
	var project models.Project
	var endDate sql.NullTime
	var shortDescription, githubURL, liveURL, imageURL sql.NullString
	var technologiesJSON string

	err := rows.Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		&shortDescription,
		&technologiesJSON,
		&githubURL,
		&liveURL,
		&imageURL,
		&project.StartDate,
		&endDate,
		&project.Status,
		&project.Featured,
		&project.SortOrder,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return r.populateProject(&project, endDate, shortDescription, githubURL, liveURL, imageURL, technologiesJSON)
}

// scanProjectRow scans a project from sql.Row
func (r *MySQLProjectRepository) scanProjectRow(row *sql.Row) (*models.Project, error) {
	var project models.Project
	var endDate sql.NullTime
	var shortDescription, githubURL, liveURL, imageURL sql.NullString
	var technologiesJSON string

	err := row.Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		&shortDescription,
		&technologiesJSON,
		&githubURL,
		&liveURL,
		&imageURL,
		&project.StartDate,
		&endDate,
		&project.Status,
		&project.Featured,
		&project.SortOrder,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return r.populateProject(&project, endDate, shortDescription, githubURL, liveURL, imageURL, technologiesJSON)
}

// populateProject populates nullable fields and parses JSON
func (r *MySQLProjectRepository) populateProject(project *models.Project, endDate sql.NullTime, shortDescription, githubURL, liveURL, imageURL sql.NullString, technologiesJSON string) (*models.Project, error) {
	// Handle nullable fields
	if endDate.Valid {
		project.EndDate = &endDate.Time
	}
	if shortDescription.Valid {
		project.ShortDescription = &shortDescription.String
	}
	if githubURL.Valid {
		project.GitHubURL = &githubURL.String
	}
	if liveURL.Valid {
		project.LiveURL = &liveURL.String
	}
	if imageURL.Valid {
		project.ImageURL = &imageURL.String
	}

	// Parse technologies JSON
	if err := json.Unmarshal([]byte(technologiesJSON), &project.Technologies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal technologies JSON: %w", err)
	}

	return project, nil
}