package logging

import (
	"context"

	gobit "github.com/pot-code/gobit/pkg"
	"go.uber.org/zap"
)

var DefaultLoggingContextKey = gobit.AppContextKey("logger")

// InjectContext set logger into target context
func InjectContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, DefaultLoggingContextKey, logger)
}

// ExtractFromContext try to extract logger from context
func ExtractFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(DefaultLoggingContextKey).(*zap.Logger)
}
