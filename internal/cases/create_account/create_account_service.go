package create_account

import (
	"fmt"

	"gofin/internal/models"
	"gofin/pkg/money"
)

type CreateAccountService struct {
	accountRepo models.AccountRepository
	projectRepo models.ProjectRepository
}

func NewCreateAccountService(accountRepo models.AccountRepository, projectRepo models.ProjectRepository) *CreateAccountService {
	return &CreateAccountService{
		accountRepo: accountRepo,
		projectRepo: projectRepo,
	}
}

func (s *CreateAccountService) CreateAccount(projectSlug, name, currencyCode string) (*models.Account, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	currency, err := money.ParseCurrency(currencyCode)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	project, err := s.projectRepo.GetBySlug(projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	exists, err := s.accountRepo.ExistsByName(project.ID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check account name existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("account with name '%s' already exists for this project", name)
	}

	account := models.NewAccount(project.ID, name, currency)

	if err := s.accountRepo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}
