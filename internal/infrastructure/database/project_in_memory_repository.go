package database

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type ProjectInMemoryRepository struct {
	projects map[string]*models.Project
	mu       sync.RWMutex
}

func NewProjectInMemoryRepository() models.ProjectRepository {
	return &ProjectInMemoryRepository{
		projects: make(map[string]*models.Project),
	}
}

func (r *ProjectInMemoryRepository) Create(project *models.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.projects[project.Slug]; exists {
		return fmt.Errorf("project with slug '%s' already exists", project.Slug)
	}

	r.projects[project.Slug] = project
	return nil
}

func (r *ProjectInMemoryRepository) GetBySlug(slug string) (*models.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	project, exists := r.projects[slug]
	if !exists {
		return nil, fmt.Errorf("project not found")
	}

	return project, nil
}

func (r *ProjectInMemoryRepository) ExistsBySlug(slug string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.projects[slug]
	return exists, nil
}

func (r *ProjectInMemoryRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, project := range r.projects {
		if project.ID == id {
			return project, nil
		}
	}

	return nil, fmt.Errorf("project with ID '%s' not found", id.String())
}
