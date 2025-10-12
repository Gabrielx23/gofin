package database

import (
	"fmt"
	"sort"
	"sync"
	"time"

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

func (r *TransactionInMemoryRepository) GetByAccountIDWithDateRange(accountID uuid.UUID, startDate, endDate *time.Time) ([]*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*models.Transaction
	for _, transaction := range r.transactions {
		if transaction.AccountID == accountID {
			if r.isTransactionInDateRange(transaction, startDate, endDate) {
				transactions = append(transactions, transaction)
			}
		}
	}

	return transactions, nil
}

func (r *TransactionInMemoryRepository) GetByProjectIDWithDateRange(projectID uuid.UUID, startDate, endDate *time.Time) ([]*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*models.Transaction
	for _, transaction := range r.transactions {
		if r.isTransactionInDateRange(transaction, startDate, endDate) {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func (r *TransactionInMemoryRepository) isTransactionInDateRange(transaction *models.Transaction, startDate, endDate *time.Time) bool {
	if startDate != nil && transaction.TransactionDate.Before(*startDate) {
		return false
	}
	if endDate != nil && transaction.TransactionDate.After(*endDate) {
		return false
	}
	return true
}

func (r *TransactionInMemoryRepository) isTransactionInDateRangeWithFutureFilter(transaction *models.Transaction, startDate, endDate *time.Time, excludeFuture bool) bool {
	if startDate != nil && transaction.TransactionDate.Before(*startDate) {
		return false
	}
	if endDate != nil && transaction.TransactionDate.After(*endDate) {
		return false
	}
	if excludeFuture && endDate == nil && transaction.TransactionDate.After(time.Now()) {
		return false
	}
	return true
}

func (r *TransactionInMemoryRepository) GetTransactionsWithFilters(query models.TransactionQuery) ([]*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*models.Transaction
	for _, transaction := range r.transactions {
		if r.matchesFilters(transaction, query) {
			transactions = append(transactions, transaction)
		}
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].TransactionDate.After(transactions[j].TransactionDate)
	})

	return transactions, nil
}

func (r *TransactionInMemoryRepository) matchesFilters(transaction *models.Transaction, query models.TransactionQuery) bool {
	if query.ProjectID != nil {
		return r.isTransactionInDateRangeWithFutureFilter(transaction, query.StartDate, query.EndDate, query.ExcludeFutureTransactions)
	}

	if query.AccountID != nil && transaction.AccountID != *query.AccountID {
		return false
	}

	return r.isTransactionInDateRangeWithFutureFilter(transaction, query.StartDate, query.EndDate, query.ExcludeFutureTransactions)
}
