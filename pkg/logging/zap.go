package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig options used in creating zap logger
type LoggerConfig struct {
	FilePath string // log file path
	Level    string // global logging level
	Env      string // app environment
	AppID    string
}

// NewLogger returns a zap logger
func NewLogger(cfg *LoggerConfig) (*zap.Logger, error) {
	var (
		core zapcore.Core
		err  error
	)
	core, err = createProductionLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger core::%w", err)
	}

	// logger := zap.New(core, zap.AddCaller())
	logger := zap.New(core)
	return logger.With(
		zap.String("labels.application", cfg.AppID),
		zap.String("labels.env", string(cfg.Env)),
	), nil
}

func getZapLoggingLevel(level string) (zlevel zapcore.Level) {
	switch strings.ToLower(level) {
	case "debug":
		zlevel = zap.DebugLevel
	case "info":
		zlevel = zap.InfoLevel
	case "warn":
		zlevel = zap.WarnLevel
	case "error":
		zlevel = zap.ErrorLevel
	case "fatal":
		zlevel = zap.FatalLevel
	default:
		log.Fatal(fmt.Errorf("unknown logging level: %s", level))
	}
	return
}

func createProductionLogger(cfg *LoggerConfig) (zapcore.Core, error) {
	logEnabler := getLevelEnabler(cfg)
	ecsEncoderConfig := zap.NewProductionEncoderConfig()
	ecsEncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
	})
	ecsEncoderConfig.TimeKey = "@timestamp"
	ecsEncoderConfig.MessageKey = "message"
	ecsEncoderConfig.LevelKey = "log.level"
	// ecsEncoderConfig.CallerKey = "log.position"
	ecsEncoderConfig.StacktraceKey = "error.stack_trace"
	ecsEncoder := zapcore.NewJSONEncoder(ecsEncoderConfig)

	if cfg.FilePath != "" {
		elkOutput, err := getFileSyncer(cfg)
		return zapcore.NewCore(ecsEncoder, elkOutput, logEnabler), err
	}
	return zapcore.NewCore(ecsEncoder, os.Stderr, logEnabler), nil
}

func getFileSyncer(cfg *LoggerConfig) (zapcore.WriteSyncer, error) {
	fd, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return fd, err
}

func getLevelEnabler(cfg *LoggerConfig) zapcore.LevelEnabler {
	level := getZapLoggingLevel(cfg.Level)
	return zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= level
	})
}

// InjectContext set logger into target context
func InjectContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, gobit.DefaultLoggingContextKey, logger)
}

// ExtractFromContext try to extract logger from context
func ExtractFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(gobit.DefaultLoggingContextKey).(*zap.Logger)
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
