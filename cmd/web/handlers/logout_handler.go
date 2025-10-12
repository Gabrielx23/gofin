package handlers

import (
	"net/http"

	"gofin/internal/container"
	webcontext "gofin/pkg/web"
	webpkg "gofin/pkg/web"
	"gofin/pkg/session"
)

type LogoutHandler struct {
	container *container.Container
}

func NewLogoutHandler(container *container.Container) *LogoutHandler {
	return &LogoutHandler{
		container: container,
	}
}

func (h *LogoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	project, _ := webcontext.GetProject(r.Context())

	session.ClearSessionCookie(w)

	webpkg.RedirectToProjectLogin(w, r, project.Slug)
}
