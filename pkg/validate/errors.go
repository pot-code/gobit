package validate

import (
	"strings"
)

type ValidationError []string

// NewFieldError create new field error
func NewValidationError() *ValidationError {
	return new(ValidationError)
}

func (ve *ValidationError) AddFieldError(e string) {
	*ve = append(*ve, e)
}

func (ve ValidationError) Error() string {
	var msg []string
	for _, e := range ve {
		msg = append(msg, e)
	}
	return strings.Join(msg, "\n")
}
