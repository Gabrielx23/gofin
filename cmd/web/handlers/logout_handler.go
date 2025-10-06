package handlers

import (
	"net/http"

	"gofin/cmd/web/middleware"
	"gofin/pkg/session"
	"gofin/web"
)

type LogoutHandler struct{}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{}
}

func (h *LogoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	projectSlug, ok := middleware.GetProjectSlugFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session.ClearSessionCookie(w)

	http.Redirect(w, r, "/"+projectSlug+web.RouteLogin, http.StatusSeeOther)
}
