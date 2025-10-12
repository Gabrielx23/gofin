package handlers

import (
	"net/http"

	"gofin/pkg/session"
	webpkg "gofin/pkg/web"
	"gofin/web"
	"gofin/web/components"
)

type LoginFormHandler struct {
	loginComponent *components.LoginComponent
	sessionManager *session.SessionManager
}

func NewLoginFormHandler(loginComponent *components.LoginComponent, sessionManager *session.SessionManager) *LoginFormHandler {
	return &LoginFormHandler{
		loginComponent: loginComponent,
		sessionManager: sessionManager,
	}
}

func (h *LoginFormHandler) Handle(w http.ResponseWriter, r *http.Request) {
	project, _ := webpkg.GetProject(r.Context())

	cookie, err := r.Cookie(web.SessionTokenCookie)
	if err != nil || cookie.Value == web.EmptyString {
		h.loginComponent.RenderLoginPage(w, r, project.Slug, web.EmptyString)
		return
	}

	token, valid := h.sessionManager.ValidateSessionToken(cookie.Value)
	if !valid || token.ProjectID != project.ID.String() {
		h.loginComponent.RenderLoginPage(w, r, project.Slug, web.EmptyString)
		return
	}

	webpkg.RedirectToProjectHomeWithSuccess(w, r, project.Slug, web.SuccessKeyLoginSuccessful)
}
