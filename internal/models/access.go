package models

import (
	"time"

	"github.com/google/uuid"
)

type Access struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	UID       string    `json:"uid" db:"uid"`
	PIN       string    `json:"-" db:"pin"`
	Name      string    `json:"name" db:"name"`
	ReadOnly  bool      `json:"readonly" db:"readonly"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AccessRepository interface {
	Create(access *Access) error
	GetByProjectID(projectID uuid.UUID) ([]*Access, error)
	GetByUID(projectID uuid.UUID, uid string) (*Access, error)
	ExistsByUID(projectID uuid.UUID, uid string) (bool, error)
}

func NewAccess(projectID uuid.UUID, uid, pin, name string, readonly bool) *Access {
	now := time.Now()
	return &Access{
		ID:        uuid.New(),
		ProjectID: projectID,
		UID:       uid,
		PIN:       pin,
		Name:      name,
		ReadOnly:  readonly,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
