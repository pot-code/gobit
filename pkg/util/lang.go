package util

import (
	"net/http"

	"golang.org/x/text/language"
)

func ParseLangFromHttpRequest(r *http.Request, matcher language.Matcher) string {
	lang := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(matcher, lang)
	return GetLangFromTag(tag)
}

func GetLangFromTag(tag language.Tag) string {
	base, _ := tag.Base()
	return base.String()
}
