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

func createProductionLogger(cfg *LoggerConfig) (zapcore.Core, error) {
	logEnabler, err := getLevelEnabler(cfg)
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

func getZapLoggingLevel(level string) (zlevel zapcore.Level, err error) {
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
		err = fmt.Errorf("unknown logging level: %s", level)
	}
	return
}

func getLevelEnabler(cfg *LoggerConfig) (zapcore.LevelEnabler, error) {
	level, err := getZapLoggingLevel(cfg.Level)
	if err != nil {
		return nil, err
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
