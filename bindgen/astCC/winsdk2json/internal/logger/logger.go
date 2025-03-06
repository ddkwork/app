// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package log provides context-aware and structured logging capabilities.
package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
	// With returns a logger based off the root logger and decorates it with
	// the given context and arguments.
	With(ctx context.Context, args ...any) Logger

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...any)
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...any)
	// Fatalf uses fmt.Sprintf to construct and log a FATAL message.
	Error(args ...any)
	// Fatalf logs a FATAL message.
	Fatal(args ...any)

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...any)
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...any)
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...any)
	// Fatalf uses fmt.Sprintf to construct and log a FATAL message.
	Fatalf(format string, args ...any)
}

type logger struct {
	*zap.SugaredLogger
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
)

// New creates a new logger using the default configuration.
func New() Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// NewCustom creates a new logger using a custom configuration
// given a log level.
func NewCustom(level string) Logger {
	// NewProductionConfig is a reasonable production logging configuration
	// Uses JSON, writes to standard error, and enables sampling.
	// Stacktraces are automatically included on logs of ErrorLevel and above.
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Encoding = "console"
	config.DisableCaller = true
	config.Level = getLoggingLevel(level)
	logger, _ := config.Build()
	return NewWithZap(logger)
}

// NewWithZap creates a new logger using the pre-configured zap logger.
func NewWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}

// NewForTest returns a new logger and the corresponding observed logs which
// can be used in unit tests to verify log entries.
func NewForTest() (Logger, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	return NewWithZap(zap.New(core)), recorded
}

// With returns a logger based off the root logger and decorates it with
// the given context and arguments.
//
// If the context contains request ID and/or correlation ID information
// (recorded via WithRequestID() and WithCorrelationID()), they will be
// added to every log message generated by the new logger.
//
// The arguments should be specified as a sequence of name, value pairs
// with names being strings.
// The arguments will also be added to every log message generated by the logger.
func (l *logger) With(ctx context.Context, args ...any) Logger {
	if ctx != nil {
		if id, ok := ctx.Value(requestIDKey).(string); ok {
			args = append(args, zap.String("request_id", id))
		}
		if id, ok := ctx.Value(correlationIDKey).(string); ok {
			args = append(args, zap.String("correlation_id", id))
		}
	}
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}
