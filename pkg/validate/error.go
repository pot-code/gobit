package validate

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ValidationResult struct {
	field  string
	reason string
}

func NewValidationResult(field string, reason string) *ValidationResult {
	return &ValidationResult{field, reason}
}

func (vr *ValidationResult) Field() string {
	return vr.field
}

func (vr *ValidationResult) Reason() string {
	return vr.reason
}

func (vr *ValidationResult) String() string {
	return fmt.Sprintf("%s: %s", vr.field, vr.reason)
}

type ValidationError []*ValidationResult

func FromValidatorErrors(errs validator.ValidationErrors, translator ut.Translator) ValidationError {
	var ve []*ValidationResult
	for _, err := range errs {
		reason := err.Translate(translator)
		ve = append(ve, NewValidationResult(err.Field(), reason))
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
