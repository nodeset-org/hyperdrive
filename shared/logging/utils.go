package logging

import (
	"context"
	"log/slog"
	"time"
)

// Prints an error to a log line
func Err(err error) slog.Attr {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return slog.String("error", msg)
}

// Replaces the default time formatting (RFC3339) in a logger with an easier to read format
func ReplaceTime(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		t := a.Value.Time()
		return slog.String(slog.TimeKey, t.UTC().Format(time.DateTime))
	}
	return a
}

// Removes the time and level from the message
func WithoutTimeAndLevel(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
		return slog.Attr{}
	}
	return a
}

// Logs a debug message. If the logger is nil, this function does nothing.
func SafeDebug(logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.Debug(msg, attrs...)
}

// Logs a debug message using the provided context. If the logger is nil, this function does nothing.
func SafeDebugWithContext(ctx context.Context, logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.DebugContext(ctx, msg, attrs...)
}

// Logs an info message. If the logger is nil, this function does nothing.
func SafeInfo(logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.Info(msg, attrs...)
}

// Logs an info message using the provided context. If the logger is nil, this function does nothing.
func SafeInfoWithContext(ctx context.Context, logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.InfoContext(ctx, msg, attrs...)
}

// Logs a warning message. If the logger is nil, this function does nothing.
func SafeWarn(logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.Warn(msg, attrs...)
}

// Logs a warning message using the provided context. If the logger is nil, this function does nothing.
func SafeWarnWithContext(ctx context.Context, logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.WarnContext(ctx, msg, attrs...)
}

// Logs an error message. If the logger is nil, this function does nothing.
func SafeError(logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.Error(msg, attrs...)
}

// Logs an error message using the provided context. If the logger is nil, this function does nothing.
func SafeErrorWithContext(ctx context.Context, logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.ErrorContext(ctx, msg, attrs...)
}

// Logs a message with the given level and attributes. If the logger is nil, this function does nothing.
func SafeLog(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.Log(ctx, level, msg, attrs...)
}
