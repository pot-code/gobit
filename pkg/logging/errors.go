package logging

import (
	"github.com/pkg/errors"
	"github.com/pot-code/gobit/pkg/util"
	"go.uber.org/zap/zapcore"
)

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
	if ste, ok := err.(util.StackTracer); ok {
		trace := util.GetVerboseStackTrace(te.depth, ste)
		enc.AddString("stack_trace", trace)

		cause := errors.Cause(err)
		enc.AddString("message", cause.Error())
	} else {
		enc.AddString("message", err.Error())
	}
	return nil
}
