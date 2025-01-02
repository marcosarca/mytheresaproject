package mocks

import (
	"context"
	"mytheresa/internal/database"
	"mytheresa/pkg/product"

	"github.com/stretchr/testify/mock"
)

type Service struct {
	mock.Mock
}

func (s *Service) CreateProduct(ctx context.Context, p product.ProductRequest) (product.Product, error) {
	args := s.Called(ctx, p)
	return args.Get(0).(product.Product), args.Error(1)
}

func (s *Service) GetProduct(ctx context.Context, id string) (product.Product, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(product.Product), args.Error(1)
}

func (s *Service) ListProducts(ctx context.Context, filters ...database.Filter) ([]product.ProductResponse, error) {
	args := s.Called(ctx, filters)
	return args.Get(0).([]product.ProductResponse), args.Error(1)
}
