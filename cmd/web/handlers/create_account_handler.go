package handlers

import (
	"encoding/json"
	"net/http"

	"gofin/internal/cases/create_account"
	"gofin/pkg/money"
	webpkg "gofin/pkg/web"
)

type CreateAccountHandler struct {
	createAccountService *create_account.CreateAccountService
}

func NewCreateAccountHandler(createAccountService *create_account.CreateAccountService) *CreateAccountHandler {
	return &CreateAccountHandler{
		createAccountService: createAccountService,
	}
}

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type CreateAccountResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Error    string `json:"error,omitempty"`
}

func (h *CreateAccountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	project, _ := webpkg.GetProject(r.Context())

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := CreateAccountResponse{
			Error: "Invalid request data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	currency, err := money.ParseCurrency(req.Currency)
	if err != nil {
		response := CreateAccountResponse{
			Error: "Invalid currency. Please select PLN or EUR",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	account, err := h.createAccountService.CreateAccount(create_account.CreateAccountData{
		ProjectID: project.ID,
		Name:      req.Name,
		Currency:  currency,
	})
	if err != nil {
		response := CreateAccountResponse{
			Error: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CreateAccountResponse{
		ID:       account.ID.String(),
		Name:     account.Name,
		Currency: account.Currency.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
