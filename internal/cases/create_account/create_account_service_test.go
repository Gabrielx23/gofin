package create_account

import (
	"testing"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/pkg/money"
)

func TestCreateAccountService_CreateAccount(t *testing.T) {
	accountRepo := database.NewAccountInMemoryRepository()
	service := NewCreateAccountService(accountRepo)

	projectID := uuid.New()

	tests := []struct {
		name        string
		data        CreateAccountData
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful account creation with PLN",
			data: CreateAccountData{
				ProjectID: projectID,
				Name:      "Test Account PLN",
				Currency:  money.Currency("PLN"),
			},
			expectError: false,
		},
		{
			name: "successful account creation with EUR",
			data: CreateAccountData{
				ProjectID: projectID,
				Name:      "Test Account EUR",
				Currency:  money.Currency("EUR"),
			},
			expectError: false,
		},
		{
			name: "error when account name is empty",
			data: CreateAccountData{
				ProjectID: projectID,
				Name:      "",
				Currency:  money.Currency("PLN"),
			},
			expectError: true,
			errorMsg:    "account name is required",
		},
		{
			name: "error when account name already exists",
			data: CreateAccountData{
				ProjectID: projectID,
				Name:      "Duplicate Account",
				Currency:  money.Currency("PLN"),
			},
			expectError: true,
			errorMsg:    "account with name 'Duplicate Account' already exists for this project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "error when account name already exists" {
				existingAccount := models.NewAccount(projectID, "Duplicate Account", money.Currency("PLN"))
				accountRepo.Create(existingAccount)
			}

			account, err := service.CreateAccount(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if account != nil {
					t.Errorf("Expected no account to be returned on error, got %v", account)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
					return
				}
				if account == nil {
					t.Errorf("Expected account to be returned but got nil")
					return
				}
				if account.Name != tt.data.Name {
					t.Errorf("Expected account name '%s', got '%s'", tt.data.Name, account.Name)
				}
				if account.Currency != tt.data.Currency {
					t.Errorf("Expected account currency '%s', got '%s'", tt.data.Currency, account.Currency)
				}
				if account.ProjectID != tt.data.ProjectID {
					t.Errorf("Expected account project ID '%s', got '%s'", tt.data.ProjectID, account.ProjectID)
				}
			}
		})
	}
}
