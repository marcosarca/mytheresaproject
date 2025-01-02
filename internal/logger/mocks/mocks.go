package mocks

import (
	"context"
	"mytheresa/internal/logger"
)

type NoopLogger struct{}

func (m *NoopLogger) Info(ctx context.Context, message string) {}

func (m *NoopLogger) Error(ctx context.Context, message string) {}

func (m *NoopLogger) Debug(ctx context.Context, message string) {}

func (m *NoopLogger) Warn(ctx context.Context, message string) {}

func (m *NoopLogger) WithField(key string, value interface{}) logger.Logger {
	return m
}

func (m *NoopLogger) WithError(err error) logger.Logger {
	return m
}

func (m *NoopLogger) Sync() error {
	return nil
}
