package context

import "context"

type AppContextKey string

// GetContextValue return value, found
func Value(ctx context.Context, key string) (interface{}, bool) {
	v := ctx.Value(AppContextKey(key))
	if v == nil {
		return nil, false
	}
	return v, true
}

func WithContext(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, AppContextKey(key), value)
}
