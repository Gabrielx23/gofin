package database

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type AccountInMemoryRepository struct {
	accounts map[string]*models.Account
	mu       sync.RWMutex
}

func NewAccountInMemoryRepository() *AccountInMemoryRepository {
	return &AccountInMemoryRepository{
		accounts: make(map[string]*models.Account),
	}
}

func (r *AccountInMemoryRepository) Create(account *models.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.getKey(account.ProjectID, account.Name)
	if _, exists := r.accounts[key]; exists {
		return fmt.Errorf("account with name '%s' already exists for project", account.Name)
	}

	r.accounts[key] = account
	return nil
}

func (r *AccountInMemoryRepository) GetByProjectID(projectID uuid.UUID) ([]*models.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var accounts []*models.Account
	for _, account := range r.accounts {
		if account.ProjectID == projectID {
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}

func (r *AccountInMemoryRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, account := range r.accounts {
		if account.ID == id {
			return account, nil
		}
	}

	return nil, fmt.Errorf("account with ID '%s' not found", id.String())
}

func (r *AccountInMemoryRepository) ExistsByName(projectID uuid.UUID, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.getKey(projectID, name)
	_, exists := r.accounts[key]
	return exists, nil
}

func (r *AccountInMemoryRepository) getKey(projectID uuid.UUID, name string) string {
	return fmt.Sprintf("%s:%s", projectID.String(), name)
}
