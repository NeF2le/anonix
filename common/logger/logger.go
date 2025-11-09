package logger

import (
	"context"
	"log/slog"
	"os"
)

type key string

const (
	KeyForLogger    key = "logger"
	KeyForRequestID key = "request_id"
	KeyForLogLevel  key = "log_level"
)

type Logger struct {
	l *slog.Logger
}

func NewLogger(logLevel string) *Logger {
	var logger *slog.Logger

	switch logLevel {
	case "debug":
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "info":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return &Logger{logger}
}

func New(ctx context.Context) context.Context {
	logLevel, ok := ctx.Value(KeyForLogLevel).(string)
	if !ok {
		logLevel = "debug"
	}

	loggerStruct := NewLogger(logLevel)
	ctx = context.WithValue(ctx, KeyForLogger, loggerStruct)

	return ctx
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	v, ok := ctx.Value(KeyForLogger).(*Logger)
	if !ok {
		ctx = New(ctx)
		return ctx.Value(KeyForLogger).(*Logger)
	}
	return v
}

func TryAppendRequestIDFromCtx(ctx context.Context, args []any) []any {
	if ctx.Value(KeyForRequestID) != nil {
		args = append(args, string(KeyForRequestID), ctx.Value(KeyForRequestID).(string))
	}
	return args
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	args = TryAppendRequestIDFromCtx(ctx, args)
	l.l.Info(msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	args = TryAppendRequestIDFromCtx(ctx, args)
	l.l.Error(msg, args...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	args = TryAppendRequestIDFromCtx(ctx, args)
	l.l.Debug(msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	args = TryAppendRequestIDFromCtx(ctx, args)
	l.l.Warn(msg, args...)
}
