package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gofin/internal/container"
	"gofin/internal/models"
)

var (
	transactionAccountID string
	transactionValue     string
	transactionName      string
	transactionType      string
	transactionDate      string
	transactionGroup     string
)

var createTransactionCmd = &cobra.Command{
	Use:   "create-transaction",
	Short: "Create one or more transactions",
	Long:  `Create a single transaction or multiple grouped transactions. Use --group to create grouped transactions.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := createTransaction(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	createTransactionCmd.Flags().StringVarP(&transactionAccountID, "account", "a", "", "Account ID (required for single transaction)")
	createTransactionCmd.Flags().StringVarP(&transactionValue, "value", "v", "", "Transaction value (required for single transaction)")
	createTransactionCmd.Flags().StringVarP(&transactionName, "name", "n", "", "Transaction name (required for single transaction)")
	createTransactionCmd.Flags().StringVarP(&transactionType, "type", "t", "", "Transaction type: debit or top-up (required for single transaction)")
	createTransactionCmd.Flags().StringVarP(&transactionDate, "date", "d", "", "Transaction date (optional, format: 2006-01-02 15:04:05 or 2006-01-02)")
	createTransactionCmd.Flags().StringVarP(&transactionGroup, "group", "g", "", "Group transactions together (optional, format: 'account1:value1:name1:type1:date1,account2:value2:name2:type2:date2')")
}

func createTransaction() error {
	container, err := container.NewContainerWithDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer container.DB.Close()

	if transactionGroup != "" {
		return createGroupedTransactions(container)
	}

	if transactionAccountID == "" || transactionValue == "" || transactionName == "" || transactionType == "" {
		return fmt.Errorf("for single transaction, all flags are required: --account, --value, --name, --type")
	}

	return createSingleTransaction(container)
}

func createSingleTransaction(container *container.Container) error {
	accountID, err := parseUUID(transactionAccountID)
	if err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	value, err := strconv.ParseFloat(transactionValue, 64)
	if err != nil {
		return fmt.Errorf("invalid value: %w", err)
	}

	transactionType, err := models.ParseTransactionType(transactionType)
	if err != nil {
		return fmt.Errorf("invalid transaction type: %w", err)
	}

	var parsedTransactionDate *time.Time
	if transactionDate != "" {
		parsedDate, err := parseTransactionDate(transactionDate)
		if err != nil {
			return fmt.Errorf("invalid transaction date: %w", err)
		}
		parsedTransactionDate = &parsedDate
	}

	data := models.TransactionData{
		AccountID:       accountID,
		Value:           value,
		Name:            transactionName,
		Type:            transactionType,
		TransactionDate: parsedTransactionDate,
	}

	transaction, err := container.CreateTransactionService.CreateSingleTransactionFromData(data)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Transaction created successfully!\n")
	fmt.Printf("   ID: %s\n", transaction.ID)
	fmt.Printf("   Account: %s\n", transaction.AccountID)
	fmt.Printf("   Value: %.2f\n", transaction.Value)
	fmt.Printf("   Name: %s\n", transaction.Name)
	fmt.Printf("   Type: %s\n", transaction.Type)
	fmt.Printf("   Date: %s\n", transaction.TransactionDate.Format("2006-01-02 15:04:05"))

	return nil
}

func createGroupedTransactions(container *container.Container) error {
	transactions, err := parseGroupedTransactions(transactionGroup)
	if err != nil {
		return fmt.Errorf("failed to parse grouped transactions: %w", err)
	}

	createdTransactions, err := container.CreateTransactionService.CreateGroupedTransactions(transactions)
	if err != nil {
		return err
	}

	fmt.Printf("✅ %d grouped transactions created successfully!\n", len(createdTransactions))
	fmt.Printf("   Group ID: %s\n", *createdTransactions[0].GroupID)

	for i, transaction := range createdTransactions {
		fmt.Printf("   Transaction %d:\n", i+1)
		fmt.Printf("     ID: %s\n", transaction.ID)
		fmt.Printf("     Account: %s\n", transaction.AccountID)
		fmt.Printf("     Value: %.2f\n", transaction.Value)
		fmt.Printf("     Name: %s\n", transaction.Name)
		fmt.Printf("     Type: %s\n", transaction.Type)
	}

	return nil
}

func parseGroupedTransactions(groupStr string) ([]models.TransactionData, error) {
	parts := strings.Split(groupStr, ",")
	var transactions []models.TransactionData

	for i, part := range parts {
		fields := strings.Split(part, ":")
		if len(fields) < 4 || len(fields) > 5 {
			return nil, fmt.Errorf("transaction %d must have format 'account:value:name:type[:date]'", i+1)
		}

		accountID, err := parseUUID(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid account ID in transaction %d: %w", i+1, err)
		}

		value, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in transaction %d: %w", i+1, err)
		}

		transactionType, err := models.ParseTransactionType(fields[3])
		if err != nil {
			return nil, fmt.Errorf("invalid transaction type in transaction %d: %w", i+1, err)
		}

		var transactionDate *time.Time
		if len(fields) == 5 && fields[4] != "" {
			parsedDate, err := parseTransactionDate(fields[4])
			if err != nil {
				return nil, fmt.Errorf("invalid transaction date in transaction %d: %w", i+1, err)
			}
			transactionDate = &parsedDate
		}

		transactions = append(transactions, models.TransactionData{
			AccountID:       accountID,
			Value:           value,
			Name:            fields[2],
			Type:            transactionType,
			TransactionDate: transactionDate,
		})
	}

	return transactions, nil
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func parseTransactionDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date '%s', expected formats: 2006-01-02 15:04:05, 2006-01-02, etc.", dateStr)
}
