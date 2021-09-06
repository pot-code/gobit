package middleware

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type LoggingConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
	LogFn   func(echo.Context)
}

// Logging create a logging middleware
func Logging(option LoggingConfig) echo.MiddlewareFunc {
	cfg := &LoggingConfig{
		Skipper: middleware.DefaultSkipper,
	}
	logFn := func(c echo.Context) {
		code := c.Response().Status
		rid := c.Response().Header().Get(echo.HeaderXRequestID)
		log.Printf("status=%d trace_id=%s url=%s referrer=%s method=%s route_param_names=%v, route_param_values=%v",
			code,
			rid,
			c.Path(),
			c.Request().Referer(),
			c.Request().Method,
			c.ParamNames(),
			c.ParamValues(),
		)
	}
	if option.LogFn != nil {
		logFn = option.LogFn
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}
			err := next(c)
			logFn(c)
			return err
		}
	}
}
