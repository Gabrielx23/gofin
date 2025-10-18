package components

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"gofin/internal/cases/get_project_balance"
	"gofin/internal/container"
	"gofin/internal/models"
	webhelpers "gofin/pkg/web"
	"gofin/web"
)

const (
	dashboardTemplateFile = "dashboard.html"
	dashboardBodyClass    = "dashboard-page"
	dashboardTitle        = "Dashboard"
)

type AccountBalanceDisplay struct {
	Name       string
	Balance    string
	IsPositive bool
}

type CurrencyTotalDisplay struct {
	Currency   string
	Balance    string
	IsPositive bool
}

type TransactionDisplay struct {
	ID              string
	AccountName     string
	AccountCurrency string
	Value           string
	FormattedValue  string
	Name            string
	TransactionDate string
	Type            string
	IsDebit         bool
	IsTopUp         bool
}

type DashboardComponent struct {
	container *container.Container
	template  *template.Template
}

func NewDashboardComponent(container *container.Container) (*DashboardComponent, error) {
	tmpl, err := template.ParseFiles(
		web.BaseTemplate,
		webhelpers.GetTemplatePath(dashboardTemplateFile),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dashboard template: %w", err)
	}

	return &DashboardComponent{
		container: container,
		template:  tmpl,
	}, nil
}

func (c *DashboardComponent) RenderDashboard(w http.ResponseWriter, r *http.Request, project *models.Project, access *models.Access, projectSlug, successKey string, year, month int, transactions []*models.Transaction, balanceData *get_project_balance.ProjectBalanceData) {
	successMessage := c.getSuccessMessage(successKey)

	data := struct {
		Title           string
		BodyClass       string
		ProjectID       string
		ProjectSlug     string
		ProjectName     string
		AccessName      string
		ReadOnly        bool
		SuccessMsg      string
		AccountBalances []AccountBalanceDisplay
		CurrencyTotals  []CurrencyTotalDisplay
		Transactions    []TransactionDisplay
		SelectedYear    int
		SelectedMonth   int
		Years           []int
		Months          []int
	}{
		Title:           project.Name,
		BodyClass:       dashboardBodyClass,
		ProjectID:       project.ID.String(),
		ProjectSlug:     projectSlug,
		ProjectName:     project.Name,
		AccessName:      access.Name,
		ReadOnly:        access.ReadOnly,
		SuccessMsg:      successMessage,
		AccountBalances: c.formatAccountBalances(balanceData.AccountBalances),
		CurrencyTotals:  c.formatCurrencyTotals(balanceData.CurrencyTotals),
		Transactions:    c.formatTransactions(transactions),
		SelectedYear:    year,
		SelectedMonth:   month,
		Years:           c.getYears(),
		Months:          c.getMonths(),
	}

	if err := c.template.Execute(w, data); err != nil {
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
	}
}

func (c *DashboardComponent) getSuccessMessage(successKey string) string {
	successMessages := map[string]string{
		web.SuccessKeyTransactionsCreated: web.SuccessTransactionsCreated,
		web.SuccessKeyLoginSuccessful:     web.SuccessLoginSuccessful,
	}

	if message, exists := successMessages[successKey]; exists {
		return message
	}
	return ""
}

func (c *DashboardComponent) formatAccountBalances(balances []models.AccountSummary) []AccountBalanceDisplay {
	var displayBalances []AccountBalanceDisplay

	for _, balance := range balances {
		formattedBalance := c.formatBalance(balance.Balance, balance.Currency)

		displayBalances = append(displayBalances, AccountBalanceDisplay{
			Name:       balance.Name,
			Balance:    formattedBalance,
			IsPositive: balance.Balance >= 0,
		})
	}

	return displayBalances
}

func (c *DashboardComponent) formatBalance(balance float64, currency string) string {
	return fmt.Sprintf("%.2f %s", balance, currency)
}

func (c *DashboardComponent) formatCurrencyTotals(currencyTotals []models.CurrencyTotal) []CurrencyTotalDisplay {
	var displayTotals []CurrencyTotalDisplay

	for _, total := range currencyTotals {
		formattedBalance := c.formatBalance(total.Balance, total.Currency)

		displayTotals = append(displayTotals, CurrencyTotalDisplay{
			Currency:   total.Currency,
			Balance:    formattedBalance,
			IsPositive: total.IsPositive,
		})
	}

	return displayTotals
}

func (c *DashboardComponent) getYears() []int {
	currentYear := time.Now().Year()
	years := make([]int, 0, 11)
	for year := currentYear - 5; year <= currentYear+5; year++ {
		years = append(years, year)
	}
	return years
}

func (c *DashboardComponent) getMonths() []int {
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
}

func (c *DashboardComponent) formatTransactions(transactions []*models.Transaction) []TransactionDisplay {
	var displayTransactions []TransactionDisplay

	for _, transaction := range transactions {
		account, err := c.container.AccountRepository.GetByID(transaction.AccountID)
		if err != nil {
			continue
		}

		formattedValue := c.formatTransactionValue(transaction.Value, account.Currency.String(), transaction.Type)
		formattedDate := transaction.TransactionDate.Format("2006-01-02")

		displayTransactions = append(displayTransactions, TransactionDisplay{
			ID:              transaction.ID.String(),
			AccountName:     account.Name,
			AccountCurrency: account.Currency.String(),
			Value:           fmt.Sprintf("%.2f", transaction.Value),
			FormattedValue:  formattedValue,
			Name:            transaction.Name,
			TransactionDate: formattedDate,
			Type:            transaction.Type.String(),
			IsDebit:         transaction.Type == models.Debit,
			IsTopUp:         transaction.Type == models.TopUp,
		})
	}

	return displayTransactions
}

func (c *DashboardComponent) formatTransactionValue(value float64, currency string, transactionType models.TransactionType) string {
	if transactionType == models.Debit {
		return fmt.Sprintf("-%.2f %s", value, currency)
	}
	return fmt.Sprintf("+%.2f %s", value, currency)
}
