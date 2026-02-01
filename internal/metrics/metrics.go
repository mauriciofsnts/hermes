package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// EmailsSentTotal tracks the total number of emails sent, labeled by status and app.
var EmailsSentTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "hermes_emails_sent_total",
		Help: "Total number of emails sent",
	},
	[]string{"status", "app"},
)

// EmailsFailedTotal tracks the total number of failed email attempts.
var EmailsFailedTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "hermes_emails_failed_total",
		Help: "Total number of failed email attempts",
	},
	[]string{"reason", "app"},
)

// QueueDepth tracks the current depth of the notification queue.
var QueueDepth = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "hermes_queue_depth",
		Help: "Current depth of the notification queue",
	},
	[]string{"queue_type"},
)

// CircuitBreakerState tracks the state of circuit breakers (0=closed, 1=open, 2=half-open).
var CircuitBreakerState = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "hermes_circuit_breaker_state",
		Help: "State of the circuit breaker (0=closed, 1=open, 2=half-open)",
	},
	[]string{"provider"},
)

// CircuitBreakerFailures tracks the number of failures recorded by circuit breakers.
var CircuitBreakerFailures = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "hermes_circuit_breaker_failures_total",
		Help: "Total number of failures recorded by circuit breakers",
	},
	[]string{"provider"},
)

// SMTPConnectionDuration tracks the latency of SMTP operations in seconds.
var SMTPConnectionDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "hermes_smtp_duration_seconds",
		Help:    "SMTP operation duration in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"operation"},
)

// QueueProcessingDuration tracks the latency of queue processing in seconds.
var QueueProcessingDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "hermes_queue_processing_duration_seconds",
		Help:    "Queue processing duration in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"queue_type"},
)

// RequestDuration tracks the latency of HTTP requests in seconds.
var RequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "hermes_http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)

// APIKeyRateLimit tracks rate limit events by API key.
var APIKeyRateLimit = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "hermes_rate_limit_exceeded_total",
		Help: "Total number of rate limit exceeded events",
	},
	[]string{"api_key", "limit_type"},
)
