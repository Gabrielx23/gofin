package create_transaction

import (
	"testing"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
)

func TestCreateTransactionService_CreateGroupedTransactions(t *testing.T) {
	tests := []struct {
		name         string
		transactions []models.TransactionData
		repoSetup    func(models.AccountRepository, models.TransactionRepository, models.ProjectRepository, []uuid.UUID, uuid.UUID)
		wantErr      bool
		wantCount    int
	}{
		{
			name: "success with multiple transactions",
			transactions: []models.TransactionData{
				{AccountID: uuid.New(), Value: 50.0, Name: "Transaction 1", Type: models.Debit, TransactionDate: nil},
				{AccountID: uuid.New(), Value: 100.0, Name: "Transaction 2", Type: models.TopUp, TransactionDate: nil},
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, projectRepo models.ProjectRepository, accountIDs []uuid.UUID, projectID uuid.UUID) {
				account1 := models.NewAccount(projectID, "Account 1", "USD")
				account1.ID = accountIDs[0]
				accountRepo.Create(account1)

				account2 := models.NewAccount(projectID, "Account 2", "USD")
				account2.ID = accountIDs[1]
				accountRepo.Create(account2)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:         "error when no transactions provided",
			transactions: []models.TransactionData{},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, projectRepo models.ProjectRepository, accountIDs []uuid.UUID, projectID uuid.UUID) {
			},
			wantErr: true,
		},
		{
			name: "error when account not found",
			transactions: []models.TransactionData{
				{AccountID: uuid.New(), Value: 50.0, Name: "Transaction 1", Type: models.Debit, TransactionDate: nil},
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, projectRepo models.ProjectRepository, accountIDs []uuid.UUID, projectID uuid.UUID) {
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := database.NewAccountInMemoryRepository()
			transactionRepo := database.NewTransactionInMemoryRepository()
			projectRepo := database.NewProjectInMemoryRepository()
			service := NewCreateTransactionService(transactionRepo, accountRepo, projectRepo)

			var accountIDs []uuid.UUID
			for _, tx := range tt.transactions {
				accountIDs = append(accountIDs, tx.AccountID)
			}
			var projectID uuid.UUID
			if len(tt.transactions) > 0 {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
				projectID = project.ID
			}
			
			tt.repoSetup(accountRepo, transactionRepo, projectRepo, accountIDs, projectID)
			transactions, err := service.CreateGroupedTransactions(projectID, tt.transactions)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateGroupedTransactions() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateGroupedTransactions() unexpected error: %v", err)
				return
			}

			if len(transactions) != tt.wantCount {
				t.Errorf("CreateGroupedTransactions() count = %v, want %v", len(transactions), tt.wantCount)
			}

			if len(transactions) > 0 {
				groupID := transactions[0].GroupID
				if groupID == nil {
					t.Errorf("CreateGroupedTransactions() GroupID should not be nil for grouped transactions")
				}

				for i, transaction := range transactions {
					if transaction.GroupID == nil || *transaction.GroupID != *groupID {
						t.Errorf("CreateGroupedTransactions() transaction %d should have same GroupID", i)
					}
				}
			}
		})
	}
}