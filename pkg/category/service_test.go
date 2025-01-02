package category_test

import (
	"context"
	"errors"
	"mytheresa/internal/apierror"
	databasemocks "mytheresa/internal/database/mocks"
	loggermocks "mytheresa/internal/logger/mocks"
	"mytheresa/pkg/category"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	dbmock := databasemocks.Database{}
	logmock := loggermocks.NoopLogger{}

	s := category.NewService(&dbmock, &logmock)

	assert.NotNil(t, s)
}

func TestService_CreateCategory_Success(t *testing.T) {
	dbmock := databasemocks.Database{}
	logmock := loggermocks.NoopLogger{}

	catReq := category.CategoryRequest{
		Name: "Test Category",
	}

	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s := category.NewService(&dbmock, &logmock)

	result, err := s.CreateCategory(context.Background(), catReq)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, catReq.Name, result.Name)
}

func TestService_CreateCategory_Error(t *testing.T) {
	dbmock := databasemocks.Database{}
	logmock := loggermocks.NoopLogger{}
	catReq := category.CategoryRequest{
		Name: "Test Category",
	}

	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some DB error"))

	s := category.NewService(&dbmock, &logmock)
	_, err := s.CreateCategory(context.Background(), catReq)

	assert.NotNil(t, err)

	apierr, ok := err.(*apierror.ApiError)

	assert.True(t, ok)
	assert.Equal(t, "there was an error saving the category", apierr.Error())
	assert.Equal(t, http.StatusInternalServerError, apierr.Code())
}
