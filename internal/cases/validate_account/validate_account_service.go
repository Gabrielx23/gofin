package validate_account

import (
	"fmt"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type ValidateAccountService struct {
	accountRepo models.AccountRepository
}

func NewValidateAccountService(accountRepo models.AccountRepository) *ValidateAccountService {
	return &ValidateAccountService{
		accountRepo: accountRepo,
	}
}

func (s *ValidateAccountService) ValidateAccountExists(accountID uuid.UUID) error {
	_, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	return nil
}

func (s *ValidateAccountService) ValidateAccountForProject(projectID uuid.UUID, accountID uuid.UUID) error {
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if account.ProjectID != projectID {
		return fmt.Errorf("account does not belong to the specified project")
	}

	return nil
}
