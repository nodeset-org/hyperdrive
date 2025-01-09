package adapter

import (
	"context"
	"log/slog"
	"strings"
	"testing"
)

type testLogger struct {
	t     *testing.T
	group string
	attrs []slog.Attr
}

func (l *testLogger) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (l *testLogger) Handle(ctx context.Context, record slog.Record) error {
	records := []string{}
	for _, attr := range l.attrs {
		records = append(records, attr.String())
	}
	record.Attrs(func(a slog.Attr) bool {
		records = append(records, a.String())
		return true
	})

	switch record.Level {
	case slog.LevelDebug:
		l.t.Logf("DEBUG: %s, %s", record.Message, strings.Join(records, ", "))
	case slog.LevelInfo:
		l.t.Logf("INFO: %s, %s", record.Message, strings.Join(records, ", "))
	case slog.LevelWarn:
		l.t.Logf("WARN: %s, %s", record.Message, strings.Join(records, ", "))
	case slog.LevelError:
		l.t.Errorf("ERROR: %s, %s", record.Message, strings.Join(records, ", "))
	default:
		if record.Level > slog.LevelError {
			l.t.Errorf("UNKNOWN: %s, %s", record.Message, strings.Join(records, ", "))
		} else {
			l.t.Logf("UNKNOWN: %s, %s", record.Message, strings.Join(records, ", "))
		}
	}
	return nil
}

func (l *testLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &testLogger{
		t:     l.t,
		group: l.group,
		attrs: append(l.attrs, attrs...),
	}
}

func (l *testLogger) WithGroup(name string) slog.Handler {
	if name == "" {
		return l
	}
	return &testLogger{
		t:     l.t,
		group: name,
		attrs: l.attrs,
	}
}

func CreateLogger(t *testing.T) *slog.Logger {
	return slog.New(&testLogger{
		t: t,
	})
}
