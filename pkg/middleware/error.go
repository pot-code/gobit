package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/db"
	"github.com/pot-code/gobit/pkg/util"
	"go.uber.org/zap"
)

// ErrorHandlingOption options for error handling
//
// the original error frame is always at the top.
type ErrorHandlingOption struct {
	Handler    func(c echo.Context, err error) // error handler
	StackDepth int                             // maximum frames to log counting from the top
}

// ErrorHandling auto recovery and handle the errors returned from handlers
//
// the default depth is infinite (-1)
func ErrorHandling(logger *zap.Logger, options ...ErrorHandlingOption) echo.MiddlewareFunc {
	depth := -1

	// default error handler
	handler := func(c echo.Context, e error) {
		traceID := c.Response().Header().Get(echo.HeaderXRequestID)

		cause := errors.Cause(e)
		msg := gobit.ErrInternalError.Error()
		if err, ok := cause.(*db.SqlDBError); ok {
			logger.Error(e.Error(), zap.String("trace.id", traceID), zap.Object("db", err), zap.Object("error", util.NewZapErrorWrapper(e, depth)))
			msg = gobit.ErrDBError.Error()
		} else {
			logger.Error(e.Error(), zap.String("trace.id", traceID), zap.Object("error", util.NewZapErrorWrapper(e, depth)))
		}
		c.JSON(http.StatusInternalServerError,
			util.NewRESTStandardError(msg).SetTraceID(traceID),
		)
	}

	if len(options) > 0 {
		option := options[0]
		if option.Handler != nil {
			handler = option.Handler
		}
		if option.StackDepth > 0 {
			depth = option.StackDepth
		}
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
