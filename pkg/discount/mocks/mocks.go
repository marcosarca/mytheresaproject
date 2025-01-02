package mocks

import (
	"context"
	"mytheresa/pkg/discount"

	"github.com/stretchr/testify/mock"
)

type Service struct {
	mock.Mock
}

func (s *Service) CreateDiscountType(ctx context.Context, discountType discount.DiscountTypeRequest) (discount.DiscountType, error) {
	args := s.Called(ctx, discountType)
	return args.Get(0).(discount.DiscountType), args.Error(1)
}

func (s *Service) CreateDiscount(ctx context.Context, d discount.DiscountRequest) (discount.Discount, error) {
	args := s.Called(ctx, d)
	return args.Get(0).(discount.Discount), args.Error(1)
}

func (s *Service) GetDiscounts(ctx context.Context) ([]discount.Discount, error) {
	args := s.Called(ctx)
	return args.Get(0).([]discount.Discount), args.Error(1)
}
