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

func CursorPagination(v *validate.GoValidatorV10, option PaginationOption) echo.MiddlewareFunc {
	key := DefaultPaginationEchoKey
	langKey := DefaultLangContextKey

	if option.EchoKey != "" {
		key = option.EchoKey
	}
	if option.LangKey != "" {
		langKey = option.LangKey
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lang := c.Request().Context().Value(langKey).(string)
			pagination := new(api.CursorPaginationReq)
			if err := c.Bind(pagination); err != nil {
				return api.BadRequestResponse(c, err)
			}

			if err := v.Struct(pagination, lang); err != nil {
				return api.ValidateFailedResponse(c, err)
			}
			c.Set(key, pagination)
			return next(c)
		}
	}
}
