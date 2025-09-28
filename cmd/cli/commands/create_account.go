package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"gofin/internal/container"
	"gofin/pkg/money"
)

var (
	accountProjectSlug string
	accountName        string
	accountCurrency    string
)

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account for a project",
	Long:  `Create a new account with a name and currency for a project.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := createAccount(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	createAccountCmd.Flags().StringVarP(&accountProjectSlug, "project", "p", "", "Project slug (required)")
	createAccountCmd.Flags().StringVarP(&accountName, "name", "n", "", "Account name (required)")
	createAccountCmd.Flags().StringVarP(&accountCurrency, "currency", "c", "USD", "Currency code (default: USD)")
	createAccountCmd.MarkFlagRequired("project")
	createAccountCmd.MarkFlagRequired("name")
}

func createAccount() error {
	if accountProjectSlug == "" {
		return fmt.Errorf("project slug is required")
	}

	if accountName == "" {
		return fmt.Errorf("name is required")
	}

	_, err := money.ParseCurrency(accountCurrency)
	if err != nil {
		return fmt.Errorf("invalid currency '%s': %w", accountCurrency, err)
	}

	container, err := container.NewContainerWithDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer container.DB.Close()

	account, err := container.CreateAccountService.CreateAccount(accountProjectSlug, accountName, accountCurrency)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Account created successfully!\n")
	fmt.Printf("   Project: %s\n", accountProjectSlug)
	fmt.Printf("   Name: %s\n", account.Name)
	fmt.Printf("   Currency: %s (%s)\n", account.Currency.String(), money.GetCurrencySymbol(account.Currency))
	fmt.Printf("   ID: %s\n", account.ID)

	return nil
}
