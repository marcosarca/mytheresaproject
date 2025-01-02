package discount

import (
	"encoding/json"
	"mytheresa/internal/apierror"
	"mytheresa/internal/logger"
	"mytheresa/internal/response"
	"net/http"
)

type Handler interface {
	CreateDiscount(w http.ResponseWriter, r *http.Request)
	GetDiscounts(w http.ResponseWriter, r *http.Request)
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

// CreateDiscount godoc
// @Summary Create a new discount
// @Description Create a new discount with the provided details
// @Accept  json
// @Produce  json
// @Param discount body DiscountRequest true "Discount details"
// @Success 201 {object} DiscountResponse
// @Failure 400 {object} apierror.ApiError "Wrong body"
// @Failure 500 {object} apierror.ApiError "Internal server error"
// @Router /v1/discounts [post]
func (h handler) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var discount DiscountRequest
	err := json.NewDecoder(r.Body).Decode(&discount)
	if err != nil {
		h.logger.WithError(err).Error(ctx, "Error decoding request while trying to create discount")
		response.RespondWithError(w, apierror.BadRequest("Wrong body"))
		return
	}

	d, err := h.service.CreateDiscount(ctx, discount)
	if err != nil {
		h.logger.WithError(err).Error(ctx, "Error creating discount")
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithData(w, http.StatusCreated, d.ToDiscountResponse())
}

// GetDiscounts godoc
// @Summary Get all discounts
// @Description Retrieve a list of all available discounts
// @Produce  json
// @Success 200 {array} GeneralDiscount
// @Failure 500 {object} apierror.ApiError "Internal server error"
// @Router /v1/discounts [get]
func (h handler) GetDiscounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	discounts, err := h.service.GetDiscounts(ctx)
	if err != nil {
		h.logger.WithError(err).Error(ctx, "Error getting discounts")
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithData(w, http.StatusOK, discounts)
}
