package models

import (
	"time"

	"github.com/google/uuid"
	"gofin/pkg/money"
)

type Account struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	ProjectID uuid.UUID      `json:"project_id" db:"project_id"`
	Name      string         `json:"name" db:"name"`
	Currency  money.Currency `json:"currency" db:"currency"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

type AccountRepository interface {
	Create(account *Account) error
	GetByProjectID(projectID uuid.UUID) ([]*Account, error)
	GetByID(id uuid.UUID) (*Account, error)
	ExistsByName(projectID uuid.UUID, name string) (bool, error)
}

func NewAccount(projectID uuid.UUID, name string, currency money.Currency) *Account {
	now := time.Now()
	return &Account{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      name,
		Currency:  currency,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func ParseCurrency(s string) (money.Currency, error) {
	return money.ParseCurrency(s)
}
