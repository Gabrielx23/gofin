package create_transaction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gofin/internal/cases/validate_account"
	"gofin/internal/models"
)

type CreateTransactionService struct {
	transactionRepo    models.TransactionRepository
	accountRepo        models.AccountRepository
	projectRepo        models.ProjectRepository
	validateAccountSvc *validate_account.ValidateAccountService
}

func NewCreateTransactionService(transactionRepo models.TransactionRepository, accountRepo models.AccountRepository, projectRepo models.ProjectRepository) *CreateTransactionService {
	return &CreateTransactionService{
		transactionRepo:    transactionRepo,
		accountRepo:        accountRepo,
		projectRepo:        projectRepo,
		validateAccountSvc: validate_account.NewValidateAccountService(accountRepo),
	}
}

func (s *CreateTransactionService) CreateGroupedTransactions(projectID uuid.UUID, transactions []models.TransactionData) ([]*models.Transaction, error) {
	if len(transactions) == 0 {
		return nil, fmt.Errorf("at least one transaction is required")
	}

	for _, txData := range transactions {
		if err := s.validateAccountSvc.ValidateAccountForProject(projectID, txData.AccountID); err != nil {
			return nil, err
		}
	}

	groupID := uuid.New()
	var createdTransactions []*models.Transaction

	for _, txData := range transactions {
		if err := s.validateTransactionData(txData); err != nil {
			return nil, err
		}

		transaction := models.NewTransaction(txData, groupID)

		if err := s.transactionRepo.Create(transaction); err != nil {
			return nil, fmt.Errorf("failed to create transaction: %w", err)
		}

		createdTransactions = append(createdTransactions, transaction)
	}

	return createdTransactions, nil
}

func (s *CreateTransactionService) validateTransactionData(data models.TransactionData) error {
	if data.Name == "" {
		return fmt.Errorf("name is required")
	}

	if data.Value <= 0 {
		return fmt.Errorf("value must be positive")
	}

	if !data.Type.IsValid() {
		return fmt.Errorf("invalid transaction type: %s", data.Type)
	}

	if data.TransactionDate != nil {
		now := time.Now()
		transactionDate := *data.TransactionDate

		if transactionDate.Before(now.AddDate(-10, 0, 0)) {
			return fmt.Errorf("transaction date cannot be more than 10 years in the past")
		}
	}

	return nil
}
