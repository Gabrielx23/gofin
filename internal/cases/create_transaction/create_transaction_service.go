package create_transaction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type CreateTransactionService struct {
	transactionRepo models.TransactionRepository
	accountRepo     models.AccountRepository
}

func NewCreateTransactionService(transactionRepo models.TransactionRepository, accountRepo models.AccountRepository) *CreateTransactionService {
	return &CreateTransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (s *CreateTransactionService) CreateSingleTransaction(accountID uuid.UUID, value float64, name string, transactionType models.TransactionType) (*models.Transaction, error) {
	return s.CreateSingleTransactionFromData(models.TransactionData{
		AccountID: accountID,
		Value:     value,
		Name:      name,
		Type:      transactionType,
	})
}

func (s *CreateTransactionService) CreateSingleTransactionFromData(data models.TransactionData) (*models.Transaction, error) {
	if err := s.validateAccount(data.AccountID); err != nil {
		return nil, err
	}

	if err := s.validateTransactionData(data); err != nil {
		return nil, err
	}

	transaction := models.NewTransaction(data)

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (s *CreateTransactionService) CreateGroupedTransactions(transactions []models.TransactionData) ([]*models.Transaction, error) {
	if len(transactions) == 0 {
		return nil, fmt.Errorf("at least one transaction is required")
	}

	groupID := uuid.New()
	var createdTransactions []*models.Transaction

	for _, txData := range transactions {
		if err := s.validateAccount(txData.AccountID); err != nil {
			return nil, err
		}

		if err := s.validateTransactionData(txData); err != nil {
			return nil, err
		}

		transaction := models.NewGroupedTransaction(txData, groupID)

		if err := s.transactionRepo.Create(transaction); err != nil {
			return nil, fmt.Errorf("failed to create transaction: %w", err)
		}

		createdTransactions = append(createdTransactions, transaction)
	}

	return createdTransactions, nil
}

func (s *CreateTransactionService) validateAccount(accountID uuid.UUID) error {
	_, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	return nil
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

	if data.TransactionDate == nil {
		return nil
	}

	now := time.Now()
	if data.TransactionDate.After(now) {
		return fmt.Errorf("transaction date cannot be in the future")
	}

	tenYearsAgo := now.AddDate(-10, 0, 0)
	if data.TransactionDate.Before(tenYearsAgo) {
		return fmt.Errorf("transaction date cannot be more than 10 years in the past")
	}

	return nil
}
