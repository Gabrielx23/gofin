package components

import (
	"fmt"
	"html/template"
	"net/http"

	"gofin/internal/container"
	"gofin/web"
)

type LoginComponent struct {
	container *container.Container
	template  *template.Template
}

func NewLoginComponent(container *container.Container) (*LoginComponent, error) {
	tmpl, err := template.ParseFiles(
		web.BaseTemplate,
		web.GetTemplatePath("login.html"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse login template: %w", err)
	}
	
	return &LoginComponent{
		container: container,
		template:  tmpl,
	}, nil
}

func (c *LoginComponent) RenderLoginPage(w http.ResponseWriter, r *http.Request, projectSlug string, errorMsg string) {
	data := struct {
		Title     string
		BodyClass string
		ProjectSlug string
		ErrorMsg  string
	}{
		Title:     "Login",
		BodyClass: "login-page",
		ProjectSlug: projectSlug,
		ErrorMsg:  errorMsg,
	}

	if err := c.template.Execute(w, data); err != nil {
		http.Error(w, "Failed to render login page", http.StatusInternalServerError)
	}
}

