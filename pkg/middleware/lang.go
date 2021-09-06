package middleware

import (
	"github.com/labstack/echo/v4"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/api"
	"golang.org/x/text/language"
)

var DefaultLangContextKey = gobit.AppContextKey("lang")

type ParseAcceptLanguageOption struct {
	ContextKey gobit.AppContextKey
	Lang       []language.Tag
}

func ParseAcceptLanguage(option ParseAcceptLanguageOption) echo.MiddlewareFunc {
	key := DefaultLangContextKey
	tags := []language.Tag{
		language.English,
		language.Chinese,
	}
	if option.ContextKey != "" {
		key = option.ContextKey
	}
	if option.Lang != nil {
		tags = option.Lang
	}
	matcher := language.NewMatcher(tags)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lang := c.Request().Header.Get("Accept-Language")
			tag, _ := language.MatchStrings(matcher, lang)
			base, _ := tag.Base()
			api.WithContextValue(c, key, base.String())
			return next(c)
		}
	}
}
