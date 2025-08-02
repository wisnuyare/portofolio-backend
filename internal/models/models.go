package models

import (
	"time"
)

// Profile represents the user's profile information
type Profile struct {
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Title     string    `json:"title" db:"title" validate:"required,min=2,max=200"`
	Location  string    `json:"location" db:"location" validate:"required,min=2,max=100"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Phone     *string   `json:"phone,omitempty" db:"phone" validate:"omitempty,min=10,max=20"`
	LinkedIn  *string   `json:"linkedin,omitempty" db:"linkedin" validate:"omitempty,url"`
	Summary   string    `json:"summary" db:"summary" validate:"required,min=10,max=1000"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateProfileRequest represents the request payload for updating profile
type UpdateProfileRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=100"`
	Title    string  `json:"title" validate:"required,min=2,max=200"`
	Location string  `json:"location" validate:"required,min=2,max=100"`
	Email    string  `json:"email" validate:"required,email"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	LinkedIn *string `json:"linkedin,omitempty" validate:"omitempty,url"`
	Summary  string  `json:"summary" validate:"required,min=10,max=1000"`
}

// Experience represents work experience
type Experience struct {
	ID          int       `json:"id" db:"id"`
	Company     string    `json:"company" db:"company" validate:"required,min=2,max=100"`
	Position    string    `json:"position" db:"position" validate:"required,min=2,max=100"`
	StartDate   time.Time `json:"start_date" db:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	Description string    `json:"description" db:"description" validate:"required,min=10,max=2000"`
	Location    string    `json:"location" db:"location" validate:"required,min=2,max=100"`
	IsCurrent   bool      `json:"is_current" db:"is_current"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Skill represents a technical skill
type Skill struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Category    string    `json:"category" db:"category" validate:"required"`
	Level       string    `json:"level" db:"level" validate:"required,oneof=Beginner Intermediate Advanced Expert"`
	YearsOfExp  *int      `json:"years_of_experience,omitempty" db:"years_of_experience" validate:"omitempty,min=0,max=50"`
	Description *string   `json:"description,omitempty" db:"description" validate:"omitempty,max=500"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SkillCategory represents skill categories for grouping
type SkillCategory struct {
	Category string  `json:"category"`
	Skills   []Skill `json:"skills"`
}

// Education represents educational background
type Education struct {
	ID          int       `json:"id" db:"id"`
	Institution string    `json:"institution" db:"institution" validate:"required,min=2,max=200"`
	Degree      string    `json:"degree" db:"degree" validate:"required,min=2,max=100"`
	Field       string    `json:"field" db:"field" validate:"required,min=2,max=100"`
	StartDate   time.Time `json:"start_date" db:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	GPA         *float64  `json:"gpa,omitempty" db:"gpa" validate:"omitempty,min=0.0,max=4.0"`
	Description *string   `json:"description,omitempty" db:"description" validate:"omitempty,max=1000"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Certification represents professional certifications
type Certification struct {
	ID           int        `json:"id" db:"id"`
	Name         string     `json:"name" db:"name" validate:"required,min=2,max=200"`
	Issuer       string     `json:"issuer" db:"issuer" validate:"required,min=2,max=200"`
	IssueDate    time.Time  `json:"issue_date" db:"issue_date" validate:"required"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	CredentialID *string    `json:"credential_id,omitempty" db:"credential_id" validate:"omitempty,max=100"`
	URL          *string    `json:"url,omitempty" db:"url" validate:"omitempty,url"`
	Description  *string    `json:"description,omitempty" db:"description" validate:"omitempty,max=1000"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Message *string     `json:"message,omitempty"`
}

// APIError represents API error responses
type APIError struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Project represents a portfolio project
type Project struct {
	ID               int        `json:"id" db:"id"`
	Title            string     `json:"title" db:"title" validate:"required,min=2,max=200"`
	Description      string     `json:"description" db:"description" validate:"required,min=10,max=2000"`
	ShortDescription *string    `json:"short_description,omitempty" db:"short_description" validate:"omitempty,max=500"`
	Technologies     []string   `json:"technologies" db:"technologies" validate:"required,min=1"`
	GitHubURL        *string    `json:"github_url,omitempty" db:"github_url" validate:"omitempty,url"`
	LiveURL          *string    `json:"live_url,omitempty" db:"live_url" validate:"omitempty,url"`
	ImageURL         *string    `json:"image_url,omitempty" db:"image_url" validate:"omitempty,url"`
	StartDate        time.Time  `json:"start_date" db:"start_date" validate:"required"`
	EndDate          *time.Time `json:"end_date,omitempty" db:"end_date"`
	Status           string     `json:"status" db:"status" validate:"required,oneof=Planning 'In Progress' Completed 'On Hold' Cancelled"`
	Featured         bool       `json:"featured" db:"featured"`
	SortOrder        int        `json:"sort_order" db:"sort_order"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status     string            `json:"status"`
	Timestamp  time.Time         `json:"timestamp"`
	Version    string            `json:"version"`
	Components map[string]string `json:"components"`
}