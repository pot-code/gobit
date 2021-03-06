package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/api"
)

// ErrorHandlingOption options for error handling
//
// the original error frame is always at the top.
type ErrorHandlingOption struct {
	Handler func(c echo.Context, err error) // error handler
}

// ErrorHandling auto recovery and handle the errors returned from handlers
func ErrorHandling(option ErrorHandlingOption) echo.MiddlewareFunc {
	handler := func(c echo.Context, e error) {
		msg := api.ErrInternalError.Error()
		c.JSON(http.StatusInternalServerError, api.NewRESTStandardError(msg))
	}

	if option.Handler != nil {
		handler = option.Handler
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				if v, ok := err.(*echo.HTTPError); ok {
					c.String(v.Code, v.Error())
				} else {
					handler(c, err)
				}
			}
			return nil
		}
	}
}
