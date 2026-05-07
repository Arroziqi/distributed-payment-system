package observability

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(service string) gin.HandlerFunc {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "docker"
	}

	return func(c *gin.Context) {
		start := time.Now()
		reqID := uuid.New().String()
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)

		path := c.Request.URL.Path

		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"service_name", service,
					"environment", environment,
					"error", err,
					"path", path,
					"request_id", reqID,
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()

		latency := time.Since(start).String()
		status := c.Writer.Status()

		if path == "/metrics" || path == "/healthz" || strings.HasPrefix(path, "/swagger/") {
			return
		}

		args := []any{
			"service_name", service,
			"environment", environment,
			"request_id", reqID,
			"method", c.Request.Method,
			"path", path,
			"status_code", status,
			"latency", latency,
		}

		if status >= 500 {
			slog.Error("request failed", args...)
		} else if status >= 400 {
			slog.Warn("client error", args...)
		} else {
			slog.Info("request processed", args...)
		}
	}
}
