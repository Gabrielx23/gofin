package commands

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gofin/internal/container"
	"gofin/internal/models"
)

var (
	transactionsProjectID string
	transactionsAccountID string
	transactionsStartDate string
	transactionsEndDate   string
)

var getTransactionsCmd = &cobra.Command{
	Use:   "get-transactions",
	Short: "Get transactions with optional filtering",
	Long:  `Get transactions with optional filtering by project/account and date range. Results are ordered by transaction date descending.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := getTransactions(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	getTransactionsCmd.Flags().StringVarP(&transactionsProjectID, "project", "p", "", "Project ID (optional if account is provided)")
	getTransactionsCmd.Flags().StringVarP(&transactionsAccountID, "account", "a", "", "Account ID (optional if project is provided)")
	getTransactionsCmd.Flags().StringVarP(&transactionsStartDate, "start-date", "s", "", "Start date for filtering transactions (format: 2006-01-02 or 2006-01-02 15:04:05)")
	getTransactionsCmd.Flags().StringVarP(&transactionsEndDate, "end-date", "e", "", "End date for filtering transactions (format: 2006-01-02 or 2006-01-02 15:04:05)")
}

func getTransactions() error {
	container, err := container.NewContainerWithDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer container.DB.Close()

	query, err := parseTransactionQuery()
	if err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}

	transactions, err := container.GetTransactionsService.GetTransactions(query)
	if err != nil {
		return err
	}

	jsonOutput, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	fmt.Println(string(jsonOutput))
	return nil
}

func parseTransactionQuery() (models.TransactionQuery, error) {
	query := models.TransactionQuery{}

	if transactionsProjectID != "" {
		projectID, err := uuid.Parse(transactionsProjectID)
		if err != nil {
			return query, fmt.Errorf("invalid project ID: %w", err)
		}
		query.ProjectID = &projectID
	}

	if transactionsAccountID != "" {
		accountID, err := uuid.Parse(transactionsAccountID)
		if err != nil {
			return query, fmt.Errorf("invalid account ID: %w", err)
		}
		query.AccountID = &accountID
	}

	if transactionsStartDate != "" {
		startDate, err := parseTransactionDate(transactionsStartDate)
		if err != nil {
			return query, fmt.Errorf("invalid start date: %w", err)
		}
		query.StartDate = &startDate
	}

	if transactionsEndDate != "" {
		endDate, err := parseTransactionDate(transactionsEndDate)
		if err != nil {
			return query, fmt.Errorf("invalid end date: %w", err)
		}
		query.EndDate = &endDate
	}

	return query, nil
}
