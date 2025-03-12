package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is a simple wrapper for a slog Logger that writes to a file on disk.
type Logger struct {
	*slog.Logger
	logFile *lumberjack.Logger
	path    string
}

// Creates a new logger that writes out to a log file on disk.
func NewLogger(logFilePath string, options LoggerOptions) (*Logger, error) {
	// Make the file
	err := os.MkdirAll(filepath.Dir(logFilePath), logDirMode)
	if err != nil {
		return nil, fmt.Errorf("error creating API log directory for \"%s\": %w", logFilePath, err)
	}
	handle, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR, logFileMode)
	handle.Close()

	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    options.MaxSize,
		MaxBackups: options.MaxBackups,
		MaxAge:     options.MaxAge,
		LocalTime:  options.LocalTime,
		Compress:   options.Compress,
	}

	// Create the logging options
	logOptions := &slog.HandlerOptions{
		AddSource:   options.AddSource,
		Level:       options.Level,
		ReplaceAttr: ReplaceTime,
	}

	// Make the logger
	var handler slog.Handler
	switch options.Format {
	case LogFormat_Json:
		handler = slog.NewJSONHandler(logFile, logOptions)
	case LogFormat_Logfmt:
		handler = slog.NewTextHandler(logFile, logOptions)
	}
	return &Logger{
		Logger:  slog.New(handler),
		logFile: logFile,
		path:    logFilePath,
	}, nil
}

// Creates a new logger that uses the slog default logger, which writes to the terminal instead of a file.
// Operations like rotation don't apply to this logger.
func NewDefaultLogger() *Logger {
	return &Logger{
		Logger: slog.Default(),
	}
}

// Get the path of the file this logger is writing to
func (l *Logger) GetFilePath() string {
	return l.path
}

// Rotate the log file, migrating the current file to an old backup and starting a new one
func (l *Logger) Rotate() error {
	if l.logFile != nil {
		return l.logFile.Rotate()
	}
	return nil
}

// Closes the log file
func (l *Logger) Close() {
	if l.logFile != nil {
		l.Info("Shutting down.")
		l.logFile.Close()
		l.logFile = nil
	}
}

// Create a clone of the logger that prints each message with the "origin" attribute.
// The underlying file handle isn't copied, so calling Close() on the sublogger won't do anything.
func (l *Logger) CreateSubLogger(origin string) *Logger {
	return &Logger{
		Logger:  l.With(slog.String("origin", origin)),
		logFile: nil,
	}
}

// Creates a copy of the parent context with the logger put into the ContextLogKey value
func (l *Logger) CreateContextWithLogger(parent context.Context) context.Context {
	return context.WithValue(parent, ContextLogKey, l)
}

// Retrieves the logger from the context
func FromContext(ctx context.Context) (*Logger, bool) {
	log, ok := ctx.Value(ContextLogKey).(*Logger)
	return log, ok
}
