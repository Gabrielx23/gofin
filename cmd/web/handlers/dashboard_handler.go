package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gofin/internal/container"
	webcontext "gofin/pkg/web"
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
	year, month := h.parseAndValidateFilterParams(r)

	transactions, err := h.container.GetProjectTransactionsService.GetProjectTransactions(project.ID, year, month)
	if err != nil {
		http.Error(w, "Failed to get project transactions", http.StatusInternalServerError)
		return
	}

	balanceData, err := h.container.GetProjectBalanceService.GetProjectBalancesFromTransactions(project.ID, transactions)
	if err != nil {
		http.Error(w, "Failed to get project balances", http.StatusInternalServerError)
		return
	}

	h.dashboardComponent.RenderDashboard(w, r, project, access, project.Slug, successMsg, year, month, transactions, balanceData)
}

func (h *DashboardHandler) parseAndValidateFilterParams(r *http.Request) (int, int) {
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	year := h.parseYear(yearStr)
	month := h.parseMonth(monthStr)

	return year, month
}

func (h *DashboardHandler) parseYear(yearStr string) int {
	defaultYear := time.Now().Year()
	minYear := defaultYear - 10
	maxYear := defaultYear + 10

	if yearStr == "" {
		return defaultYear
	}

	parsedYear, err := strconv.Atoi(yearStr)
	if err != nil {
		return defaultYear
	}

	if parsedYear < minYear || parsedYear > maxYear {
		return defaultYear
	}

	return parsedYear
}

func (h *DashboardHandler) parseMonth(monthStr string) int {
	defaultMonth := int(time.Now().Month())

	if monthStr == "" {
		return defaultMonth
	}

	parsedMonth, err := strconv.Atoi(monthStr)
	if err != nil {
		return defaultMonth
	}

	if parsedMonth < 1 || parsedMonth > 12 {
		return defaultMonth
	}

	return parsedMonth
}
