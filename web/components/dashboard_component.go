package components

import (
	"fmt"
	"html/template"
	"net/http"

	"gofin/internal/container"
	"gofin/internal/models"
	"gofin/web"
)

type DashboardComponent struct {
	container *container.Container
	template  *template.Template
}

func NewDashboardComponent(container *container.Container) (*DashboardComponent, error) {
	tmpl, err := template.ParseFiles(
		web.BaseTemplate,
		web.GetTemplatePath("dashboard.html"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dashboard template: %w", err)
	}
	
	return &DashboardComponent{
		container: container,
		template:  tmpl,
	}, nil
}

func (c *DashboardComponent) RenderDashboard(w http.ResponseWriter, r *http.Request, project *models.Project, access *models.Access, projectSlug string) {
	data := struct {
		Title       string
		BodyClass   string
		ProjectID   string
		ProjectSlug string
		ProjectName string
		AccessName  string
		ReadOnly    bool
	}{
		Title:       project.Name,
		BodyClass:   "dashboard-page",
		ProjectID:   project.ID.String(),
		ProjectSlug: projectSlug,
		ProjectName: project.Name,
		AccessName:  access.Name,
		ReadOnly:    access.ReadOnly,
	}

	if err := c.template.Execute(w, data); err != nil {
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
	}
}
