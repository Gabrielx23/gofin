package get_project_balance

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/pkg/money"
)

func TestGetProjectBalanceService_GetProjectBalance(t *testing.T) {
	t.Run("success with project ID - all accounts", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetProjectBalanceService(transactionRepo, accountRepo, projectRepo)

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
			Name:      "Top-up 1",
			Type:      models.TopUp,
		})
		transactionRepo.Create(transaction1)

		transaction2 := models.NewTransaction(models.TransactionData{
			AccountID: account1.ID,
			Value:     50.0,
			Name:      "Debit 1",
			Type:      models.Debit,
		})
		transactionRepo.Create(transaction2)

		transaction3 := models.NewTransaction(models.TransactionData{
			AccountID: account2.ID,
			Value:     200.0,
			Name:      "Top-up 2",
			Type:      models.TopUp,
		})
		transactionRepo.Create(transaction3)

		query := models.BalanceQuery{
			ProjectID: &projectID,
		}

		summary, err := service.GetProjectBalance(query)
		if err != nil {
			t.Errorf("GetProjectBalance() unexpected error: %v", err)
			return
		}

		if summary == nil {
			t.Errorf("GetProjectBalance() returned nil summary")
			return
		}

		if len(summary.Summaries) != 2 {
			t.Errorf("GetProjectBalance() got %d summaries, want 2", len(summary.Summaries))
		}

		account1Balance := 0.0
		account2Balance := 0.0
		for _, accountSummary := range summary.Summaries {
			if accountSummary.AccountID == account1.ID {
				account1Balance = accountSummary.Balance
			} else if accountSummary.AccountID == account2.ID {
				account2Balance = accountSummary.Balance
			}
		}

		if account1Balance != 50.0 { // 100 - 50
			t.Errorf("GetProjectBalance() account1 balance = %f, want 50.0", account1Balance)
		}

		if account2Balance != 200.0 { // 200 - 0
			t.Errorf("GetProjectBalance() account2 balance = %f, want 200.0", account2Balance)
		}
	})

	t.Run("success with account ID - single account", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetProjectBalanceService(transactionRepo, accountRepo, projectRepo)

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
			Name:      "Top-up 1",
			Type:      models.TopUp,
		})
		transactionRepo.Create(transaction1)

		transaction2 := models.NewTransaction(models.TransactionData{
			AccountID: accountID,
			Value:     30.0,
			Name:      "Debit 1",
			Type:      models.Debit,
		})
		transactionRepo.Create(transaction2)

		query := models.BalanceQuery{
			AccountID: &accountID,
		}

		summary, err := service.GetProjectBalance(query)
		if err != nil {
			t.Errorf("GetProjectBalance() unexpected error: %v", err)
			return
		}

		if summary == nil {
			t.Errorf("GetProjectBalance() returned nil summary")
			return
		}

		if len(summary.Summaries) != 1 {
			t.Errorf("GetProjectBalance() got %d summaries, want 1 (only the queried account)", len(summary.Summaries))
		}

		accountSummary := summary.Summaries[0]
		if accountSummary.AccountID != accountID {
			t.Errorf("GetProjectBalance() returned wrong account ID: %s, want %s", accountSummary.AccountID, accountID)
		}

		if accountSummary.Balance != 70.0 { // 100 - 30
			t.Errorf("GetProjectBalance() account balance = %f, want 70.0", accountSummary.Balance)
		}
	})

	t.Run("error when neither project_id nor account_id provided", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetProjectBalanceService(transactionRepo, accountRepo, projectRepo)

		query := models.BalanceQuery{}

		_, err := service.GetProjectBalance(query)
		if err == nil {
			t.Errorf("GetProjectBalance() expected error, got nil")
		}
	})

	t.Run("error when both project_id and account_id provided", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetProjectBalanceService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		accountID := uuid.New()
		query := models.BalanceQuery{
			ProjectID: &projectID,
			AccountID: &accountID,
		}

		_, err := service.GetProjectBalance(query)
		if err == nil {
			t.Errorf("GetProjectBalance() expected error, got nil")
		}
	})

	t.Run("error when end_date before start_date", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetProjectBalanceService(transactionRepo, accountRepo, projectRepo)

		projectID := uuid.New()
		startDate := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		query := models.BalanceQuery{
			ProjectID: &projectID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		_, err := service.GetProjectBalance(query)
		if err == nil {
			t.Errorf("GetProjectBalance() expected error, got nil")
		}
	})

	t.Run("error when account not found", func(t *testing.T) {
		projectRepo := database.NewProjectInMemoryRepository()
		accountRepo := database.NewAccountInMemoryRepository()
		transactionRepo := database.NewTransactionInMemoryRepository()
		service := NewGetProjectBalanceService(transactionRepo, accountRepo, projectRepo)

		accountID := uuid.New()
		query := models.BalanceQuery{
			AccountID: &accountID,
		}

		_, err := service.GetProjectBalance(query)
		if err == nil {
			t.Errorf("GetProjectBalance() expected error, got nil")
		}
	})
}

func TestBalanceQuery_Validate(t *testing.T) {
	tests := []struct {
		name    string
		query   models.BalanceQuery
		wantErr bool
	}{
		{
			name: "valid with project ID",
			query: models.BalanceQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			wantErr: false,
		},
		{
			name: "valid with account ID",
			query: models.BalanceQuery{
				AccountID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			wantErr: false,
		},
		{
			name: "valid with date range",
			query: models.BalanceQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				StartDate: func() *time.Time { t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
				EndDate:   func() *time.Time { t := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC); return &t }(),
			},
			wantErr: false,
		},
		{
			name:    "error when neither project_id nor account_id",
			query:   models.BalanceQuery{},
			wantErr: true,
		},
		{
			name: "error when both project_id and account_id",
			query: models.BalanceQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				AccountID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			wantErr: true,
		},
		{
			name: "error when end_date before start_date",
			query: models.BalanceQuery{
				ProjectID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				StartDate: func() *time.Time { t := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC); return &t }(),
				EndDate:   func() *time.Time { t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceQuery.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
