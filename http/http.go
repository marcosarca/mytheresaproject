package transport

import (
	"context"
	"fmt"
	"mytheresa/internal/logger"
	"mytheresa/pkg/discount"
	"mytheresa/pkg/product"
	"net/http"

	_ "mytheresa/docs"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHTTPRouter(l logger.Logger, ps product.Service, ds discount.Service) *mux.Router {

	ph := product.NewHandler(ps, l)
	dh := discount.NewHandler(ds, l)

	r := mux.NewRouter()
	r.Use(requestIDMiddleware)
	r.Use(contentTypeMiddleware)

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}).Methods(http.MethodGet)

	//Documentation
	r.HandleFunc("/swagger/{any:.*}", httpSwagger.WrapHandler).Methods(http.MethodGet)

	v1 := r.PathPrefix("/v1").Subrouter()

	//Product endpoints
	v1.HandleFunc("/product", ph.CreateProduct).Methods(http.MethodPost)
	v1.HandleFunc("/product/{id}", ph.GetProduct).Methods(http.MethodGet)
	v1.HandleFunc("/products", ph.ListProducts).Methods(http.MethodGet)
	//Discount endpoints
	v1.HandleFunc("/discount", dh.CreateDiscount).Methods(http.MethodPost)
	v1.HandleFunc("/discounts", dh.GetDiscounts).Methods(http.MethodGet)

	return r
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			// generate new version 4 uuid
			id = uuid.New().String()
		}
		// set the id to the request context
		ctx = context.WithValue(ctx, "request_id", id)
		r = r.WithContext(ctx)

		// set the response header
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}
