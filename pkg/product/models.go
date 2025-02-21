package product

import (
	"mytheresa/internal/database"
	"mytheresa/pkg/category"
)

type Product struct {
	SKU        string            `gorm:"primaryKey" json:"sku"`
	Name       string            `gorm:"not null" json:"name"`
	Category   category.Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	CategoryID int               `gorm:"not null" json:"category_id"`
	Price      int               `gorm:"not null" json:"price"`
}

// ProductRequest represents the body for creating a product
// @Description ProductRequest is the input for creating a new product
// @Accept json
// @Produce json
// @Param product body ProductRequest true "Product details"
type ProductRequest struct {
	SKU        string `json:"sku" example:"000005"`
	Name       string `json:"name" example:"Legendary Boots"`
	Price      int    `json:"price" example:"10000"`
	CategoryID int    `json:"category_id" example:"1"`
}

func (p *ProductRequest) ToProduct() Product {
	return Product{
		SKU:        p.SKU,
		Name:       p.Name,
		Price:      p.Price,
		CategoryID: p.CategoryID,
	}
}

func (p *Product) ToProductResponse() ProductResponse {
	return ProductResponse{
		SKU:      p.SKU,
		Name:     p.Name,
		Category: p.Category.Name,
		Price: PriceResponse{
			Original:           p.Price,
			Final:              p.Price,
			DiscountPercentage: nil,
			Currency:           "EUR",
		},
	}
}

func (p *Product) GetIdentifier() string {
	return p.SKU
}

// ProductResponse represents a product with its details
// @Description ProductResponse is the output when retrieving product details
// @Accept json
// @Produce json
// @Success 200 {object} ProductResponse
type ProductResponse struct {
	SKU      string        `json:"sku" example:"000005"`
	Name     string        `json:"name" example:"Legendary boots"`
	Category string        `json:"category" example:"Boots"`
	Price    PriceResponse `json:"price"`
}

// PriceResponse represents the price details of a product
// @Description PriceResponse includes the original and final price of a product, along with any discounts
// @Accept json
// @Produce json
// @Success 200 {object} PriceResponse
type PriceResponse struct {
	Original           int     `json:"original" example:"10000"`
	Final              int     `json:"final" example:"8000"`
	DiscountPercentage *string `json:"discount_percentage,omitempty" example:"20"`
	Currency           string  `json:"currency" example:"EUR"`
}

type categoryFilter struct {
	field   string
	Value   string
	Operand string
}

func (f *categoryFilter) GetColumnName() string {
	return f.field
}

func (f *categoryFilter) GetValue() interface{} {
	return f.Value
}

func (f *categoryFilter) GetOperand() string {
	return f.Operand
}

func NewCategoryFilter(value string, operand string) database.Filter {
	return &categoryFilter{
		field:   "category_id",
		Value:   value,
		Operand: operand,
	}
}

type priceFilter struct {
	field   string
	Value   string
	Operand string
}

func (f *priceFilter) GetColumnName() string {
	return f.field
}

func (f *priceFilter) GetValue() interface{} {
	return f.Value
}

func (f *priceFilter) GetOperand() string {
	return f.Operand
}

func NewPriceFilter(value string, operand string) database.Filter {
	return &priceFilter{
		field:   "price",
		Value:   value,
		Operand: operand,
	}
}
