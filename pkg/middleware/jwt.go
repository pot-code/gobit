package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/api"
	"github.com/pot-code/gobit/pkg/auth"
)

var DefaultTokenContextKey = gobit.AppContextKey("user")

var DefaultRefreshTokenEchoKey = "refresh"

// RefreshTokenOption ...
type RefreshTokenOption struct {
	Skipper        func(uri string) bool
	EchoContextKey string
	TokenName      string
}

// VerifyRefreshToken validate refresh JWT
func VerifyRefreshToken(jp *auth.JwtAuth, options RefreshTokenOption) echo.MiddlewareFunc {
	skipper := func(string) bool { return false }
	echoKey := DefaultRefreshTokenEchoKey
	tokenName := "refresh_token"

	if options.EchoContextKey != "" {
		echoKey = options.EchoContextKey
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
				c.Set(echoKey, token)
				return next(c)
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}

// ValidateTokenOption ...
type ValidateTokenOption struct {
	Skipper    func(uri string) bool
	ContextKey gobit.AppContextKey
}

// VerifyRefreshToken validate normal JWT
func VerifyAccessToken(jp *auth.JwtAuth, options ValidateTokenOption) echo.MiddlewareFunc {
	skipper := func(string) bool { return false }
	contextKey := DefaultTokenContextKey

	if options.ContextKey != "" {
		contextKey = options.ContextKey
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
				api.WithContextValue(c, contextKey, token)
				return next(c)
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}
