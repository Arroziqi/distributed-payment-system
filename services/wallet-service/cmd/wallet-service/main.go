package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wallet-service/internal/config"
	delivery "wallet-service/internal/delivery/http"
	pgrepo "wallet-service/internal/infrastructure/postgres"
	redisrepo "wallet-service/internal/infrastructure/redis"
	"wallet-service/internal/observability"
	"wallet-service/internal/usecase"

	sharedMiddleware "distributed-payment-system/shared/middleware"

	_ "wallet-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	goredis "github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate sh -c "cd ../.. && swag init -g cmd/wallet-service/main.go -o docs --parseInternal --parseDependency"

// @title Wallet Service API
// @version 1.0
// @description Wallet service API documentation.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(slogHandler))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.Error("connect postgres failed", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		slog.Error("ping redis failed", "error", err)
		os.Exit(1)
	}

	walletRepo := pgrepo.NewWalletRepository(dbPool)
	cacheRepo := redisrepo.NewBalanceCacheRepository(redisClient)
	walletUC := usecase.NewWalletUsecase(walletRepo, cacheRepo, cfg.CacheTTL)

	if os.Getenv("ENV") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(sharedMiddleware.CORS())
	r.Use(observability.LoggingMiddleware(cfg.ServiceName), observability.GinMiddleware(cfg.ServiceName))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	observability.SetServiceUp(cfg.ServiceName)
	h := delivery.NewHandler(walletUC)
	h.RegisterRoutes(r, cfg.JWTAccessSecret)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		fmt.Printf("\n=================================\n")
		fmt.Printf("%s started\n", cfg.ServiceName)
		fmt.Printf("HTTP Port: %s\n", cfg.HTTPPort)
		fmt.Printf("Environment: %s\n", os.Getenv("ENV"))
		fmt.Printf("=================================\n\n")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen and serve failed", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("server gracefully stopped")
	}
}
