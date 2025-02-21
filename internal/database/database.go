package database

import (
	"context"
)

type Database interface {
	Save(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, here interface{}) error
	GetWithFilters(ctx context.Context, here interface{}, filters ...Filter) error
	ErrRecordNotFound() error
	MigrateModels(models ...interface{}) error
}

type Filter interface {
	GetColumnName() string
	GetValue() interface{}
	GetOperand() string
}
