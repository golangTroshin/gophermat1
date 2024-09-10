package main

import (
	"log"
	"net/http"
	"time"

	"github.com/golangTroshin/gophermat/internal/config"
	"github.com/golangTroshin/gophermat/internal/db"
	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/repository"
	"github.com/golangTroshin/gophermat/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	internal_middleware "github.com/golangTroshin/gophermat/internal/middleware"
)

func main() {
	if err := config.ParseFlags(); err != nil {
		log.Fatalf("error ocured while parsing flags: %v", err)
	}

	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	go service.StartPolling(db, 5*time.Second, 2)

	userRepo := repository.NewUserRepository(db)
	balanceRepo := repository.NewUserBalanceRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	withdrawRepo := repository.NewWithdrawRepository(db)

	userService := service.NewUserService(userRepo, balanceRepo, withdrawRepo)
	orderService := service.NewOrderService(orderRepo)

	authHandler := handler.NewAuthHandler(userService)
	orderHandler := handler.NewOrderHandler(orderService)
	balanceHandler := handler.NewBalanceHandler(userService)
	withdrawHandler := handler.NewWithdrawHandler(userService, orderService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * 1000000000))

	r.Use(middleware.Compress(5))

	r.Post("/api/user/register", authHandler.Register)
	r.Post("/api/user/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(internal_middleware.AuthMiddleware)

		r.Post("/api/user/orders", orderHandler.UploadOrder)
		r.Get("/api/user/orders", orderHandler.GetOrders)
		r.Get("/api/user/balance", balanceHandler.GetUserBalance)
		r.Post("/api/user/balance/withdraw", withdrawHandler.Withdraw)
		r.Get("/api/user/withdrawals", withdrawHandler.GetWithdrawals)
	})

	log.Fatal(http.ListenAndServe(config.Options.ServerAddress, r))
}
