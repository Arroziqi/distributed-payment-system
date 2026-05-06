package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	ServiceName     string
	HTTPPort        string
	PostgresDSN     string
	RabbitMQURL     string
	Exchange        string
	Queue           string
	RoutingKey      string
	ShutdownTimeout time.Duration
	JWTAccessSecret string
}

func Load() (Config, error) {
	shutdownTimeout, err := time.ParseDuration(env("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg := Config{
		ServiceName:     env("SERVICE_NAME", "notification-service"),
		HTTPPort:        env("HTTP_PORT", "8080"),
		PostgresDSN:     env("POSTGRES_DSN", ""),
		RabbitMQURL:     env("RABBITMQ_URL", ""),
		Exchange:        env("RABBITMQ_EXCHANGE", "payments.events"),
		Queue:           env("RABBITMQ_QUEUE", "notification.transaction.success"),
		RoutingKey:      env("RABBITMQ_ROUTING_KEY", "transaction.completed"),
		ShutdownTimeout: shutdownTimeout,
		JWTAccessSecret: env("JWT_ACCESS_SECRET", "dev-access-secret"),
	}
	if cfg.PostgresDSN == "" {
		return Config{}, fmt.Errorf("POSTGRES_DSN is required")
	}
	if cfg.RabbitMQURL == "" {
		return Config{}, fmt.Errorf("RABBITMQ_URL is required")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
