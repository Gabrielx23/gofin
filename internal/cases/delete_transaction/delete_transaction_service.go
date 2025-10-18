package delete_transaction

import (
	"fmt"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type DeleteTransactionService struct {
	transactionRepo models.TransactionRepository
}

func NewDeleteTransactionService(transactionRepo models.TransactionRepository) *DeleteTransactionService {
	return &DeleteTransactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *DeleteTransactionService) DeleteTransaction(transactionID uuid.UUID) error {
	_, err := s.transactionRepo.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	err = s.transactionRepo.DeleteByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}
