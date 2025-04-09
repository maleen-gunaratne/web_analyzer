package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Requests = prometheus.NewCounterVec(

		prometheus.CounterOpts{
			Namespace: "web_analyzer",
			Name:      "requests_total",
			Help:      "Total number of page analysis requests",
		},
		[]string{"status"},
	)

	LinksProcessed = promauto.NewCounter(prometheus.CounterOpts{

		Namespace: "web_analyzer",
		Name:      "links_processed_total",
		Help:      "Total number of links processed",
	})

	BrokenLinks = promauto.NewCounter(prometheus.CounterOpts{

		Namespace: "web_analyzer",
		Name:      "broken_links_total",
		Help:      "Total number of broken links found",
	})

	AnalysisTime = promauto.NewHistogram(prometheus.HistogramOpts{

		Namespace: "web_analyzer",
		Name:      "analysis_duration_seconds",
		Help:      "Time taken to analyze pages",
		Buckets:   prometheus.DefBuckets,
	})

	HTTPResponseCodes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "web_analyzer",
			Name:      "http_response_codes_total",
			Help:      "HTTP response codes encountered",
		},
		[]string{"code"},
	)

	ActiveRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "web_analyzer",
		Name:      "active_requests",
		Help:      "Number of currently processing requests",
	})
)

func InitMetrics() { //  registers all metrics with Prometheus
	prometheus.MustRegister(Requests)
}

func IncrementActiveRequests() {
	ActiveRequests.Inc()
}

func DecrementActiveRequests() {
	ActiveRequests.Dec()
}
