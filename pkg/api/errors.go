package api

import (
	"errors"

	"github.com/pot-code/gobit/pkg/validate"
)

// i18n error
var (
	ErrFailedBinding  = errors.New("errors.bind")
	ErrFailedValidate = errors.New("errors.validate")
	ErrInternalError  = errors.New("errors.internal")
)

// RESTStandardError response error
type RESTStandardError struct {
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}

func NewRESTStandardError(msg string) *RESTStandardError {
	return &RESTStandardError{
		Message: msg,
	}
}

func (re RESTStandardError) Error() string {
	return re.Message
}

func (re RESTStandardError) SetTraceID(id string) RESTStandardError {
	re.TraceID = id
	return re
}

// RESTValidationError standard validation error
type RESTValidationError struct {
	RESTStandardError
	Errors *validate.ValidationError `json:"errors"`
}

func NewRESTValidationError(msg string, ve *validate.ValidationError) *RESTValidationError {
	return &RESTValidationError{
		RESTStandardError: RESTStandardError{
			Message: msg,
		},
		Errors: ve,
	}
}

func (rve RESTValidationError) Error() string {
	return rve.Message
}

func (rve RESTValidationError) SetTraceID(id string) RESTValidationError {
	rve.RESTStandardError.TraceID = id
	return rve
}

// RESTValidationError standard validation error
type RESTBindingError struct {
	RESTStandardError
	Reason interface{} `json:"reason,omitempty"`
}

func NewRESTBindingError(msg string, reason interface{}) *RESTBindingError {
	return &RESTBindingError{
		RESTStandardError: RESTStandardError{
			Message: msg,
		},
		Reason: reason,
	}
}

func (rbe RESTBindingError) Error() string {
	return rbe.Message
}
