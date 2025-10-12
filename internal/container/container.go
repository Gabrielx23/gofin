package container

import (
	"fmt"
	"path/filepath"

	"gofin/internal/cases/create_access"
	"gofin/internal/cases/create_account"
	"gofin/internal/cases/create_project"
	"gofin/internal/cases/create_transaction"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
	"gofin/web"
)

type Container struct {
	ProjectRepository        models.ProjectRepository
	AccessRepository         models.AccessRepository
	AccountRepository        models.AccountRepository
	TransactionRepository    models.TransactionRepository
	CreateProjectService     *create_project.CreateProjectService
	CreateAccessService      *create_access.CreateAccessService
	CreateAccountService     *create_account.CreateAccountService
	CreateTransactionService *create_transaction.CreateTransactionService
	DB                       database.Database
}

func NewContainer(dbPath string) (*Container, error) {
	db, err := database.NewDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	projectRepo := database.NewProjectSqliteRepository(db.GetConnection())
	accessRepo := database.NewAccessSqliteRepository(db.GetConnection())
	accountRepo := database.NewAccountSqliteRepository(db.GetConnection())
	transactionRepo := database.NewTransactionSqliteRepository(db.GetConnection())
	createProjectService := create_project.NewCreateProjectService(projectRepo)
	createAccessService := create_access.NewCreateAccessService(accessRepo, projectRepo)
	createAccountService := create_account.NewCreateAccountService(accountRepo)
	createTransactionService := create_transaction.NewCreateTransactionService(transactionRepo, accountRepo, projectRepo)

	return &Container{
		ProjectRepository:        projectRepo,
		AccessRepository:         accessRepo,
		AccountRepository:        accountRepo,
		TransactionRepository:    transactionRepo,
		CreateProjectService:     createProjectService,
		CreateAccessService:      createAccessService,
		CreateAccountService:     createAccountService,
		CreateTransactionService: createTransactionService,
		DB:                       db,
	}, nil
}

func NewContainerWithDefaultConfig() (*Container, error) {
	return NewContainer(filepath.Join(".", web.DatabaseFile))
}
