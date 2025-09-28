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
	AccountID uuid.UUID
	Value     float64
	Name      string
	Type      TransactionType
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
}

func NewTransaction(data TransactionData) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:              uuid.New(),
		AccountID:       data.AccountID,
		Value:           data.Value,
		Name:            data.Name,
		TransactionDate: now,
		Type:            data.Type,
		GroupID:         nil,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func NewGroupedTransaction(data TransactionData, groupID uuid.UUID) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:              uuid.New(),
		AccountID:       data.AccountID,
		Value:           data.Value,
		Name:            data.Name,
		TransactionDate: now,
		Type:            data.Type,
		GroupID:         &groupID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
