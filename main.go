package main

import (
	"context"
	"fmt"
	"log"
	transport "mytheresa/http"
	"mytheresa/internal/config"
	"mytheresa/internal/database/sqlite"
	"mytheresa/internal/logger"
	"mytheresa/pkg/category"
	"mytheresa/pkg/discount"
	"mytheresa/pkg/product"
	"net/http"
	"os"
	"os/signal"
	"time"

	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "mytheresa/docs"
)

func main() {

	conf := config.New()

	l := logger.NewLogger("mytheresa")
	defer l.Sync()

	// Cleaning DB for a fresh start everytime
	_ = os.Remove(conf.DbFile)

	db, err := gorm.Open(gormsqlite.Open(conf.DbFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sql := sqlite.NewSQLiteDB(db, l)

	err = sql.MigrateModels(
		&product.Product{},
		&category.Category{},
		&discount.DiscountType{},
		&discount.GeneralDiscount{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	cs := category.NewService(sql, l)
	ds := discount.NewService(sql, l)
	ps := product.NewService(sql, l, ds)

	insertInitialData(cs, ps, ds)

	httpTransportRouter := transport.NewHTTPRouter(l, ps, ds)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", conf.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      httpTransportRouter,
	}

	l.WithField("transport", "http").WithField("port", conf.Port).
		Info(context.Background(), "Transport Start")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.WithField(
				"transport", "http").
				WithError(err).
				Info(context.Background(), "Transport Stopped")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	l.Info(context.Background(), "Service gracefully shutted down")
	os.Exit(0)
}

func insertInitialData(cs category.Service, ps product.Service, ds discount.Service) {
	ctx := context.Background()
	c1, _ := cs.CreateCategory(ctx, category.CategoryRequest{
		Name: "boots",
	})

	c2, _ := cs.CreateCategory(ctx, category.CategoryRequest{
		Name: "sandals",
	})

	c3, _ := cs.CreateCategory(ctx, category.CategoryRequest{
		Name: "sneakers",
	})

	p1 := product.ProductRequest{
		SKU:        "000001",
		Name:       "BV Lean leather ankle boots",
		CategoryID: c1.ID,
		Price:      89000,
	}
	_, _ = ps.CreateProduct(ctx, p1)

	p2 := product.ProductRequest{
		SKU:        "000002",
		Name:       "BV Lean leather ankle boots",
		CategoryID: c1.ID,
		Price:      99000,
	}
	_, _ = ps.CreateProduct(ctx, p2)

	p3 := product.ProductRequest{
		SKU:        "000003",
		Name:       "Ashlington leather ankle boots",
		CategoryID: c1.ID,
		Price:      71000,
	}
	_, _ = ps.CreateProduct(ctx, p3)

	p4 := product.ProductRequest{
		SKU:        "000004",
		Name:       "Naima embellished suede sandals",
		CategoryID: c2.ID,
		Price:      79500,
	}
	_, _ = ps.CreateProduct(ctx, p4)

	p5 := product.ProductRequest{
		SKU:        "000005",
		Name:       "Nathane leather sneakers",
		CategoryID: c3.ID,
		Price:      59000,
	}
	_, _ = ps.CreateProduct(ctx, p5)

	dt1, _ := ds.CreateDiscountType(ctx, discount.DiscountTypeRequest{
		Type: "category",
	})

	dt2, _ := ds.CreateDiscountType(ctx, discount.DiscountTypeRequest{
		Type: "sku",
	})

	dt3, _ := ds.CreateDiscountType(ctx, discount.DiscountTypeRequest{
		Type: "general",
	})

	//-----Discounts-----
	_, _ = ds.CreateDiscount(ctx, discount.DiscountRequest{
		DiscountTypeID: dt1.ID,
		Target:         fmt.Sprint(c1.ID),
		Percentage:     30,
	})

	_, _ = ds.CreateDiscount(ctx, discount.DiscountRequest{
		DiscountTypeID: dt2.ID,
		Target:         p3.SKU,
		Percentage:     15,
	})

	_, _ = ds.CreateDiscount(ctx, discount.DiscountRequest{
		DiscountTypeID: dt3.ID,
		Percentage:     0,
	})

}
