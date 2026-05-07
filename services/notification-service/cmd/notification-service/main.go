package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification-service/internal/config"
	"notification-service/internal/consumer"
	delivery "notification-service/internal/delivery/http"
	"notification-service/internal/observability"
	"notification-service/internal/service"

	_ "notification-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate sh -c "cd ../.. && swag init -g cmd/notification-service/main.go -o docs --parseInternal --parseDependency"

// @title Notification Service API
// @version 1.0
// @description Notification service API documentation.
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.Error("connect postgres failed", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	rabbitConn, err := dialRabbitWithRetry(cfg.RabbitMQURL, 10, 2*time.Second)
	if err != nil {
		slog.Error("connect rabbitmq failed", "error", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		slog.Error("create rabbitmq channel failed", "error", err)
		os.Exit(1)
	}
	defer rabbitCh.Close()

	svc := service.NewNotificationService(dbPool)
	cons := consumer.NewRabbitConsumer(rabbitCh, cfg.Queue, cfg.Exchange, cfg.RoutingKey, svc)
	if err := cons.Start(ctx); err != nil {
		slog.Error("start consumer failed", "error", err)
		os.Exit(1)
	}

	if os.Getenv("ENV") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(observability.LoggingMiddleware(cfg.ServiceName), observability.GinMiddleware(cfg.ServiceName))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	observability.SetServiceUp(cfg.ServiceName)
	h := delivery.NewHandler()
	h.RegisterRoutes(router, cfg.JWTAccessSecret)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		fmt.Printf("\n=================================\n")
		fmt.Printf("%s started\n", cfg.ServiceName)
		fmt.Printf("HTTP Port: %s\n", cfg.HTTPPort)
		fmt.Printf("RabbitMQ: connected\n")
		fmt.Printf("Environment: %s\n", os.Getenv("ENV"))
		fmt.Printf("=================================\n\n")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen and serve failed", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("server gracefully stopped")
	}
}

func dialRabbitWithRetry(url string, maxAttempts int, delay time.Duration) (*amqp.Connection, error) {
	var lastErr error
	slog.Info("connecting to rabbitmq...")
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err := amqp.Dial(url)
		if err == nil {
			slog.Info("rabbitmq connected successfully")
			return conn, nil
		}
		lastErr = err
		slog.Warn(fmt.Sprintf("rabbitmq not ready, retrying in %v...", delay), "attempt", attempt, "maxAttempts", maxAttempts, "error", err)
		time.Sleep(delay)
		delay = time.Duration(float64(delay) * 1.5)
	}
	return nil, fmt.Errorf("rabbitmq dial failed after %d attempts: %w", maxAttempts, lastErr)
}
