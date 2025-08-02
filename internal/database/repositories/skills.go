package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"portfolio-backend/internal/models"
)

type SkillRepository interface {
	GetAllSkills(ctx context.Context) ([]models.Skill, error)
	GetSkillsByCategory(ctx context.Context) ([]models.SkillCategory, error)
}

type MySQLSkillRepository struct {
	db *sql.DB
}

func NewSkillRepository(db *sql.DB) SkillRepository {
	return &MySQLSkillRepository{db: db}
}

func (r *MySQLSkillRepository) GetAllSkills(ctx context.Context) ([]models.Skill, error) {
	query := `
		SELECT id, name, category, level, years_of_experience, description, created_at, updated_at
		FROM skills 
		ORDER BY category, name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query skills: %w", err)
	}
	defer rows.Close()

	var skills []models.Skill

	for rows.Next() {
		var skill models.Skill
		var yearsOfExp sql.NullInt32
		var description sql.NullString

		err := rows.Scan(
			&skill.ID,
			&skill.Name,
			&skill.Category,
			&skill.Level,
			&yearsOfExp,
			&description,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		// Handle nullable fields
		if yearsOfExp.Valid {
			years := int(yearsOfExp.Int32)
			skill.YearsOfExp = &years
		}
		if description.Valid {
			skill.Description = &description.String
		}

		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over skills: %w", err)
	}

	return skills, nil
}

func (r *MySQLSkillRepository) GetSkillsByCategory(ctx context.Context) ([]models.SkillCategory, error) {
	// First get all skills
	skills, err := r.GetAllSkills(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	// Group skills by category
	categoryMap := make(map[string][]models.Skill)
	for _, skill := range skills {
		categoryMap[skill.Category] = append(categoryMap[skill.Category], skill)
	}

	// Convert map to slice of SkillCategory
	var categories []models.SkillCategory
	for category, categorySkills := range categoryMap {
		categories = append(categories, models.SkillCategory{
			Category: category,
			Skills:   categorySkills,
		})
	}

	return categories, nil
}