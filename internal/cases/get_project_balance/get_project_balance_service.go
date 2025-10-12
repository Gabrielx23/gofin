package get_project_balance

import (
	"fmt"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type GetProjectBalanceService struct {
	transactionRepo models.TransactionRepository
	accountRepo     models.AccountRepository
}

func NewGetProjectBalanceService(transactionRepo models.TransactionRepository, accountRepo models.AccountRepository) *GetProjectBalanceService {
	return &GetProjectBalanceService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

type ProjectBalanceData struct {
	ProjectID       uuid.UUID               `json:"project_id"`
	AccountBalances []models.AccountSummary `json:"account_balances"`
}

func (s *GetProjectBalanceService) GetProjectBalances(projectID uuid.UUID) (*ProjectBalanceData, error) {
	accounts, err := s.accountRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project accounts: %w", err)
	}

	var accountBalances []models.AccountSummary
	for _, account := range accounts {
		balance, err := s.calculateAccountBalance(account.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate balance for account %s: %w", account.Name, err)
		}

		accountBalances = append(accountBalances, models.AccountSummary{
			AccountID: account.ID,
			Name:      account.Name,
			Currency:  account.Currency.String(),
			Balance:   balance,
		})
	}

	return &ProjectBalanceData{
		ProjectID:       projectID,
		AccountBalances: accountBalances,
	}, nil
}

func (s *GetProjectBalanceService) calculateAccountBalance(accountID uuid.UUID) (float64, error) {
	query := models.TransactionQuery{
		AccountID:                 &accountID,
		ExcludeFutureTransactions: true,
	}

	transactions, err := s.transactionRepo.GetTransactionsWithFilters(query)
	if err != nil {
		return 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	var balance float64

	for _, transaction := range transactions {
		if transaction.Type == models.Debit {
			balance -= transaction.Value
		} else {
			balance += transaction.Value
		}
	}

	return balance, nil
}

func (s *GetProjectBalanceService) FormatBalance(balance float64, currency string) string {
	return fmt.Sprintf("%.2f %s", balance, currency)
}
