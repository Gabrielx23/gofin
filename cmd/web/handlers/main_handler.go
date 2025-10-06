package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"gofin/cmd/web/middleware"
	"gofin/web"
)

type MainHandler struct{}

func NewMainHandler() *MainHandler {
	return &MainHandler{}
}

func (h *MainHandler) Handle(w http.ResponseWriter, req *http.Request) {
	projectID, ok := middleware.GetProjectIDFromContext(req.Context())

	if !ok || projectID == uuid.Nil {
		http.NotFound(w, req)
		return
	}

	projectSlug, ok := middleware.GetProjectSlugFromContext(req.Context())
	if !ok {
		http.Error(w, web.ProjectSlugNotFoundError, http.StatusInternalServerError)
		return
	}

	cookie, err := req.Cookie(web.SessionTokenCookie)

	if err != nil || cookie.Value == web.EmptyString {
		http.Redirect(w, req, "/"+projectSlug+web.RouteLogin, http.StatusSeeOther)
		return
	}

	http.Redirect(w, req, "/"+projectSlug+web.RouteDashboard, http.StatusSeeOther)
}
