package get_transactions

import (
	"fmt"

	"gofin/internal/models"
)

type GetTransactionsService struct {
	transactionRepo models.TransactionRepository
	accountRepo     models.AccountRepository
	projectRepo     models.ProjectRepository
}

func NewGetTransactionsService(transactionRepo models.TransactionRepository, accountRepo models.AccountRepository, projectRepo models.ProjectRepository) *GetTransactionsService {
	return &GetTransactionsService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		projectRepo:     projectRepo,
	}
}

func (s *GetTransactionsService) GetTransactions(query models.TransactionQuery) ([]*models.Transaction, error) {
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	if query.AccountID != nil {
		_, err := s.accountRepo.GetByID(*query.AccountID)
		if err != nil {
			return nil, fmt.Errorf("account not found: %w", err)
		}
	}

	if query.ProjectID != nil {
		_, err := s.projectRepo.GetByID(*query.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("project not found: %w", err)
		}
	}

	transactions, err := s.transactionRepo.GetTransactionsWithFilters(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
