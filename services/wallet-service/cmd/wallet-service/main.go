package main

import (
	"context"
	"log"
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

	_ "wallet-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	goredis "github.com/redis/go-redis/v9"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer dbPool.Close()

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("ping redis: %v", err)
	}

	walletRepo := pgrepo.NewWalletRepository(dbPool)
	cacheRepo := redisrepo.NewBalanceCacheRepository(redisClient)
	walletUC := usecase.NewWalletUsecase(walletRepo, cacheRepo, cfg.CacheTTL)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), observability.GinMiddleware(cfg.ServiceName))
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
		log.Printf("%s listening on :%s", cfg.ServiceName, cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
