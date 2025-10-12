package components

import (
	"fmt"
	"html/template"
	"net/http"

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

func (c *DashboardComponent) RenderDashboard(w http.ResponseWriter, r *http.Request, project *models.Project, access *models.Access, projectSlug, successKey string) {
	successMessage := c.getSuccessMessage(successKey)

	balanceData, err := c.container.GetProjectBalanceService.GetProjectBalances(project.ID)
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
		formattedBalance := c.container.GetProjectBalanceService.FormatBalance(
			balance.Balance,
			balance.Currency,
		)

		displayBalances = append(displayBalances, AccountBalanceDisplay{
			Name:       balance.Name,
			Balance:    formattedBalance,
			IsPositive: balance.Balance >= 0,
		})
	}

	return displayBalances
}
