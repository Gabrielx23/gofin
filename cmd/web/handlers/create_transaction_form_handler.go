package handlers

import (
	"net/http"

	"gofin/internal/container"
	webcontext "gofin/pkg/web"
	"gofin/web/components"
)

type CreateTransactionFormHandler struct {
	container            *container.Container
	transactionComponent *components.TransactionCreationComponent
}

func NewCreateTransactionFormHandler(container *container.Container, transactionComponent *components.TransactionCreationComponent) *CreateTransactionFormHandler {
	return &CreateTransactionFormHandler{
		container:            container,
		transactionComponent: transactionComponent,
	}
}

func (h *CreateTransactionFormHandler) Handle(w http.ResponseWriter, r *http.Request) {
	project, _ := webcontext.GetProject(r.Context())

	accounts, err := h.container.AccountRepository.GetByProjectID(project.ID)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	h.transactionComponent.RenderCreateTransactionPage(w, r, project.Slug, accounts, "")
}
