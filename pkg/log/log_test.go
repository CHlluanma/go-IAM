package log

import (
	"context"
	"testing"
)

func TestLog_std(t *testing.T) {
	defer Flush()
	// info log
	Info("This is a info message", String("key1", "value1"))
	Infof("This is a formatted %s message", "info")
	Infow("Message printed with Infow", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	// debug log
	Debug("This is a debug message", String("key1", "value1"))
	Debugf("This is a formatted %s message", "debug")
	Debugw("Message printed with debugw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	// warn log
	Warn("This is a info message", String("key1", "value1"))
	Warnf("This is a formatted %s message", "info")
	Warnw("Message printed with Infow", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	// error log
	Error("This is a warn message", String("key1", "value1"))
	Errorf("This is a formatted %s message", "warn")
	Errorw("Message printed with warnw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	// panic log
	// Panic("This is a panic message", String("key1", "value1"))
	// Panicf("This is a formatted %s message", "panic")
	// Panicw("Message printed with panicw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	// faltal log
	// Fatal("This is a faltal message", String("key1", "value1"))
	// Fatalf("This is a formatted %s message", "faltal")
	// Fatalw("Message printed with faltalw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")
}

func TestLog_std_V(t *testing.T) {
	// V Level 通过整型数值来灵活指定日志级别，数值越大，优先级越低
	V(1).Info("This is a V level message")
	V(2).Infof("This is a %s V level message", "formatted")
	V(3).Infow("This is a V level message with fields", "X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")
}

func TestLog_std_WithValues(t *testing.T) {
	// 返回一个携带指定key-value的Logger
	lv := WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")
	lv.Infow("Info message printed with [WithValues] logger")
	lv.Warnf("Warning message printed with [WithValues] logger")
}

func TestLog_std_WithContext(t *testing.T) {
	lv := WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	// 返回一个携带指定上下文的Logger
	ctx := lv.WithContext(context.Background())
	// 从上下文获取Logger
	lc := FromContext(ctx)
	lc.Infow("Info message printed with [WithContext] logger")
	lv.Warnf("Warning message printed with [WithContext] logger")
}

func TestLog_std_WithName(t *testing.T) {
	// 返回一个携带指定名称的Logger
	lv := WithName("test")
	lv.Infow("Info message printed with [WithName] logger")
	lv.Warnf("Warning message printed with [WithName] logger")
}
