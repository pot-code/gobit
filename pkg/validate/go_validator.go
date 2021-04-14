package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// GoValidatorV10 Validator implementation using go-playground
type GoValidatorV10 struct {
	backend *validator.Validate
	uni     *ut.UniversalTranslator
}

// NewValidator create a new Validator
func NewValidator() *GoValidatorV10 {
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, en, zh)

	validate := validator.New()
	ent, _ := uni.GetTranslator("en")
	zht, _ := uni.GetTranslator("zh")
	en_translations.RegisterDefaultTranslations(validate, ent)
	zh_translations.RegisterDefaultTranslations(validate, zht)
	validate.RegisterTagNameFunc(func(sf reflect.StructField) string {
		name := sf.Tag.Get("json")
		if name != "-" && name != "" {
			return name
		}
		name = sf.Tag.Get("yaml")
		if name != "-" && name != "" {
			return name
		}
		return sf.Name
	})
	return &GoValidatorV10{
		backend: validate,
		uni:     uni,
	}
}

// Struct validate struct
func (gv GoValidatorV10) Struct(s interface{}, lang string) *ValidationError {
	result := NewValidationError()
	validate := gv.backend
	trans, _ := gv.uni.GetTranslator(lang)
	if err := validate.Struct(s); err != nil {
		for _, item := range err.(validator.ValidationErrors) {
			result.AddFieldError(NewFieldError(item.Field(), item.Translate(trans)))
		}
		return result
	}
	return nil
}

func (gv GoValidatorV10) Email(name, v string, lang string) *ValidationError {
	return gv.validateWithTag(v, name, "email", lang)
}

func (gv GoValidatorV10) Required(name string, v interface{}, lang string) *ValidationError {
	return gv.validateWithTag(v, name, "required", lang)
}

func (gv GoValidatorV10) Tags(name string, v interface{}, tags []string, lang string) *ValidationError {
	tag := strings.Join(tags, ",")
	return gv.validateWithTag(v, name, tag, lang)
}

func (gv GoValidatorV10) validateWithTag(v interface{}, name, tag, lang string) *ValidationError {
	result := NewValidationError()
	validate := gv.backend
	trans, _ := gv.uni.GetTranslator(lang)
	if err := validate.Var(v, tag); err != nil {
		for _, item := range err.(validator.ValidationErrors) {
			result.AddFieldError(NewFieldError(name, item.Translate(trans)))
		}
		return result
	}
	return nil
}
