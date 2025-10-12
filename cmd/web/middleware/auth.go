package middleware

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	webcontext "gofin/pkg/web"
	webpkg "gofin/pkg/web"
	"gofin/internal/container"
	"gofin/pkg/session"
	"gofin/web"
)

func AuthRequired(container *container.Container, sessionManager *session.SessionManager) func(http.HandlerFunc) http.HandlerFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			project, ok := webcontext.GetProject(r.Context())
			if !ok {
				http.Error(w, web.ProjectIDNotFoundError, http.StatusInternalServerError)
				return
			}

			sessionToken, err := getSessionTokenFromCookie(r)
			if err != nil {
				redirectToLogin(w, r, container)
				return
			}

			token, valid := sessionManager.ValidateSessionToken(sessionToken)
			if !valid {
				clearInvalidCookie(w)
				redirectToLogin(w, r, container)
				return
			}

			if token.ProjectID != project.ID.String() {
				clearInvalidCookie(w)
				redirectToLogin(w, r, container)
				return
			}

			access, err := container.AccessRepository.GetByID(uuid.MustParse(token.AccessID))
			if err != nil || access.ProjectID != project.ID {
				clearInvalidCookie(w)
				redirectToLogin(w, r, container)
				return
			}

			ctx := webcontext.SetAccess(r.Context(), access)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func getSessionTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(web.SessionTokenCookie)
	if err != nil {
		return "", err
	}

	if cookie.Value == web.EmptyString {
		return "", fmt.Errorf(web.EmptySessionTokenError)
	}

	return cookie.Value, nil
}

func redirectToLogin(w http.ResponseWriter, r *http.Request, container *container.Container) {
	project, ok := webcontext.GetProject(r.Context())
	if !ok {
		http.Error(w, web.ProjectIDNotFoundError, http.StatusInternalServerError)
		return
	}

	webpkg.RedirectToProjectLogin(w, r, project.Slug)
}



func clearInvalidCookie(w http.ResponseWriter) {
	session.ClearSessionCookie(w)
}

