package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/api"
	"github.com/pot-code/gobit/pkg/context"
	"github.com/pot-code/gobit/pkg/util"
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
			lang := util.ParseLangFromHttpRequest(c.Request(), matcher)
			api.WithContextValue(c, key, lang)
			return next(c)
		}
	}
}
