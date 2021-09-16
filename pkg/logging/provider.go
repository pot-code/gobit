package logging

import (
	"context"
	"log"
	"os"

	"github.com/pot-code/gobit/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLoggerProvider(lc *LoggingConfig, lm *util.LifecycleManager) *zap.Logger {
	if lc == nil {
		panic("LoggingConfig is nil")
	}

	ec := NewEcsConfig()
	enabler, err := GetLevelEnabler(lc.Level)
	util.HandlePanicError("failed to create logger", err)

	var zc zapcore.Encoder
	if lc.Format == "json" {
		zc = zapcore.NewJSONEncoder(ec)
	} else {
		zc = zapcore.NewConsoleEncoder(ec)
	}

	var out zapcore.WriteSyncer
	p := lc.FilePath
	if p == "" {
		out = os.Stderr
	} else {
		out, err = NewFileSyncer(p)
		util.HandlePanicError("failed to create logger", err)
	}

	core := zapcore.NewCore(zc, out, enabler)
	logger := zap.New(core)

	lm.OnExit(func(ctx context.Context) {
		log.Println("[zap.Logger] sync logger")
		logger.Sync()
	})

	return logger
}
