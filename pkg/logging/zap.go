package logging

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var DefaultLoggingContextKey = gobit.AppContextKey("logger")

func NewEcsConfig() zapcore.EncoderConfig {
	ec := zap.NewProductionEncoderConfig()

	ec.EncodeTime = zapcore.TimeEncoder(zapcore.ISO8601TimeEncoder)
	ec.TimeKey = "@timestamp"
	ec.MessageKey = "message"
	ec.LevelKey = "log.level"
	ec.CallerKey = "log.origin.file.line"
	ec.StacktraceKey = "error.stack_trace"
	return ec
}

func NewFileSyncer(p string) (zapcore.WriteSyncer, error) {
	fd, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return fd, err
}

var zapLoggingLevelMap map[string]zapcore.Level = map[string]zapcore.Level{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"fatal":  zap.FatalLevel,
	"panic":  zap.PanicLevel,
	"dpanic": zap.DPanicLevel,
}

func GetLevelEnabler(l string) (zapcore.LevelEnabler, error) {
	level, ok := zapLoggingLevelMap[strings.ToLower(l)]
	if !ok {
		return nil, fmt.Errorf("unsupported logging level '%s'", l)
	}
	return zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= level
	}), nil
}

// InjectContext set logger into target context
func InjectContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, DefaultLoggingContextKey, logger)
}

// ExtractFromContext try to extract logger from context
func ExtractFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(DefaultLoggingContextKey).(*zap.Logger)
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
