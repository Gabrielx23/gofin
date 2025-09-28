package container

import (
	"fmt"
	"path/filepath"

	"gofin/internal/cases/create_access"
	"gofin/internal/cases/create_project"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
)

type Container struct {
	ProjectRepository    models.ProjectRepository
	AccessRepository     models.AccessRepository
	CreateProjectService *create_project.CreateProjectService
	CreateAccessService  *create_access.CreateAccessService
	DB                   database.Database
}

func NewContainer(dbPath string) (*Container, error) {
	db, err := database.NewDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	projectRepo := database.NewProjectSqliteRepository(db.GetConnection())
	accessRepo := database.NewAccessSqliteRepository(db.GetConnection())
	createProjectService := create_project.NewCreateProjectService(projectRepo)
	createAccessService := create_access.NewCreateAccessService(accessRepo, projectRepo)

	return &Container{
		ProjectRepository:    projectRepo,
		AccessRepository:     accessRepo,
		CreateProjectService: createProjectService,
		CreateAccessService:  createAccessService,
		DB:                   db,
	}, nil
}

func NewContainerWithDefaultConfig() (*Container, error) {
	return NewContainer(filepath.Join(".", "database.db"))
}
