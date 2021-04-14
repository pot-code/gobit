package util

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/pot-code/gobit/pkg/validate"
	"go.uber.org/zap/zapcore"
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

func (re RESTStandardError) SetTraceID(traceID string) RESTStandardError {
	re.TraceID = traceID
	return re
}

// RESTValidationError standard validation error
type RESTValidationError struct {
	RESTStandardError
	Errors *validate.ValidationError `json:"errors"`
}

func NewRESTValidationError(msg string, internal *validate.ValidationError) *RESTValidationError {
	return &RESTValidationError{
		RESTStandardError: RESTStandardError{
			Message: msg,
		},
		Errors: internal,
	}
}

func (rve RESTValidationError) Error() string {
	return rve.Message
}

func (rve RESTValidationError) SetTraceID(traceID string) RESTValidationError {
	rve.RESTStandardError.TraceID = traceID
	return rve
}

// RESTValidationError standard validation error
type RESTBindingError struct {
	RESTStandardError
	Reason interface{}        `json:"reason,omitempty"`
	Errors *echo.BindingError `json:"errors,omitempty"`
}

func NewRESTBindingError(msg string, reason interface{}, err *echo.BindingError) *RESTBindingError {
	return &RESTBindingError{
		RESTStandardError: RESTStandardError{
			Message: msg,
		},
		Reason: reason,
		Errors: err,
	}
}

func (rbe RESTBindingError) Error() string {
	return rbe.Message
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type ZapErrorWrapper struct {
	depth int
	err   error
}

// NewZapErrorWrapper create wrapper object that implements the `MarshalLogObject` protocol
//
// depth: set stack trace depth if the error type supports it
//
// err: the error to be wrapped
func NewZapErrorWrapper(err error, depth int) *ZapErrorWrapper {
	return &ZapErrorWrapper{depth, err}
}

func (te ZapErrorWrapper) Unwrap() error {
	return te.err
}

func (te ZapErrorWrapper) Error() string {
	if te.err != nil {
		return te.err.Error()
	}
	return "traceable error"
}

func (te ZapErrorWrapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	err := te.err
	if ste, ok := err.(StackTracer); ok {
		trace := GetVerboseStackTrace(te.depth, ste)
		enc.AddString("stack_trace", trace)

		cause := errors.Cause(err)
		enc.AddString("message", cause.Error())
	} else {
		enc.AddString("message", err.Error())
	}
	return nil
}
