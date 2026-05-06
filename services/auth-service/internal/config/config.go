package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServiceName          string
	HTTPPort             string
	PostgresDSN          string
	RedisAddr            string
	RedisPassword        string
	RedisDB              int
	JWTAccessSecret      string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	ShutdownTimeout      time.Duration
	RefreshTokenKeyPrefix string
}

func Load() (Config, error) {
	redisDB, err := strconv.Atoi(env("REDIS_DB", "0"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	accessTTL, err := time.ParseDuration(env("ACCESS_TOKEN_TTL", "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid ACCESS_TOKEN_TTL: %w", err)
	}

	refreshTTL, err := time.ParseDuration(env("REFRESH_TOKEN_TTL", "168h"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid REFRESH_TOKEN_TTL: %w", err)
	}

	shutdownTimeout, err := time.ParseDuration(env("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}

	cfg := Config{
		ServiceName:          env("SERVICE_NAME", "auth-service"),
		HTTPPort:             env("HTTP_PORT", "8080"),
		PostgresDSN:          env("POSTGRES_DSN", ""),
		RedisAddr:            env("REDIS_ADDR", "localhost:6379"),
		RedisPassword:        env("REDIS_PASSWORD", ""),
		RedisDB:              redisDB,
		JWTAccessSecret:      env("JWT_ACCESS_SECRET", ""),
		AccessTokenTTL:       accessTTL,
		RefreshTokenTTL:      refreshTTL,
		ShutdownTimeout:       shutdownTimeout,
		RefreshTokenKeyPrefix: env("REFRESH_TOKEN_KEY_PREFIX", "rt:"),
	}

	if cfg.PostgresDSN == "" {
		return Config{}, fmt.Errorf("POSTGRES_DSN is required")
	}
	if cfg.JWTAccessSecret == "" {
		return Config{}, fmt.Errorf("JWT_ACCESS_SECRET is required")
	}

	return cfg, nil
}

func env(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
