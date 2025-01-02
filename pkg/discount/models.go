package discount

import (
	"strconv"
)

const (
	CATEGORY = 1 //applies to a whole category
	SKU      = 2 //applies to a single product SKU
	GENERAL  = 3 //applies to all products
)

type Discount interface {
	IsApplicableFor(item DiscountConditions) bool
	Apply(original int) int
	GetPercentage() int
	ToDiscountResponse() DiscountResponse
}

type DiscountConditions struct {
	CategoryID string
	SKU        string
}

// DiscountType represents the type of discount
type DiscountType struct {
	ID   int    `gorm:"primaryKey" json:"id" example:"1"`
	Type string `gorm:"unique;not null" json:"type" example:"category"`
}

// DiscountTypeRequest represents the body for creating a discount type
// @Description DiscountTypeRequest is the input for creating a new discount type
// @Accept json
// @Produce json
// @Param discount_type body DiscountTypeRequest true "Discount Type details"
type DiscountTypeRequest struct {
	Type string `json:"type" example:"category"`
}

// DiscountTypeResponse represents the response for a discount type
// @Description DiscountTypeResponse is the response structure for discount type details
// @Accept json
// @Produce json
// @Success 200 {object} DiscountTypeResponse
type DiscountTypeResponse struct {
	ID   string `json:"id" example:"1"`
	Type string `json:"type" example:"category"`
}

func (d *DiscountTypeRequest) ToDiscountType() DiscountType {
	return DiscountType{
		Type: d.Type,
	}
}

func (d *DiscountType) ToDiscountTypeResponse() DiscountTypeResponse {
	return DiscountTypeResponse{
		ID:   strconv.Itoa(d.ID),
		Type: d.Type,
	}
}

func (d *DiscountType) GetIdentifier() string {
	return strconv.Itoa(d.ID)
}

// GeneralDiscount represents the structure of a general discount
// @Description GeneralDiscount defines the fields for a general discount, including percentage and target
// @Accept json
// @Produce json
// @Success 200 {object} GeneralDiscount
type GeneralDiscount struct {
	ID             int          `gorm:"primaryKey" json:"id" example:"1"`
	Percentage     int          `gorm:"not null" json:"percentage" example:"10"`
	DiscountTypeID int          `gorm:"not null" json:"discount_type_id" example:"1"`
	DiscountType   DiscountType `gorm:"foreignKey:DiscountTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"discount_type"`
	Target         string       `gorm:"not null" json:"target" example:"boots"`
}

// DiscountRequest represents the body for creating a discount
// @Description DiscountRequest is the input for creating a new discount
// @Accept json
// @Produce json
// @Param discount body DiscountRequest true "Discount details"
type DiscountRequest struct {
	Percentage     int    `json:"percentage" example:"10"`
	DiscountTypeID int    `json:"discount_type_id" example:"1"`
	Target         string `json:"target" example:"boots"`
}

// DiscountResponse represents the output when retrieving discount details
// @Description DiscountResponse is the response structure when fetching discounts
// @Accept json
// @Produce json
// @Success 200 {object} DiscountResponse
type DiscountResponse struct {
	ID           string       `json:"id" example:"1"`
	Target       string       `json:"target" example:"boots"`
	DiscountType DiscountType `json:"discount_type"`
	Percentage   int          `json:"percentage" example:"10"`
}

func (d *DiscountRequest) ToDiscount() GeneralDiscount {
	return GeneralDiscount{
		Percentage:     d.Percentage,
		DiscountTypeID: d.DiscountTypeID,
		Target:         d.Target,
	}
}

func (d *GeneralDiscount) ToDiscountResponse() DiscountResponse {
	return DiscountResponse{
		ID:           strconv.Itoa(d.ID),
		Percentage:   d.Percentage,
		DiscountType: d.DiscountType,
		Target:       d.Target,
	}
}

func (d *GeneralDiscount) GetIdentifier() string {
	return strconv.Itoa(d.ID)
}

func (d *GeneralDiscount) Apply(original int) int {
	return original - (original * d.Percentage / 100)
}

func (d *GeneralDiscount) GetPercentage() int {
	return d.Percentage
}

func (d *GeneralDiscount) IsApplicableFor(item DiscountConditions) bool {
	return true
}

type CategoryDiscount struct {
	GeneralDiscount
}

func (d *CategoryDiscount) IsApplicableFor(item DiscountConditions) bool {
	return item.CategoryID == d.Target
}

type SkuDiscount struct {
	GeneralDiscount
}

func (d *SkuDiscount) IsApplicableFor(item DiscountConditions) bool {
	return item.SKU == d.Target
}
