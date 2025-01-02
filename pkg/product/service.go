package product

import (
	"context"
	"errors"
	"fmt"
	"mytheresa/internal/apierror"
	"mytheresa/internal/database"
	"mytheresa/internal/logger"
	"mytheresa/pkg/discount"

	"gorm.io/gorm"
)

type Service interface {
	CreateProduct(ctx context.Context, product ProductRequest) (Product, error)
	GetProduct(ctx context.Context, id string) (Product, error)
	ListProducts(ctx context.Context, filters ...database.Filter) ([]ProductResponse, error)
}

type service struct {
	db              database.Database
	logger          logger.Logger
	discountService discount.Service
}

func NewService(db database.Database, logger logger.Logger, ds discount.Service) Service {
	return &service{
		db:              db,
		logger:          logger,
		discountService: ds,
	}
}

func (s *service) CreateProduct(ctx context.Context, req ProductRequest) (Product, error) {
	product := req.ToProduct()
	err := s.db.Save(ctx, product.GetIdentifier(), &product)
	if err != nil {
		msg := fmt.Sprintf("Error creating product: %s", product.Name)
		s.logger.Error(ctx, msg)
		return Product{}, apierror.InternalServerError(msg)
	}
	return product, nil
}

func (s *service) GetProduct(ctx context.Context, id string) (Product, error) {
	var product Product
	err := s.db.Get(ctx, id, &product)
	if err != nil {
		s.logger.
			WithField("id", id).
			WithError(err).
			Error(ctx, "error getting product from DB")

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Product{}, apierror.NotFound("Product not found")
		}

		return Product{}, apierror.InternalServerError(fmt.Sprintf("Error getting Product with ID %s", id))
	}
	return product, nil
}

func (s *service) ListProducts(ctx context.Context, filters ...database.Filter) ([]ProductResponse, error) {
	var products []Product

	s.logger.WithField("filters", filters).Info(ctx, "Listing products")

	err := s.db.GetWithFilters(ctx, &products, filters...)
	if err != nil {
		s.logger.WithError(err).Error(ctx, "Failed to get products from database")
		return nil, apierror.InternalServerError(fmt.Sprintf("Failed to get products from database"))
	}

	response, err := s.getProductResponseWithDiscounts(ctx, products)
	if err != nil {
		return nil, err
	}

	s.logger.WithField("quantity", len(products)).Info(ctx, "Successfully retrieved products")
	return response, nil
}

func (s *service) getProductResponseWithDiscounts(ctx context.Context, products []Product) ([]ProductResponse, error) {
	discounts, err := s.discountService.GetDiscounts(ctx)
	if err != nil {
		s.logger.WithError(err).Error(ctx, "Failed to get discounts from database")
		return nil, err
	}

	response := []ProductResponse{}
	for _, p := range products {
		pr := p.ToProductResponse()

		item := discount.DiscountConditions{
			CategoryID: fmt.Sprint(p.CategoryID),
			SKU:        p.SKU,
		}

		//apply greater discount
		for _, d := range discounts {
			if d.IsApplicableFor(item) {
				candidate := d.Apply(pr.Price.Original)
				if candidate < pr.Price.Final {
					percentaje := fmt.Sprint(d.GetPercentage())
					pr.Price = PriceResponse{
						Original:           pr.Price.Original,
						Final:              candidate,
						DiscountPercentage: &percentaje,
						Currency:           "EUR",
					}
				}
			}
		}
		response = append(response, pr)
	}
	return response, nil
}
