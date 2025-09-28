package create_transaction

import (
	"testing"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
)

func TestCreateTransactionService_CreateSingleTransaction(t *testing.T) {
	tests := []struct {
		name            string
		accountID       uuid.UUID
		value           float64
		transactionName string
		transactionType models.TransactionType
		repoSetup       func(models.AccountRepository, models.TransactionRepository, uuid.UUID)
		wantErr         bool
	}{
		{
			name:            "success with debit transaction",
			accountID:       uuid.New(),
			value:           50.0,
			transactionName: "Grocery Shopping",
			transactionType: models.Debit,
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: false,
		},
		{
			name:            "success with top-up transaction",
			accountID:       uuid.New(),
			value:           100.0,
			transactionName: "Salary Deposit",
			transactionType: models.TopUp,
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: false,
		},
		{
			name:            "error when account not found",
			accountID:       uuid.New(),
			value:           50.0,
			transactionName: "Test Transaction",
			transactionType: models.Debit,
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
			},
			wantErr: true,
		},
		{
			name:            "error when name is empty",
			accountID:       uuid.New(),
			value:           50.0,
			transactionName: "",
			transactionType: models.Debit,
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: true,
		},
		{
			name:            "error when value is zero or negative",
			accountID:       uuid.New(),
			value:           0.0,
			transactionName: "Test Transaction",
			transactionType: models.Debit,
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := database.NewAccountInMemoryRepository()
			transactionRepo := database.NewTransactionInMemoryRepository()
			service := NewCreateTransactionService(transactionRepo, accountRepo)
			tt.repoSetup(accountRepo, transactionRepo, tt.accountID)

			transaction, err := service.CreateSingleTransaction(tt.accountID, tt.value, tt.transactionName, tt.transactionType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateSingleTransaction() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateSingleTransaction() unexpected error: %v", err)
				return
			}

			if transaction == nil {
				t.Errorf("CreateSingleTransaction() expected transaction, got nil")
				return
			}

			if transaction.AccountID != tt.accountID {
				t.Errorf("CreateSingleTransaction() accountID = %v, want %v", transaction.AccountID, tt.accountID)
			}

			if transaction.Value != tt.value {
				t.Errorf("CreateSingleTransaction() value = %v, want %v", transaction.Value, tt.value)
			}

			if transaction.Name != tt.transactionName {
				t.Errorf("CreateSingleTransaction() name = %v, want %v", transaction.Name, tt.transactionName)
			}

			if transaction.Type != tt.transactionType {
				t.Errorf("CreateSingleTransaction() type = %v, want %v", transaction.Type, tt.transactionType)
			}

			if transaction.GroupID != nil {
				t.Errorf("CreateSingleTransaction() GroupID should be nil for single transaction")
			}
		})
	}
}

func TestCreateTransactionService_CreateGroupedTransactions(t *testing.T) {
	tests := []struct {
		name         string
		transactions []models.TransactionData
		repoSetup    func(models.AccountRepository, models.TransactionRepository, []uuid.UUID)
		wantErr      bool
		wantCount    int
	}{
		{
			name: "success with multiple transactions",
			transactions: []models.TransactionData{
				{AccountID: uuid.New(), Value: 50.0, Name: "Transaction 1", Type: models.Debit},
				{AccountID: uuid.New(), Value: 100.0, Name: "Transaction 2", Type: models.TopUp},
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountIDs []uuid.UUID) {
				account1 := models.NewAccount(uuid.New(), "Account 1", "USD")
				account1.ID = accountIDs[0]
				accountRepo.Create(account1)

				account2 := models.NewAccount(uuid.New(), "Account 2", "USD")
				account2.ID = accountIDs[1]
				accountRepo.Create(account2)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:         "error when no transactions provided",
			transactions: []models.TransactionData{},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountIDs []uuid.UUID) {
			},
			wantErr: true,
		},
		{
			name: "error when account not found",
			transactions: []models.TransactionData{
				{AccountID: uuid.New(), Value: 50.0, Name: "Transaction 1", Type: models.Debit},
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountIDs []uuid.UUID) {
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := database.NewAccountInMemoryRepository()
			transactionRepo := database.NewTransactionInMemoryRepository()
			service := NewCreateTransactionService(transactionRepo, accountRepo)

			var accountIDs []uuid.UUID
			for _, tx := range tt.transactions {
				accountIDs = append(accountIDs, tx.AccountID)
			}
			tt.repoSetup(accountRepo, transactionRepo, accountIDs)

			transactions, err := service.CreateGroupedTransactions(tt.transactions)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateGroupedTransactions() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateGroupedTransactions() unexpected error: %v", err)
				return
			}

			if len(transactions) != tt.wantCount {
				t.Errorf("CreateGroupedTransactions() count = %v, want %v", len(transactions), tt.wantCount)
			}

			if len(transactions) > 0 {
				groupID := transactions[0].GroupID
				if groupID == nil {
					t.Errorf("CreateGroupedTransactions() GroupID should not be nil for grouped transactions")
				}

				for i, transaction := range transactions {
					if transaction.GroupID == nil || *transaction.GroupID != *groupID {
						t.Errorf("CreateGroupedTransactions() transaction %d should have same GroupID", i)
					}
				}
			}
		})
	}
}

func TestCreateTransactionService_CreateSingleTransactionFromData(t *testing.T) {
	tests := []struct {
		name      string
		data      models.TransactionData
		repoSetup func(models.AccountRepository, models.TransactionRepository, uuid.UUID)
		wantErr   bool
	}{
		{
			name: "success with debit transaction",
			data: models.TransactionData{
				AccountID: uuid.New(),
				Value:     50.0,
				Name:      "Grocery Shopping",
				Type:      models.Debit,
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: false,
		},
		{
			name: "success with top-up transaction",
			data: models.TransactionData{
				AccountID: uuid.New(),
				Value:     100.0,
				Name:      "Salary Deposit",
				Type:      models.TopUp,
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: false,
		},
		{
			name: "error when account not found",
			data: models.TransactionData{
				AccountID: uuid.New(),
				Value:     50.0,
				Name:      "Test Transaction",
				Type:      models.Debit,
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
			},
			wantErr: true,
		},
		{
			name: "error when name is empty",
			data: models.TransactionData{
				AccountID: uuid.New(),
				Value:     50.0,
				Name:      "",
				Type:      models.Debit,
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: true,
		},
		{
			name: "error when value is zero or negative",
			data: models.TransactionData{
				AccountID: uuid.New(),
				Value:     0.0,
				Name:      "Test Transaction",
				Type:      models.Debit,
			},
			repoSetup: func(accountRepo models.AccountRepository, transactionRepo models.TransactionRepository, accountID uuid.UUID) {
				account := models.NewAccount(uuid.New(), "Test Account", "USD")
				account.ID = accountID
				accountRepo.Create(account)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := database.NewAccountInMemoryRepository()
			transactionRepo := database.NewTransactionInMemoryRepository()
			service := NewCreateTransactionService(transactionRepo, accountRepo)
			tt.repoSetup(accountRepo, transactionRepo, tt.data.AccountID)

			transaction, err := service.CreateSingleTransactionFromData(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateSingleTransactionFromData() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateSingleTransactionFromData() unexpected error: %v", err)
				return
			}

			if transaction == nil {
				t.Errorf("CreateSingleTransactionFromData() expected transaction, got nil")
				return
			}

			if transaction.AccountID != tt.data.AccountID {
				t.Errorf("CreateSingleTransactionFromData() accountID = %v, want %v", transaction.AccountID, tt.data.AccountID)
			}

			if transaction.Value != tt.data.Value {
				t.Errorf("CreateSingleTransactionFromData() value = %v, want %v", transaction.Value, tt.data.Value)
			}

			if transaction.Name != tt.data.Name {
				t.Errorf("CreateSingleTransactionFromData() name = %v, want %v", transaction.Name, tt.data.Name)
			}

			if transaction.Type != tt.data.Type {
				t.Errorf("CreateSingleTransactionFromData() type = %v, want %v", transaction.Type, tt.data.Type)
			}

			if transaction.GroupID != nil {
				t.Errorf("CreateSingleTransactionFromData() GroupID should be nil for single transaction")
			}
		})
	}
}
