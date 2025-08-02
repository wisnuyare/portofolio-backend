package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"portfolio-backend/internal/models"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context) (*models.Profile, error)
	UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) (*models.Profile, error)
}

type MySQLProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) ProfileRepository {
	return &MySQLProfileRepository{db: db}
}

func (r *MySQLProfileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
	query := `
		SELECT name, title, location, email, phone, linkedin, summary, updated_at
		FROM profiles 
		LIMIT 1`

	var profile models.Profile
	var phone, linkedin sql.NullString

	err := r.db.QueryRowContext(ctx, query).Scan(
		&profile.Name,
		&profile.Title,
		&profile.Location,
		&profile.Email,
		&phone,
		&linkedin,
		&profile.Summary,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Handle nullable fields
	if phone.Valid {
		profile.Phone = &phone.String
	}
	if linkedin.Valid {
		profile.LinkedIn = &linkedin.String
	}

	return &profile, nil
}

func (r *MySQLProfileRepository) UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) (*models.Profile, error) {
	query := `
		UPDATE profiles 
		SET name = ?, title = ?, location = ?, email = ?, phone = ?, linkedin = ?, summary = ?, updated_at = NOW()
		WHERE id = 1`

	var phone, linkedin interface{}
	if req.Phone != nil {
		phone = *req.Phone
	}
	if req.LinkedIn != nil {
		linkedin = *req.LinkedIn
	}

	result, err := r.db.ExecContext(ctx, query,
		req.Name,
		req.Title,
		req.Location,
		req.Email,
		phone,
		linkedin,
		req.Summary,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("profile not found or no changes made")
	}

	// Return the updated profile
	return r.GetProfile(ctx)
}