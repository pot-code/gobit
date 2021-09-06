package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// AbortRequestOption option for AbortRequest
type AbortRequestOption struct {
	Timeout time.Duration // timeout for request, negative value mean never timeout
}

// AbortRequest handle request abortion
func AbortRequest(option AbortRequestOption) echo.MiddlewareFunc {
	timeout := 30 * time.Second
	if option.Timeout > 0 {
		timeout = option.Timeout
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if timeout > 0 {
				req := c.Request()
				ctx, cancel := context.WithTimeout(req.Context(), timeout)
				defer cancel()
				c.SetRequest(req.WithContext(ctx))
			}
			err := next(c)
			if errors.Is(err, context.DeadlineExceeded) {
				return echo.NewHTTPError(http.StatusRequestTimeout)
			} else if errors.Is(err, context.Canceled) { // if request is canceled
				return nil
			}
			return err
		}
	}
}
