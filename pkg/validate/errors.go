package validate

import (
	"fmt"
	"strings"
)

// FieldError field error to be nested by other errors
type FieldError struct {
	Domain string `json:"domain"`
	Reason string `json:"reason"`
}

// NewFieldError create new field error
func NewFieldError(domain string, reason string) *FieldError {
	return &FieldError{domain, reason}
}

func (fe FieldError) Error() string {
	return fmt.Sprintf("[FieldError]%s: %s", fe.Domain, fe.Reason)
}

type ValidationError []*FieldError

// NewFieldError create new field error
func NewValidationError() *ValidationError {
	return new(ValidationError)
}

func (ve *ValidationError) AddFieldError(fe *FieldError) {
	*ve = append(*ve, fe)
}

func (ve ValidationError) Error() string {
	var msg []string
	for _, fe := range ve {
		msg = append(msg, fe.Error())
	}
	return strings.Join(msg, "\n")
}
