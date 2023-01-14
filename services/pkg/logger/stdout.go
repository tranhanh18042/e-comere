package logger

import (
	"context"
	"io"
	"os"
)

// StdoutLogger implements Logger interface.
// It writes the log to stdout.
type StdoutLogger struct {
	w io.Writer
}

func NewStdoutLogger() (*StdoutLogger, error) {
	return &StdoutLogger{
		w: os.Stdout,
	}, nil
}

func (l *StdoutLogger) Info(ctx context.Context, msg ...any) {
	write(ctx, l.w, LogPrefixInfo, msg...)
}

func (l *StdoutLogger) Error(ctx context.Context, msg ...any) {
	write(ctx, l.w, LogPrefixError, msg...)
}

func (l *StdoutLogger) Debug(ctx context.Context, msg ...any) {
	write(ctx, l.w, LogPrefixDebug, msg...)
}
