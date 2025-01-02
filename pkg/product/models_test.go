package product_test

import (
	"mytheresa/pkg/category"
	"mytheresa/pkg/product"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductRequest_ToProduct(t *testing.T) {
	request := product.ProductRequest{
		SKU:        "000005",
		Name:       "Legendary Boots",
		Price:      11000,
		CategoryID: 1,
	}
	p := request.ToProduct()

	assert.Equal(t, "000005", p.SKU)
	assert.Equal(t, "Legendary Boots", p.Name)
	assert.Equal(t, 11000, p.Price)
	assert.Equal(t, 1, p.CategoryID)
}

func TestProduct_ToProductResponse(t *testing.T) {
	p := product.Product{
		SKU:  "000005",
		Name: "Epic Sandals",
		Category: category.Category{
			ID:   2,
			Name: "Sandals",
		},
		Price: 500,
	}

	response := p.ToProductResponse()

	assert.Equal(t, "000005", response.SKU)
	assert.Equal(t, "Epic Sandals", response.Name)
	assert.Equal(t, "Sandals", response.Category)
	assert.Equal(t, 500, response.Price.Original)
	assert.Equal(t, 500, response.Price.Final)
	assert.Equal(t, "EUR", response.Price.Currency)
	assert.Nil(t, response.Price.DiscountPercentage)
}

func TestProduct_GetIdentifier(t *testing.T) {
	p := product.Product{SKU: "000001"}
	identifier := p.GetIdentifier()

	assert.Equal(t, "000001", identifier)
}

func TestNewCategoryFilter(t *testing.T) {
	filter := product.NewCategoryFilter("1", "=")

	field, _ := reflect.TypeOf(product.Product{}).FieldByName("CategoryID")
	assert.Equal(t, field.Name, filter.GetField().Name)
	assert.Equal(t, "1", filter.GetValue())
	assert.Equal(t, "=", filter.GetOperand())
}

func TestNewPriceFilter(t *testing.T) {
	filter := product.NewPriceFilter("100", ">")

	field, _ := reflect.TypeOf(product.Product{}).FieldByName("Price")
	assert.Equal(t, field.Name, filter.GetField().Name)
	assert.Equal(t, "100", filter.GetValue())
	assert.Equal(t, ">", filter.GetOperand())
}
