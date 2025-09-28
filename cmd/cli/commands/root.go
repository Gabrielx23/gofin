package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gofin",
	Short: "A simple financial management CLI",
	Long:  `Gofin is a CLI tool for managing financial projects and transactions.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createProjectCmd)
	rootCmd.AddCommand(createAccessCmd)
	rootCmd.AddCommand(createAccountCmd)
	rootCmd.AddCommand(createTransactionCmd)
	rootCmd.AddCommand(getBalanceCmd)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
