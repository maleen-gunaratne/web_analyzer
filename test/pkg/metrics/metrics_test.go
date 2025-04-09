package metrics_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"

	"web-analyzer/pkg/metrics"
)

func TestMetricsInitialization(t *testing.T) {

	metrics.InitMetrics()

	metrics.Requests.WithLabelValues("success").Inc()
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.Requests.WithLabelValues("success")))

	metrics.LinksProcessed.Add(5)
	assert.Equal(t, float64(5), testutil.ToFloat64(metrics.LinksProcessed))

	metrics.BrokenLinks.Inc()
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.BrokenLinks))

	metrics.ActiveRequests.Set(3)
	assert.Equal(t, float64(3), testutil.ToFloat64(metrics.ActiveRequests))

	metrics.IncrementActiveRequests()
	assert.Equal(t, float64(4), testutil.ToFloat64(metrics.ActiveRequests))

	metrics.DecrementActiveRequests()
	assert.Equal(t, float64(3), testutil.ToFloat64(metrics.ActiveRequests))

}
