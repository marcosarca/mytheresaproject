package database

import (
	"context"
	"reflect"
)

type Database interface {
	Save(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, here interface{}) error
	GetWithFilters(ctx context.Context, here interface{}, filters ...Filter) error
	ErrRecordNotFound() error
	MigrateModels(models ...interface{}) error
}

type Filter interface {
	GetField() reflect.StructField
	GetValue() interface{}
	GetOperand() string
}
