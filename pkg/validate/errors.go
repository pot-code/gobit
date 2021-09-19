package validate

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field  string               `json:"field"`
	Reason string               `json:"reason"`
	err    validator.FieldError `json:"-"`
}

func newFieldError(msg string, err validator.FieldError) *FieldError {
	if err == nil {
		panic("err is nil")
	}
	return &FieldError{Reason: msg, err: err, Field: err.Field()}
}

func (fe *FieldError) Translate(translator ut.Translator) *FieldError {
	return &FieldError{err: fe.err, Field: fe.Field, Reason: fe.err.Translate(translator)}
}

func (fe *FieldError) Error() string {
	return fe.Reason
}

type ValidationErrors []*FieldError

// NewFieldError create new field error
func newValidationErrors() *ValidationErrors {
	return new(ValidationErrors)
}

func (ve *ValidationErrors) addFieldError(fe *FieldError) {
	*ve = append(*ve, fe)
}

func (ve ValidationErrors) Error() string {
	var msg []string
	for _, e := range ve {
		msg = append(msg, e.Error())
	}
	return strings.Join(msg, "\n")
}
