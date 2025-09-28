package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"gofin/internal/container"
)

var (
	projectName string
	projectSlug string
)

var createProjectCmd = &cobra.Command{
	Use:   "create-project",
	Short: "Create a new financial project",
	Long:  `Create a new financial project with a unique slug and name.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := createProject(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	createProjectCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (required)")
	createProjectCmd.Flags().StringVarP(&projectSlug, "slug", "s", "", "Project slug (optional, will be generated from name if not provided)")

	createProjectCmd.MarkFlagRequired("name")
}

func createProject() error {
	container, err := container.NewContainerWithDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer container.Close()

	project, err := container.CreateProjectService.CreateProject(projectName, projectSlug)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Project created successfully!\n")
	fmt.Printf("   Name: %s\n", project.Name)
	fmt.Printf("   Slug: %s\n", project.Slug)
	fmt.Printf("   ID: %s\n", project.ID.String())

	return nil
}
