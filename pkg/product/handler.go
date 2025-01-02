package product

import (
	"encoding/json"
	"mytheresa/internal/apierror"
	"mytheresa/internal/database"
	"mytheresa/internal/logger"
	"mytheresa/internal/response"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProduct(http.ResponseWriter, *http.Request)
	ListProducts(http.ResponseWriter, *http.Request)
}

type handler struct {
	service Service
	logger  logger.Logger
}

func NewHandler(service Service, logger logger.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Accept  json
// @Produce  json
// @Param product body ProductRequest true "Product details"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} apierror.ApiError "Wrong body"
// @Failure 500 {object} apierror.ApiError "Internal server error"
// @Router /v1/products [post]
func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var product ProductRequest
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		h.logger.WithError(err).Error(ctx, "Error decoding request while trying to create product")
		response.RespondWithError(w, apierror.BadRequest("Wrong body"))
		return
	}

	p, err := h.service.CreateProduct(ctx, product)
	if err != nil {
		h.logger.WithError(err).Error(ctx, "Error creating product")
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithData(w, http.StatusCreated, p)
}

// GetProduct godoc
// @Summary Get a product by SKU
// @Description Get the details of a product by its SKU
// @Produce  json
// @Param id path string true "Product SKU"
// @Success 200 {object} ProductResponse
// @Failure 404 {object} apierror.ApiError "Product not found"
// @Failure 500 {object} apierror.ApiError "Internal server error"
// @Router /v1/products/{id} [get]
func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	product, err := h.service.GetProduct(ctx, mux.Vars(r)["id"])
	if err != nil {
		h.logger.
			WithField("product_id", mux.Vars(r)["id"]).
			WithError(err).
			Error(ctx, "Error getting product")
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithData(w, http.StatusOK, product)
}

// ListProducts godoc
// @Summary List all products
// @Description Retrieve a list of products, with optional filtering by category and price range
// @Produce  json
// @Param limit query int false "Limit the number of products" default(5)
// @Param category query string false "Filter products by category ID"
// @Param priceLessThan query int false "Filter products with price less than"
// @Param priceGreaterThan query int false "Filter products with price greater than"
// @Success 200 {array} ProductResponse
// @Failure 500 {object} apierror.ApiError "Internal server error"
// @Router /v1/products [get]
func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := 5

	queryParams := r.URL.Query()
	if l, err := strconv.Atoi(queryParams.Get("limit")); queryParams.Get("limit") != "" && err == nil {
		limit = l
	}
	filters := createFilters(queryParams)

	products, err := h.service.ListProducts(ctx, filters...)
	if err != nil {
		h.logger.
			WithError(err).
			Error(ctx, "Error getting list of products")
		response.RespondWithError(w, err)
		return
	}
	if len(products) < limit {
		limit = len(products)
	}
	response.RespondWithData(w, http.StatusOK, products[:limit])
}

func createFilters(params url.Values) []database.Filter {
	var filters []database.Filter

	if p := params.Get("category"); p != "" {
		filters = append(filters, NewCategoryFilter(p, "="))
	}
	if p := params.Get("priceLessThan"); p != "" {
		filters = append(filters, NewPriceFilter(p, "<="))
	}
	if p := params.Get("priceGreaterThan"); p != "" {
		filters = append(filters, NewPriceFilter(p, ">="))
	}

	return filters
}
