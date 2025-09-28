package get_transactions

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/pkg/money"
)

func TestGetTransactionsService_GetTransactions(t *testing.T) {
	t.Run("success with project ID - all transactions", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		project := models.NewProject("Test Project", "test-project")
		project.ID = projectID
		projectRepo.Create(project)

		account1 := models.NewAccount(projectID, "Account 1", money.USD)
		account1.ID = uuid.New()
		accountRepo.Create(account1)

		account2 := models.NewAccount(projectID, "Account 2", money.EUR)
		account2.ID = uuid.New()
		accountRepo.Create(account2)

		transaction1 := models.NewTransaction(models.TransactionData{
			AccountID: account1.ID,
			Value:     100.0,
			Name:      "Transaction 1",
			Type:      models.TopUp,
		})
		transactionRepo.Create(transaction1)

		transaction2 := models.NewTransaction(models.TransactionData{
			AccountID: account2.ID,
			Value:     50.0,
			Name:      "Transaction 2",
			Type:      models.Debit,
		})
		transactionRepo.Create(transaction2)

		query := models.TransactionQuery{
			ProjectID: &projectID,
		}

		transactions, err := service.GetTransactions(query)
		if err != nil {
			t.Errorf("GetTransactions() unexpected error: %v", err)
			return
		}

		if len(transactions) != 2 {
			t.Errorf("GetTransactions() got %d transactions, want 2", len(transactions))
		}
	})

	t.Run("success with account ID - single account transactions", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		project := models.NewProject("Test Project", "test-project")
		project.ID = projectID
		projectRepo.Create(project)

		accountID := uuid.New()
		account := models.NewAccount(projectID, "Account 1", money.USD)
		account.ID = accountID
		accountRepo.Create(account)

		account2 := models.NewAccount(projectID, "Account 2", money.EUR)
		account2.ID = uuid.New()
		accountRepo.Create(account2)

		transaction1 := models.NewTransaction(models.TransactionData{
			AccountID: accountID,
			Value:     100.0,
			Name:      "Transaction 1",
			Type:      models.TopUp,
		})
		transactionRepo.Create(transaction1)

		transaction2 := models.NewTransaction(models.TransactionData{
			AccountID: account2.ID,
			Value:     50.0,
			Name:      "Transaction 2",
			Type:      models.Debit,
		})
		transactionRepo.Create(transaction2)

		query := models.TransactionQuery{
			AccountID: &accountID,
		}

		transactions, err := service.GetTransactions(query)
		if err != nil {
			t.Errorf("GetTransactions() unexpected error: %v", err)
			return
		}

		if len(transactions) != 1 {
			t.Errorf("GetTransactions() got %d transactions, want 1", len(transactions))
		}

		if transactions[0].AccountID != accountID {
			t.Errorf("GetTransactions() got transaction for account %s, want %s", transactions[0].AccountID, accountID)
		}
	})

	t.Run("success with date range filtering", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		project := models.NewProject("Test Project", "test-project")
		project.ID = projectID
		projectRepo.Create(project)

		accountID := uuid.New()
		account := models.NewAccount(projectID, "Account 1", money.USD)
		account.ID = accountID
		accountRepo.Create(account)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(-12 * time.Hour)

		transaction1 := models.NewTransaction(models.TransactionData{
			AccountID:       accountID,
			Value:           100.0,
			Name:            "Transaction 1",
			Type:            models.TopUp,
			TransactionDate: &startDate,
		})
		transactionRepo.Create(transaction1)

		transaction2 := models.NewTransaction(models.TransactionData{
			AccountID:       accountID,
			Value:           50.0,
			Name:            "Transaction 2",
			Type:            models.Debit,
			TransactionDate: &endDate,
		})
		transactionRepo.Create(transaction2)

		transaction3 := models.NewTransaction(models.TransactionData{
			AccountID:       accountID,
			Value:           25.0,
			Name:            "Transaction 3",
			Type:            models.TopUp,
			TransactionDate: &now,
		})
		transactionRepo.Create(transaction3)

		query := models.TransactionQuery{
			AccountID: &accountID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		transactions, err := service.GetTransactions(query)
		if err != nil {
			t.Errorf("GetTransactions() unexpected error: %v", err)
			return
		}

		if len(transactions) != 2 {
			t.Errorf("GetTransactions() got %d transactions, want 2", len(transactions))
		}
	})

	t.Run("error when neither project_id nor account_id provided", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		query := models.TransactionQuery{}

		_, err := service.GetTransactions(query)
		if err == nil {
			t.Errorf("GetTransactions() expected error, got nil")
		}
	})

	t.Run("error when both project_id and account_id provided", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		accountID := uuid.New()
		query := models.TransactionQuery{
			ProjectID: &projectID,
			AccountID: &accountID,
		}

		_, err := service.GetTransactions(query)
		if err == nil {
			t.Errorf("GetTransactions() expected error, got nil")
		}
	})

	t.Run("error when account not found", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		accountID := uuid.New()
		query := models.TransactionQuery{
			AccountID: &accountID,
		}

		_, err := service.GetTransactions(query)
		if err == nil {
			t.Errorf("GetTransactions() expected error, got nil")
		}
	})

	t.Run("error when project not found", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		query := models.TransactionQuery{
			ProjectID: &projectID,
		}

		_, err := service.GetTransactions(query)
		if err == nil {
			t.Errorf("GetTransactions() expected error, got nil")
		}
	})

	t.Run("error when end_date before start_date", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetTransactionsService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		startDate := time.Now()
		endDate := startDate.Add(-24 * time.Hour)
		query := models.TransactionQuery{
			ProjectID: &projectID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		_, err := service.GetTransactions(query)
		if err == nil {
			t.Errorf("GetTransactions() expected error, got nil")
		}
	})
}

func TestTransactionQuery_Validate(t *testing.T) {
	tests := []struct {
		name    string
		query   models.TransactionQuery
		wantErr bool
	}{
		{
			name: "valid with project ID",
			query: models.TransactionQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			wantErr: false,
		},
		{
			name: "valid with account ID",
			query: models.TransactionQuery{
				AccountID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			wantErr: false,
		},
		{
			name: "valid with date range",
			query: models.TransactionQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				StartDate: func() *time.Time { t := time.Now(); return &t }(),
				EndDate:   func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
			},
			wantErr: false,
		},
		{
			name:    "error when neither project_id nor account_id",
			query:   models.TransactionQuery{},
			wantErr: true,
		},
		{
			name: "error when both project_id and account_id",
			query: models.TransactionQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				AccountID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			wantErr: true,
		},
		{
			name: "error when end_date before start_date",
			query: models.TransactionQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				StartDate: func() *time.Time { t := time.Now(); return &t }(),
				EndDate:   func() *time.Time { t := time.Now().Add(-24 * time.Hour); return &t }(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionQuery.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
