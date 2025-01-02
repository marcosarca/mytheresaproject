package mocks

import (
	"context"
	"mytheresa/pkg/category"

	"github.com/stretchr/testify/mock"
)

type Service struct {
	mock.Mock
}

func (s *Service) CreateCategory(ctx context.Context, c category.CategoryRequest) (category.Category, error) {
	args := s.Called(ctx, c)
	return args.Get(0).(category.Category), args.Error(1)
}
