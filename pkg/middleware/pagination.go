package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/api"
	"github.com/pot-code/gobit/pkg/context"
	"github.com/pot-code/gobit/pkg/validate"
)

const PaginationKey context.AppContextKey = "pagination"

type PaginationOption struct {
	ContextKey context.AppContextKey
}

func CursorPagination(v *validate.ValidatorV10, option PaginationOption) echo.MiddlewareFunc {
	key := PaginationKey

	if option.ContextKey != "" {
		key = option.ContextKey
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			pagination := new(api.CursorPaginationReq)
			if err := c.Bind(pagination); err != nil {
				return api.BadRequestResponse(c, err)
			}

			if err := v.Struct(pagination); err != nil {
				return api.ValidateFailedResponse(c, err)
			}
			api.WithContextValue(c, key, pagination)
			return next(c)
		}
	}
}
