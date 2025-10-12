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
	
	data := struct {
		Title        string
		BodyClass    string
		ProjectID    string
		ProjectSlug  string
		ProjectName  string
		AccessName   string
		ReadOnly     bool
		SuccessMsg   string
	}{
		Title:        project.Name,
		BodyClass:    dashboardBodyClass,
		ProjectID:    project.ID.String(),
		ProjectSlug:  projectSlug,
		ProjectName:  project.Name,
		AccessName:   access.Name,
		ReadOnly:     access.ReadOnly,
		SuccessMsg:   successMessage,
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
