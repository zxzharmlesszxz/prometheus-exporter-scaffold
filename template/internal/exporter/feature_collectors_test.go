package exporter

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestFeatureRegistersCollector(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	registry := prometheus.NewRegistry()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v", err)
	}

	exportertest.WaitForMetricValue(t, registry, metricLastCollectionSuccess, nil, 1)
}

func TestFeatureReportsCollectorRegistrationError(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	exportertest.Register(t, registry, NewCollector(defaultMetricNamespace, testFeatureContext().Logger, SnapshotGatherer{}, time.Minute))

	feature := NewFeature()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err == nil {
		t.Fatal("RegisterCollectors() error = nil, want duplicate registration error")
	}
}
