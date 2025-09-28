package create_project

import (
	"fmt"
	"testing"

	"gofin/internal/infrastructure/database"
	"gofin/internal/models"
)

func TestCreateProjectService_CreateProject(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		customSlug  string
		repoSetup   func(models.ProjectRepository)
		wantErr     bool
		wantSlug    string
	}{
		{
			name:        "success with auto-generated slug",
			projectName: "My Test Project",
			customSlug:  "",
			repoSetup:   func(r models.ProjectRepository) {},
			wantErr:     false,
			wantSlug:    "my-test-project",
		},
		{
			name:        "success with custom slug",
			projectName: "My Test Project",
			customSlug:  "custom-slug",
			repoSetup:   func(r models.ProjectRepository) {},
			wantErr:     false,
			wantSlug:    "custom-slug",
		},
		{
			name:        "error when name is empty",
			projectName: "",
			customSlug:  "",
			repoSetup:   func(r models.ProjectRepository) {},
			wantErr:     true,
		},
		{
			name:        "error when custom slug is invalid",
			projectName: "My Test Project",
			customSlug:  "Invalid_Slug",
			repoSetup:   func(r models.ProjectRepository) {},
			wantErr:     true,
		},
		{
			name:        "error when slug already exists",
			projectName: "My Test Project",
			customSlug:  "existing-slug",
			repoSetup: func(r models.ProjectRepository) {
				existingProject := models.NewProject("Existing Project", "existing-slug")
				r.Create(existingProject)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := database.NewProjectInMemoryRepository()
			service := NewCreateProjectService(repository)
			tt.repoSetup(repository)

			project, err := service.CreateProject(tt.projectName, tt.customSlug)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateProject() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateProject() unexpected error: %v", err)
				return
			}

			if project.Name != tt.projectName {
				t.Errorf("CreateProject() project name = %v, want %v", project.Name, tt.projectName)
			}

			if project.Slug != tt.wantSlug {
				t.Errorf("CreateProject() project slug = %v, want %v", project.Slug, tt.wantSlug)
			}

			if project.ID.String() == "" {
				t.Errorf("CreateProject() project ID should not be empty")
			}
		})
	}
}

func TestCreateProjectService_ensureUniqueSlug(t *testing.T) {
	tests := []struct {
		name          string
		baseSlug      string
		existingSlugs []string
		wantSlug      string
		wantErr       bool
	}{
		{
			name:          "unique slug returned when not exists",
			baseSlug:      "test-project",
			existingSlugs: []string{},
			wantSlug:      "test-project",
			wantErr:       false,
		},
		{
			name:          "numbered slug when base exists",
			baseSlug:      "test-project",
			existingSlugs: []string{"test-project"},
			wantSlug:      "test-project-1",
			wantErr:       false,
		},
		{
			name:          "incremented slug when multiple exist",
			baseSlug:      "test-project",
			existingSlugs: []string{"test-project", "test-project-1"},
			wantSlug:      "test-project-2",
			wantErr:       false,
		},
		{
			name:          "error when max attempts reached",
			baseSlug:      "test-project",
			existingSlugs: generateExistingSlugs("test-project", maxSlugAttempts+1),
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := database.NewProjectInMemoryRepository()
			service := NewCreateProjectService(repository)

			for _, slug := range tt.existingSlugs {
				project := models.NewProject("Test Project", slug)
				repository.Create(project)
			}

			slug, err := service.ensureUniqueSlug(tt.baseSlug)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ensureUniqueSlug() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ensureUniqueSlug() unexpected error: %v", err)
				return
			}

			if slug != tt.wantSlug {
				t.Errorf("ensureUniqueSlug() = %v, want %v", slug, tt.wantSlug)
			}
		})
	}
}

func TestCreateProjectService_determineSlug(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		customSlug  string
		wantSlug    string
		wantErr     bool
	}{
		{
			name:        "auto-generated slug from name",
			projectName: "My Test Project",
			customSlug:  "",
			wantSlug:    "my-test-project",
			wantErr:     false,
		},
		{
			name:        "custom slug when provided",
			projectName: "My Test Project",
			customSlug:  "custom-slug",
			wantSlug:    "custom-slug",
			wantErr:     false,
		},
		{
			name:        "error when custom slug is invalid",
			projectName: "My Test Project",
			customSlug:  "Invalid_Slug",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := database.NewProjectInMemoryRepository()
			service := NewCreateProjectService(repository)
			slug, err := service.determineSlug(tt.projectName, tt.customSlug)

			if tt.wantErr {
				if err == nil {
					t.Errorf("determineSlug() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("determineSlug() unexpected error: %v", err)
				return
			}

			if slug != tt.wantSlug {
				t.Errorf("determineSlug() = %v, want %v", slug, tt.wantSlug)
			}
		})
	}
}

func generateExistingSlugs(baseSlug string, count int) []string {
	slugs := make([]string, count)
	for i := 0; i < count; i++ {
		if i == 0 {
			slugs[i] = baseSlug
		} else {
			slugs[i] = fmt.Sprintf("%s-%d", baseSlug, i)
		}
	}
	return slugs
}
