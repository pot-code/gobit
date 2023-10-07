package validate

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(any) ValidationError
}

type V10Validator struct {
	v          *validator.Validate
	utt        *ut.UniversalTranslator
	translator ut.Translator
}

func (vv V10Validator) Struct(s any) ValidationError {
	if err := vv.v.Struct(s); err != nil {
		return FromValidatorErrors(err.(validator.ValidationErrors), vv.translator)
	}
	return nil
}
