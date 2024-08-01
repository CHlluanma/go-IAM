package log

import (
	"context"
)

type key int

const (
	logContextKey key = iota
)

func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

// FromContext 函数根据上下文提取 Logger 实例
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		// 尝试从上下文中获取 Logger 实例，
		// 如果找到了就返回一个新的 zapLogger 实例，其内部指向找到的 zap.Logger 实例
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(*zapLogger)
		}
	}
	return WithName("Unknown-Context")
}
