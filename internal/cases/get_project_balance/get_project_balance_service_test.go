package get_project_balance

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/pkg/money"
)

func TestGetProjectBalanceService_GetProjectBalancesFromTransactions(t *testing.T) {
	accountRepo := database.NewAccountInMemoryRepository()
	service := NewGetProjectBalanceService(accountRepo)

	projectID := uuid.New()
	account1 := models.NewAccount(projectID, "Savings", money.PLN)
	account2 := models.NewAccount(projectID, "Checking", money.EUR)

	accountRepo.Create(account1)
	accountRepo.Create(account2)

	now := time.Now()
	transactionDate1 := now.Add(-2 * time.Hour)
	transaction1 := models.NewTransaction(models.TransactionData{
		AccountID:       account1.ID,
		Value:           100.0,
		Name:            "Initial deposit",
		TransactionDate: &transactionDate1,
		Type:            models.TopUp,
	}, uuid.New())

	transactionDate2 := now.Add(-1 * time.Hour)
	transaction2 := models.NewTransaction(models.TransactionData{
		AccountID:       account1.ID,
		Value:           30.0,
		Name:            "Withdrawal",
		TransactionDate: &transactionDate2,
		Type:            models.Debit,
	}, uuid.New())

	transactionDate3 := now.Add(-30 * time.Minute)
	transaction3 := models.NewTransaction(models.TransactionData{
		AccountID:       account2.ID,
		Value:           200.0,
		Name:            "Salary",
		TransactionDate: &transactionDate3,
		Type:            models.TopUp,
	}, uuid.New())

	transactions := []*models.Transaction{transaction1, transaction2, transaction3}

	result, err := service.GetProjectBalancesFromTransactions(projectID, transactions)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ProjectID != projectID {
		t.Errorf("Expected project ID %s, got %s", projectID, result.ProjectID)
	}

	if len(result.AccountBalances) != 2 {
		t.Fatalf("Expected 2 account balances, got %d", len(result.AccountBalances))
	}

	balances := make(map[string]models.AccountSummary)
	for _, balance := range result.AccountBalances {
		balances[balance.Name] = balance
	}

	savingsBalance := balances["Savings"]
	if savingsBalance.Balance != 70.0 {
		t.Errorf("Expected savings balance 70.0, got %.2f", savingsBalance.Balance)
	}
	if savingsBalance.Currency != "PLN" {
		t.Errorf("Expected savings currency PLN, got %s", savingsBalance.Currency)
	}

	checkingBalance := balances["Checking"]
	if checkingBalance.Balance != 200.0 {
		t.Errorf("Expected checking balance 200.0, got %.2f", checkingBalance.Balance)
	}
	if checkingBalance.Currency != "EUR" {
		t.Errorf("Expected checking currency EUR, got %s", checkingBalance.Currency)
	}

	if len(result.CurrencyTotals) != 2 {
		t.Fatalf("Expected 2 currency totals, got %d", len(result.CurrencyTotals))
	}

	currencyTotals := make(map[string]models.CurrencyTotal)
	for _, total := range result.CurrencyTotals {
		currencyTotals[total.Currency] = total
	}

	plnTotal := currencyTotals["PLN"]
	if plnTotal.Balance != 70.0 {
		t.Errorf("Expected PLN total 70.0, got %.2f", plnTotal.Balance)
	}
	if !plnTotal.IsPositive {
		t.Errorf("Expected PLN total to be positive")
	}

	eurTotal := currencyTotals["EUR"]
	if eurTotal.Balance != 200.0 {
		t.Errorf("Expected EUR total 200.0, got %.2f", eurTotal.Balance)
	}
	if !eurTotal.IsPositive {
		t.Errorf("Expected EUR total to be positive")
	}
}

func TestGetProjectBalanceService_GetProjectBalancesFromTransactions_EmptyProject(t *testing.T) {
	accountRepo := database.NewAccountInMemoryRepository()
	service := NewGetProjectBalanceService(accountRepo)

	projectID := uuid.New()
	transactions := []*models.Transaction{}

	result, err := service.GetProjectBalancesFromTransactions(projectID, transactions)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ProjectID != projectID {
		t.Errorf("Expected project ID %s, got %s", projectID, result.ProjectID)
	}

	if len(result.AccountBalances) != 0 {
		t.Errorf("Expected 0 account balances, got %d", len(result.AccountBalances))
	}

	if len(result.CurrencyTotals) != 0 {
		t.Errorf("Expected 0 currency totals, got %d", len(result.CurrencyTotals))
	}
}
