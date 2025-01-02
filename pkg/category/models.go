package category

import (
	"fmt"
	"strconv"
)

type Category struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *CategoryRequest) ToCategory() Category {
	return Category{
		Name: c.Name,
	}
}

func (c *Category) ToCategoryResponse() CategoryResponse {
	return CategoryResponse{
		ID:   strconv.Itoa(c.ID),
		Name: c.Name,
	}
}

func (c *Category) GetIdentifier() string {
	return fmt.Sprint(c.ID)
}
