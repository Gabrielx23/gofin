package database

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type TransactionInMemoryRepository struct {
	transactions map[string]*models.Transaction
	mu           sync.RWMutex
}

func NewTransactionInMemoryRepository() *TransactionInMemoryRepository {
	return &TransactionInMemoryRepository{
		transactions: make(map[string]*models.Transaction),
	}
}

func (r *TransactionInMemoryRepository) Create(transaction *models.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := transaction.ID.String()
	if _, exists := r.transactions[key]; exists {
		return fmt.Errorf("transaction with ID '%s' already exists", transaction.ID.String())
	}

	r.transactions[key] = transaction
	return nil
}

func (r *TransactionInMemoryRepository) GetByAccountID(accountID uuid.UUID) ([]*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*models.Transaction
	for _, transaction := range r.transactions {
		if transaction.AccountID == accountID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func (r *TransactionInMemoryRepository) GetByGroupID(groupID uuid.UUID) ([]*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*models.Transaction
	for _, transaction := range r.transactions {
		if transaction.GroupID != nil && *transaction.GroupID == groupID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func (r *TransactionInMemoryRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transaction, exists := r.transactions[id.String()]
	if !exists {
		return nil, fmt.Errorf("transaction with ID '%s' not found", id.String())
	}

	return transaction, nil
}
