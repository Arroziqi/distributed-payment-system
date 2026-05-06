package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"auth-service/internal/config"
	delivery "auth-service/internal/delivery/http"
	"auth-service/internal/observability"
	pgrepo "auth-service/internal/infrastructure/postgres"
	redisrepo "auth-service/internal/infrastructure/redis"
	"auth-service/internal/infrastructure/security"
	"auth-service/internal/usecase"

	_ "auth-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	goredis "github.com/redis/go-redis/v9"
)

//go:generate sh -c "cd ../.. && swag init -g cmd/auth-service/main.go -o docs --parseInternal --parseDependency"

// @title Auth Service API
// @version 1.0
// @description Authentication service API documentation.
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

	userRepo := pgrepo.NewUserRepository(dbPool)
	refreshRepo := redisrepo.NewRefreshTokenRepository(redisClient, cfg.RefreshTokenKeyPrefix)
	hasher := security.NewPasswordHasher(0)
	tokenManager := security.NewTokenManager(cfg.JWTAccessSecret, cfg.AccessTokenTTL)

	authUC := usecase.NewAuthUsecase(
		userRepo,
		refreshRepo,
		hasher,
		tokenManager,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), observability.GinMiddleware(cfg.ServiceName))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	observability.SetServiceUp(cfg.ServiceName)
	handler := delivery.NewAuthHandler(authUC)
	handler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
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
