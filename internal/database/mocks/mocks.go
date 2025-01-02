package mocks

import (
	"context"
	"mytheresa/internal/database"

	"github.com/stretchr/testify/mock"
)

type Database struct {
	mock.Mock
}

func (d *Database) Save(ctx context.Context, key string, value interface{}) error {
	args := d.Called(ctx, key, value)
	return args.Error(0)
}

func (d *Database) Get(ctx context.Context, key string, here interface{}) error {
	args := d.Called(ctx, key, here)
	return args.Error(0)
}

func (d *Database) GetWithFilters(ctx context.Context, here interface{}, filters ...database.Filter) error {
	args := d.Called(ctx, here, filters)
	return args.Error(0)
}

func (d *Database) ErrRecordNotFound() error {
	args := d.Called()
	return args.Error(0)
}

func (d *Database) MigrateModels(models ...interface{}) error {
	args := d.Called(models)
	return args.Error(0)
}
