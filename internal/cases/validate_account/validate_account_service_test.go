package validate_account

import (
	"testing"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/pkg/money"
)

func TestValidateAccountService_ValidateAccountExists(t *testing.T) {
	accountID := uuid.New()

	tests := []struct {
		name      string
		accountID uuid.UUID
		repoSetup func(models.AccountRepository)
		wantErr   bool
		errorMsg  string
	}{
		{
			name:      "success when account exists",
			accountID: accountID,
			repoSetup: func(accountRepo models.AccountRepository) {
				project := models.NewProject("Test Project", "test-project")
				account := models.NewAccount(project.ID, "Test Account", money.USD)
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: false,
		},
		{
			name:      "error when account does not exist",
			accountID: uuid.New(),
			repoSetup: func(accountRepo models.AccountRepository) {
			},
			wantErr:  true,
			errorMsg: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := database.NewAccountInMemoryRepository()
			service := NewValidateAccountService(accountRepo)
			tt.repoSetup(accountRepo)

			err := service.ValidateAccountExists(tt.accountID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAccountExists() expected error, got nil")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateAccountExists() error message = %v, want to contain %v", err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateAccountExists() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAccountService_ValidateAccountForProject(t *testing.T) {
	projectID := uuid.New()
	otherProjectID := uuid.New()
	accountID := uuid.New()

	tests := []struct {
		name      string
		projectID uuid.UUID
		accountID uuid.UUID
		repoSetup func(models.AccountRepository)
		wantErr   bool
		errorMsg  string
	}{
		{
			name:      "success when account belongs to project",
			projectID: projectID,
			accountID: accountID,
			repoSetup: func(accountRepo models.AccountRepository) {
				project := models.NewProject("Test Project", "test-project")
				project.ID = projectID
				account := models.NewAccount(project.ID, "Test Account", money.USD)
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: false,
		},
		{
			name:      "error when account does not exist",
			projectID: projectID,
			accountID: accountID,
			repoSetup: func(accountRepo models.AccountRepository) {
			},
			wantErr:  true,
			errorMsg: "account not found",
		},
		{
			name:      "error when account belongs to different project",
			projectID: projectID,
			accountID: accountID,
			repoSetup: func(accountRepo models.AccountRepository) {
				project := models.NewProject("Other Project", "other-project")
				project.ID = otherProjectID
				account := models.NewAccount(project.ID, "Test Account", money.USD)
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr:  true,
			errorMsg: "account does not belong to the specified project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := database.NewAccountInMemoryRepository()
			service := NewValidateAccountService(accountRepo)
			tt.repoSetup(accountRepo)

			err := service.ValidateAccountForProject(tt.projectID, tt.accountID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAccountForProject() expected error, got nil")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateAccountForProject() error message = %v, want to contain %v", err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateAccountForProject() unexpected error: %v", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
