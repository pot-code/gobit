package api

import (
	"net/http"
	"time"
)

// NewRefreshTokenCookie create token cookie
func NewRefreshTokenCookie(name, token string, exp time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(exp),
	}
}
