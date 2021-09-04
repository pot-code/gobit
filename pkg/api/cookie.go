package api

import (
	"net/http"
	"time"
)

// NewHttpCookie create new cookie
func NewHttpCookie(name, value string, exp time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(exp),
	}
}
