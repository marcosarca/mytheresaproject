package discount_test

import (
	"bytes"
	"encoding/json"
	"errors"
	loggermocks "mytheresa/internal/logger/mocks"
	"mytheresa/pkg/discount"
	discountmocks "mytheresa/pkg/discount/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHandler(t *testing.T) {
	logMock := loggermocks.NoopLogger{}
	smock := discountmocks.Service{}

	h := discount.NewHandler(&smock, &logMock)

	assert.NotNil(t, h)
}

func TestHandlerCreateDiscount_OK(t *testing.T) {
	logMock := loggermocks.NoopLogger{}
	smock := discountmocks.Service{}
	smock.On("CreateDiscount", mock.Anything, mock.Anything).Return(&discount.GeneralDiscount{
		ID:             1,
		Percentage:     10,
		DiscountTypeID: 1,
		DiscountType:   discount.DiscountType{},
		Target:         "000005",
	}, nil)

	h := discount.NewHandler(&smock, &logMock)

	w := httptest.NewRecorder()

	body, _ := json.Marshal(discount.DiscountRequest{
		Percentage:     10,
		DiscountTypeID: 1,
		Target:         "000005",
	})
	r := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	h.CreateDiscount(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandlerCreateDiscount_WrongBody(t *testing.T) {
	logMock := loggermocks.NoopLogger{}
	smock := discountmocks.Service{}

	h := discount.NewHandler(&smock, &logMock)

	w := httptest.NewRecorder()

	body, _ := json.Marshal(struct {
		Percentage     string `json:"percentage"`
		DiscountTypeID string `json:"discount_type_id"`
	}{
		"15",
		"1",
	})
	r := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	h.CreateDiscount(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerCreateDiscount_ServiceError(t *testing.T) {

	logMock := loggermocks.NoopLogger{}
	smock := discountmocks.Service{}
	smock.On("CreateDiscount", mock.Anything, mock.Anything).Return(&discount.GeneralDiscount{}, errors.New("Some Error"))

	h := discount.NewHandler(&smock, &logMock)

	w := httptest.NewRecorder()

	body, _ := json.Marshal(discount.DiscountRequest{
		Percentage:     10,
		DiscountTypeID: 1,
		Target:         "000005",
	})
	r := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	h.CreateDiscount(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandlerGetDiscounts_OK(t *testing.T) {
	logMock := loggermocks.NoopLogger{}
	smock := discountmocks.Service{}
	var discounts []discount.Discount
	smock.On("GetDiscounts", mock.Anything).Return(discounts, nil)

	h := discount.NewHandler(&smock, &logMock)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.GetDiscounts(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerGetDiscounts_ServiceError(t *testing.T) {
	logMock := loggermocks.NoopLogger{}
	smock := discountmocks.Service{}
	var discounts []discount.Discount
	smock.On("GetDiscounts", mock.Anything).Return(discounts, errors.New("Some Error"))

	h := discount.NewHandler(&smock, &logMock)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.GetDiscounts(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
