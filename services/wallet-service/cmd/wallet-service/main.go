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

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	goredis "github.com/redis/go-redis/v9"
)

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
	observability.SetServiceUp(cfg.ServiceName)
	h := delivery.NewHandler(walletUC)
	h.RegisterRoutes(r)

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
