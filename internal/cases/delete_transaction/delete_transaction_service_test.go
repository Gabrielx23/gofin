package delete_transaction

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/pkg/money"
)

func TestDeleteTransactionService_DeleteTransaction(t *testing.T) {
	transactionRepo := database.NewTransactionInMemoryRepository()
	service := NewDeleteTransactionService(transactionRepo)

	projectID := uuid.New()
	account := models.NewAccount(projectID, "Test Account", money.PLN)
	accountRepo := database.NewAccountInMemoryRepository()
	accountRepo.Create(account)

	transaction := models.NewTransaction(models.TransactionData{
		AccountID: account.ID,
		Value:     100.0,
		Name:      "Test Transaction",
		Type:      models.Debit,
	}, uuid.New())

	transactionRepo.Create(transaction)

	err := service.DeleteTransaction(transaction.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = transactionRepo.GetByID(transaction.ID)
	if err == nil {
		t.Errorf("Expected transaction to be deleted, but it still exists")
	}
}

func TestDeleteTransactionService_DeleteTransaction_NotFound(t *testing.T) {
	transactionRepo := database.NewTransactionInMemoryRepository()
	service := NewDeleteTransactionService(transactionRepo)

	nonExistentID := uuid.New()

	err := service.DeleteTransaction(nonExistentID)
	if err == nil {
		t.Fatalf("Expected error for non-existent transaction, got nil")
	}

	expectedError := "transaction not found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDeleteTransactionService_DeleteTransaction_RepositoryError(t *testing.T) {
	transactionRepo := database.NewTransactionInMemoryRepository()
	service := NewDeleteTransactionService(transactionRepo)

	projectID := uuid.New()
	account := models.NewAccount(projectID, "Test Account", money.PLN)
	accountRepo := database.NewAccountInMemoryRepository()
	accountRepo.Create(account)

	transaction := models.NewTransaction(models.TransactionData{
		AccountID: account.ID,
		Value:     100.0,
		Name:      "Test Transaction",
		Type:      models.Debit,
	}, uuid.New())

	transactionRepo.Create(transaction)

	err := service.DeleteTransaction(transaction.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = service.DeleteTransaction(transaction.ID)
	if err == nil {
		t.Fatalf("Expected error when deleting already deleted transaction, got nil")
	}

	expectedError := "transaction not found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}
