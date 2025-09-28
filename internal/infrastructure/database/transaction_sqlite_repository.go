package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type TransactionSqliteRepository struct {
	db *sql.DB
}

func NewTransactionSqliteRepository(db *sql.DB) *TransactionSqliteRepository {
	return &TransactionSqliteRepository{db: db}
}

func (r *TransactionSqliteRepository) Create(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (id, account_id, value, name, transaction_date, type, group_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var groupID *string
	if transaction.GroupID != nil {
		groupIDStr := transaction.GroupID.String()
		groupID = &groupIDStr
	}

	_, err := r.db.Exec(
		query,
		transaction.ID.String(),
		transaction.AccountID.String(),
		transaction.Value,
		transaction.Name,
		transaction.TransactionDate,
		transaction.Type.String(),
		groupID,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *TransactionSqliteRepository) GetByAccountID(accountID uuid.UUID) ([]*models.Transaction, error) {
	query := `
		SELECT id, account_id, value, name, transaction_date, type, group_id, created_at, updated_at
		FROM transactions
		WHERE account_id = ?
		ORDER BY transaction_date DESC, created_at DESC
	`

	rows, err := r.db.Query(query, accountID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by account_id: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

func (r *TransactionSqliteRepository) GetByGroupID(groupID uuid.UUID) ([]*models.Transaction, error) {
	query := `
		SELECT id, account_id, value, name, transaction_date, type, group_id, created_at, updated_at
		FROM transactions
		WHERE group_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, groupID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by group_id: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

func (r *TransactionSqliteRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	query := `
		SELECT id, account_id, value, name, transaction_date, type, group_id, created_at, updated_at
		FROM transactions
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id.String())
	return r.scanTransaction(row)
}

func (r *TransactionSqliteRepository) scanTransaction(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Transaction, error) {
	var id, accountID, name, transactionType string
	var value float64
	var transactionDate, createdAt, updatedAt time.Time
	var groupIDStr sql.NullString

	err := scanner.Scan(&id, &accountID, &value, &name, &transactionDate, &transactionType, &groupIDStr, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to scan transaction row: %w", err)
	}

	transactionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	parsedType, err := models.ParseTransactionType(transactionType)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction type: %w", err)
	}

	var groupID *uuid.UUID
	if groupIDStr.Valid && groupIDStr.String != "" {
		groupUUID, err := uuid.Parse(groupIDStr.String)
		if err != nil {
			return nil, fmt.Errorf("invalid group ID: %w", err)
		}
		groupID = &groupUUID
	}

	return &models.Transaction{
		ID:              transactionID,
		AccountID:       accountUUID,
		Value:           value,
		Name:            name,
		TransactionDate: transactionDate,
		Type:            parsedType,
		GroupID:         groupID,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

func (r *TransactionSqliteRepository) GetByAccountIDWithDateRange(accountID uuid.UUID, startDate, endDate *time.Time) ([]*models.Transaction, error) {
	query := `
		SELECT t.id, t.account_id, t.value, t.name, t.transaction_date, t.type, t.group_id, t.created_at, t.updated_at
		FROM transactions t
		WHERE t.account_id = ?
	`
	args := []interface{}{accountID.String()}

	if startDate != nil {
		query += " AND t.transaction_date >= ?"
		args = append(args, *startDate)
	}

	if endDate != nil {
		query += " AND t.transaction_date <= ?"
		args = append(args, *endDate)
	}

	query += " ORDER BY t.transaction_date ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by account_id with date range: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

func (r *TransactionSqliteRepository) GetByProjectIDWithDateRange(projectID uuid.UUID, startDate, endDate *time.Time) ([]*models.Transaction, error) {
	query := `
		SELECT t.id, t.account_id, t.value, t.name, t.transaction_date, t.type, t.group_id, t.created_at, t.updated_at
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.project_id = ?
	`
	args := []interface{}{projectID.String()}

	if startDate != nil {
		query += " AND t.transaction_date >= ?"
		args = append(args, *startDate)
	}

	if endDate != nil {
		query += " AND t.transaction_date <= ?"
		args = append(args, *endDate)
	}

	query += " ORDER BY t.transaction_date ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by project_id with date range: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}
