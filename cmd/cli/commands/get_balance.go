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
	balanceProjectID string
	balanceAccountID string
	balanceStartDate string
	balanceEndDate   string
)

var getBalanceCmd = &cobra.Command{
	Use:   "get-balance",
	Short: "Get balance summary for a project or account",
	Long:  `Get balance summary for a project (all accounts) or specific account with optional date filtering.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := getBalance(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	getBalanceCmd.Flags().StringVarP(&balanceProjectID, "project", "p", "", "Project ID (optional if account is provided)")
	getBalanceCmd.Flags().StringVarP(&balanceAccountID, "account", "a", "", "Account ID (optional if project is provided)")
	getBalanceCmd.Flags().StringVarP(&balanceStartDate, "start-date", "s", "", "Start date for filtering transactions (format: 2006-01-02 or 2006-01-02 15:04:05)")
	getBalanceCmd.Flags().StringVarP(&balanceEndDate, "end-date", "e", "", "End date for filtering transactions (format: 2006-01-02 or 2006-01-02 15:04:05)")
}

func getBalance() error {
	container, err := container.NewContainerWithDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer container.DB.Close()

	query, err := parseBalanceQuery()
	if err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}

	summary, err := container.GetProjectBalanceService.GetProjectBalance(query)
	if err != nil {
		return err
	}

	jsonOutput, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(jsonOutput))
	return nil
}

func parseBalanceQuery() (models.BalanceQuery, error) {
	query := models.BalanceQuery{}

	if balanceProjectID != "" {
		projectID, err := uuid.Parse(balanceProjectID)
		if err != nil {
			return query, fmt.Errorf("invalid project ID: %w", err)
		}
		query.ProjectID = &projectID
	}

	if balanceAccountID != "" {
		accountID, err := uuid.Parse(balanceAccountID)
		if err != nil {
			return query, fmt.Errorf("invalid account ID: %w", err)
		}
		query.AccountID = &accountID
	}

	if balanceStartDate != "" {
		startDate, err := parseTransactionDate(balanceStartDate)
		if err != nil {
			return query, fmt.Errorf("invalid start date: %w", err)
		}
		query.StartDate = &startDate
	}

	if balanceEndDate != "" {
		endDate, err := parseTransactionDate(balanceEndDate)
		if err != nil {
			return query, fmt.Errorf("invalid end date: %w", err)
		}
		query.EndDate = &endDate
	}

	return query, nil
}
