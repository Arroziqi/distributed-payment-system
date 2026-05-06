package observability

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_http_requests_total",
		Help: "Total number of HTTP requests processed.",
	}, []string{"service", "method", "route", "status"})
	requestsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_http_requests_failed_total",
		Help: "Total number of failed HTTP requests (status >= 500).",
	}, []string{"service", "method", "route"})
	requestLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "service_http_request_latency_seconds",
		Help:    "HTTP request latency distribution in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method", "route"})
	serviceUp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "service_up",
		Help: "Service liveness status (1 = up).",
	}, []string{"service"})
)

func SetServiceUp(service string) {
	serviceUp.WithLabelValues(service).Set(1)
}

func GinMiddleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		status := strconv.Itoa(c.Writer.Status())
		latency := time.Since(start).Seconds()

		requestsTotal.WithLabelValues(service, c.Request.Method, route, status).Inc()
		requestLatency.WithLabelValues(service, c.Request.Method, route).Observe(latency)
		if c.Writer.Status() >= 500 {
			requestsFailed.WithLabelValues(service, c.Request.Method, route).Inc()
		}
	}
}
