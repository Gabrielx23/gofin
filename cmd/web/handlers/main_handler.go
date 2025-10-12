package handlers

import (
	"net/http"

	"gofin/internal/container"
	webcontext "gofin/pkg/web"
	webpkg "gofin/pkg/web"
	"gofin/web"
)

type MainHandler struct {
	container *container.Container
}

func NewMainHandler(container *container.Container) *MainHandler {
	return &MainHandler{
		container: container,
	}
}

func (h *MainHandler) Handle(w http.ResponseWriter, req *http.Request) {
	project, _ := webcontext.GetProject(req.Context())

	cookie, err := req.Cookie(web.SessionTokenCookie)

	if err != nil || cookie.Value == web.EmptyString {
		webpkg.RedirectToProjectLogin(w, req, project.Slug)
		return
	}

	webpkg.RedirectToProjectDashboard(w, req, project.Slug)
}
