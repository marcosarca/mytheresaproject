package sqlite

import (
	"context"
	"fmt"
	"mytheresa/internal/database"
	"mytheresa/internal/logger"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

const (
	FOREIGNKEY_TAG = "foreignKey"
)

type sqliteDB struct {
	*gorm.DB
	logger logger.Logger
}

func NewSQLiteDB(db *gorm.DB, logger logger.Logger) database.Database {
	//Initial data from problem description
	s := &sqliteDB{db, logger}

	return s
}

func (db *sqliteDB) ErrRecordNotFound() error {
	return gorm.ErrRecordNotFound
}

func (db *sqliteDB) Save(ctx context.Context, key string, value interface{}) error {
	t := getActualType(value)

	db.logger.WithField("key", key).Info(ctx, fmt.Sprintf("creating %v ", t))

	err := db.Create(value).Error
	if err != nil {
		db.logger.WithError(err).Error(ctx, fmt.Sprintf("error creating %v ", t))
		return err
	}

	return nil
}

func (db *sqliteDB) Get(ctx context.Context, key string, here interface{}) error {
	t := getActualType(here)
	db.logger.WithField("key", key).Info(ctx, fmt.Sprintf("getting %v ", t))

	err := preloadTables(db.DB, t).First(here, key).Error
	if err != nil {
		db.logger.WithField("key", key).WithError(err).
			Error(ctx, fmt.Sprintf("error getting %v ", t))
	}

	return err
}

func (db *sqliteDB) GetWithFilters(ctx context.Context, here interface{}, filters ...database.Filter) error {
	t := getActualType(here)

	query := applyFilters(preloadTables(db.DB, t), filters...)

	err := query.Find(here).Error

	db.logger.WithField("found in db", here).Info(ctx, fmt.Sprintf("getting %v ", t))
	return err
}

func preloadTables(query *gorm.DB, t reflect.Type) *gorm.DB {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		gormTag := field.Tag.Get("gorm")

		if strings.Contains(gormTag, FOREIGNKEY_TAG) {
			query = query.Preload(field.Name)
		}
	}
	return query
}

func applyFilters(query *gorm.DB, filters ...database.Filter) *gorm.DB {

	for _, filter := range filters {
		q := fmt.Sprintf("%s %s ?", filter.GetColumnName(), filter.GetOperand())
		query = query.Where(q, filter.GetValue())
	}
	return query
}

func getActualType(val interface{}) reflect.Type {
	t := reflect.TypeOf(val)

	switch t.Kind() {
	case reflect.Ptr:
		// If it's a pointer, check the type it points to
		elemType := t.Elem()
		if elemType.Kind() == reflect.Slice || elemType.Kind() == reflect.Array {
			return elemType.Elem() // Return the struct type inside the slice/array
		}
		return elemType

	case reflect.Slice, reflect.Array:
		// If it's a slice or array, return the element type
		return t.Elem()

	default:
		// If it's not a pointer, slice, or array, return nil or handle differently
		return t
	}

}

func (db *sqliteDB) MigrateModels(models ...interface{}) error {
	// Auto-migrate the models
	return db.AutoMigrate(
		models...,
	)
}
