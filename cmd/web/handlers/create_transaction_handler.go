package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gofin/internal/cases/create_transaction"
	"gofin/internal/container"
	"gofin/internal/models"
	"gofin/pkg/config"
	webpkg "gofin/pkg/web"
	"gofin/web"
	"gofin/web/components"
)

type CreateTransactionHandler struct {
	container           *container.Container
	transactionComponent *components.TransactionCreationComponent
	createTransactionSvc *create_transaction.CreateTransactionService
}

func NewCreateTransactionHandler(container *container.Container, transactionComponent *components.TransactionCreationComponent, createTransactionSvc *create_transaction.CreateTransactionService) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		container:            container,
		transactionComponent: transactionComponent,
		createTransactionSvc: createTransactionSvc,
	}
}

const (
	formParseError = "Failed to parse form data"
	noTransactionsError = "At least one transaction is required"
	invalidValueError = "Invalid value for group %d"
	invalidAccountError = "Invalid account for group %d"
	invalidTypeError = "Invalid type for group %d"
	invalidDateError = "Invalid date for group %d"
	createTransactionError = "Failed to create transactions: %v"
)

type TransactionGroupData struct {
	Name      string
	Value     float64
	Type      string
	AccountID string
	Date      time.Time
}

func (h *CreateTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	project, _ := webpkg.GetProject(r.Context())

	accounts, err := h.container.AccountRepository.GetByProjectID(project.ID)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

		if err := r.ParseForm(); err != nil {
		h.renderCreateTransactionForm(w, r, accounts, project.Slug, formParseError)
		return
	}

	var groups []TransactionGroupData
	
	for key, values := range r.Form {
		if len(values) == 0 || values[0] == "" {
			continue
		}
		
		if !strings.HasPrefix(key, "groups[") || !strings.Contains(key, "].value") {
			continue
		}
		
		indexStr := strings.TrimPrefix(key, "groups[")
		indexStr = strings.TrimSuffix(indexStr, "].value")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}
		
		valueStr := r.FormValue(fmt.Sprintf("groups[%d].value", index))
		if valueStr == "" {
			continue
		}
		
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			h.renderCreateTransactionForm(w, r, accounts, project.Slug, fmt.Sprintf(invalidValueError, index+1))
			return
		}
		
		accountIDStr := r.FormValue(fmt.Sprintf("groups[%d].account_id", index))
		if accountIDStr == "" {
			h.renderCreateTransactionForm(w, r, accounts, project.Slug, fmt.Sprintf("Account is required for group %d", index+1))
			return
		}
		
		_, err = uuid.Parse(accountIDStr)
		if err != nil {
			h.renderCreateTransactionForm(w, r, accounts, project.Slug, fmt.Sprintf(invalidAccountError, index+1))
			return
		}
		
		typeStr := r.FormValue(fmt.Sprintf("groups[%d].type", index))
		if typeStr == "" {
			h.renderCreateTransactionForm(w, r, accounts, project.Slug, fmt.Sprintf("Type is required for group %d", index+1))
			return
		}
		
		_, err = models.ParseTransactionType(typeStr)
		if err != nil {
			h.renderCreateTransactionForm(w, r, accounts, project.Slug, fmt.Sprintf(invalidTypeError, index+1))
			return
		}
		
		date := time.Now()
		dateStr := r.FormValue(fmt.Sprintf("groups[%d].date", index))
		if dateStr != "" {
			date, err = time.Parse(config.DateTimeFormat, dateStr)
			if err != nil {
				h.renderCreateTransactionForm(w, r, accounts, project.Slug, fmt.Sprintf(invalidDateError, index+1))
				return
			}
		}
		
		groups = append(groups, TransactionGroupData{
			Name:      r.FormValue(fmt.Sprintf("groups[%d].name", index)),
			Value:     value,
			Type:      typeStr,
			AccountID: accountIDStr,
			Date:      date,
		})
	}
	
	if len(groups) == 0 {
		h.renderCreateTransactionForm(w, r, accounts, project.Slug, noTransactionsError)
		return
	}
	
	var transactionData []models.TransactionData
	for _, group := range groups {
		accountID, _ := uuid.Parse(group.AccountID)
		transactionType, _ := models.ParseTransactionType(group.Type)
		
		transactionData = append(transactionData, models.TransactionData{
			AccountID:       accountID,
			Value:           group.Value,
			Name:            group.Name,
			Type:            transactionType,
			TransactionDate: &group.Date,
		})
	}
	
	_, err = h.createTransactionSvc.CreateGroupedTransactions(project.ID, transactionData)
	if err != nil {
		h.renderCreateTransactionForm(w, r, accounts, project.Slug, fmt.Sprintf(createTransactionError, err))
		return
	}
	
	webpkg.RedirectWithSuccess(w, r, "/"+project.Slug+web.RouteDashboard, web.SuccessKeyTransactionsCreated)
}

func (h *CreateTransactionHandler) renderCreateTransactionForm(w http.ResponseWriter, r *http.Request, accounts []*models.Account, projectSlug, errorMsg string) {
	h.transactionComponent.RenderCreateTransactionPage(w, r, projectSlug, accounts, errorMsg)
}
