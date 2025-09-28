package get_project_balance

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type GetProjectBalanceService struct {
	transactionRepo models.TransactionRepository
	accountRepo     models.AccountRepository
	projectRepo     models.ProjectRepository
}

func NewGetProjectBalanceService(transactionRepo models.TransactionRepository, accountRepo models.AccountRepository, projectRepo models.ProjectRepository) *GetProjectBalanceService {
	return &GetProjectBalanceService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		projectRepo:     projectRepo,
	}
}

func (s *GetProjectBalanceService) GetProjectBalance(query models.BalanceQuery) (*models.ProjectBalanceSummary, error) {
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	if query.ProjectID != nil {
		return s.getProjectBalance(*query.ProjectID, query.StartDate, query.EndDate)
	}

	return s.getAccountBalance(*query.AccountID, query.StartDate, query.EndDate)
}

func (s *GetProjectBalanceService) getProjectBalance(projectID uuid.UUID, startDate, endDate *time.Time) (*models.ProjectBalanceSummary, error) {
	transactions, err := s.transactionRepo.GetByProjectIDWithDateRange(projectID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for project: %w", err)
	}

	accounts, err := s.accountRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts for project: %w", err)
	}

	accountBalances := make(map[uuid.UUID]float64)
	for _, transaction := range transactions {
		balance := accountBalances[transaction.AccountID]
		if transaction.Type == models.TopUp {
			balance += transaction.Value
		}
		if transaction.Type == models.Debit {
			balance -= transaction.Value
		}
		accountBalances[transaction.AccountID] = balance
	}

	var summaries []models.AccountSummary
	for _, account := range accounts {
		balance := accountBalances[account.ID]
		summaries = append(summaries, models.AccountSummary{
			AccountID: account.ID,
			Name:      account.Name,
			Currency:  account.Currency.String(),
			Balance:   balance,
		})
	}

	return &models.ProjectBalanceSummary{
		ProjectID: projectID,
		Summaries: summaries,
	}, nil
}

func (s *GetProjectBalanceService) getAccountBalance(accountID uuid.UUID, startDate, endDate *time.Time) (*models.ProjectBalanceSummary, error) {
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	transactions, err := s.transactionRepo.GetByAccountIDWithDateRange(accountID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for account: %w", err)
	}

	balance := 0.0
	for _, transaction := range transactions {
		if transaction.Type == models.TopUp {
			balance += transaction.Value
		}
		if transaction.Type == models.Debit {
			balance -= transaction.Value
		}
	}

	summaries := []models.AccountSummary{
		{
			AccountID: account.ID,
			Name:      account.Name,
			Currency:  account.Currency.String(),
			Balance:   balance,
		},
	}

	return &models.ProjectBalanceSummary{
		ProjectID: account.ProjectID,
		Summaries: summaries,
	}, nil
}
