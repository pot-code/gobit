package middleware

import (
	"github.com/labstack/echo/v4"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/api"
	"github.com/pot-code/gobit/pkg/util"
	"github.com/pot-code/gobit/pkg/validate"
)

type PaginationOption struct {
	EchoKey string
	LangKey gobit.AppContextKey
}

func CursorPagination(v *validate.GoValidatorV10, options ...PaginationOption) echo.MiddlewareFunc {
	key := gobit.DefaultPaginationEchoKey
	langKey := gobit.DefaultLangContextKey
	if len(options) > 0 {
		option := options[0]
		if option.EchoKey != "" {
			key = option.EchoKey
		}
		if option.LangKey != "" {
			langKey = option.LangKey
		}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lang := c.Request().Context().Value(langKey).(string)
			pagination := new(util.CursorPaginationReq)
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
