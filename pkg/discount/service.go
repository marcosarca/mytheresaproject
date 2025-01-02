package discount

import (
	"context"
	"mytheresa/internal/apierror"
	"mytheresa/internal/database"
	"mytheresa/internal/logger"
)

type Service interface {
	CreateDiscountType(ctx context.Context, discountType DiscountTypeRequest) (DiscountType, error)
	CreateDiscount(ctx context.Context, discount DiscountRequest) (Discount, error)
	GetDiscounts(ctx context.Context) ([]Discount, error)
}

type service struct {
	db     database.Database
	logger logger.Logger
}

func NewService(db database.Database, logger logger.Logger) Service {
	return &service{
		db,
		logger,
	}
}

func (s *service) CreateDiscountType(ctx context.Context, req DiscountTypeRequest) (DiscountType, error) {
	discountType := req.ToDiscountType()
	err := s.db.Save(ctx, discountType.GetIdentifier(), &discountType)
	if err != nil {
		s.logger.WithError(err).Error(ctx, "error creating discount type")
		return DiscountType{}, apierror.InternalServerError("error creating discount type")
	}
	return discountType, nil
}

func (s *service) CreateDiscount(ctx context.Context, req DiscountRequest) (Discount, error) {
	discount := req.ToDiscount()
	err := s.db.Save(ctx, discount.GetIdentifier(), &discount)
	if err != nil {
		s.logger.
			WithError(err).
			Error(ctx, "error creating discount")
		return &GeneralDiscount{}, apierror.InternalServerError("error creating discount")
	}

	return &discount, nil
}

func (s *service) GetDiscounts(ctx context.Context) ([]Discount, error) {
	var discounts []GeneralDiscount
	results := []Discount{}

	err := s.db.GetWithFilters(ctx, &discounts)
	if err != nil {
		s.logger.WithError(err).Error(ctx, "error getting discounts")
		return nil, apierror.InternalServerError("error getting discounts")
	}

	for _, d := range discounts {
		switch d.DiscountTypeID {
		case CATEGORY:
			results = append(results, &CategoryDiscount{d})
			break
		case SKU:
			results = append(results, &SkuDiscount{d})
			break
		default:
			results = append(results, &d)

		}
	}

	return results, nil
}
