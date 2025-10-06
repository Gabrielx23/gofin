package database

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"gofin/internal/models"
)

type AccessInMemoryRepository struct {
	accesses map[string]*models.Access
	mu       sync.RWMutex
}

func NewAccessInMemoryRepository() *AccessInMemoryRepository {
	return &AccessInMemoryRepository{
		accesses: make(map[string]*models.Access),
	}
}

func (r *AccessInMemoryRepository) Create(access *models.Access) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.getKey(access.ProjectID, access.UID)
	if _, exists := r.accesses[key]; exists {
		return fmt.Errorf("access with UID '%s' already exists for project", access.UID)
	}

	r.accesses[key] = access
	return nil
}

func (r *AccessInMemoryRepository) GetByProjectID(projectID uuid.UUID) ([]*models.Access, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var accesses []*models.Access
	for _, access := range r.accesses {
		if access.ProjectID == projectID {
			accesses = append(accesses, access)
		}
	}

	return accesses, nil
}

func (r *AccessInMemoryRepository) GetByUID(projectID uuid.UUID, uid string) (*models.Access, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.getKey(projectID, uid)
	access, exists := r.accesses[key]
	if !exists {
		return nil, fmt.Errorf("access with UID '%s' not found for project", uid)
	}

	return access, nil
}

func (r *AccessInMemoryRepository) ExistsByUID(projectID uuid.UUID, uid string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.getKey(projectID, uid)
	_, exists := r.accesses[key]
	return exists, nil
}

func (r *AccessInMemoryRepository) GetByID(id uuid.UUID) (*models.Access, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, access := range r.accesses {
		if access.ID == id {
			return access, nil
		}
	}

	return nil, fmt.Errorf("access with ID '%s' not found", id.String())
}

func (r *AccessInMemoryRepository) getKey(projectID uuid.UUID, uid string) string {
	return fmt.Sprintf("%s:%s", projectID.String(), uid)
}
