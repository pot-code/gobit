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
	ec := NewEcsConfig()
	enabler, err := GetLevelEnabler(lc.Level)
	util.HandleFatalError("failed to create logger", err)

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
		if err != nil {
			log.Fatalf("failed to create logger: %s", err)
		}
	}

	core := zapcore.NewCore(zc, out, enabler)
	logger := zap.New(core)

	lm.OnExit(func(ctx context.Context) {
		log.Println("sync zap logger")
		logger.Sync()
	})

	return logger
}
