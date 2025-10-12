package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Debit TransactionType = "debit"
	TopUp TransactionType = "top-up"
)

func (t TransactionType) String() string {
	return string(t)
}

func (t TransactionType) IsValid() bool {
	return t == Debit || t == TopUp
}

func ParseTransactionType(s string) (TransactionType, error) {
	switch s {
	case "debit":
		return Debit, nil
	case "top-up", "topup", "top_up":
		return TopUp, nil
	default:
		return "", fmt.Errorf("invalid transaction type: %s", s)
	}
}

type TransactionData struct {
	AccountID       uuid.UUID
	Value           float64
	Name            string
	Type            TransactionType
	TransactionDate *time.Time
}

type Transaction struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	AccountID       uuid.UUID       `json:"account_id" db:"account_id"`
	Value           float64         `json:"value" db:"value"`
	Name            string          `json:"name" db:"name"`
	TransactionDate time.Time       `json:"transaction_date" db:"transaction_date"`
	Type            TransactionType `json:"type" db:"type"`
	GroupID         *uuid.UUID      `json:"group_id,omitempty" db:"group_id"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByAccountID(accountID uuid.UUID) ([]*Transaction, error)
	GetByGroupID(groupID uuid.UUID) ([]*Transaction, error)
	GetByID(id uuid.UUID) (*Transaction, error)
	GetByAccountIDWithDateRange(accountID uuid.UUID, startDate, endDate *time.Time) ([]*Transaction, error)
	GetByProjectIDWithDateRange(projectID uuid.UUID, startDate, endDate *time.Time) ([]*Transaction, error)
	GetTransactionsWithFilters(query TransactionQuery) ([]*Transaction, error)
}

type TransactionQuery struct {
	ProjectID           *uuid.UUID
	AccountID           *uuid.UUID
	StartDate           *time.Time
	EndDate             *time.Time
	ExcludeFutureTransactions bool
}

func (q *TransactionQuery) Validate() error {
	if q.ProjectID == nil && q.AccountID == nil {
		return fmt.Errorf("either project_id or account_id must be provided")
	}

	if q.ProjectID != nil && q.AccountID != nil {
		return fmt.Errorf("cannot specify both project_id and account_id")
	}

	if q.StartDate != nil && q.EndDate != nil {
		if q.EndDate.Before(*q.StartDate) {
			return fmt.Errorf("end_date cannot be before start_date")
		}
	}

	return nil
}

func NewTransaction(data TransactionData, groupID ...uuid.UUID) *Transaction {
	now := time.Now()
	transactionDate := now
	if data.TransactionDate != nil {
		transactionDate = *data.TransactionDate
	}

	var groupIDPtr *uuid.UUID
	if len(groupID) > 0 {
		groupIDPtr = &groupID[0]
	}

	return &Transaction{
		ID:              uuid.New(),
		AccountID:       data.AccountID,
		Value:           data.Value,
		Name:            data.Name,
		TransactionDate: transactionDate,
		Type:            data.Type,
		GroupID:         groupIDPtr,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
