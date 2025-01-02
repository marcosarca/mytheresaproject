package product_test

import (
	"bytes"
	"encoding/json"
	"mytheresa/internal/apierror"
	loggermocks "mytheresa/internal/logger/mocks"
	"mytheresa/pkg/category"
	"mytheresa/pkg/product"
	productmocks "mytheresa/pkg/product/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHandler(t *testing.T) {
	ps := productmocks.Service{}
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	assert.NotNil(t, h)
}

func TestHandlerCreateProduct_OK(t *testing.T) {
	p := product.ProductRequest{
		SKU:        "000001",
		Name:       "Test product",
		Price:      95000,
		CategoryID: 1,
	}
	ps := productmocks.Service{}
	ps.On("CreateProduct", mock.Anything, mock.Anything).Return(p.ToProduct(), nil)
	logMock := loggermocks.NoopLogger{}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(p)
	r := httptest.NewRequest("POST", "/products", bytes.NewReader(body))

	h := product.NewHandler(&ps, &logMock)
	h.CreateProduct(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response product.Product
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Equal(t, p.Name, response.Name)
	assert.Equal(t, p.Price, response.Price)
	assert.Equal(t, p.CategoryID, response.CategoryID)
}

func TestHandlerCreateProduct_WrongBody(t *testing.T) {
	ps := productmocks.Service{}
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid body")))
	w := httptest.NewRecorder()

	h.CreateProduct(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	var apierr apierror.ApiError
	err := json.NewDecoder(w.Body).Decode(&apierr)
	assert.NoError(t, err)
	assert.Equal(t, "Wrong body", apierr.Error())
}

func TestHandlerCreateProduct_ServiceError(t *testing.T) {
	productRequest := product.ProductRequest{
		SKU:        "12345",
		Name:       "Laptop",
		Price:      1000,
		CategoryID: 1,
	}

	ps := productmocks.Service{}
	ps.On("CreateProduct", mock.Anything, productRequest).Return(product.Product{}, apierror.InternalServerError("service error"))
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	body, _ := json.Marshal(productRequest)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateProduct(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	var apierr apierror.ApiError
	err := json.NewDecoder(w.Body).Decode(&apierr)
	assert.NoError(t, err)
	assert.Equal(t, "service error", apierr.Error())
}

func TestHandlerGetProduct_OK(t *testing.T) {
	productID := "000001"
	expectedProduct := product.Product{
		SKU:        productID,
		Name:       "Test Product",
		Price:      95000,
		CategoryID: 1,
		Category:   category.Category{Name: "Boots"},
	}

	ps := productmocks.Service{}
	ps.On("GetProduct", mock.Anything, productID).Return(expectedProduct, nil)
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest("GET", "/products/"+productID, nil)
	r = mux.SetURLVars(r, map[string]string{"id": productID})
	w := httptest.NewRecorder()

	h.GetProduct(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response product.Product
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, response)
}

func TestHandlerGetProduct_NotFound(t *testing.T) {
	productID := "000002"

	ps := productmocks.Service{}
	ps.On("GetProduct", mock.Anything, productID).Return(product.Product{}, apierror.NotFound("product not found"))
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest("GET", "/products/"+productID, nil)
	r = mux.SetURLVars(r, map[string]string{"id": productID})
	w := httptest.NewRecorder()

	h.GetProduct(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var apierr apierror.ApiError
	err := json.NewDecoder(w.Body).Decode(&apierr)
	assert.NoError(t, err)
	assert.Equal(t, "product not found", apierr.Error())
}

func TestHandlerGetProduct_ServiceError(t *testing.T) {
	productID := "000003"

	ps := productmocks.Service{}
	ps.On("GetProduct", mock.Anything, productID).Return(product.Product{}, apierror.InternalServerError("service error"))
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest("GET", "/products/"+productID, nil)
	r = mux.SetURLVars(r, map[string]string{"id": productID})
	w := httptest.NewRecorder()

	h.GetProduct(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var apierr apierror.ApiError
	err := json.NewDecoder(w.Body).Decode(&apierr)
	assert.NoError(t, err)
	assert.Equal(t, "service error", apierr.Error())
}

func TestHandlerListProducts_OK(t *testing.T) {
	products := []product.ProductResponse{
		{
			SKU:      "000001",
			Name:     "Product 1",
			Category: "Boots",
			Price:    product.PriceResponse{},
		},
		{
			SKU:      "000002",
			Name:     "Product 2",
			Category: "Boots",
			Price:    product.PriceResponse{},
		},
		{
			SKU:      "000003",
			Name:     "Product 3",
			Category: "Sneakers",
			Price:    product.PriceResponse{},
		},
	}
	ps := productmocks.Service{}
	ps.On("ListProducts", mock.Anything, mock.Anything).Return(products, nil)
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest("GET", "/products?limit=2", nil)
	w := httptest.NewRecorder()

	h.ListProducts(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []product.ProductResponse
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, products[:2], response)
}

func TestHandlerListProducts_DefaultLimit(t *testing.T) {
	products := []product.ProductResponse{
		{
			SKU:      "000001",
			Name:     "Product 1",
			Category: "Boots",
			Price:    product.PriceResponse{},
		},
		{
			SKU:      "000002",
			Name:     "Product 2",
			Category: "Boots",
			Price:    product.PriceResponse{},
		},
		{
			SKU:      "000003",
			Name:     "Product 3",
			Category: "Sandals",
			Price:    product.PriceResponse{},
		},
		{
			SKU:      "000004",
			Name:     "Product 4",
			Category: "Sandals",
			Price:    product.PriceResponse{},
		},
		{
			SKU:      "000005",
			Name:     "Product 5",
			Category: "Sneakers",
			Price:    product.PriceResponse{},
		},
		{
			SKU:      "000006",
			Name:     "Product 6",
			Category: "Sneakers",
			Price:    product.PriceResponse{},
		},
	}
	ps := productmocks.Service{}
	ps.On("ListProducts", mock.Anything, mock.Anything).Return(products, nil)
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	h.ListProducts(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []product.ProductResponse
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Len(t, response, 5)
	assert.Equal(t, products[:5], response)
}

func TestHandlerListProducts_WithFilters(t *testing.T) {
	products := []product.ProductResponse{
		{
			SKU:      "000001",
			Name:     "Product 1",
			Category: "Boots",
			Price:    product.PriceResponse{},
		},
	}
	ps := productmocks.Service{}
	ps.On("ListProducts", mock.Anything, mock.Anything).Return(products, nil)
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest("GET", "/products?category=Boots&priceLessThan=90000", nil)
	w := httptest.NewRecorder()

	h.ListProducts(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []product.ProductResponse
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, products, response)
}

func TestHandlerListProducts_ServiceError(t *testing.T) {
	ps := productmocks.Service{}
	ps.On("ListProducts", mock.Anything, mock.Anything).Return([]product.ProductResponse{}, apierror.InternalServerError("service error"))
	logMock := loggermocks.NoopLogger{}

	h := product.NewHandler(&ps, &logMock)

	r := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	h.ListProducts(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var apierr apierror.ApiError
	err := json.NewDecoder(w.Body).Decode(&apierr)
	assert.NoError(t, err)
	assert.Equal(t, "service error", apierr.Error())
}
