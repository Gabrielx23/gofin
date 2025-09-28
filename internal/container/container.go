package container

import (
	"fmt"
	"path/filepath"

	"gofin/internal/cases/create-project"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
)

type Container struct {
	DB                   *database.DB
	ProjectRepository    models.ProjectRepository
	CreateProjectService *createproject.CreateProjectService
}

func NewContainer(dbPath string) (*Container, error) {
	db, err := database.NewDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	projectRepo := database.NewProjectSqliteRepository(db.GetConnection())
	createProjectService := createproject.NewCreateProjectService(projectRepo)

	return &Container{
		DB:                   db,
		ProjectRepository:    projectRepo,
		CreateProjectService: createProjectService,
	}, nil
}

func NewContainerWithDefaultConfig() (*Container, error) {
	return NewContainer(filepath.Join(".", "database.db"))
}

func (c *Container) Close() error {
	return c.DB.Close()
}
