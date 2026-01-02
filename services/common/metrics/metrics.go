package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// Business metrics - can be customized per service
	transactionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transactions_total",
			Help: "Total number of transactions",
		},
		[]string{"type", "status"},
	)

	transactionAmount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_amount",
			Help:    "Transaction amount distribution",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		},
		[]string{"type", "currency"},
	)

	activeUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of active users (with session)",
		},
	)

	dbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	externalAPILatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_latency_seconds",
			Help:    "External API call latency in seconds",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"api", "endpoint", "status"},
	)
)

// PrometheusMiddleware returns a Gin middleware for collecting HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics endpoint itself
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		httpRequestsInFlight.Inc()

		// Process request
		c.Next()

		httpRequestsInFlight.Dec()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Record metrics
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
	}
}

// MetricsHandler returns the Prometheus HTTP handler for /metrics endpoint
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Business metrics functions - call these from your service code

// RecordTransaction records a transaction metric
func RecordTransaction(txType, status string, amount float64, currency string) {
	transactionsTotal.WithLabelValues(txType, status).Inc()
	transactionAmount.WithLabelValues(txType, currency).Observe(amount)
}

// SetActiveUsers sets the current number of active users
func SetActiveUsers(count float64) {
	activeUsers.Set(count)
}

// SetDBConnections sets the current number of active DB connections
func SetDBConnections(count float64) {
	dbConnectionsActive.Set(count)
}

// RecordExternalAPICall records an external API call latency
func RecordExternalAPICall(api, endpoint, status string, duration time.Duration) {
	externalAPILatency.WithLabelValues(api, endpoint, status).Observe(duration.Seconds())
}

// Custom service-specific metrics

// NewCounter creates a new counter metric
func NewCounter(name, help string, labels []string) *prometheus.CounterVec {
	return promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
}

// NewGauge creates a new gauge metric
func NewGauge(name, help string) prometheus.Gauge {
	return promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
	)
}

// NewHistogram creates a new histogram metric
func NewHistogram(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: buckets,
		},
		labels,
	)
}
