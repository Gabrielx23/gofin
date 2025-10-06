package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gofin/web"
	"net/http"
	"strings"
	"time"
)

type SessionToken struct {
	AccessID  string
	ProjectID string
	ExpiresAt int64
	Nonce     string
}

type SessionManager struct {
	secretKey []byte
}

func NewSessionManager() *SessionManager {
	secretKey := make([]byte, 32)
	rand.Read(secretKey)

	return &SessionManager{
		secretKey: secretKey,
	}
}

func (sm *SessionManager) GenerateSessionToken(accessID, projectID string) (string, error) {
	nonce := make([]byte, 16)
	rand.Read(nonce)

	token := SessionToken{
		AccessID:  accessID,
		ProjectID: projectID,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Nonce:     base64.URLEncoding.EncodeToString(nonce),
	}

	tokenData := fmt.Sprintf("%s:%s:%d:%s", token.AccessID, token.ProjectID, token.ExpiresAt, token.Nonce)

	h := hmac.New(sha256.New, sm.secretKey)
	h.Write([]byte(tokenData))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	signedToken := fmt.Sprintf("%s.%s", tokenData, signature)

	return base64.URLEncoding.EncodeToString([]byte(signedToken)), nil
}

func (sm *SessionManager) ValidateSessionToken(token string) (*SessionToken, bool) {
	tokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, false
	}

	signedToken := string(tokenBytes)

	parts := strings.Split(signedToken, ".")
	if len(parts) != 2 {
		return nil, false
	}

	tokenData, signature := parts[0], parts[1]

	h := hmac.New(sha256.New, sm.secretKey)
	h.Write([]byte(tokenData))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, false
	}

	tokenParts := strings.Split(tokenData, ":")
	if len(tokenParts) != 4 {
		return nil, false
	}

	sessionToken := &SessionToken{
		AccessID:  tokenParts[0],
		ProjectID: tokenParts[1],
		ExpiresAt: 0,
		Nonce:     tokenParts[3],
	}

	if _, err := fmt.Sscanf(tokenParts[2], "%d", &sessionToken.ExpiresAt); err != nil {
		return nil, false
	}

	if time.Now().Unix() > sessionToken.ExpiresAt {
		return nil, false
	}

	return sessionToken, true
}

func SetSessionCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     web.SessionTokenCookie,
		Value:    value,
		Path:     web.CookiePath,
		MaxAge:   web.CookieMaxAge,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     web.SessionTokenCookie,
		Value:    "",
		Path:     web.CookiePath,
		MaxAge:   web.CookieMaxAgeClear,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
}
