package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/pot-code/gobit/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewEcsConfig() zapcore.EncoderConfig {
	ec := zap.NewProductionEncoderConfig()

	ec.EncodeTime = zapcore.TimeEncoder(zapcore.ISO8601TimeEncoder)
	ec.TimeKey = "@timestamp"
	ec.MessageKey = "message"
	ec.LevelKey = "log.level"
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

func GetEncoderByFormat(format string, cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
	switch format {
	case LogFormatJson:
		return zapcore.NewJSONEncoder(cfg), nil
	case LogFormatConsole:
		return zapcore.NewConsoleEncoder(cfg), nil
	}
	return nil, fmt.Errorf("unsupported log format '%s'", format)
}

type ZapStacktraceError struct {
	depth int
	err   error
}

// NewZapStacktraceError create wrapper object that implements the `MarshalLogObject` protocol
//
// depth: set stack trace depth if the error type supports it
//
// err: the error to be wrapped
func NewZapStacktraceError(err error, depth int) *ZapStacktraceError {
	return &ZapStacktraceError{depth, err}
}

func (te ZapStacktraceError) Unwrap() error {
	return te.err
}

func (te ZapStacktraceError) Error() string {
	if te.err != nil {
		return te.err.Error()
	}
	return "zap with stacktrace error"
}

func (te ZapStacktraceError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
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
