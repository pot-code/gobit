package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/api"
	"github.com/pot-code/gobit/pkg/auth"
	"github.com/pot-code/gobit/pkg/context"
)

const RefreshTokenKey context.AppContextKey = "refresh_token"

// RefreshTokenOption ...
type RefreshTokenOption struct {
	Skipper    func(uri string) bool
	ContextKey string
	TokenName  string
}

// VerifyRefreshToken validate refresh JWT
func VerifyRefreshToken(jp *auth.JwtAuth, options RefreshTokenOption) echo.MiddlewareFunc {
	skipper := func(string) bool { return false }
	ctxKey := RefreshTokenKey
	tokenName := "refresh_token"

	if options.ContextKey != "" {
		ctxKey = context.AppContextKey(options.ContextKey)
	}
	if options.Skipper != nil {
		skipper = options.Skipper
	}
	if options.TokenName != "" {
		tokenName = options.TokenName
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c.Request().RequestURI) {
				return next(c)
			}

			tokenStr, err := c.Cookie(tokenName)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			token, err := jp.Validate(tokenStr.Value)
			if err == nil {
				api.WithContextValue(c, ctxKey, token)
				return next(c)
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}

const TokenKey context.AppContextKey = "token"

// ValidateTokenOption ...
type ValidateTokenOption struct {
	Skipper    func(uri string) bool
	ContextKey string
}

// VerifyRefreshToken validate normal JWT
func VerifyAccessToken(jp *auth.JwtAuth, options ValidateTokenOption) echo.MiddlewareFunc {
	skipper := func(string) bool { return false }
	ctxKey := TokenKey

	if options.ContextKey != "" {
		ctxKey = context.AppContextKey(options.ContextKey)
	}
	if options.Skipper != nil {
		skipper = options.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c.Request().RequestURI) {
				return next(c)
			}

			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return c.NoContent(http.StatusUnauthorized)
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			token, err := jp.Validate(tokenStr)
			if err == nil {
				api.WithContextValue(c, ctxKey, token)
				return next(c)
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}
