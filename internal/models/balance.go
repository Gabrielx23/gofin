package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccountSummary struct {
	AccountID uuid.UUID `json:"account_id"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
}

type ProjectBalanceSummary struct {
	ProjectID uuid.UUID        `json:"project_id"`
	Summaries []AccountSummary `json:"summaries"`
}

type BalanceQuery struct {
	ProjectID *uuid.UUID
	AccountID *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
}

func (q *BalanceQuery) Validate() error {
	if q.ProjectID == nil && q.AccountID == nil {
		return fmt.Errorf("either project_id or account_id must be provided")
	}

	if q.ProjectID != nil && q.AccountID != nil {
		return fmt.Errorf("cannot specify both project_id and account_id")
	}

	if q.StartDate != nil && q.EndDate != nil {
		if q.EndDate.Before(*q.StartDate) {
			return fmt.Errorf("end_date cannot be before start_date")
		}
	}

	return nil
}
