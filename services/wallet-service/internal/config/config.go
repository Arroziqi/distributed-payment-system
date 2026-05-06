package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServiceName     string
	HTTPPort        string
	PostgresDSN     string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	CacheTTL        time.Duration
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	redisDB, err := strconv.Atoi(env("REDIS_DB", "0"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	cacheTTL, err := time.ParseDuration(env("CACHE_TTL", "5m"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid CACHE_TTL: %w", err)
	}
	shutdownTimeout, err := time.ParseDuration(env("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}

	cfg := Config{
		ServiceName:     env("SERVICE_NAME", "wallet-service"),
		HTTPPort:        env("HTTP_PORT", "8080"),
		PostgresDSN:     env("POSTGRES_DSN", ""),
		RedisAddr:       env("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   env("REDIS_PASSWORD", ""),
		RedisDB:         redisDB,
		CacheTTL:        cacheTTL,
		ShutdownTimeout: shutdownTimeout,
	}
	if cfg.PostgresDSN == "" {
		return Config{}, fmt.Errorf("POSTGRES_DSN is required")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
