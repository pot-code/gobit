package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/text/language"
)

type LocaleRegister func(v *validator.Validate, t ut.Translator) error

type ValidatorBuilder struct {
	localeMap     map[string]LocaleRegister
	localeTrans   []locales.Translator
	defaultLocale language.Tag
}

func NewBuilder() *ValidatorBuilder {
	localeTrans := []locales.Translator{en.New()}
	localeMap := map[string]LocaleRegister{
		getLangTagString(language.English): func(v *validator.Validate, utt ut.Translator) error {
			return en_translations.RegisterDefaultTranslations(v, utt)
		},
	}
	return &ValidatorBuilder{
		localeMap:     localeMap,
		localeTrans:   localeTrans,
		defaultLocale: language.English,
	}
}

func (s *ValidatorBuilder) RegisterLocale(locale language.Tag, localeTranslator locales.Translator, r LocaleRegister) *ValidatorBuilder {
	s.localeTrans = append(s.localeTrans, localeTranslator)
	s.localeMap[getLangTagString(locale)] = r
	return s
}

func (s *ValidatorBuilder) DefaultLocale(locale language.Tag) *ValidatorBuilder {
	if locale != s.defaultLocale {
		s.defaultLocale = locale
	}
	return s
}

func (s *ValidatorBuilder) Build() Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(sf reflect.StructField) string {
		name := sf.Tag.Get("json")
		if name != "-" && name != "" {
			s := strings.Split(name, ",")
			return s[0]
		}
		name = sf.Tag.Get("yaml")
		if name != "-" && name != "" {
			s := strings.Split(name, ",")
			return s[0]
		}
		return sf.Name
	})

	utt := ut.New(s.localeTrans[0], s.localeTrans...)
	for locale, register := range s.localeMap {
		t, _ := utt.GetTranslator(locale)
		err := register(v, t)
		if err != nil {
			panic(fmt.Errorf("failed to register locale '%s' for validator: %w", locale, err))
		}
	}

	locale := getLangTagString(s.defaultLocale)
	translator, ok := utt.GetTranslator(locale)
	if !ok {
		panic(fmt.Errorf("failed to get default translator for locale '%s'", locale))
	}

	return &V10Validator{
		v:          v,
		translator: translator,
		utt:        utt,
	}
}

func getLangTagString(tag language.Tag) string {
	base, _ := tag.Base()
	return base.String()
}
