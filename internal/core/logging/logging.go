package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggerKey = "logger"

var defaultLogger = zap.New(zapcore.NewCore(
	zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}),
	zapcore.AddSync(os.Stdout),
	zap.NewAtomicLevelAt(zapcore.InfoLevel),
), zap.AddCaller(), zap.AddCallerSkip(1))

func From(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return l
	}
	return defaultLogger
}

func With(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}
	return With(ctx, From(ctx).With(fields...))
}

func WithOptions(ctx context.Context, opts ...zap.Option) context.Context {
	if len(opts) == 0 {
		return ctx
	}
	return With(ctx, From(ctx).WithOptions(opts...))
}

func Sync(ctx context.Context) error {
	return From(ctx).Sync()
}
