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
	CurrencyTotals  []models.CurrencyTotal  `json:"currency_totals"`
}

func (s *GetProjectBalanceService) GetProjectBalances(projectID uuid.UUID, year int, month int) (*ProjectBalanceData, error) {
	accounts, err := s.accountRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project accounts: %w", err)
	}

	var accountBalances []models.AccountSummary
	for _, account := range accounts {
		balance, err := s.calculateAccountBalance(account.ID, year, month)
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

	currencyTotals := s.calculateCurrencyTotals(accountBalances)

	return &ProjectBalanceData{
		ProjectID:       projectID,
		AccountBalances: accountBalances,
		CurrencyTotals:  currencyTotals,
	}, nil
}

func (s *GetProjectBalanceService) calculateAccountBalance(accountID uuid.UUID, year int, month int) (float64, error) {
	startDate, endDate := s.calculateDateRange(year, month)
	
	query := models.TransactionQuery{
		AccountID:                 &accountID,
		StartDate:                 startDate,
		EndDate:                   endDate,
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

func (s *GetProjectBalanceService) calculateCurrencyTotals(balances []models.AccountSummary) []models.CurrencyTotal {
	currencyTotals := make(map[string]float64)

	for _, balance := range balances {
		currencyTotals[balance.Currency] += balance.Balance
	}

	var currencyTotalsList []models.CurrencyTotal
	for currency, total := range currencyTotals {
		currencyTotalsList = append(currencyTotalsList, models.CurrencyTotal{
			Currency:   currency,
			Balance:    total,
			IsPositive: total >= 0,
		})
	}

	return currencyTotalsList
}

func (s *GetProjectBalanceService) calculateDateRange(year int, month int) (*time.Time, *time.Time) {
	var startDate, endDate time.Time

	if month > 0 {
		startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	}

	return &startDate, &endDate
}

func (s *GetProjectBalanceService) FormatBalance(balance float64, currency string) string {
	return fmt.Sprintf("%.2f %s", balance, currency)
}
