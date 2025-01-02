package category_test

import (
	"fmt"
	"mytheresa/pkg/category"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryRequest_ToCategory(t *testing.T) {
	request := category.CategoryRequest{Name: "Sandals"}
	c := request.ToCategory()

	assert.Equal(t, request.Name, c.Name)
}

func TestCategory_ToCategoryResponse(t *testing.T) {
	c := category.Category{ID: 1, Name: "Boots"}
	response := c.ToCategoryResponse()

	assert.Equal(t, fmt.Sprint(c.ID), response.ID)
	assert.Equal(t, c.Name, response.Name)
}

func TestCategory_GetIdentifier(t *testing.T) {
	c := category.Category{ID: 42}
	identifier := c.GetIdentifier()

	assert.Equal(t, fmt.Sprint(c.ID), identifier)
}
