package createproject

import (
	"fmt"

	"gofin/internal/models"
	"gofin/pkg/slug"
)

type CreateProjectService struct {
	projectRepo models.ProjectRepository
}

func NewCreateProjectService(projectRepo models.ProjectRepository) *CreateProjectService {
	return &CreateProjectService{
		projectRepo: projectRepo,
	}
}

func (s *CreateProjectService) CreateProject(name, customSlug string) (*models.Project, error) {
	if name == "" {
		return nil, fmt.Errorf("project name is required")
	}

	finalSlug, err := s.determineSlug(name, customSlug)
	if err != nil {
		return nil, err
	}

	project := models.NewProject(name, finalSlug)
	
	err = s.projectRepo.Create(project)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

func (s *CreateProjectService) determineSlug(name, customSlug string) (string, error) {
	if customSlug != "" {
		err := slug.Validate(customSlug)
		if err != nil {
			return "", fmt.Errorf("invalid slug: %w", err)
		}
		return customSlug, nil
	}

	baseSlug := slug.Generate(name)
	return s.ensureUniqueSlug(baseSlug)
}

const maxSlugAttempts = 50

func (s *CreateProjectService) ensureUniqueSlug(baseSlug string) (string, error) {
	exists, err := s.projectRepo.ExistsBySlug(baseSlug)
	if err != nil {
		return "", fmt.Errorf("failed to check slug availability: %w", err)
	}

	if !exists {
		return baseSlug, nil
	}

	for counter := 1; counter <= maxSlugAttempts; counter++ {
		candidateSlug := fmt.Sprintf("%s-%d", baseSlug, counter)
		exists, err := s.projectRepo.ExistsBySlug(candidateSlug)
		if err != nil {
			return "", fmt.Errorf("failed to check slug availability: %w", err)
		}
		if !exists {
			return candidateSlug, nil
		}
	}

	return "", fmt.Errorf("unable to generate unique slug after %d attempts", maxSlugAttempts)
}
