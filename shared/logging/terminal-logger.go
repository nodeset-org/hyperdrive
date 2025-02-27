package logging

import (
	"log/slog"

	"github.com/fatih/color"
)

// ======================
// === TerminalLogger ===
// ======================

// Simple logger that logs colored messages directly to the terminal (stdout) without time or level.
// Useful for CLI applications that only log debug messages, where such things aren't relevant.
type TerminalLogger struct {
	*slog.Logger
	colorWriter *colorWriter
}

// Creates a new TerminalLogger instance
func NewTerminalLogger(debugEnabled bool, logColor color.Attribute) *TerminalLogger {
	// Create the logger options
	opts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: WithoutTimeAndLevel,
	}
	if debugEnabled {
		opts.Level = slog.LevelDebug
	}

	// Create the logger
	cw := newColorWriter(logColor)
	return &TerminalLogger{
		colorWriter: cw,
		Logger:      slog.New(slog.NewTextHandler(cw, opts)),
	}
}

// ===================
// === ColorWriter ===
// ===================

// Simple struct for printing colored log messages to the terminal
type colorWriter struct {
	impl *color.Color
}

// Creates a new ColorWriter
func newColorWriter(logColor color.Attribute) *colorWriter {
	return &colorWriter{
		impl: color.New(logColor),
	}
}

// Prints the logged message to the console, coloring the message with the handler's color
func (w *colorWriter) Write(p []byte) (n int, err error) {
	return w.impl.Println(string(p))
}
