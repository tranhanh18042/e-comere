package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

/*
logger package should be initialized before using it by calling Init().
And then, the loggers can be used by calling the Info(), Error(), and Debug() functions globally.
*/

// Init initializes the loggers with given configuration.
// Use should be called once and only once in the application.
func Init() error {
	lgs, err := createLoggers()
	if err != nil {
		return fmt.Errorf("failed to create loggers: %v", err)
	}

	loggers = lgs
	return nil
}

// LogPrefix is the prefix of log.
type LogPrefix string

const (
	LogPrefixInfo  LogPrefix = "[INFO]"
	LogPrefixError LogPrefix = "[ERROR]"
	LogPrefixDebug LogPrefix = "[DEBUG]"
)

var bufPool = sync.Pool{
	New: func() any {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

// Info writes the log to the logger for INFO level.
func Info(ctx context.Context, msg ...any) {
	for _, l := range loggers {
		l.Info(ctx, msg...)
	}
}

// Error writes the log to the logger for ERROR level.
func Error(ctx context.Context, msg ...any) {
	for _, l := range loggers {
		l.Error(ctx, msg...)
	}
}

// Debug writes the log to the logger for DEBUG level.
func Debug(ctx context.Context, msg ...any) {
	for _, l := range loggers {
		l.Debug(ctx, msg...)
	}
}

// internalLogger represents a internal logger,
// and all loggers should implement this interface.
// Because writing logs is a side work,
// so we don't care if the log cannot be write successfully.
type internalLogger interface {
	Info(ctx context.Context, msg ...any)
	Error(ctx context.Context, msg ...any)
	Debug(ctx context.Context, msg ...any)
}

var loggers []internalLogger

func createLoggers() ([]internalLogger, error) {
	loggers := make([]internalLogger, 0, 1)
	stdLgr, err := NewStdoutLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout logger: %v", err)
	}
	loggers = append(loggers, stdLgr)

	return loggers, nil
}

// write writes the log to the logger under the hood.
func write(ctx context.Context, w io.Writer, prefix LogPrefix, msg ...any) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(string(prefix))
	b.WriteByte(' ')
	b.WriteByte('[')
	b.WriteString(time.Now().Format(time.RFC3339))
	b.WriteByte(']')
	b.WriteByte(' ')
	b.WriteByte(' ')
	for i, m := range msg {
		if i > 0 { // don't write ":" for the first message
			b.WriteByte(':')
		}
		b.WriteString(fmt.Sprint(m))
	}
	b.WriteString("\n")
	// Writing the log is just a side work, so we don't care if the log cannot be write successfully
	_, _ = w.Write(b.Bytes())
	bufPool.Put(b)
}
