package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gofin/internal/container"
	"gofin/pkg/session"
	"gofin/web"
)

func AuthRequired(container *container.Container, sessionManager *session.SessionManager) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
	projectID, ok := GetProjectIDFromContext(r.Context())
	if !ok {
		http.Error(w, web.ProjectIDNotFoundError, http.StatusInternalServerError)
		return
	}

			sessionToken, err := getSessionTokenFromCookie(r)
			if err != nil {
				redirectToLogin(w, r)
				return
			}

			token, valid := sessionManager.ValidateSessionToken(sessionToken)
			if !valid {
				clearInvalidCookie(w)
				redirectToLogin(w, r)
				return
			}

			if token.ProjectID != projectID.String() {
				clearInvalidCookie(w)
				redirectToLogin(w, r)
				return
			}

			if !validateAccessID(container, projectID, token.AccessID) {
				clearInvalidCookie(w)
				redirectToLogin(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), web.ContextAccessID, token.AccessID)
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

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	projectSlug, ok := GetProjectSlugFromContext(r.Context())
	if !ok {
		http.Error(w, web.ProjectSlugNotFoundError, http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/"+projectSlug+web.RouteLogin, http.StatusSeeOther)
}


func GetAccessIDFromContext(ctx context.Context) (string, bool) {
	accessID, ok := ctx.Value(web.ContextAccessID).(string)
	return accessID, ok
}

func validateAccessID(container *container.Container, projectID uuid.UUID, accessID string) bool {
	accessUUID, err := uuid.Parse(accessID)
	if err != nil {
		return false
	}

	access, err := container.AccessRepository.GetByID(accessUUID)
	if err != nil {
		return false
	}

	return access.ProjectID == projectID
}

func clearInvalidCookie(w http.ResponseWriter) {
	session.ClearSessionCookie(w)
}
