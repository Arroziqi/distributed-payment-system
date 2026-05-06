package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"transaction-service/internal/config"
	delivery "transaction-service/internal/delivery/http"
	"transaction-service/internal/observability"
	pgrepo "transaction-service/internal/infrastructure/postgres"
	rmq "transaction-service/internal/infrastructure/rabbitmq"
	"transaction-service/internal/usecase"

	_ "transaction-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate sh -c "cd ../.. && swag init -g cmd/transaction-service/main.go -o docs --parseInternal --parseDependency"

// @title Transaction Service API
// @version 1.0
// @description Transaction service API documentation.
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

	rabbitConn, err := dialRabbitWithRetry(cfg.RabbitMQURL, 20, 2*time.Second)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer rabbitConn.Close()
	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("create rabbitmq channel: %v", err)
	}
	defer rabbitCh.Close()

	if err := rabbitCh.ExchangeDeclare(cfg.RabbitExchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("declare exchange: %v", err)
	}

	txRepo := pgrepo.NewTransactionRepository(dbPool)
	pub := rmq.NewPublisher(rabbitCh, cfg.RabbitExchange)
	txUC := usecase.NewTransactionUsecase(txRepo, pub, cfg.IdempotencyTTL)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), observability.GinMiddleware(cfg.ServiceName))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	observability.SetServiceUp(cfg.ServiceName)
	h := delivery.NewHandler(txUC)
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

func dialRabbitWithRetry(url string, maxAttempts int, delay time.Duration) (*amqp.Connection, error) {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err := amqp.Dial(url)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		log.Printf("rabbitmq dial attempt %d/%d failed: %v", attempt, maxAttempts, err)
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("rabbitmq dial failed after %d attempts: %w", maxAttempts, lastErr)
}
