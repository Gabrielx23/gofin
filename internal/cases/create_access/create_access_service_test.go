package create_access

import (
	"testing"

	"github.com/google/uuid"
	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
)

func TestCreateAccessService_CreateAccess(t *testing.T) {
	tests := []struct {
		name        string
		projectSlug string
		accessName  string
		readonly    bool
		repoSetup   func(models.ProjectRepository, models.AccessRepository)
		wantErr     bool
	}{
		{
			name:        "success with read-write access",
			projectSlug: "test-project",
			accessName:  "Test Access",
			readonly:    false,
			repoSetup: func(projectRepo models.ProjectRepository, accessRepo models.AccessRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
			},
			wantErr: false,
		},
		{
			name:        "success with read-only access",
			projectSlug: "test-project",
			accessName:  "Read Only Access",
			readonly:    true,
			repoSetup: func(projectRepo models.ProjectRepository, accessRepo models.AccessRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
			},
			wantErr: false,
		},
		{
			name:        "error when project not found",
			projectSlug: "non-existent-project",
			accessName:  "Test Access",
			readonly:    false,
			repoSetup:   func(projectRepo models.ProjectRepository, accessRepo models.AccessRepository) {},
			wantErr:     true,
		},
		{
			name:        "error when name is empty",
			projectSlug: "test-project",
			accessName:  "",
			readonly:    false,
			repoSetup: func(projectRepo models.ProjectRepository, accessRepo models.AccessRepository) {
				project := models.NewProject("Test Project", "test-project")
				projectRepo.Create(project)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRepo := database.NewProjectInMemoryRepository()
			accessRepo := database.NewAccessInMemoryRepository()
			service := NewCreateAccessService(accessRepo, projectRepo)
			tt.repoSetup(projectRepo, accessRepo)

			access, plainPIN, err := service.CreateAccess(tt.projectSlug, tt.accessName, tt.readonly)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateAccess() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateAccess() unexpected error: %v", err)
				return
			}

			if access == nil {
				t.Errorf("CreateAccess() expected access, got nil")
				return
			}

			if access.Name != tt.accessName {
				t.Errorf("CreateAccess() name = %v, want %v", access.Name, tt.accessName)
			}

			if access.ReadOnly != tt.readonly {
				t.Errorf("CreateAccess() readonly = %v, want %v", access.ReadOnly, tt.readonly)
			}

			if access.UID == "" {
				t.Errorf("CreateAccess() UID should not be empty")
			}

			if access.PinHash == "" {
				t.Errorf("CreateAccess() PinHash should not be empty")
			}

			if len(access.UID) != 2 {
				t.Errorf("CreateAccess() UID length = %v, want 2", len(access.UID))
			}

			if len(plainPIN) != 8 {
				t.Errorf("CreateAccess() plain PIN length = %v, want 8", len(plainPIN))
			}

			if access.PinHash == plainPIN {
				t.Errorf("CreateAccess() PinHash should be hashed, not plain text")
			}

			if len(access.PinHash) < 50 {
				t.Errorf("CreateAccess() hashed PinHash should be much longer than plain PIN")
			}
		})
	}
}

func TestCreateAccessService_generateUniqueUID(t *testing.T) {
	tests := []struct {
		name           string
		projectID      uuid.UUID
		existingUIDs   []string
		repoSetup      func(models.AccessRepository, uuid.UUID, []string)
		wantErr        bool
		expectedLength int
	}{
		{
			name:         "unique UID when none exist",
			projectID:    uuid.New(),
			existingUIDs: []string{},
			repoSetup: func(accessRepo models.AccessRepository, projectID uuid.UUID, existingUIDs []string) {
				for _, uid := range existingUIDs {
					access := models.NewAccess(projectID, uid, "12345678", "Test", false)
					accessRepo.Create(access)
				}
			},
			wantErr:        false,
			expectedLength: 2,
		},
		{
			name:         "unique UID when some exist",
			projectID:    uuid.New(),
			existingUIDs: []string{"01", "02", "03"},
			repoSetup: func(accessRepo models.AccessRepository, projectID uuid.UUID, existingUIDs []string) {
				for _, uid := range existingUIDs {
					access := models.NewAccess(projectID, uid, "12345678", "Test", false)
					accessRepo.Create(access)
				}
			},
			wantErr:        false,
			expectedLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessRepo := database.NewAccessInMemoryRepository()
			projectRepo := database.NewProjectInMemoryRepository()
			service := NewCreateAccessService(accessRepo, projectRepo)
			tt.repoSetup(accessRepo, tt.projectID, tt.existingUIDs)

			uid, err := service.generateUniqueUID(tt.projectID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("generateUniqueUID() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("generateUniqueUID() unexpected error: %v", err)
				return
			}

			if len(uid) != tt.expectedLength {
				t.Errorf("generateUniqueUID() UID length = %v, want %v", len(uid), tt.expectedLength)
			}

			exists, err := accessRepo.ExistsByUID(tt.projectID, uid)
			if err != nil {
				t.Errorf("generateUniqueUID() failed to check UID existence: %v", err)
				return
			}

			if exists {
				t.Errorf("generateUniqueUID() generated UID that already exists")
			}
		})
	}
}
