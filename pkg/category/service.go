package category

import (
	"context"
	"mytheresa/internal/apierror"
	"mytheresa/internal/database"
	"mytheresa/internal/logger"
)

type Service interface {
	CreateCategory(ctx context.Context, category CategoryRequest) (Category, error)
}

type service struct {
	db     database.Database
	logger logger.Logger
}

func NewService(db database.Database, logger logger.Logger) Service {
	return &service{db: db, logger: logger}
}

func (s *service) CreateCategory(ctx context.Context, req CategoryRequest) (Category, error) {
	category := req.ToCategory()
	err := s.db.Save(ctx, category.GetIdentifier(), &category)
	if err != nil {
		s.logger.WithError(err).Error(ctx, "failed to save category")
		return Category{}, apierror.InternalServerError("there was an error saving the category")
	}

	return category, nil
}
