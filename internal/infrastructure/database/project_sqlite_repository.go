package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type ProjectSqliteRepository struct {
	db *sql.DB
}

func NewProjectSqliteRepository(db *sql.DB) models.ProjectRepository {
	return &ProjectSqliteRepository{db: db}
}

func (r *ProjectSqliteRepository) Create(project *models.Project) error {
	query := `
		INSERT INTO projects (id, slug, name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		project.ID.String(),
		project.Slug,
		project.Name,
		project.CreatedAt,
		project.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

func (r *ProjectSqliteRepository) GetBySlug(slug string) (*models.Project, error) {
	query := `
		SELECT id, slug, name, created_at, updated_at
		FROM projects
		WHERE slug = ?
	`

	var project models.Project
	var idStr string

	err := r.db.QueryRow(query, slug).Scan(
		&idStr,
		&project.Slug,
		&project.Name,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	project.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project ID: %w", err)
	}

	return &project, nil
}

func (r *ProjectSqliteRepository) ExistsBySlug(slug string) (bool, error) {
	query := `SELECT COUNT(*) FROM projects WHERE slug = ?`

	var count int
	err := r.db.QueryRow(query, slug).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check project existence: %w", err)
	}

	return count > 0, nil
}

func (r *ProjectSqliteRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	query := `SELECT id, slug, name, created_at, updated_at FROM projects WHERE id = ?`
	row := r.db.QueryRow(query, id.String())

	var project models.Project
	err := row.Scan(&project.ID, &project.Slug, &project.Name, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with ID '%s' not found", id.String())
		}
		return nil, fmt.Errorf("failed to get project by ID: %w", err)
	}

	return &project, nil
}
