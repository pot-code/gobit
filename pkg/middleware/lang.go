package middleware

import (
	"github.com/labstack/echo/v4"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/util"
	"golang.org/x/text/language"
)

type ParseAcceptLanguageOption struct {
	ContextKey gobit.AppContextKey
}

func ParseAcceptLanguage(lang []language.Tag, options ...ParseAcceptLanguageOption) echo.MiddlewareFunc {
	key := gobit.DefaultLangContextKey
	matcher := language.NewMatcher(lang)
	if len(options) > 0 {
		option := options[0]
		if option.ContextKey != "" {
			key = option.ContextKey
		}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lang := c.Request().Header.Get("Accept-Language")
			tag, _ := language.MatchStrings(matcher, lang)
			base, _ := tag.Base()
			util.WithContextValue(c, key, base.String())
			return next(c)
		}
	}
}
