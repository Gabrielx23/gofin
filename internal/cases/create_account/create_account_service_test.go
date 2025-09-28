package create_account

import (
	"testing"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/pkg/money"
)

func TestCreateAccountService_CreateAccount(t *testing.T) {
	tests := []struct {
		name         string
		projectSlug  string
		accountName  string
		currencyCode string
		repoSetup    func(models.ProjectRepository, models.AccountRepository)
		wantErr      bool
		wantCurrency money.Currency
	}{
		{
			name:         "success with USD currency",
			projectSlug:  "test-project",
			accountName:  "Checking Account",
			currencyCode: "USD",
			repoSetup: func(projectRepo models.ProjectRepository, accountRepo models.AccountRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
			},
			wantErr:      false,
			wantCurrency: money.USD,
		},
		{
			name:         "success with EUR currency",
			projectSlug:  "test-project",
			accountName:  "Savings Account",
			currencyCode: "EUR",
			repoSetup: func(projectRepo models.ProjectRepository, accountRepo models.AccountRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
			},
			wantErr:      false,
			wantCurrency: money.EUR,
		},
		{
			name:         "error when project not found",
			projectSlug:  "non-existent-project",
			accountName:  "Test Account",
			currencyCode: "USD",
			repoSetup:    func(projectRepo models.ProjectRepository, accountRepo models.AccountRepository) {},
			wantErr:      true,
		},
		{
			name:         "error when name is empty",
			projectSlug:  "test-project",
			accountName:  "",
			currencyCode: "USD",
			repoSetup: func(projectRepo models.ProjectRepository, accountRepo models.AccountRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
			},
			wantErr: true,
		},
		{
			name:         "error when currency is invalid",
			projectSlug:  "test-project",
			accountName:  "Test Account",
			currencyCode: "INVALID",
			repoSetup: func(projectRepo models.ProjectRepository, accountRepo models.AccountRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
			},
			wantErr: true,
		},
		{
			name:         "error when account name already exists",
			projectSlug:  "test-project",
			accountName:  "Existing Account",
			currencyCode: "USD",
			repoSetup: func(projectRepo models.ProjectRepository, accountRepo models.AccountRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)

				existingAccount := models.NewAccount(project.ID, "Existing Account", money.USD)
				accountRepo.Create(existingAccount)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRepo := database.NewProjectInMemoryRepository()
			accountRepo := database.NewAccountInMemoryRepository()
			service := NewCreateAccountService(accountRepo, projectRepo)
			tt.repoSetup(projectRepo, accountRepo)

			account, err := service.CreateAccount(tt.projectSlug, tt.accountName, tt.currencyCode)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateAccount() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateAccount() unexpected error: %v", err)
				return
			}

			if account == nil {
				t.Errorf("CreateAccount() expected account, got nil")
				return
			}

			if account.Name != tt.accountName {
				t.Errorf("CreateAccount() name = %v, want %v", account.Name, tt.accountName)
			}

			if account.Currency != tt.wantCurrency {
				t.Errorf("CreateAccount() currency = %v, want %v", account.Currency, tt.wantCurrency)
			}

			if account.ID == uuid.Nil {
				t.Errorf("CreateAccount() ID should not be nil")
			}
		})
	}
}
