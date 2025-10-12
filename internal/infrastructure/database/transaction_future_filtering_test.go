package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gofin/internal/models"
)

func TestTransactionInMemoryRepository_FutureTransactionFiltering(t *testing.T) {
	repo := NewTransactionInMemoryRepository()
	accountID := uuid.New()

	now := time.Now()
	pastTime := now.AddDate(0, 0, -1)
	futureTime := now.AddDate(0, 0, 1)

	transactions := []*models.Transaction{
		{
			ID:              uuid.New(),
			AccountID:       accountID,
			Value:           100.0,
			Name:            "Past Transaction",
			TransactionDate: pastTime,
			Type:            models.Debit,
		},
		{
			ID:              uuid.New(),
			AccountID:       accountID,
			Value:           200.0,
			Name:            "Current Transaction",
			TransactionDate: now,
			Type:            models.TopUp,
		},
		{
			ID:              uuid.New(),
			AccountID:       accountID,
			Value:           300.0,
			Name:            "Future Transaction",
			TransactionDate: futureTime,
			Type:            models.Debit,
		},
	}

	for _, tx := range transactions {
		if err := repo.Create(tx); err != nil {
			t.Fatalf("Failed to create transaction: %v", err)
		}
	}

	tests := []struct {
		name                     string
		query                    models.TransactionQuery
		expectedTransactionCount int
		expectedTransactionNames []string
	}{
		{
			name: "ExcludeFutureTransactions=true, no end date - should exclude future transactions",
			query: models.TransactionQuery{
				AccountID:                   &accountID,
				ExcludeFutureTransactions:   true,
			},
			expectedTransactionCount: 2,
			expectedTransactionNames: []string{"Current Transaction", "Past Transaction"},
		},
		{
			name: "ExcludeFutureTransactions=false, no end date - should include all transactions",
			query: models.TransactionQuery{
				AccountID:                   &accountID,
				ExcludeFutureTransactions:   false,
			},
			expectedTransactionCount: 3,
			expectedTransactionNames: []string{"Future Transaction", "Current Transaction", "Past Transaction"},
		},
		{
			name: "ExcludeFutureTransactions=true, with end date in future - should include future transactions",
			query: models.TransactionQuery{
				AccountID:                   &accountID,
				EndDate:                     &futureTime,
				ExcludeFutureTransactions:   true,
			},
			expectedTransactionCount: 3,
			expectedTransactionNames: []string{"Future Transaction", "Current Transaction", "Past Transaction"},
		},
		{
			name: "ExcludeFutureTransactions=false, with end date in past - should exclude future transactions",
			query: models.TransactionQuery{
				AccountID:                   &accountID,
				EndDate:                     &now,
				ExcludeFutureTransactions:   false,
			},
			expectedTransactionCount: 2,
			expectedTransactionNames: []string{"Current Transaction", "Past Transaction"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions, err := repo.GetTransactionsWithFilters(tt.query)
			if err != nil {
				t.Fatalf("Failed to get transactions: %v", err)
			}

			if len(transactions) != tt.expectedTransactionCount {
				t.Errorf("Expected %d transactions, got %d", tt.expectedTransactionCount, len(transactions))
			}

			transactionNames := make([]string, len(transactions))
			for i, tx := range transactions {
				transactionNames[i] = tx.Name
			}

			for i, expectedName := range tt.expectedTransactionNames {
				if i >= len(transactionNames) || transactionNames[i] != expectedName {
					t.Errorf("Expected transaction %d to be '%s', got '%s'", i, expectedName, transactionNames[i])
				}
			}
		})
	}
}
