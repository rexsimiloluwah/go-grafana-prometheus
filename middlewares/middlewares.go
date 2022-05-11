package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func RegisterPrometheusMetrics() {
	prometheus.Register(latency)
}

// Middleware for recording request latency
func RecordRequestLatency(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		res := next(c)
		elapsed := time.Since(start)

		latency.WithLabelValues(
			c.Request().Method,
			c.Request().URL.Path,
		).Observe(elapsed.Seconds())
		return res
	}
}

// Summary to record the latency for each request.
// The observed requests are labelled by the method and path.
var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "api",
		Name:       "latency_seconds",
		Help:       "RPC latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)
