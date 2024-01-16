package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
)

// PrometheusMiddlewar , Update metrics
func PrometheusMiddleware(requestsTotal prometheus.Counter, cpuUsage prometheus.Gauge) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestsTotal.Inc()
		cpuUsage.Set(float64(runtime.GOMAXPROCS(0)))
		return c.Next()
	}
}
