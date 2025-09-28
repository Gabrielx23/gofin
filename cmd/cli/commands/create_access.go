package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"gofin/internal/container"
)

var (
	accessProjectSlug string
	accessName        string
	accessReadonly    bool
)

var createAccessCmd = &cobra.Command{
	Use:   "create-access",
	Short: "Create a new access credential for a project",
	Long:  `Create a new access credential with auto-generated UID and PIN for a project.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := createAccess(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	createAccessCmd.Flags().StringVarP(&accessProjectSlug, "project", "p", "", "Project slug (required)")
	createAccessCmd.Flags().StringVarP(&accessName, "name", "n", "", "Access name (required)")
	createAccessCmd.Flags().BoolVarP(&accessReadonly, "readonly", "r", false, "Create read-only access")
	createAccessCmd.MarkFlagRequired("project")
	createAccessCmd.MarkFlagRequired("name")
}

func createAccess() error {
	if accessProjectSlug == "" {
		return fmt.Errorf("project slug is required")
	}

	if accessName == "" {
		return fmt.Errorf("name is required")
	}

	container, err := container.NewContainerWithDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer container.DB.Close()

	access, plainPIN, err := container.CreateAccessService.CreateAccess(accessProjectSlug, accessName, accessReadonly)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Access created successfully!\n")
	fmt.Printf("   Project: %s\n", accessProjectSlug)
	fmt.Printf("   Name: %s\n", access.Name)
	fmt.Printf("   UID: %s\n", access.UID)
	fmt.Printf("   PIN: %s\n", plainPIN)
	fmt.Printf("   Read-only: %t\n", access.ReadOnly)
	fmt.Printf("   ID: %s\n", access.ID)

	return nil
}
