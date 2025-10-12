package create_account

import (
	"fmt"

	"github.com/google/uuid"
	"gofin/internal/models"
	"gofin/pkg/money"
)

type CreateAccountService struct {
	accountRepo models.AccountRepository
}

func NewCreateAccountService(accountRepo models.AccountRepository) *CreateAccountService {
	return &CreateAccountService{
		accountRepo: accountRepo,
	}
}

type CreateAccountData struct {
	ProjectID uuid.UUID
	Name      string
	Currency  money.Currency
}

func (s *CreateAccountService) CreateAccount(data CreateAccountData) (*models.Account, error) {
	if data.Name == "" {
		return nil, fmt.Errorf("account name is required")
	}

	exists, err := s.accountRepo.ExistsByName(data.ProjectID, data.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if account exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("account with name '%s' already exists for this project", data.Name)
	}

	account := models.NewAccount(data.ProjectID, data.Name, data.Currency)

	if err := s.accountRepo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}
