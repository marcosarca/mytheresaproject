package discount_test

import (
	"context"
	"errors"
	"mytheresa/internal/apierror"
	dbmocks "mytheresa/internal/database/mocks"
	loggermocks "mytheresa/internal/logger/mocks"
	"mytheresa/pkg/discount"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	s := discount.NewService(&dbmock, &logMock)

	assert.NotNil(t, s)
}

func TestCreateDiscountType_OK(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s := discount.NewService(&dbmock, &logMock)
	req := discount.DiscountTypeRequest{
		Type: "Test",
	}

	result, err := s.CreateDiscountType(context.Background(), req)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Type, result.Type)
}

func TestCreateDiscountType_Error(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some DB error"))

	s := discount.NewService(&dbmock, &logMock)
	req := discount.DiscountTypeRequest{
		Type: "Test",
	}

	_, err := s.CreateDiscountType(context.Background(), req)
	assert.NotNil(t, err)

	apierr, ok := err.(*apierror.ApiError)

	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apierr.Code())
	assert.Equal(t, "error creating discount type", apierr.Error())
}

func TestCreateDiscount_OK(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s := discount.NewService(&dbmock, &logMock)

	req := discount.DiscountRequest{
		Percentage:     10,
		DiscountTypeID: 1,
		Target:         "000005",
	}

	result, err := s.CreateDiscount(context.Background(), req)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Percentage, result.GetPercentage())
}

func TestCreateDiscount_Error(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some DB error"))
	s := discount.NewService(&dbmock, &logMock)
	req := discount.DiscountRequest{
		Percentage:     10,
		DiscountTypeID: 1,
		Target:         "000005",
	}

	_, err := s.CreateDiscount(context.Background(), req)

	assert.NotNil(t, err)

	apierr, ok := err.(*apierror.ApiError)

	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apierr.Code())
	assert.Equal(t, "error creating discount", apierr.Error())
}

func TestGetDiscounts_OK_OneOfEachDiscountOnDB(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	dbData := []discount.GeneralDiscount{
		{
			ID:             1,
			Percentage:     10,
			DiscountTypeID: discount.CATEGORY,
			DiscountType:   discount.DiscountType{},
			Target:         "1",
		},
		{
			ID:             2,
			Percentage:     20,
			DiscountTypeID: discount.SKU,
			DiscountType:   discount.DiscountType{},
			Target:         "000005",
		},
		{
			ID:             3,
			Percentage:     30,
			DiscountTypeID: discount.GENERAL,
			DiscountType:   discount.DiscountType{},
			Target:         "000005",
		},
	}

	dbmock.On("GetWithFilters", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if h, ok := args.Get(1).(*[]discount.GeneralDiscount); ok {
			*h = dbData
		}
	}).Return(nil)

	s := discount.NewService(&dbmock, &logMock)

	results, err := s.GetDiscounts(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 3)

	_, ok := results[0].(*discount.CategoryDiscount)
	assert.True(t, ok)

	_, ok = results[1].(*discount.SkuDiscount)
	assert.True(t, ok)
}

func TestGetDiscounts_OK_NoDiscountsOnDB(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	dbmock.On("GetWithFilters", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s := discount.NewService(&dbmock, &logMock)
	results, err := s.GetDiscounts(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 0)
}

func TestGetDiscounts_DatabaseError(t *testing.T) {
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	dbmock.On("GetWithFilters", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some DB error"))

	s := discount.NewService(&dbmock, &logMock)
	_, err := s.GetDiscounts(context.Background())

	assert.NotNil(t, err)
	apierr, ok := err.(*apierror.ApiError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apierr.Code())
	assert.Equal(t, "error getting discounts", apierr.Error())
}
