package middleware

import (
	"net/http"
	"strings"

	"gofin/internal/container"
	"gofin/internal/models"
	webcontext "gofin/pkg/web"
)

func ProjectBased(container *container.Container) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path[1:]

			if path == "" {
				http.NotFound(w, r)
				return
			}

			parts := strings.Split(path, "/")
			projectSlug := parts[0]
			if projectSlug == "" {
				http.NotFound(w, r)
				return
			}

			project, err := container.ProjectRepository.GetBySlug(projectSlug)
			if err != nil {
				http.NotFound(w, r)
				return
			}

			serveWithProject(w, r, next, project)
		})
	}
}

func serveWithProject(w http.ResponseWriter, r *http.Request, next http.Handler, project *models.Project) {
	ctx := r.Context()
	ctx = webcontext.SetProject(ctx, project)
	next.ServeHTTP(w, r.WithContext(ctx))
}
