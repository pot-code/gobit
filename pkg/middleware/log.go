package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type LoggingConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

// Logging create a logging middleware with zap logger
func Logging(logger *zap.Logger, options ...LoggingConfig) echo.MiddlewareFunc {
	cfg := &LoggingConfig{
		Skipper: middleware.DefaultSkipper,
	}
	if len(options) > 0 {
		option := options[0]
		if option.Skipper != nil {
			cfg.Skipper = option.Skipper
		}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}
			err := next(c)
			code := c.Response().Status
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			logger.Info(
				http.StatusText(code),
				zap.String("trace", rid),
				zap.String("url.original", c.Path()),
				zap.String("http.request.referrer", c.Request().Referer()),
				zap.String("http.request.method", c.Request().Method),
				zap.Int("http.request.status_code", code),
				zap.Strings("route.params.name", c.ParamNames()),
				zap.Strings("route.params.value", c.ParamValues()),
				zap.Int64("http.response.body.bytes", c.Response().Size),
			)
			return err
		}
	}
}
