package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"auth-service/internal/config"
	delivery "auth-service/internal/delivery/http"
	pgrepo "auth-service/internal/infrastructure/postgres"
	redisrepo "auth-service/internal/infrastructure/redis"
	"auth-service/internal/infrastructure/security"
	"auth-service/internal/observability"
	"auth-service/internal/usecase"
	"auth-service/migrations"

	"distributed-payment-system/shared/database"
	sharedMiddleware "distributed-payment-system/shared/middleware"

	_ "auth-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	goredis "github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate sh -c "cd ../.. && swag init -g cmd/auth-service/main.go -o docs --parseInternal --parseDependency"

// @title Auth Service API
// @version 1.0
// @description Authentication service API documentation.
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

	// Run migrations
	if err := database.RunMigrations(cfg.PostgresDSN, migrations.FS, "."); err != nil {
		slog.Error("run migrations failed", "error", err)
		os.Exit(1)
	}

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

	if os.Getenv("ENV") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(sharedMiddleware.CORS())
	router.Use(observability.LoggingMiddleware(cfg.ServiceName), observability.GinMiddleware(cfg.ServiceName))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	observability.SetServiceUp(cfg.ServiceName)
	handler := delivery.NewAuthHandler(authUC)
	handler.RegisterRoutes(router, cfg.JWTAccessSecret)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
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
