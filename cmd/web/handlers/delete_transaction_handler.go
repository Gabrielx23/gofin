package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"gofin/internal/container"
	webcontext "gofin/pkg/web"
	"gofin/web"
)

type DeleteTransactionHandler struct {
	container *container.Container
}

func NewDeleteTransactionHandler(container *container.Container) *DeleteTransactionHandler {
	return &DeleteTransactionHandler{
		container: container,
	}
}

func (h *DeleteTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	project, _ := webcontext.GetProject(r.Context())
	access, _ := webcontext.GetAccess(r.Context())

	if access.ReadOnly {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	transactionIDStr := r.URL.Query().Get("id")
	if transactionIDStr == "" {
		http.Error(w, "Transaction ID is required", http.StatusBadRequest)
		return
	}

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err = h.container.DeleteTransactionService.DeleteTransaction(transactionID)
	if err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+project.Slug+web.RouteDashboard+"?success="+web.SuccessKeyTransactionDeleted, http.StatusSeeOther)
}
