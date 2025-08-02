package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"portfolio-backend/internal/models"
)

type CertificationRepository interface {
	GetAllCertifications(ctx context.Context) ([]models.Certification, error)
}

type MySQLCertificationRepository struct {
	db *sql.DB
}

func NewCertificationRepository(db *sql.DB) CertificationRepository {
	return &MySQLCertificationRepository{db: db}
}

func (r *MySQLCertificationRepository) GetAllCertifications(ctx context.Context) ([]models.Certification, error) {
	query := `
		SELECT id, name, issuer, issue_date, expiry_date, credential_id, url, description, created_at, updated_at
		FROM certifications 
		ORDER BY issue_date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query certifications: %w", err)
	}
	defer rows.Close()

	var certifications []models.Certification

	for rows.Next() {
		var cert models.Certification
		var expiryDate sql.NullTime
		var credentialID, url, description sql.NullString

		err := rows.Scan(
			&cert.ID,
			&cert.Name,
			&cert.Issuer,
			&cert.IssueDate,
			&expiryDate,
			&credentialID,
			&url,
			&description,
			&cert.CreatedAt,
			&cert.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan certification: %w", err)
		}

		// Handle nullable fields
		if expiryDate.Valid {
			cert.ExpiryDate = &expiryDate.Time
		}
		if credentialID.Valid {
			cert.CredentialID = &credentialID.String
		}
		if url.Valid {
			cert.URL = &url.String
		}
		if description.Valid {
			cert.Description = &description.String
		}

		certifications = append(certifications, cert)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over certifications: %w", err)
	}

	return certifications, nil
}