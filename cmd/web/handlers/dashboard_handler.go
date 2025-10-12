package handlers

import (
	"net/http"

	webcontext "gofin/pkg/web"
	"gofin/internal/container"
	"gofin/web"
	"gofin/web/components"
)

type DashboardHandler struct {
	container          *container.Container
	dashboardComponent *components.DashboardComponent
}

func NewDashboardHandler(container *container.Container, dashboardComponent *components.DashboardComponent) *DashboardHandler {
	return &DashboardHandler{
		container:          container,
		dashboardComponent: dashboardComponent,
	}
}

func (h *DashboardHandler) Handle(w http.ResponseWriter, r *http.Request) {
	project, _ := webcontext.GetProject(r.Context())
	access, _ := webcontext.GetAccess(r.Context())

	successMsg := r.URL.Query().Get(web.SuccessQueryParam)
	h.dashboardComponent.RenderDashboard(w, r, project, access, project.Slug, successMsg)
}
