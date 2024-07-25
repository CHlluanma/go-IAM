package log

import (
	"context"
	"go.uber.org/zap"
)

type key int

const (
	logContextKey key = iota
)

func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

func (l *ZapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return &ZapLogger{zapL: logger.(*zap.Logger)}
		}
	}
	return WithName("Unknown-Context")
}
