package get_project_balance

import (
	"fmt"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type GetProjectBalanceService struct {
	accountRepo models.AccountRepository
}

func NewGetProjectBalanceService(accountRepo models.AccountRepository) *GetProjectBalanceService {
	return &GetProjectBalanceService{
		accountRepo: accountRepo,
	}
}

type ProjectBalanceData struct {
	ProjectID       uuid.UUID               `json:"project_id"`
	AccountBalances []models.AccountSummary `json:"account_balances"`
	CurrencyTotals  []models.CurrencyTotal  `json:"currency_totals"`
}

func (s *GetProjectBalanceService) GetProjectBalancesFromTransactions(projectID uuid.UUID, transactions []*models.Transaction) (*ProjectBalanceData, error) {
	accounts, err := s.accountRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project accounts: %w", err)
	}

	accountBalances := s.calculateAccountBalancesFromTransactions(accounts, transactions)
	currencyTotals := s.calculateCurrencyTotals(accountBalances)

	return &ProjectBalanceData{
		ProjectID:       projectID,
		AccountBalances: accountBalances,
		CurrencyTotals:  currencyTotals,
	}, nil
}

func (s *GetProjectBalanceService) calculateAccountBalancesFromTransactions(accounts []*models.Account, transactions []*models.Transaction) []models.AccountSummary {
	accountBalances := make(map[uuid.UUID]float64)

	for _, transaction := range transactions {
		if transaction.Type == models.Debit {
			accountBalances[transaction.AccountID] -= transaction.Value
		} else {
			accountBalances[transaction.AccountID] += transaction.Value
		}
	}

	var accountSummaries []models.AccountSummary
	for _, account := range accounts {
		balance := accountBalances[account.ID]
		accountSummaries = append(accountSummaries, models.AccountSummary{
			AccountID: account.ID,
			Name:      account.Name,
			Currency:  account.Currency.String(),
			Balance:   balance,
		})
	}

	return accountSummaries
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
