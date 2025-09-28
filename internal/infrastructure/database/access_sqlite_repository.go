package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type AccessSqliteRepository struct {
	db *sql.DB
}

func NewAccessSqliteRepository(db *sql.DB) *AccessSqliteRepository {
	return &AccessSqliteRepository{db: db}
}

func (r *AccessSqliteRepository) Create(access *models.Access) error {
	query := `
		INSERT INTO access (id, project_id, uid, pin, name, readonly, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		access.ID.String(),
		access.ProjectID.String(),
		access.UID,
		access.PIN,
		access.Name,
		access.ReadOnly,
		access.CreatedAt,
		access.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create access: %w", err)
	}

	return nil
}

func (r *AccessSqliteRepository) GetByProjectID(projectID uuid.UUID) ([]*models.Access, error) {
	query := `
		SELECT id, project_id, uid, pin, name, readonly, created_at, updated_at
		FROM access
		WHERE project_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, projectID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query access by project_id: %w", err)
	}
	defer rows.Close()

	var accesses []*models.Access
	for rows.Next() {
		access, err := r.scanAccess(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan access: %w", err)
		}
		accesses = append(accesses, access)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating access rows: %w", err)
	}

	return accesses, nil
}

func (r *AccessSqliteRepository) GetByUID(projectID uuid.UUID, uid string) (*models.Access, error) {
	query := `
		SELECT id, project_id, uid, pin, name, readonly, created_at, updated_at
		FROM access
		WHERE project_id = ? AND uid = ?
	`

	row := r.db.QueryRow(query, projectID.String(), uid)
	return r.scanAccess(row)
}

func (r *AccessSqliteRepository) ExistsByUID(projectID uuid.UUID, uid string) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM access
		WHERE project_id = ? AND uid = ?
	`

	var count int
	err := r.db.QueryRow(query, projectID.String(), uid).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check access existence: %w", err)
	}

	return count > 0, nil
}

func (r *AccessSqliteRepository) scanAccess(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Access, error) {
	var id, projectID, uid, pin, name string
	var readonly bool
	var createdAt, updatedAt time.Time

	err := scanner.Scan(&id, &projectID, &uid, &pin, &name, &readonly, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("access not found")
		}
		return nil, fmt.Errorf("failed to scan access row: %w", err)
	}

	accessID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid access ID: %w", err)
	}

	projID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	return &models.Access{
		ID:        accessID,
		ProjectID: projID,
		UID:       uid,
		PIN:       pin,
		Name:      name,
		ReadOnly:  readonly,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
