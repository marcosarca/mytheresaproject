package logger

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, message string)
	Error(ctx context.Context, message string)
	Debug(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Sync() error
}

type logger struct {
	internal *zap.SugaredLogger
}

// NewLogger initializes the logger and allows the usage of the Logger Interface
// user MUST defer the call to Sync() immediately afterwards.

func NewLogger(service string) Logger {
	l, err := zap.NewProduction()
	if err != nil {
		panic("unable to initialize logger: " + err.Error())
	}

	return &logger{
		internal: l.WithOptions(zap.AddCallerSkip(1)).Sugar().
			With("service", service),
	}
}

// Sync MUST be deferred to flush any buffered logs prior to shutting down the application.
func (l *logger) Sync() error {
	return l.internal.Sync()
}

// Info allows for a message with info lever to be logged
func (l *logger) Info(ctx context.Context, message string) {
	l.injectTracing(ctx).Info(message)
}

// Error allows for a message with error lever to be logged
func (l *logger) Error(ctx context.Context, message string) {
	l.injectTracing(ctx).Error(message)
}

// Debug allows for a message with debug lever to be logged
func (l *logger) Debug(ctx context.Context, message string) {
	l.injectTracing(ctx).Debug(message)
}

// Warn allows for a message with warn lever to be logged
func (l *logger) Warn(ctx context.Context, message string) {
	l.injectTracing(ctx).Warn(message)
}

// WithField allows for the inclusion of a key-value into the log
func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{
		internal: l.internal.With(key, value),
	}
}

// WithError allows for the inclusion of an error into the log. It also populates DDog stack field
func (l *logger) WithError(err error) Logger {
	newLogger := l.internal.With(
		"error", err,
	)

	newLogger = newLogger.With("dd.error.stack", string(debug.Stack()))

	return &logger{
		internal: newLogger,
	}
}

// injectTracing enters request ID, if any
func (l *logger) injectTracing(ctx context.Context) *zap.SugaredLogger {
	//add our request id if present
	rid := ctx.Value("request_id")
	entry := l.internal
	if rid != nil {
		entry = entry.With("request_id", ctx.Value("request_id"))
	}

	return entry
}
