package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/api"
)

var DefaultUserContextKey = gobit.AppContextKey("user")

var DefaultRefreshTokenEchoKey = "refresh"

// RefreshTokenOption ...
type RefreshTokenOption struct {
	Skipper        func(uri string) bool
	EchoContextKey string
	Secret         []byte
	TokenName      string
	Algorithm      jwt.SigningMethod
}

// VerifyRefreshToken validate refresh JWT
func VerifyRefreshToken(options RefreshTokenOption) echo.MiddlewareFunc {
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

	var (
		secret interface{} = options.Secret
		err    error
	)
	if options.Algorithm == jwt.SigningMethodRS256 {
		secret, err = jwt.ParseRSAPublicKeyFromPEM(options.Secret)
		if err != nil {
			log.Fatal(err)
		}
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

			token, err := jwt.Parse(tokenStr.Value, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err == nil {
				c.Set(echoKey, token.Claims)
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
	Secret     []byte
	Algorithm  jwt.SigningMethod
}

// VerifyRefreshToken validate normal JWT
func VerifyAccessToken(options ValidateTokenOption) echo.MiddlewareFunc {
	skipper := func(string) bool { return false }
	contextKey := DefaultUserContextKey
	if options.ContextKey != "" {
		contextKey = options.ContextKey
	}
	if options.Skipper != nil {
		skipper = options.Skipper
	}

	var (
		secret interface{} = options.Secret
		err    error
	)
	if options.Algorithm == jwt.SigningMethodRS256 {
		secret, err = jwt.ParseRSAPublicKeyFromPEM(options.Secret)
		if err != nil {
			log.Fatal(err)
		}
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
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err == nil {
				api.WithContextValue(c, contextKey, token.Claims)
				return next(c)
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}
