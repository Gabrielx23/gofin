package create_access

import (
	"fmt"

	"github.com/google/uuid"
	"gofin/internal/models"
	"gofin/pkg/password"
	"gofin/pkg/random"
)

type CreateAccessService struct {
	accessRepo  models.AccessRepository
	projectRepo models.ProjectRepository
}

func NewCreateAccessService(accessRepo models.AccessRepository, projectRepo models.ProjectRepository) *CreateAccessService {
	return &CreateAccessService{
		accessRepo:  accessRepo,
		projectRepo: projectRepo,
	}
}

func (s *CreateAccessService) CreateAccess(projectSlug, name string, readonly bool) (*models.Access, string, error) {
	if name == "" {
		return nil, "", fmt.Errorf("name is required")
	}

	project, err := s.projectRepo.GetBySlug(projectSlug)
	if err != nil {
		return nil, "", fmt.Errorf("project not found: %w", err)
	}

	uid, err := s.generateUniqueUID(project.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate unique UID: %w", err)
	}

	if err := s.validateUID(uid); err != nil {
		return nil, "", fmt.Errorf("invalid generated UID: %w", err)
	}

	pin := random.GenerateRandomNumber(8)
	if err := s.validatePIN(pin); err != nil {
		return nil, "", fmt.Errorf("invalid generated PIN: %w", err)
	}

	hashedPIN, err := password.Hash(pin)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash PIN: %w", err)
	}

	accessRecord := models.NewAccess(project.ID, uid, hashedPIN, name, readonly)

	if err := s.accessRepo.Create(accessRecord); err != nil {
		return nil, "", fmt.Errorf("failed to create access: %w", err)
	}

	return accessRecord, pin, nil
}

func (s *CreateAccessService) generateUniqueUID(projectID uuid.UUID) (string, error) {
	const maxAttempts = 100

	for i := 0; i < maxAttempts; i++ {
		uid := random.GenerateRandomNumber(2)

		exists, err := s.accessRepo.ExistsByUID(projectID, uid)
		if err != nil {
			return "", fmt.Errorf("failed to check UID existence: %w", err)
		}

		if !exists {
			return uid, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique UID after %d attempts", maxAttempts)
}

func (s *CreateAccessService) validateUID(uid string) error {
	if len(uid) != 2 {
		return fmt.Errorf("UID must be exactly 2 characters")
	}

	for _, char := range uid {
		if char < '0' || char > '9' {
			return fmt.Errorf("UID must contain only digits")
		}
	}

	return nil
}

func (s *CreateAccessService) validatePIN(pin string) error {
	if len(pin) != 8 {
		return fmt.Errorf("PIN must be exactly 8 characters")
	}

	for _, char := range pin {
		if char < '0' || char > '9' {
			return fmt.Errorf("PIN must contain only digits")
		}
	}

	return nil
}
