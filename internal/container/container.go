package container

import (
	"fmt"
	"path/filepath"

	"gofin/internal/cases/create_project"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
)

type Container struct {
	ProjectRepository    models.ProjectRepository
	CreateProjectService *create_project.CreateProjectService
	DB                   database.Database
}

func NewContainer(dbPath string) (*Container, error) {
	db, err := database.NewDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	projectRepo := database.NewProjectSqliteRepository(db.GetConnection())
	createProjectService := create_project.NewCreateProjectService(projectRepo)

	return &Container{
		ProjectRepository:    projectRepo,
		CreateProjectService: createProjectService,
		DB:                   db,
	}, nil
}

func NewContainerWithDefaultConfig() (*Container, error) {
	return NewContainer(filepath.Join(".", "database.db"))
}
