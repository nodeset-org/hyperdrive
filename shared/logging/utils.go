package logging

import (
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
