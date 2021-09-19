package middleware

import (
	"github.com/labstack/echo/v4"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/api"
	"github.com/pot-code/gobit/pkg/validate"
)

var DefaultPaginationEchoKey = "pagination"

type PaginationOption struct {
	EchoKey string
	LangKey gobit.AppContextKey
}

func CursorPagination(v *validate.ValidatorV10, option PaginationOption) echo.MiddlewareFunc {
	key := DefaultPaginationEchoKey

	if option.EchoKey != "" {
		key = option.EchoKey
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
			c.Set(key, pagination)
			return next(c)
		}
	}
}
