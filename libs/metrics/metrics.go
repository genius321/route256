package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ozon",
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   []float64{},
	},
		[]string{
			"status",
			"method",
		},
	)
)
