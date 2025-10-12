package components

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"gofin/internal/container"
	"gofin/internal/models"
	"gofin/pkg/config"
	webhelpers "gofin/pkg/web"
	"gofin/web"
)

const (
	transactionTemplateFile = "create_transaction.html"
	pageTitle              = "Create Transaction"
	bodyClass              = "dashboard-page"
	templateError          = "Failed to render transaction page"
)

type TransactionCreationComponent struct {
	container *container.Container
	template  *template.Template
}

func NewTransactionCreationComponent(container *container.Container) (*TransactionCreationComponent, error) {
	tmpl, err := template.ParseFiles(
		web.BaseTemplate,
		webhelpers.GetTemplatePath(transactionTemplateFile),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction template: %w", err)
	}

	return &TransactionCreationComponent{
		container: container,
		template:  tmpl,
	}, nil
}

func (c *TransactionCreationComponent) RenderCreateTransactionPage(w http.ResponseWriter, r *http.Request, projectSlug string, accounts []*models.Account, errorMsg string) {
	data := struct {
		Title       string
		BodyClass   string
		ProjectSlug string
		Accounts    []*models.Account
		DefaultDate string
		ErrorMsg    string
	}{
		Title:       pageTitle,
		BodyClass:   bodyClass,
		ProjectSlug: projectSlug,
		Accounts:    accounts,
		DefaultDate: time.Now().Format(config.DateTimeFormat),
		ErrorMsg:    errorMsg,
	}

	if err := c.template.Execute(w, data); err != nil {
		http.Error(w, templateError, http.StatusInternalServerError)
	}
}
