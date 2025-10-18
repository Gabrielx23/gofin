package components

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

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

func (c *DashboardComponent) RenderDashboard(w http.ResponseWriter, r *http.Request, project *models.Project, access *models.Access, projectSlug, successKey string, year, month int) {
	successMessage := c.getSuccessMessage(successKey)

	balanceData, err := c.container.GetProjectBalanceService.GetProjectBalances(project.ID, year, month)
	if err != nil {
		http.Error(w, "Failed to get project balances", http.StatusInternalServerError)
		return
	}

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

