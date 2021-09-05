package logging

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	gobit "github.com/pot-code/gobit/pkg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var DefaultLoggingContextKey = gobit.AppContextKey("logger")

// LoggerConfig options used in creating zap logger
type LoggerConfig struct {
	FilePath  string // log file path
	Level     string // global logging level
	AddCaller bool
}

// NewLogger returns a zap logger
func NewLogger(cfg *LoggerConfig) (*zap.Logger, error) {
	var (
		core zapcore.Core
		err  error
	)
	core, err = NewProductionCore(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger core::%w", err)
	}

	logger := zap.New(core)
	if cfg.AddCaller {
		logger.WithOptions(zap.AddCaller())
	}
	return logger, nil
}

func NewProductionCore(cfg *LoggerConfig) (zapcore.Core, error) {
	logEnabler, err := getLevelEnabler(cfg.Level)
	if err != nil {
		return nil, err
	}

	ecsEncoderConfig := zap.NewProductionEncoderConfig()
	ecsEncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
	})
	ecsEncoderConfig.TimeKey = "@timestamp"
	ecsEncoderConfig.MessageKey = "message"
	ecsEncoderConfig.LevelKey = "log.level"
	ecsEncoderConfig.CallerKey = "log.position"
	ecsEncoderConfig.StacktraceKey = "error.stack_trace"
	ecsEncoder := zapcore.NewJSONEncoder(ecsEncoderConfig)

	if cfg.FilePath != "" {
		elkOutput, err := getFileSyncer(cfg.FilePath)
		return zapcore.NewCore(ecsEncoder, elkOutput, logEnabler), err
	}
	return zapcore.NewCore(ecsEncoder, os.Stderr, logEnabler), nil
}

func getFileSyncer(p string) (zapcore.WriteSyncer, error) {
	fd, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return fd, err
}

func getZapLoggingLevel(l string) (level zapcore.Level) {
	switch strings.ToLower(l) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return
}

func getLevelEnabler(l string) (zapcore.LevelEnabler, error) {
	level := getZapLoggingLevel(l)
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
