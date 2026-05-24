package exporter

import (
	"testing"
)

func TestFeatureServesMetricsThroughTemplate(t *testing.T) {
	t.Parallel()

	info := ExporterInfo()
	wantMetrics := append([]string{info.Metrics.BuildInfo}, info.Smoke.WantMetrics...)

	handler := newTestHandler(t)
	body := waitForHandlerMetrics(t, handler, wantMetrics)
	if body == "" {
		t.Fatal("waitForHandlerMetrics() returned empty body")
	}
}
