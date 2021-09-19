package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/api"
	"github.com/pot-code/gobit/pkg/context"
	"golang.org/x/text/language"
)

const LangKey context.AppContextKey = "lang"

type ParseAcceptLanguageOption struct {
	ContextKey string
	Lang       []language.Tag
}

func ParseAcceptLanguage(option ParseAcceptLanguageOption) echo.MiddlewareFunc {
	key := LangKey
	tags := []language.Tag{
		language.English,
		language.Chinese,
	}
	if option.ContextKey != "" {
		key = context.AppContextKey(option.ContextKey)
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
