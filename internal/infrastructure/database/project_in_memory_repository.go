package database

import (
	"fmt"
	"sync"

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
