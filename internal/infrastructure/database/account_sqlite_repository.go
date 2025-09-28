package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type AccountSqliteRepository struct {
	db *sql.DB
}

func NewAccountSqliteRepository(db *sql.DB) *AccountSqliteRepository {
	return &AccountSqliteRepository{db: db}
}

func (r *AccountSqliteRepository) Create(account *models.Account) error {
	query := `
		INSERT INTO accounts (id, project_id, name, currency, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		account.ID.String(),
		account.ProjectID.String(),
		account.Name,
		account.Currency.String(),
		account.CreatedAt,
		account.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (r *AccountSqliteRepository) GetByProjectID(projectID uuid.UUID) ([]*models.Account, error) {
	query := `
		SELECT id, project_id, name, currency, created_at, updated_at
		FROM accounts
		WHERE project_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, projectID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts by project_id: %w", err)
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		account, err := r.scanAccount(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating account rows: %w", err)
	}

	return accounts, nil
}

func (r *AccountSqliteRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	query := `
		SELECT id, project_id, name, currency, created_at, updated_at
		FROM accounts
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id.String())
	return r.scanAccount(row)
}

func (r *AccountSqliteRepository) ExistsByName(projectID uuid.UUID, name string) (bool, error) {
	query := `SELECT COUNT(*) FROM accounts WHERE project_id = ? AND name = ?`

	var count int
	err := r.db.QueryRow(query, projectID.String(), name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}

	return count > 0, nil
}

func (r *AccountSqliteRepository) scanAccount(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Account, error) {
	var id, projectID, name, currency string
	var createdAt, updatedAt time.Time

	err := scanner.Scan(&id, &projectID, &name, &currency, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to scan account row: %w", err)
	}

	accountID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	projID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	currencyType, err := models.ParseCurrency(currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	return &models.Account{
		ID:        accountID,
		ProjectID: projID,
		Name:      name,
		Currency:  currencyType,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
