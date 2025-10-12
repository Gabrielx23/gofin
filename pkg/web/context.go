package web

import (
	"context"

	"gofin/internal/models"
)

type contextKey string

const (
	projectKey contextKey = "project"
	accessKey  contextKey = "access"
)

func SetProject(ctx context.Context, project *models.Project) context.Context {
	return context.WithValue(ctx, projectKey, project)
}

func GetProject(ctx context.Context) (*models.Project, bool) {
	project, ok := ctx.Value(projectKey).(*models.Project)
	return project, ok
}

func SetAccess(ctx context.Context, access *models.Access) context.Context {
	return context.WithValue(ctx, accessKey, access)
}

func GetAccess(ctx context.Context) (*models.Access, bool) {
	access, ok := ctx.Value(accessKey).(*models.Access)
	return access, ok
}
