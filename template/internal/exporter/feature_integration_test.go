package exporter

import (
	"testing"
)

func TestFeatureServesMetricsThroughTemplate(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)
	body := waitForHandlerMetrics(t, handler, []string{
		metricBuildInfo,
		metricLastCollectionSuccess + " 1",
		metricExampleValue + " 1",
	})
	if body == "" {
		t.Fatal("waitForHandlerMetrics() returned empty body")
	}
}
