package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CounterRequests = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "ozon",
		Name:      "counter_requests",
	})

	CounterRequestsByGroup = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ozon",
		Name:      "counter_requests_by_method",
	},
		[]string{
			"method",
		},
	)

	HistogramResponseServerTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ozon",
		Name:      "histogram_response_server_time_seconds",
		Buckets:   prometheus.DefBuckets,
	},
		[]string{
			"code_response",
			"method",
		},
	)

	HistogramResponseClientTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ozon",
		Name:      "histogram_response_client_time_seconds",
		Buckets:   prometheus.DefBuckets,
	},
		[]string{
			"code_response",
			"method",
		},
	)
)
