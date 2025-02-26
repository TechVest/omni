//nolint:unused // This is a work in progress.
package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	statusOffset = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "solver",
		Subsystem: "processor",
		Name:      "status_offset",
		Help:      "Last inbox offset processed by chain and status",
	}, []string{"chain", "target", "status"})

	processedEvents = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "solver",
		Subsystem: "processor",
		Name:      "processed_events_total",
		Help:      "Total number of events processed by chain and status",
	}, []string{"chain", "target", "status"})

	rejectedOrders = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "solver",
		Subsystem: "processor",
		Name:      "rejected_orders_total",
		Help:      "Total number of rejected orders by chain and reason",
	}, []string{"src_chain", "dest_chain", "target", "reason"})

	tokenBalance = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "solver",
		Subsystem: "liquidity",
		Name:      "token_balance",
		Help:      "Token balance of solver",
	}, []string{"chain", "solver_addr", "token_addr", "token_symbol", "is_native"})

	apiLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "solver",
		Subsystem: "api",
		Name:      "latency",
		Help:      "API server request latency in seconds per endpoint",
	}, []string{"endpoint"})

	apiResponses = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "solver",
		Subsystem: "api",
		Name:      "response_total",
		Help:      "Total responses served by the API server per endpoint per status code class (2xx, 4xx, 5xx)",
	}, []string{"endpoint", "class"})

	apiConcurrent = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "solver",
		Subsystem: "api",
		Name:      "concurrent_requests",
		Help:      "Number of concurrent requests being served by the API server (at scrape time)",
	})
)
