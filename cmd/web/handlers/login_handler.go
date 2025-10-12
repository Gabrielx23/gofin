package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gofin/internal/container"
	"gofin/pkg/password"
	"gofin/pkg/session"
	webpkg "gofin/pkg/web"
	"gofin/web"
	"gofin/web/components"
)

type LoginHandler struct {
	container      *container.Container
	loginComponent *components.LoginComponent
	sessionManager *session.SessionManager
}

func NewLoginHandler(container *container.Container, loginComponent *components.LoginComponent, sessionManager *session.SessionManager) *LoginHandler {
	return &LoginHandler{
		container:      container,
		loginComponent: loginComponent,
		sessionManager: sessionManager,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	project, _ := webpkg.GetProject(r.Context())
	h.handleLoginPost(w, r, project.ID, project.Slug)
}

func (h *LoginHandler) handleLoginPost(w http.ResponseWriter, r *http.Request, projectID uuid.UUID, projectSlug string) {
	if err := r.ParseForm(); err != nil {
		h.loginComponent.RenderLoginPage(w, r, projectSlug, "Invalid form data")
		return
	}

	uid := h.extractUID(r)
	pin := h.extractPIN(r)

	if len(uid) != 2 || len(pin) != 8 {
		h.loginComponent.RenderLoginPage(w, r, projectSlug, "Invalid UID or PIN format")
		return
	}

	access, err := h.container.AccessRepository.GetByUID(projectID, uid)
	if err != nil {
		h.loginComponent.RenderLoginPage(w, r, projectSlug, "Invalid credentials")
		return
	}

	valid, err := password.Verify(pin, access.PinHash)
	if err != nil || !valid {
		h.loginComponent.RenderLoginPage(w, r, projectSlug, "Invalid credentials")
		return
	}

	sessionToken, err := h.sessionManager.GenerateSessionToken(access.ID.String(), projectID.String())
	if err != nil {
		h.loginComponent.RenderLoginPage(w, r, projectSlug, "Failed to create session")
		return
	}

	session.SetSessionCookie(w, sessionToken)

	webpkg.RedirectToProjectHomeWithSuccess(w, r, projectSlug, web.SuccessKeyLoginSuccessful)
}

func (h *LoginHandler) extractUID(r *http.Request) string {
	var uidParts []string
	for i := 0; i < 2; i++ {
		part := r.FormValue(fmt.Sprintf("uid_%d", i))
		if part == "" {
			return ""
		}
		uidParts = append(uidParts, part)
	}
	return strings.Join(uidParts, "")
}

func (h *LoginHandler) extractPIN(r *http.Request) string {
	var pinParts []string
	for i := 0; i < 8; i++ {
		part := r.FormValue(fmt.Sprintf("pin_%d", i))
		if part == "" {
			return ""
		}
		pinParts = append(pinParts, part)
	}
	return strings.Join(pinParts, "")
}
