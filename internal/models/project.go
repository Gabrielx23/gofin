package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Slug      string    `json:"slug" db:"slug"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ProjectRepository interface {
	Create(project *Project) error
	GetBySlug(slug string) (*Project, error)
	ExistsBySlug(slug string) (bool, error)
}

func NewProject(name, slug string) *Project {
	now := time.Now()
	return &Project{
		ID:        uuid.New(),
		Slug:      slug,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
