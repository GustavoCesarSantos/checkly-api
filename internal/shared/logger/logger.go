package logger

import (
	"context"
	"log/slog"
)

func log(
	ctx context.Context,
	level slog.Level,
	message string,
	who string,
	operation string,
	attrs ...any,
) {
	baseAttrs := []any{
		"who", who,
		"operation", operation,
	}

	allAttrs := append(baseAttrs, attrs...)

	if ctx != nil {
		slog.Log(ctx, level, message, allAttrs...)
		return
	}

	slog.Log(context.TODO(), level, message, allAttrs...)
}

func Info(message, who, operation string, attrs ...any) {
	log(context.TODO(), slog.LevelInfo, message, who, operation, attrs...)
}

func InfoContext(ctx context.Context, message, who, operation string, attrs ...any) {
	log(ctx, slog.LevelInfo, message, who, operation, attrs...)
}

func Warn(message, who, operation string, attrs ...any) {
	log(context.TODO(), slog.LevelWarn, message, who, operation, attrs...)
}

func WarnContext(ctx context.Context, message, who, operation string, attrs ...any) {
	log(ctx, slog.LevelWarn, message, who, operation, attrs...)
}

func Error(message, who, operation string, err error, attrs ...any) {
	log(context.TODO(), slog.LevelError, message, who, operation, append(attrs, "error", err)...)
}

func ErrorContext(ctx context.Context, message, who, operation string, err error, attrs ...any) {
	log(ctx, slog.LevelError, message, who, operation, append(attrs, "error", err)...)
}