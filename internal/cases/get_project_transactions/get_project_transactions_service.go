package get_project_transactions

import (
	"time"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type GetProjectTransactionsService struct {
	transactionRepo models.TransactionRepository
}

func NewGetProjectTransactionsService(transactionRepo models.TransactionRepository) *GetProjectTransactionsService {
	return &GetProjectTransactionsService{
		transactionRepo: transactionRepo,
	}
}

func (s *GetProjectTransactionsService) GetProjectTransactions(projectID uuid.UUID, year int, month int) ([]*models.Transaction, error) {
	startDate, endDate := s.calculateDateRange(year, month)

	query := models.TransactionQuery{
		ProjectID: &projectID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	return s.transactionRepo.GetTransactionsWithFilters(query)
}

func (s *GetProjectTransactionsService) calculateDateRange(year int, month int) (*time.Time, *time.Time) {
	var startDate, endDate time.Time

	if month == 0 {
		startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)
	} else {
		startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		if month == 12 {
			endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
		} else {
			endDate = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
		}
	}

	return &startDate, &endDate
}
