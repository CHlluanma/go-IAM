package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InfoLogger 定义了日志信息记录的接口。
// 它提供了多种方式来记录日志信息，包括简单信息、格式化信息和自定义字段的信息.
type InfoLogger interface {
	Info(msg string, fields ...Field)
	Infof(format string, v ...any)
	Infow(msg string, keysAndValues ...any)

	// Enabled 返回一个布尔值，指示当前是否启用了日志记录。
	// 这个方法可以用于在记录日志之前检查日志级别是否允许。
	Enabled() bool
}

// Logger 定义了日志记录器的接口，表示能够记录消息，包括错误和非错误.
type Logger interface {
	InfoLogger
	Debug(msg string, fields ...Field)
	Debugf(format string, v ...any)
	Debugw(msg string, keysAndValues ...any)
	Warn(msg string, fields ...Field)
	Warnf(format string, v ...any)
	Warnw(msg string, keysAndValues ...any)
	Error(msg string, fields ...Field)
	Errorf(format string, v ...any)
	Errorw(msg string, keysAndValues ...any)
	Panic(msg string, fields ...Field)
	Panicf(format string, v ...any)
	Panicw(msg string, keysAndValues ...any)
	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...any)
	Fatalw(msg string, keysAndValues ...any)

	// V 返回一个具有指定日志级别的信息记录器。
	// 这个方法允许动态调整日志记录的详细程度，以根据需要查看更具体或更高级别的日志信息。
	// 返回的 InfoLogger 是一个具有指定日志级别并可用于记录信息的接口。
	V(level Level) InfoLogger
	Write(p []byte) (n int, err error)

	// WithValues 返回一个新的日志记录器，其中包含指定的键值对。
	WithValues(keysAndValues ...any) Logger

	// WithName 返回一个新的日志记录器，其中包含指定的名称。
	WithName(name string) Logger

	// WithContext  返回一个新的上下文，其中包含指定的日志记录器。
	WithContext(ctx context.Context) context.Context

	// Flush 将缓冲区中的日志数据写入输出。
	Flush() error
}

// noopInfoLogger 是一个日志记录器的实现，它不对任何信息进行记录。
// 这种实现用于在不需要日志记录或者需要关闭日志记录的场景下。
// 它充当一个"空操作"的实现，即执行操作但不产生任何效果。
type noopInfoLogger struct{}

func (n *noopInfoLogger) Info(_ string, _ ...Field) {}

func (n *noopInfoLogger) Infof(_ string, _ ...any) {}

func (n *noopInfoLogger) Infow(_ string, _ ...any) {}

func (n *noopInfoLogger) Enabled() bool { return false }

var disableInfoLogger = &noopInfoLogger{}

type infoLogger struct {
	l     *zap.Logger
	level zapcore.Level
}

func (l *infoLogger) Info(msg string, fields ...Field) {
	if checkedEntry := l.l.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

func (l *infoLogger) Infof(format string, v ...any) {
	if checkedEntry := l.l.Check(l.level, fmt.Sprintf(format, v...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

func (l *infoLogger) Infow(msg string, keysAndValues ...any) {
	if checkedEntry := l.l.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.l, keysAndValues)...)
	}
}

func (l *infoLogger) Enabled() bool {
	return true
}

// handleFields converts a bunch of arbitrary key-value pairs into Zap fields.  It takes
// additional pre-converted Zap fields, for use with automatically attached fields, like
// `error`.
func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return additional
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		// check just in case for strongly-typed Zap fields, which is illegal (since
		// it breaks implementation agnosticism), so we can give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))

			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}

type Level = zapcore.Level

type zapLogger struct {
	infoLogger
	zapL *zap.Logger
	// al   *zap.AtomicLevel
}

// Enabled implements Logger.
// Subtle: this method shadows the method (infoLogger).Enabled of zapLogger.infoLogger.
func (l *zapLogger) Enabled() bool {
	panic("unimplemented")
}

// Info implements Logger.
// Subtle: this method shadows the method (infoLogger).Info of zapLogger.infoLogger.
func (l *zapLogger) Info(msg string, fields ...zapcore.Field) {
	l.zapL.Info(msg, fields...)
}

// Infof implements Logger.
// Subtle: this method shadows the method (infoLogger).Infof of zapLogger.infoLogger.
func (l *zapLogger) Infof(format string, v ...any) {
	l.zapL.Sugar().Infof(format, v...)
}

// Infow implements Logger.
// Subtle: this method shadows the method (infoLogger).Infow of zapLogger.infoLogger.
func (l *zapLogger) Infow(msg string, keysAndValues ...any) {
	l.zapL.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zapL.Debug(msg, fields...)
}
func (l *zapLogger) Debugf(format string, v ...any) {
	l.zapL.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...any) {
	l.zapL.Sugar().Debugw(msg, keysAndValues...)
}
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zapL.Warn(msg, fields...)
}
func (l *zapLogger) Warnf(format string, v ...any) {
	l.zapL.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...any) {
	l.zapL.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zapL.Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, v ...any) {
	l.zapL.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...any) {
	l.zapL.Sugar().Errorw(msg, keysAndValues...)
}
func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.zapL.Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, v ...any) {
	l.zapL.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...any) {
	l.zapL.Sugar().Panicw(msg, keysAndValues...)
}
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zapL.Fatal(msg, fields...)
}
func (l *zapLogger) Fatalf(format string, v ...any) {
	l.zapL.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...any) {
	l.zapL.Sugar().Fatalw(msg, keysAndValues...)
}

// V 返回一个日志记录器，该记录器仅记录级别大于等于给定级别的日志。
func (l *zapLogger) V(level Level) InfoLogger {
	if l.zapL.Core().Enabled(level) {
		return &infoLogger{
			l:     l.zapL,
			level: level,
		}
	}
	return disableInfoLogger
}

// L 从 Context 中取出指定的keyValue， 作为上下文添加到日志输出中
func (l *zapLogger) L(ctx context.Context) *zapLogger {
	lg := l.clone()

	if requestID := ctx.Value("requestID"); requestID != nil {
		lg.zapL = lg.zapL.With(zap.Any("requestID", requestID))
	}
	if username := ctx.Value("username"); username != nil {
		lg.zapL = lg.zapL.With(zap.Any("username", username))
	}
	if watcherName := ctx.Value("watcher"); watcherName != nil {
		lg.zapL = lg.zapL.With(zap.Any("watcher", watcherName))
	}

	return lg
}

//nolint:predeclared
func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}

func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zapL.Info(string(p))

	return len(p), nil
}

func (l *zapLogger) WithValues(keysAndValues ...any) Logger {
	newLogger := l.zapL.With(handleFields(l.zapL, keysAndValues)...)

	return NewLoggerWithX(newLogger)
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapL.Named(name)

	return NewLoggerWithX(newLogger)
}

func (l *zapLogger) Flush() error {
	return l.zapL.Sync()
}

func (l *zapLogger) Sync() error {
	return l.zapL.Sync()
}

// TODO: 设置日志编码器 && 日志写入器
//func (l *zapLogger) SetEncoder(encoderType LogEncoder) {}
//func (l *zapLogger) SetWriter() zapcore.WriteSyncer {}
// func (l *zapLogger) SetLevel(level Level) {
// 	if l.al != nil {
// 		l.al.SetLevel(level)
// 	}
// }

var _ Logger = &zapLogger{}

// New 函数创建一个新的 zapLogger
func New(out io.Writer, level Level, opts ...Option) *zapLogger {
	if out == nil {
		out = os.Stdout
	}

	al := zap.NewAtomicLevelAt(level)
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(out),
		al,
	)
	return &zapLogger{
		zapL: zap.New(core, opts...),
		// al:   &al,
	}
}

func NewLoggerWithX(l *zap.Logger) Logger {
	return &zapLogger{
		zapL: l,
		infoLogger: infoLogger{
			l:     l,
			level: InfoLevel,
		},
	}
}

func ZapLogger() *zap.Logger {
	return std.zapL
}

var (
	std = New(os.Stderr, InfoLevel)
	mu  sync.Mutex
)

func Default() *zapLogger {
	return std
}

func ReplaceDefault(l *zapLogger) {
	std = l
}

// func SetLevel(level Level) {
// 	std.SetLevel(level)
// }

func Info(msg string, fields ...Field) {
	std.Info(msg, fields...)
}

func Infof(format string, v ...any) {
	std.Infof(format, v...)
}

func Infow(msg string, keysAndValues ...any) {
	std.Infow(msg, keysAndValues...)
}

func Debug(msg string, fields ...Field) {
	std.Debug(msg, fields...)
}

func Debugf(format string, v ...any) {
	std.Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...any) {
	std.Debugw(msg, keysAndValues...)
}

func Warn(msg string, fields ...Field) {
	std.Warn(msg, fields...)
}

func Warnf(format string, v ...any) {
	std.Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...any) {
	std.Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...Field) {
	std.Error(msg, fields...)
}

func Errorf(format string, v ...any) {
	std.Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...any) {
	std.Errorw(msg, keysAndValues...)
}

func Panic(msg string, fields ...Field) {
	std.Panic(msg, fields...)
}

func Panicf(format string, v ...any) {
	std.Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...any) {
	std.Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...Field) {
	std.Fatal(msg, fields...)
}

func Fatalf(format string, v ...any) {
	std.Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...any) {
	std.Fatalw(msg, keysAndValues...)
}

func Flush() error {
	return std.Flush()
}

func WithName(s string) Logger {
	return std.WithName(s)
}

func WithValues(keysAndValues ...any) Logger {
	return std.WithValues(keysAndValues...)
}

func V(level Level) InfoLogger {
	return std.V(level)
}

func L(ctx context.Context) Logger {
	return std.L(ctx)
}
