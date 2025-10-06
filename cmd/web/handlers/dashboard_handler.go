package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"gofin/internal/container"
	"gofin/web"
	"gofin/web/components"
	"gofin/cmd/web/middleware"
)

type DashboardHandler struct {
	container *container.Container
	dashboardComponent *components.DashboardComponent
}

func NewDashboardHandler(container *container.Container, dashboardComponent *components.DashboardComponent) *DashboardHandler {
	return &DashboardHandler{
		container: container,
		dashboardComponent: dashboardComponent,
	}
}

func (h *DashboardHandler) Handle(w http.ResponseWriter, r *http.Request) {
	projectID, ok := middleware.GetProjectIDFromContext(r.Context())
	if !ok {
		http.Error(w, web.ProjectIDNotFoundError, http.StatusInternalServerError)
		return
	}

	projectSlug, ok := middleware.GetProjectSlugFromContext(r.Context())
	if !ok {
		http.Error(w, web.ProjectSlugNotFoundError, http.StatusInternalServerError)
		return
	}

	accessID, ok := middleware.GetAccessIDFromContext(r.Context())
	if !ok {
		http.Error(w, web.AccessIDNotFoundError, http.StatusInternalServerError)
		return
	}

	project, err := h.container.ProjectRepository.GetByID(projectID)
	if err != nil {
		http.Error(w, web.ProjectNotFoundError, http.StatusNotFound)
		return
	}

	access, err := h.container.AccessRepository.GetByID(uuid.MustParse(accessID))
	if err != nil {
		http.Error(w, web.AccessNotFoundError, http.StatusNotFound)
		return
	}

	h.dashboardComponent.RenderDashboard(w, r, project, access, projectSlug)
}
