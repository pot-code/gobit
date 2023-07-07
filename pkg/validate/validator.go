package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/text/language"
)

type LocaleRegister func(v *validator.Validate, trans ut.Translator) error

type ValidatorBuilderSchema struct {
	localeMap     map[string]LocaleRegister
	localeTrans   []locales.Translator
	defaultLocale language.Tag
}

func NewValidator() *ValidatorBuilderSchema {
	localeTrans := []locales.Translator{en.New()}
	localeMap := map[string]LocaleRegister{
		getLangTagString(language.English): func(v *validator.Validate, utt ut.Translator) error {
			return en_translations.RegisterDefaultTranslations(v, utt)
		},
	}
	return &ValidatorBuilderSchema{
		localeMap:     localeMap,
		localeTrans:   localeTrans,
		defaultLocale: language.English,
	}
}

func (s *ValidatorBuilderSchema) RegisterLocale(locale language.Tag, localeTranslator locales.Translator, r LocaleRegister) *ValidatorBuilderSchema {
	s.localeTrans = append(s.localeTrans, localeTranslator)
	s.localeMap[getLangTagString(locale)] = r
	return s
}

func (s *ValidatorBuilderSchema) DefaultLocale(locale language.Tag) *ValidatorBuilderSchema {
	if locale != s.defaultLocale {
		s.defaultLocale = locale
	}
	return s
}

func (s *ValidatorBuilderSchema) Build() *Validator {
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
		trans, _ := utt.GetTranslator(locale)
		err := register(v, trans)
		if err != nil {
			panic(fmt.Errorf("failed to register locale %s for validator: %w", locale, err))
		}
	}

	return &Validator{
		v:          v,
		translator: utt.GetFallback(),
		utt:        utt,
	}
}

type Validator struct {
	v          *validator.Validate
	utt        *ut.UniversalTranslator
	translator ut.Translator
}

func (vv Validator) Struct(s interface{}) ValidationError {
	if err := vv.v.Struct(s); err != nil {
		return FromValidatorErrors(err.(validator.ValidationErrors), vv.translator)
	}
	return nil
}

func (vv Validator) Var(rule, name string, v interface{}) ValidationError {
	if err := vv.v.Var(v, rule); err != nil {
		return FromValidatorErrors(err.(validator.ValidationErrors), vv.translator)

	}
	return nil
}

func getLangTagString(tag language.Tag) string {
	base, _ := tag.Base()
	return base.String()
}
