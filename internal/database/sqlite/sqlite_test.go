package sqlite_test

import (
	"context"
	"errors"
	"mytheresa/internal/database"
	"mytheresa/internal/database/sqlite"
	loggermocks "mytheresa/internal/logger/mocks"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbname string = "fakeDB"

// Dummy model for testing
type dummyModel struct {
	ID   int
	Name string
}

// Dummy filter implementation that satisfies the database.Filter interface
type dummyFilter struct {
	field   string
	operand string
	value   interface{}
}

func (f *dummyFilter) GetColumnName() string {
	return f.field
}

func (f *dummyFilter) GetValue() interface{} {
	return f.value
}

func (f *dummyFilter) GetOperand() string {
	return f.operand
}

func NewDummyFilter(operand string, value interface{}) *dummyFilter {
	field := "Name"

	return &dummyFilter{
		field:   field,
		operand: operand,
		value:   value,
	}
}

func MockDB() database.Database {
	db, _ := gorm.Open(gormsqlite.Open(dbname), &gorm.Config{})
	sqliteDB := sqlite.NewSQLiteDB(db, &loggermocks.NoopLogger{})
	sqliteDB.MigrateModels(&dummyModel{})

	return sqliteDB
}

func TestNewSQLiteDB(t *testing.T) {
	sqldb := MockDB()
	defer os.Remove(dbname)

	assert.NotNil(t, sqldb)
}

// TestSave tests the Save method of the sqliteDB struct.
func TestSave(t *testing.T) {
	sqliteDB := MockDB()
	defer os.Remove(dbname)

	model := dummyModel{Name: "test"}
	err := sqliteDB.Save(context.Background(), "test_key", &model)
	assert.NoError(t, err)
	assert.Equal(t, 1, model.ID)

	model2 := dummyModel{ID: 1, Name: "test"}
	err = sqliteDB.Save(context.Background(), "test_key", &model2)

	assert.Error(t, err)
}

// TestGet tests the Get method of the sqliteDB struct.
func TestGet(t *testing.T) {

	sqliteDB := MockDB()
	defer os.Remove(dbname)

	sqliteDB.Save(context.Background(), "test_key", &dummyModel{Name: "test"})

	// Simulate success
	var result dummyModel
	err := sqliteDB.Get(context.Background(), "1", &result)
	assert.NoError(t, err)

	err = sqliteDB.Get(context.Background(), "2", &result)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestGetWithFilters(t *testing.T) {
	sqliteDB := MockDB()
	defer os.Remove(dbname)

	sqliteDB.Save(context.Background(), "test_key", &dummyModel{Name: "test"})

	// Create a dummy filter
	filter := NewDummyFilter("=", "test")

	var result []dummyModel
	err := sqliteDB.GetWithFilters(context.Background(), &result, filter)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestErrRecordNotFound(t *testing.T) {
	sqliteDB := MockDB()
	defer os.Remove(dbname)

	err := sqliteDB.ErrRecordNotFound()
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
