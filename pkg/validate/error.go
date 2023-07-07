package validate

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationResult struct {
	field  string
	reason string
}

func NewValidationResult(field string, reason string) *ValidationResult {
	return &ValidationResult{field, strings.ReplaceAll(reason, field, "")}
}

func (vr *ValidationResult) Field() string {
	return vr.field
}

func (vr *ValidationResult) Reason() string {
	return vr.reason
}

func (vr *ValidationResult) String() string {
	return fmt.Sprintf("%s %s", vr.field, vr.reason)
}

type ValidationError []*ValidationResult

func FromValidatorErrors(err validator.ValidationErrors, t ut.Translator) ValidationError {
	var ve []*ValidationResult
	for _, err := range err {
		reason := err.Translate(t)
		ve = append(ve, NewValidationResult(err.Field(), reason))
	}
	return ValidationError(ve)
}

func FromVarValidatorErrors(name string, err validator.ValidationErrors, t ut.Translator) ValidationError {
	var ve []*ValidationResult
	for _, err := range err {
		reason := err.Translate(t)
		ve = append(ve, NewValidationResult(name, reason))
	}
	return ValidationError(ve)
}

func (ve ValidationError) Error() string {
	var msg []string
	for _, e := range ve {
		msg = append(msg, e.String())
	}
	return strings.Join(msg, "\n")
}
