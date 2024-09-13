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
	"gorm.io/gorm"

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

	if err := http.ListenAndServe(config.Options.ServerAddress, getRouter(db)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getRouter(db *gorm.DB) chi.Router {
	userRepo := repository.NewUserRepository(db)
	balanceRepo := repository.NewUserBalanceRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	withdrawRepo := repository.NewWithdrawRepository(db)

	balanceService := service.NewBalanceService(userRepo, balanceRepo)
	authService := service.NewAuthService(userRepo, balanceRepo)
	orderService := service.NewOrderService(orderRepo)
	withdrawService := service.NewWithdrawService(withdrawRepo)

	authHandler := handler.NewAuthHandler(authService)
	orderHandler := handler.NewOrderHandler(orderService)
	balanceHandler := handler.NewBalanceHandler(balanceService)
	withdrawHandler := handler.NewWithdrawHandler(balanceService, withdrawService, orderService)

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

	return r
}
