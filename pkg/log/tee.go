package log

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelEnablerFunc func(level Level) bool

type TeeOption struct {
	Out io.Writer
	LevelEnablerFunc
}

func NewTee(tees []TeeOption, opts ...Option) *zapLogger {
	var cores []zapcore.Core
	for _, tee := range tees {
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.RFC3339TimeEncoder
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			zapcore.AddSync(tee.Out),
			zap.LevelEnablerFunc(tee.LevelEnablerFunc),
		)
		cores = append(cores, core)
	}
	return &zapLogger{
		zapL: zap.New(zapcore.NewTee(cores...), opts...),
		// al:   nil,
	}
}
