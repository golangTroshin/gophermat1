package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	go service.StartPolling(db, time.Duration(config.Options.AccrualSystemInterval)*time.Second, config.Options.AccrualSystemWorkers)

	server := &http.Server{
		Addr:    config.Options.ServerAddress,
		Handler: getRouter(db),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", config.Options.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully.")
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

	r.Route("/api/user", func(r chi.Router) {
		r.Use(internal_middleware.AuthMiddleware)

		r.Post("/orders", orderHandler.UploadOrder)
		r.Get("/orders", orderHandler.GetOrders)
		r.Get("/balance", balanceHandler.GetUserBalance)
		r.Post("/balance/withdraw", withdrawHandler.Withdraw)
		r.Get("/withdrawals", withdrawHandler.GetWithdrawals)
	})

	return r
}
