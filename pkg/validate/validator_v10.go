package validate

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pot-code/gobit/pkg/util"
	"golang.org/x/text/language"
)

type ValidatorOption interface {
	apply(*ValidatorV10)
}

type optionFunc func(*ValidatorV10)

func (o optionFunc) apply(v *ValidatorV10) {
	o(v)
}

type LocaleRegister func(v *validator.Validate, trans ut.Translator) error

func AddLocale(locale language.Tag, translator locales.Translator, reg LocaleRegister) ValidatorOption {
	return optionFunc(func(vv *ValidatorV10) {
		vv.translators = append(vv.translators, translator)
		vv.locales[util.GetLangFromTag(locale)] = reg
		vv.tags = append(vv.tags, locale)
	})
}

// ValidatorV10 Validator implementation using go-playground
type ValidatorV10 struct {
	v           *validator.Validate
	uni         *ut.UniversalTranslator
	locales     map[string]LocaleRegister
	tags        []language.Tag
	matcher     language.Matcher
	translators []locales.Translator
}

// NewValidator create a new Validator
func NewValidator(options ...ValidatorOption) *ValidatorV10 {
	options = append(options, AddLocale(language.English, en.New(), en_translations.RegisterDefaultTranslations))

	validate := &ValidatorV10{
		v:           validator.New(),
		locales:     make(map[string]LocaleRegister),
		translators: []locales.Translator{},
	}
	for _, o := range options {
		o.apply(validate)
	}

	uni := ut.New(en.New(), validate.translators...)
	validate.uni = uni
	validate.matcher = language.NewMatcher(validate.tags)
	for locale, reg := range validate.locales {
		t, _ := uni.GetTranslator(locale)
		reg(validate.v, t)
	}

	validate.v.RegisterTagNameFunc(func(sf reflect.StructField) string {
		name := sf.Tag.Get("json")
		if name != "-" && name != "" {
			sg := strings.Split(name, ",")
			return sg[0]
		}
		name = sf.Tag.Get("yaml")
		if name != "-" && name != "" {
			sg := strings.Split(name, ",")
			return sg[0]
		}
		return sf.Name
	})

	return validate
}

func (gv ValidatorV10) Struct(s interface{}) *ValidationErrors {
	result := newValidationErrors()
	validate := gv.v
	trans := gv.uni.GetFallback()
	if err := validate.Struct(s); err != nil {
		for _, fe := range err.(validator.ValidationErrors) {
			result.addFieldError(newFieldError(fe.Translate(trans), fe))
		}
		return result
	}
	return nil
}

func (gv ValidatorV10) translate(ves *ValidationErrors, translator ut.Translator) *ValidationErrors {
	errors := ValidationErrors(make([]*FieldError, len(*ves)))
	for i, fe := range *ves {
		errors[i] = fe.Translate(translator)
	}
	return &errors
}

func (gv ValidatorV10) Translate(ves *ValidationErrors, locale string) *ValidationErrors {
	translator, _ := gv.uni.GetTranslator(locale)
	return gv.translate(ves, translator)
}

func (gv ValidatorV10) TranslateWithHttpRequest(ves *ValidationErrors, request *http.Request) *ValidationErrors {
	lang := util.ParseLangFromHttpRequest(request, gv.matcher)
	translator, _ := gv.uni.FindTranslator(lang)
	return gv.translate(ves, translator)
}

func (gv ValidatorV10) Email(name, v string) *ValidationErrors {
	return gv.validateWithTag(v, name, "email")
}

func (gv ValidatorV10) Required(name string, v interface{}) *ValidationErrors {
	return gv.validateWithTag(v, name, "required")
}

func (gv ValidatorV10) Tags(name string, v interface{}, tags []string) *ValidationErrors {
	tag := strings.Join(tags, ",")
	return gv.validateWithTag(v, name, tag)
}

func (gv ValidatorV10) validateWithTag(v interface{}, name, tag string) *ValidationErrors {
	result := newValidationErrors()
	validate := gv.v
	trans := gv.uni.GetFallback()
	if err := validate.Var(v, tag); err != nil {
		for _, fe := range err.(validator.ValidationErrors) {
			result.addFieldError(newFieldError(fe.Translate(trans), fe))

		}
		return result
	}
	return nil
}
