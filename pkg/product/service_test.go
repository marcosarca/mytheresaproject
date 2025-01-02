package product_test

import (
	"context"
	"errors"
	"fmt"
	"mytheresa/internal/apierror"
	dbmocks "mytheresa/internal/database/mocks"
	loggermocks "mytheresa/internal/logger/mocks"
	"mytheresa/pkg/category"
	"mytheresa/pkg/discount"
	discountmocks "mytheresa/pkg/discount/mocks"
	"mytheresa/pkg/product"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	ds := discountmocks.Service{}
	dbmock := dbmocks.Database{}
	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	assert.NotNil(t, s)
}

func TestCreateProduct_OK(t *testing.T) {
	pr := product.ProductRequest{
		SKU:        "1234",
		Name:       "Test product",
		Price:      11000,
		CategoryID: 1,
	}

	ds := discountmocks.Service{}
	dbmock := dbmocks.Database{}
	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if h, ok := args.Get(2).(*product.Product); ok {
			*h = pr.ToProduct()
		}
	},
	).Return(nil)
	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	result, err := s.CreateProduct(context.Background(), pr)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, pr.SKU, result.SKU)
	assert.Equal(t, pr.Name, result.Name)
	assert.Equal(t, pr.Price, result.Price)
	assert.Equal(t, pr.CategoryID, result.CategoryID)
}

func TestCreateProduct_ErrorSavingOnDB(t *testing.T) {
	pr := product.ProductRequest{
		SKU:        "1234",
		Name:       "Test product",
		Price:      11000,
		CategoryID: 1,
	}

	ds := discountmocks.Service{}
	dbmock := dbmocks.Database{}
	dbmock.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some DB error"))
	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	_, err := s.CreateProduct(context.Background(), pr)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Error creating product: %s", pr.Name), err.Error())
}

func TestGetProduct_OK(t *testing.T) {
	p := product.Product{
		SKU:        "1234",
		Name:       "Test product",
		Category:   category.Category{},
		CategoryID: 1,
		Price:      11000,
	}
	ds := discountmocks.Service{}
	dbmock := dbmocks.Database{}
	dbmock.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if h, ok := args.Get(2).(*product.Product); ok {
			*h = p
		}
	}).Return(nil)
	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	result, err := s.GetProduct(context.Background(), "1234")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, p.SKU, result.SKU)
}

func TestGetProduct_ErrorGettingFromDB(t *testing.T) {
	ds := discountmocks.Service{}
	dbmock := dbmocks.Database{}
	dbmock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some DB error"))
	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	_, err := s.GetProduct(context.Background(), "1234")
	assert.NotNil(t, err)
	apierr, ok := err.(*apierror.ApiError)
	assert.True(t, ok)
	assert.Equal(t, "Error getting Product with ID 1234", apierr.Error())
}

func TestListProducts_OK(t *testing.T) {
	dbdata := []product.Product{
		{
			SKU:        "1234",
			Name:       "Test product",
			Category:   category.Category{},
			CategoryID: 1,
			Price:      11000,
		},
	}
	ds := discountmocks.Service{}
	ds.On("GetDiscounts", mock.Anything).Return([]discount.Discount{
		&discount.GeneralDiscount{
			ID:             1,
			Percentage:     10,
			DiscountTypeID: 1,
			DiscountType:   discount.DiscountType{},
			Target:         "1",
		},
	}, nil)

	dbmock := dbmocks.Database{}
	dbmock.On("GetWithFilters", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if h, ok := args.Get(1).(*[]product.Product); ok {
			*h = dbdata
		}
	}).Return(nil)

	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	result, err := s.ListProducts(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
}

func TestListProducts_SeveralDiscountsApplyToSameProduct(t *testing.T) {
	dbdata := []product.Product{
		{
			SKU:        "1234",
			Name:       "Test product",
			Category:   category.Category{},
			CategoryID: 1,
			Price:      11000,
		},
	}
	ds := discountmocks.Service{}
	minorDiscount := &discount.CategoryDiscount{
		GeneralDiscount: discount.GeneralDiscount{

			ID:             1,
			Percentage:     10,
			DiscountTypeID: 1,
			DiscountType:   discount.DiscountType{},
			Target:         "1",
		},
	}
	greaterDiscount := &discount.SkuDiscount{
		GeneralDiscount: discount.GeneralDiscount{
			ID:             2,
			Percentage:     50,
			DiscountTypeID: 2,
			DiscountType:   discount.DiscountType{},
			Target:         "1234",
		},
	}
	ds.On("GetDiscounts", mock.Anything).Return([]discount.Discount{
		minorDiscount,
		greaterDiscount,
	}, nil)

	dbmock := dbmocks.Database{}
	dbmock.On("GetWithFilters", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if h, ok := args.Get(1).(*[]product.Product); ok {
			*h = dbdata
		}
	}).Return(nil)

	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	result, err := s.ListProducts(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	p := result[0]
	resultPrice := dbdata[0].Price * greaterDiscount.Percentage / 100
	assert.Equal(t, resultPrice, p.Price.Final)
}

func TestListProducts_ErrorSearchingOnDBProducts(t *testing.T) {
	ds := discountmocks.Service{}

	dbmock := dbmocks.Database{}
	dbmock.On("GetWithFilters", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some DB error"))

	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	result, err := s.ListProducts(context.Background())

	assert.Nil(t, result)
	assert.NotNil(t, err)
	apierr, ok := err.(*apierror.ApiError)
	assert.True(t, ok)
	assert.Equal(t, "Failed to get products from database", apierr.Error())
}

func TestListProducts_ErrorGettingDiscounts(t *testing.T) {
	dbdata := []product.Product{
		{
			SKU:        "1234",
			Name:       "Test product",
			Category:   category.Category{},
			CategoryID: 1,
			Price:      11000,
		},
	}
	ds := discountmocks.Service{}
	discountErr := apierror.InternalServerError("error getting discounts")
	ds.On("GetDiscounts", mock.Anything).Return([]discount.Discount{}, discountErr)

	dbmock := dbmocks.Database{}
	dbmock.On("GetWithFilters", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if h, ok := args.Get(1).(*[]product.Product); ok {
			*h = dbdata
		}
	}).Return(nil)

	logMock := loggermocks.NoopLogger{}

	s := product.NewService(&dbmock, &logMock, &ds)

	result, err := s.ListProducts(context.Background())

	assert.Nil(t, result)
	assert.NotNil(t, err)
	apierr, ok := err.(*apierror.ApiError)
	assert.True(t, ok)
	assert.Equal(t, discountErr.Error(), apierr.Error())
}
