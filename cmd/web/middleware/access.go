package middleware

import (
	"net/http"

	"gofin/internal/container"
	webcontext "gofin/pkg/web"
	"gofin/web"
)

func ReadOnlyProhibited(container *container.Container) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			access, ok := webcontext.GetAccess(r.Context())
			if !ok {
				http.Error(w, web.AccessIDNotFoundError, http.StatusInternalServerError)
				return
			}

			if access.ReadOnly {
				http.Error(w, "Access denied: Write access required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
