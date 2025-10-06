package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gofin/internal/container"
	"gofin/web"
)

func ProjectBased(container *container.Container) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path[1:]
			var projectID uuid.UUID
			var projectSlug string

			if path == "" {
				serveWithProjectID(w, r, next, projectID, projectSlug)
				return
			}

			parts := strings.Split(path, "/")
			projectSlug = parts[0]
			if projectSlug == "" {
				serveWithProjectID(w, r, next, projectID, projectSlug)
				return
			}

			project, err := container.ProjectRepository.GetBySlug(projectSlug)
			if err != nil {
				http.NotFound(w, r)
				return
			}

			projectID = project.ID
			serveWithProjectID(w, r, next, projectID, projectSlug)
		})
	}
}

func serveWithProjectID(w http.ResponseWriter, r *http.Request, next http.Handler, projectID uuid.UUID, projectSlug string) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, web.ContextProjectID, projectID)
	ctx = context.WithValue(ctx, web.ContextProjectSlug, projectSlug)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func GetProjectIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	projectID, ok := ctx.Value(web.ContextProjectID).(uuid.UUID)

	return projectID, ok
}

func GetProjectSlugFromContext(ctx context.Context) (string, bool) {
	projectSlug, ok := ctx.Value(web.ContextProjectSlug).(string)

	return projectSlug, ok
}

func GetProjectIDFromContextOrFail(ctx context.Context, w http.ResponseWriter) (uuid.UUID, bool) {
	projectID, ok := ctx.Value(web.ContextProjectID).(uuid.UUID)

	if !ok {
		http.Error(w, web.ProjectIDNotFoundError, http.StatusInternalServerError)
		return uuid.Nil, false
	}

	if projectID == uuid.Nil {
		http.Error(w, web.ProjectIDNotFoundError, http.StatusInternalServerError)
		return uuid.Nil, false
	}

	return projectID, ok
}
